package tracing

type Storage interface {
	Store(trace string) error
}
