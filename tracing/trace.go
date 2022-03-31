package tracing

import (
	"time"
)

type Trace struct {
	Timestamp     time.Time
	Message       string
	CorrelationId string
	Metrics       map[string]float64
	Level         string
}
