package utils

import (
	"encoding/json"
	"fmt"
)

// ParseJSON safely parses a JSON byte slice into a map
func ParseJSON(body []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}
