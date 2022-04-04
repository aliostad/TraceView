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

func TestParser_clef_full(t *testing.T) {
	json := `{"@t":"2016-11-21T11:22:33Z","@m":"hello","@mt":"hellomtt","CorrelationId":"12345","@l":"debug","foo":"sumagh","bar":2}`
	parser := PayloadParser{}
	trc, err := parser.Parse(json)
	assert.Nil(t, err)
	assert.Equal(t, "hello", trc.Message)
	assert.Equal(t, 11, trc.Timestamp.Hour())
	assert.Equal(t, "debug", trc.Level)
	assert.Equal(t, "", trc.CorrelationId)
	assert.Equal(t, 3, len(trc.Properties))
	assert.Equal(t, 1, len(trc.Metrics))
	assert.Equal(t, "sumagh", trc.Properties["foo"])
	assert.Equal(t, 2.0, trc.Metrics["bar"]) // golang reads json number as float64
}

func TestParser_clef_doesnt_fail_no_message(t *testing.T) {
	json := `{"@t":"2016-11-21T11:22:33Z","CorrelationId":"12345","@l":"infos","foo":"sumagh","bar":2}`
	parser := PayloadParser{}
	_, err := parser.Parse(json)
	assert.Nil(t, err)
}

func TestParser_nonclef_no_config(t *testing.T) {
	json := `{"Timestamp":"2016-11-21T11:22:33Z","message": "I was here!","CorrelationId":"12345","severity":"infos","foo":"sumagh","bar":2}`
	parser := PayloadParser{}
	trc, err := parser.Parse(json)
	assert.Nil(t, err)
	assert.Equal(t, 11, trc.Timestamp.Hour())
	assert.Equal(t, "infos", trc.Level)
	assert.Equal(t, "", trc.CorrelationId)
	assert.Equal(t, "I was here!", trc.Message)
	assert.Equal(t, 2, len(trc.Properties))
	assert.Equal(t, 1, len(trc.Metrics))
	assert.Equal(t, "sumagh", trc.Properties["foo"])
	assert.Equal(t, 2.0, trc.Metrics["bar"]) // golang reads json number as float64
}

func TestParser_nonclef_full_config(t *testing.T) {
	json := `{"Timestampi":"2016-11-21T11:22:33Z","mensaje": "I was here!","corrrId":"12345","suvirity":"infos","foo":"sumagh","bar":2}`
	parser := PayloadParser{
		config: Config{
			TimestampFieldNames:     []string{"Timestampi"},
			MessageFieldNames:       []string{"mensaje"},
			LevelFieldNames:         []string{"suvirity"},
			CorrelationIdFieldNames: []string{"corrrId"},
		},
	}
	trc, err := parser.Parse(json)
	assert.Nil(t, err)
	assert.Equal(t, 11, trc.Timestamp.Hour())
	assert.Equal(t, "infos", trc.Level)
	assert.Equal(t, "12345", trc.CorrelationId)
	assert.Equal(t, "I was here!", trc.Message)
	assert.Equal(t, 1, len(trc.Properties))
	assert.Equal(t, 1, len(trc.Metrics))
	assert.Equal(t, "sumagh", trc.Properties["foo"])
	assert.Equal(t, 2.0, trc.Metrics["bar"]) // golang reads json number as float64
}

func TestParser_nonclef_full_config_message_not_found(t *testing.T) {
	json := `{"Timestampi":"2016-11-21T11:22:33Z","message": "I was here!","corrrId":"12345","suvirity":"infos","foo":"sumagh","bar":2}`
	parser := PayloadParser{
		config: Config{
			TimestampFieldNames:     []string{"Timestampi"},
			MessageFieldNames:       []string{"mensaje"},
			LevelFieldNames:         []string{"suvirity"},
			CorrelationIdFieldNames: []string{"corrrId"},
		},
	}

	_, err := parser.Parse(json)
	assert.NotNil(t, err)

}

func TestParser_nonclef_full_config_level_not_found(t *testing.T) {
	json := `{"Timestampi":"2016-11-21T11:22:33Z","message": "I was here!","corrrId":"12345","severity":"infos","foo":"sumagh","bar":2}`
	parser := PayloadParser{
		config: Config{
			TimestampFieldNames:     []string{"Timestampi"},
			MessageFieldNames:       []string{},
			LevelFieldNames:         []string{"suvirity"},
			CorrelationIdFieldNames: []string{"corrrId"},
		},
	}

	trc, err := parser.Parse(json)
	assert.Nil(t, err)
	assert.Equal(t, 11, trc.Timestamp.Hour())
	assert.Equal(t, "info", trc.Level)
	assert.Equal(t, "12345", trc.CorrelationId)
	assert.Equal(t, "I was here!", trc.Message)
	assert.Equal(t, 2, len(trc.Properties))
	assert.Equal(t, 1, len(trc.Metrics))
	assert.Equal(t, "sumagh", trc.Properties["foo"])
	assert.Equal(t, 2.0, trc.Metrics["bar"]) // golang reads json number as float64
}
func TestParser_clef_default_level_with_integer_level(t *testing.T) {
	json := `{"@t":"2016--21T11:22:33Z","CorrelationId":"12345","@l":4,"foo":"sumagh","bar":2}`
	parser := PayloadParser{}
	trc, err := parser.Parse(json)
	assert.Nil(t, err)
	assert.Equal(t, "info", trc.Level)
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
