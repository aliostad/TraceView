package tracing

type TraceStore interface {
	Store(trace Trace, originalPayload string) error
}

type InMemoryStore struct {
	config *Config
}

func (store *InMemoryStore) Store(trace Trace, originalPayload string) error {
	return nil
}
