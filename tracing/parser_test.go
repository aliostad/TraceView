package tracing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_canLoadJson(t *testing.T) {
	jsonString := `{"Timestamp":"2016-01-01T00:00:00Z","Message":"hello","CorrelationId":"12345","Metrics":{"foo":1,"bar":2},"Level":"info"}`
	parser := PayloadParser{}
	parser.Parse(jsonString)
}

func TestParser_string(t *testing.T) {
	parser := PayloadParser{}
	trace, err := parser.Parse("hello")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "hello", trace.Message)
	assert.Equal(t, "info", trace.Level)
	assert.Equal(t, time.Now().Day(), trace.Timestamp.Day())

}

func TestParser_badjson(t *testing.T) {
	parser := PayloadParser{}
	_, err := parser.Parse("{hello")

	assert.NotNil(t, err)
}

func Test_isDate(t *testing.T) {
	assert.Nil(t, getSecondParam(parseDate("2016-01-01T00:00:00Z")))
	assert.Nil(t, getSecondParam(parseDate("Fri, 01 Apr 2022 19:29:21 GMT")))
	assert.Nil(t, getSecondParam(parseDate(1351700038.0))) // golang reads json number as float64

	assert.NotNil(t, getSecondParam(parseDate("hello")))
}

func getSecondParam(one interface{}, two interface{}) interface{} {
	return two
}
