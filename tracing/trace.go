package tracing

import (
	"time"
)

type Trace struct {
	Timestamp     time.Time
	Message       string
	CorrelationId string
	Level         string
	Metrics       map[string]float64
	Properties    map[string]string
}
