package tracing

import (
	"strconv"
	"time"

	"github.com/hashicorp/go-memdb"
)

type InMemoryStore struct {
	config   *Config
	payloads map[string]string
	db       *memdb.MemDB
}

const (
	tableName             = "Trace"
	id_index              = "id"
	timestamp_index       = "timestamp_idx"
	id_column_name        = "TraceId"
	timestamp_column_name = "TimeIndex"
	max_return            = 100
)

func NewInMemoryStore(config *Config) (*InMemoryStore, error) {
	db, err := memdb.NewMemDB(getSchema())
	if err != nil {
		return nil, err
	}

	return &InMemoryStore{
		config:   config,
		payloads: make(map[string]string),
		db:       db,
	}, nil
}

func (store *InMemoryStore) Store(trace *Trace, originalPayload string) error {
	if store.config.KeepOriginalPayload {
		store.payloads[trace.TraceId] = originalPayload
	}
	txn := store.db.Txn(true)
	txn.Insert("Trace", trace)
	txn.Commit()
	return nil
}

// returns null if not found
func (store *InMemoryStore) GetById(id string) (*Trace, error) {
	txn := store.db.Txn(false)
	defer txn.Abort()
	trace, err := txn.First(tableName, id_index, id)
	if err != nil {
		return nil, err
	}

	if trace == nil {
		return nil, nil
	}

	return trace.(*Trace), nil
}

func (store *InMemoryStore) ListByTimeRange(n int, from, to *time.Time, exclusive bool) ([]*Trace, error) {
	reverse := false
	if from == nil && to != nil {
		reverse = true
	}

	if from == nil {
		fromX := time.Unix(0, 0)
		from = &fromX
	}

	if to == nil {
		toX := time.Now().UTC()
		from = &toX
	}

	txn := store.db.Txn(false)
	defer txn.Abort()
	var iter memdb.ResultIterator
	var err error
	if reverse {
		iter, err = txn.ReverseLowerBound(tableName, timestamp_index, strconv.FormatInt((*to).UnixMicro(), 10))
	} else {
		iter, err = txn.LowerBound(tableName, timestamp_index, strconv.FormatInt((*from).UnixMicro(), 10))
	}

	if err != nil {
		return nil, err
	}

	traces := make([]*Trace, 0)

	for obj := iter.Next(); obj != nil; obj = iter.Next() {
		trc := obj.(*Trace)

		if trc.Timestamp.After(*to) {
			break
		}

		if exclusive && (trc.Timestamp.Equal(*to) || trc.Timestamp.Equal(*from)) {
			continue
		}

		traces = append(traces, trc)
		if len(traces) >= min(max_return, n) {
			break
		}
	}

	if reverse {
		for i, j := 0, len(traces)-1; i < j; i, j = i+1, j-1 {
			traces[i], traces[j] = traces[j], traces[i]
		}
	}

	return traces, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getSchema() *memdb.DBSchema {

	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"Trace": {
				Name: tableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    id_index,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: id_column_name},
					},
					"timestamp_idx": {
						Name:    timestamp_index,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: timestamp_column_name},
					},
				},
			},
		},
	}
}
