package tracing

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func Parse(payload string) (Trace, error) {

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

		return ParseJson(payload, result)
	}

	return Trace{
		Message:   payload,
		Timestamp: time.Now().UTC(),
		Level:     "info",
	}, nil

}

func ParseJson(payload string, jsonMap map[string]interface{}) (Trace, error) {
	return Trace{}, nil
}

func ParseClef(payload string, jsonMap map[string]interface{}) (Trace, error) {
	return Trace{}, nil
}
