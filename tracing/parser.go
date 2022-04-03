package tracing

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println(keys)
	if slices.Contains(keys, "@t") && isDate(jsonMap["@t"].(string)) {
		return parser.parseClef(payload, jsonMap)
	}
	return Trace{}, nil
}

func (parser *PayloadParser) parseClef(payload string, jsonMap map[string]interface{}) (Trace, error) {
	return Trace{}, nil
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
