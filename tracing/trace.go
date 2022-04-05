package tracing

import (
	"time"

	"github.com/google/uuid"
)

type Trace struct {
	TraceId       string
	Timestamp     time.Time
	Message       string
	CorrelationId string
	Level         string
	Metrics       map[string]float64
	Properties    map[string]string
}

func NewTrace(ts time.Time, message string, corrId string, level string) Trace {
	return Trace{
		TraceId:       uuid.New().String(),
		Timestamp:     ts.UTC(),
		Message:       message,
		CorrelationId: corrId,
		Level:         level,
		Metrics:       make(map[string]float64),
		Properties:    make(map[string]string),
	}

}
