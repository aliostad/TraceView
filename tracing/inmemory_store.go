package tracing

import (
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
	trace, err := txn.First(tableName, id_index, id)
	if err != nil {
		return nil, err
	}

	if trace == nil {
		return nil, nil
	}

	return trace.(*Trace), nil
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
