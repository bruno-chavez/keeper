package handler

import (
	"fmt"
	"strings"
)

// parseEncryptionContext parses comma-separated "key:value" pairs into a map.
// An empty input returns a nil map (no encryption context). Only the first colon
// of each pair separates key from value, so values may contain colons.
func parseEncryptionContext(raw string) (map[string]string, error) {
	if raw == "" {
		return nil, nil
	}

	result := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		key, value, found := strings.Cut(pair, ":")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if !found || key == "" || value == "" {
			return nil, fmt.Errorf("invalid context pair %q, expected key:value", pair)
		}

		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("duplicate context key %q", key)
		}

		result[key] = value
	}

	return result, nil
}
