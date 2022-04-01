package tracing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_canLoadJson(t *testing.T) {
	jsonString := `{"Timestamp":"2016-01-01T00:00:00Z","Message":"hello","CorrelationId":"12345","Metrics":{"foo":1,"bar":2},"Level":"info"}`
	Parse(jsonString)
}

func TestParser_string(t *testing.T) {
	trace, err := Parse("hello")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "hello", trace.Message)
	assert.Equal(t, "info", trace.Level)
	assert.Equal(t, time.Now().Day(), trace.Timestamp.Day())

}

func TestParser_badjson(t *testing.T) {
	_, err := Parse("{hello")

	assert.NotNil(t, err)
}
