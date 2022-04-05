package tracing

import "time"

type TraceStore interface {
	Store(trace *Trace, originalPayload string) error
	GetById(id string) (*Trace, error)
	GetByTimeRange(from, to time.Time) ([]*Trace, error)
}
