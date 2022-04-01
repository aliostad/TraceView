package tracing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/maps"
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

		return parseJson(payload, result)
	}

	return Trace{
		Message:   payload,
		Timestamp: time.Now().UTC(),
		Level:     "info",
	}, nil

}

func parseJson(payload string, jsonMap map[string]interface{}) (Trace, error) {
	keys := maps.Keys(jsonMap)
	fmt.Println(keys)
	return Trace{}, nil
}

func parseClef(payload string, jsonMap map[string]interface{}) (Trace, error) {

	return Trace{}, nil
}
