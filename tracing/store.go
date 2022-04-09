package tracing

import "time"

type TraceStore interface {
	Store(trace *Trace, originalPayload string) error
	GetById(id string) (*Trace, error)
	ListByTimeRange(n int, from, to *time.Time) ([]*Trace, error)
}
