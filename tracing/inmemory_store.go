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

func (store *InMemoryStore) ListByTimeRange(n int, from, to *time.Time) ([]*Trace, error) {
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
	iter, err := txn.LowerBound(tableName, timestamp_index, strconv.FormatInt((*from).UnixMicro(), 10))
	if err != nil {
		return nil, err
	}

	traces := make([]*Trace, 0)

	i := 0
	for obj := iter.Next(); obj != nil; obj = iter.Next() {
		trc := obj.(*Trace)
		i++
		if i > max_return || i > n || trc.Timestamp.After(*to) {
			break
		}

		traces = append(traces, trc)
	}

	return traces, nil
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
