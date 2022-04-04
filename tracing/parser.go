package tracing

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type PayloadParser struct {
	config Config
}

func (parser *PayloadParser) Parse(payload string) (Trace, error) {
	if payload == "" {
		return Trace{}, errors.New("empty payload")
	}

	payload = strings.TrimSpace(payload)
	if strings.HasPrefix(payload, "{") {
		var result map[string]interface{}
		err := json.Unmarshal([]byte(payload), &result)
		if err != nil {
			return Trace{}, err
		}

		return parser.parseJson(payload, result)
	}

	return Trace{
		Message:   payload,
		Timestamp: time.Now().UTC(),
		Level:     "info",
	}, nil

}

func (parser *PayloadParser) parseJson(payload string, jsonMap map[string]interface{}) (Trace, error) {
	keys := maps.Keys(jsonMap)
	if slices.Contains(keys, "@t") && isDate(jsonMap["@t"].(string)) {
		return parser.parseClef(payload, jsonMap)
	}

	var timestamp time.Time
	var message string
	var level string
	var corrId string
	var err error

	// ______________________ TIMESTAMP ______________________
	// if not defined, try to guess the timestamp field
	if len(parser.config.TimestampFieldNames) == 0 {
		var timestampFieldName string
		timestampFieldName, timestamp = findTimestampField(jsonMap, "timestamp", "Timestamp", "time", "Time",
			"date", "Date", "datetime", "DateTime", "eventDate", "EventDate")
		if timestampFieldName == "" {
			timestamp = time.Now().UTC()
		} else {
			delete(jsonMap, timestampFieldName)
		}
	} else {
		// try to find the first timestamp field that has date value
		for _, fieldName := range parser.config.TimestampFieldNames {
			if value, ok := jsonMap[fieldName]; ok {
				timestamp, err = parseDate(value)
				if err == nil {
					delete(jsonMap, fieldName)
					break
				} else {
					continue
				}
			}
		}

		// if not found throw error
		if timestamp.IsZero() {
			return Trace{}, errors.New("no valid timestamp field found")
		}
	}

	// ______________________ MESSAGE ______________________
	messageFieldName, fromConfig := findStringField(jsonMap, parser.config.MessageFieldNames,
		"message", "Message", "Description", "description",
		"Text", "text", "Error", "error", "ErrorText", "errorText", "errorText")
	if messageFieldName == "" {
		if fromConfig {
			return Trace{}, errors.New("defined message field names could not be found")
		}
	} else {
		message = jsonMap[messageFieldName].(string)
		delete(jsonMap, messageFieldName)
	}

	// ______________________ LEVEL ______________________
	levelFieldName, _ := findStringField(jsonMap, parser.config.LevelFieldNames, "level", "Level", "severity", "Severity")
	if levelFieldName == "" {
		level = "info"
	} else {
		level = jsonMap[levelFieldName].(string)
		delete(jsonMap, levelFieldName)
	}

	corrIdFieldName, _ := findStringField(jsonMap, parser.config.CorrelationIdFieldNames) // we don't waste time on this to try other things
	if corrIdFieldName != "" {
		corrId = jsonMap[corrIdFieldName].(string)
		delete(jsonMap, corrIdFieldName)
	}

	trc := Trace{
		Timestamp:     timestamp,
		Message:       message,
		Level:         level,
		CorrelationId: corrId,
		Metrics:       make(map[string]float64),
		Properties:    make(map[string]string),
	}
	populatePropertiesAndMetrics(jsonMap, &trc)
	return trc, nil
}

func findStringField(jsonMap map[string]interface{}, configFieldNames []string, fieldNames ...string) (string, bool) {
	if len(configFieldNames) > 0 {
		for _, fieldName := range configFieldNames {
			if value, ok := jsonMap[fieldName]; ok {
				if value.(string) != "" {
					return fieldName, true
				}
			}
		}
	} else {
		for _, fieldName := range fieldNames {
			if value, ok := jsonMap[fieldName]; ok {
				_, success := value.(string)
				if success {
					return fieldName, false
				}
			}
		}
		return "", false
	}
	return "", true
}

