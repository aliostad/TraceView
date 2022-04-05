package tracing

type TraceStore interface {
	Store(trace Trace, originalPayload string) error
}
