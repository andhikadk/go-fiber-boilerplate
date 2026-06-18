package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

const redactedValue = "[REDACTED]"

var sensitiveKeys = map[string]struct{}{
	"password":         {},
	"old_password":     {},
	"new_password":     {},
	"confirm_password": {},
	"token":            {},
	"access_token":     {},
	"refresh_token":    {},
	"authorization":    {},
	"secret":           {},
	"api_key":          {},
	"card":             {},
	"card_number":      {},
	"cvv":              {},
	"pin":              {},
	"signature":        {},
}

func RedactJSONValue(raw []byte, maxBytes int) any {
	if len(raw) == 0 {
		return nil
	}
	if maxBytes > 0 && len(raw) > maxBytes {
		return fmt.Sprintf("[truncated %d bytes]", len(raw))
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "[unparseable json body]"
	}
	return redactValue(parsed)
}

func redactValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			if _, sensitive := sensitiveKeys[strings.ToLower(k)]; sensitive {
				val[k] = redactedValue
				continue
			}
			val[k] = redactValue(child)
		}
		return val
	case []any:
		for i, child := range val {
			val[i] = redactValue(child)
		}
		return val
	default:
		return v
	}
}