func findTimestampField(jsonMap map[string]interface{}, fieldNames ...string) (string, time.Time) {

	for _, fieldName := range fieldNames {
		if value, ok := jsonMap[fieldName]; ok {
			dt, err := parseDate(value)
			if err == nil {
				return fieldName, dt
			}
		}
	}

	return "", time.Time{}
}

/*
@t	Timestamp	An ISO 8601 timestamp	Yes
@m	Message	A fully-rendered message describing the event
@mt	Message template	Alternative to Message; specifies a message template over the event’s properties that provides for rendering into a textual description of the event
@l	Level	An implementation-specific level or severity identifier (string or number)	Absence implies “informational”
@x	Exception	A language-dependent error representation potentially including backtrace
@i	Event id	An implementation specific event id (string or number)
@r	Renderings	If @mt includes tokens with programming-language-specific formatting, an array of pre-rendered values for each such token	May be omitted; if present, the count of renderings must match the count of formatted tokens exactly

*/
func (parser *PayloadParser) parseClef(payload string, jsonMap map[string]interface{}) (Trace, error) {
	timestamp, err := parseDate(jsonMap["@t"])
	delete(jsonMap, "@t")
	if err != nil {
		panic("invalid timestamp while it was supposed to work")
	}

	message := safeGetValue(jsonMap, "@m")
	if message == "" {
		message = safeGetValue(jsonMap, "@mt")
		if message != "" {
			delete(jsonMap, "@mt")
		}
	} else {
		delete(jsonMap, "@m")
	}

	level := safeGetValue(jsonMap, "@l")
	if level == "" {
		level = "info"
	} else {
		delete(jsonMap, "@l")
	}

	var corrId string
	corrIdFieldName, _ := findStringField(jsonMap, parser.config.CorrelationIdFieldNames) // we don't waste time on this to try other things
	if corrIdFieldName != "" {
		corrId = jsonMap[corrIdFieldName].(string)
		delete(jsonMap, corrIdFieldName)
	}

	trc := Trace{
		Timestamp:     timestamp,
		Message:       message,
		CorrelationId: corrId,
		Level:         level,
		Metrics:       make(map[string]float64),
		Properties:    make(map[string]string),
	}

	populatePropertiesAndMetrics(jsonMap, &trc)

	return trc, nil
}

func populatePropertiesAndMetrics(jsonMap map[string]interface{}, trc *Trace) {
	for key, value := range jsonMap {
		s, success := value.(string)
		if success {
			trc.Properties[key] = s
		} else {
			f, success := value.(float64)
			if success {
				trc.Metrics[key] = f
			}
		}
	}
}

func safeGetValue(jsonMap map[string]interface{}, key string) string {
	if value, ok := jsonMap[key]; ok {
		s, success := value.(string)
		if success {
			return s
		}
	}
	return ""
}

func isDate(s interface{}) bool {
	_, err := parseDate(s)
	return err == nil
}

// parses date in typical formats and epoch
func parseDate(s interface{}) (time.Time, error) {
	var dt time.Time
	i, success := s.(float64)
	if success {
		return getDateFromEpoch(int64(i))
	}

	dt, err := time.Parse(time.RFC3339, s.(string))
	if err == nil {
		return dt, nil
	}

	dt, err = time.Parse(time.RFC850, s.(string))
	if err == nil {
		return dt, nil
	}

	dt, err = time.Parse(time.RFC1123, s.(string))
	if err == nil {
		return dt, nil
	}

	return time.Time{}, err
}

func getDateFromEpoch(epoch int64) (time.Time, error) {
	s := strconv.FormatInt(epoch, 10)
	switch len(s) {
	case 10:
		return time.Unix(epoch, 0), nil
	case 13:
		return time.UnixMilli(epoch), nil
	case 16:
		return time.UnixMicro(epoch), nil
	case 19:
		return time.UnixMicro(epoch / 1000), nil
	default:
		return time.Time{}, errors.New("invalid epoch")
	}
}
