package tracing

import (
	"github.com/hashicorp/go-memdb"
)

type InMemoryStore struct {
	config   *Config
	payloads map[string]Trace
	db       *memdb.MemDB
}

func NewInMemoryStore(config *Config) (*InMemoryStore, error) {
	db, err := memdb.NewMemDB(getSchema())
	if err != nil {
		return nil, err
	}

	return &InMemoryStore{
		config:   config,
		payloads: make(map[string]Trace),
		db:       db,
	}, nil
}

func (store *InMemoryStore) Store(trace *Trace, originalPayload string) error {
	txn := store.db.Txn(true)
	txn.Insert("Trace", trace)
	txn.Commit()
	return nil
}

func (store *InMemoryStore) GetById(id string) (*Trace, error) {
	txn := store.db.Txn(false)
	trace, err := txn.First("Trace", "id", id)
	if err != nil {
		return nil, err
	}

	return trace.(*Trace), nil
}

func getSchema() *memdb.DBSchema {

	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"Trace": {
				Name: "Trace",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "TraceId"},
					},
					"timestamp_idx": {
						Name:    "timestamp_idx",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "TimeIndex"},
					},
				},
			},
		},
	}
}
