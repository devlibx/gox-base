package goxJsonUtils

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var (
	piiMaskingRegistry = make(map[string]func(string) string)
	registryLock       = &sync.RWMutex{}
)

// RegisterPiiMasker registers a function to mask a PII field.
// For example, to mask an email, you can register a function like:
//
//	goxJsonUtils.RegisterPiiMasker("email", func(email string) string {
//		// Your logic to mask the email
//		return "masked-email"
//	})
func RegisterPiiMasker(key string, masker func(string) string) {
	registryLock.Lock()
	defer registryLock.Unlock()
	piiMaskingRegistry[key] = masker
}

// PrettyStringLoggingSuppressError converts an input object to a pretty-printed JSON string,
// with PII masking applied.
// It suppresses errors and returns a string representation of the input object even if marshalling or masking fails.
func PrettyStringLoggingSuppressError(in any) string {
	// 1. Marshal the input to bytes
	bytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Sprintf("%v", in)
	}

	// 2. Unmarshal into a generic interface{}
	var data interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return string(bytes) // return original json
	}

	// 3. Mask the data
	maskedData := maskPii(data)

	// 4. Marshal back to pretty json
	prettyJSON, err := json.MarshalIndent(maskedData, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v", maskedData)
	}

	return string(prettyJSON)
}

func maskPii(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			registryLock.RLock()
			masker, exists := piiMaskingRegistry[key]
			registryLock.RUnlock()

			if exists {
				if strVal, ok := val.(string); ok {
					v[key] = applyMasker(masker, strVal)
				} else {
					strVal := fmt.Sprintf("%v", val)
					v[key] = applyMasker(masker, strVal)
				}
			} else {
				v[key] = maskPii(val)
			}
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = maskPii(val)
		}
		return v
	default:
		return data
	}
}

func applyMasker(masker func(string) string, value string) (maskedValue string) {
	defer func() {
		if r := recover(); r != nil {
			// If the masker panics, replace with asterisks
			maskedValue = strings.Repeat("*", len(value))
		}
	}()
	return masker(value)
}
