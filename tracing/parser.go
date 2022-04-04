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
	return Trace{}, nil
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
	delete(jsonMap, "@m")
	if message == "" {
		message = safeGetValue(jsonMap, "@mt")
		delete(jsonMap, "@mt")
	}

	corrId := safeGetValue(jsonMap, "CorrelationId")
	delete(jsonMap, "CorrelationId")
	if corrId == "" {
		corrId = safeGetValue(jsonMap, "corrId")
		delete(jsonMap, "corrId")
	}

	trc := Trace{
		Timestamp:     timestamp,
		Message:       message,
		CorrelationId: corrId,
		Level:         safeGetValue(jsonMap, "@l"), // not acceptting integer level
		Metrics:       make(map[string]float64),
		Properties:    make(map[string]string),
	}

	delete(jsonMap, "@l")

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

	return trc, nil
}

func safeGetValue(jsonMap map[string]interface{}, key string) string {
	if value, ok := jsonMap[key]; ok {
		return value.(string)
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
