package gox

import (
	"encoding/json"
	"github.com/devlibx/gox-base/serialization"
	"reflect"
	"strconv"
)

type StringObjectMap map[string]interface{}

func (s StringObjectMap) Int(name string) (int, bool) {
	switch value := s[name].(type) {
	case int:
		return value, true
	case int32:
		return int(value), true
	case int64:
		return int(value), true
	case uint32:
		return int(value), true
	case uint64:
		return int(value), true
	case float32:
		return int(value), true
	case float64:
		return int(value), true
	case string:
		if val, err := strconv.Atoi(value); err != nil {
			return 0, false
		} else {
			return val, true
		}
	}
	return 0, false
}

func (s StringObjectMap) IntOrDefault(name string, defaultValue int) int {
	if value, ok := s.Int(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) Int32(name string) (int32, bool) {
	switch value := s[name].(type) {
	case int:
		return int32(value), true
	case int32:
		return int32(value), true
	case int64:
		return int32(value), true
	case uint32:
		return int32(value), true
	case uint64:
		return int32(value), true
	case float32:
		return int32(value), true
	case float64:
		return int32(value), true
	case string:
		if val, err := strconv.Atoi(value); err != nil {
			return 0, false
		} else {
			return int32(val), true
		}
	}
	return 0, false
}

func (s StringObjectMap) Int32OrDefault(name string, defaultValue int32) int32 {
	if value, ok := s.Int32(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) Int64(name string) (int64, bool) {
	switch value := s[name].(type) {
	case int:
		return int64(value), true
	case int32:
		return int64(value), true
	case int64:
		return int64(value), true
	case uint32:
		return int64(value), true
	case uint64:
		return int64(value), true
	case float32:
		return int64(value), true
	case float64:
		return int64(value), true
	case string:
		if val, err := strconv.Atoi(value); err != nil {
			if n, err := strconv.ParseInt(value, 10, 64); err != nil {
				return 0, false
			} else {
				return n, true
			}
		} else {
			return int64(val), true
		}
	}
	return 0, false
}

func (s StringObjectMap) Int64OrDefault(name string, defaultValue int64) int64 {
	if value, ok := s.Int64(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) Float32(name string) (float32, bool) {
	switch value := s[name].(type) {
	case int:
		return float32(value), true
	case int32:
		return float32(value), true
	case int64:
		return float32(value), true
	case uint32:
		return float32(value), true
	case uint64:
		return float32(value), true
	case float32:
		return float32(value), true
	case float64:
		return float32(value), true
	case string:
		if val, err := strconv.ParseFloat(value, 32); err != nil {
			return 0, false
		} else {
			return float32(val), true
		}
	}
	return 0, false
}

func (s StringObjectMap) Float32OrDefault(name string, defaultValue float32) float32 {
	if value, ok := s.Float32(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) Float64(name string) (float64, bool) {
	switch value := s[name].(type) {
	case int:
		return float64(value), true
	case int32:
		return float64(value), true
	case int64:
		return float64(value), true
	case uint32:
		return float64(value), true
	case uint64:
		return float64(value), true
	case float32:
		return float64(value), true
	case float64:
		return float64(value), true
	case string:
		if val, err := strconv.ParseFloat(value, 64); err != nil {
			return 0, false
		} else {
			return float64(val), true
		}
	}
	return 0, false
}

func (s StringObjectMap) Float64OrDefault(name string, defaultValue float64) float64 {
	if value, ok := s.Float64(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) Bool(name string) (bool, bool) {
	switch value := s[name].(type) {
	case int:
		return value == 1, true
	case int32:
		return value == 1, true
	case int64:
		return value == 1, true
	case uint32:
		return value == 1, true
	case uint64:
		return value == 1, true
	case float32:
		return value == 1, true
	case float64:
		return value == 1, true
	case string:
		if val, err := strconv.ParseBool(value); err != nil {
			return false, false
		} else {
			return val, true
		}
	}
	return false, false
}

func (s StringObjectMap) BoolOrDefault(name string, defaultValue bool) bool {
	if value, ok := s.Bool(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) String(name string) (string, bool) {
	switch value := s[name].(type) {
	case string:
		return value, true
	default:
		if val, err := serialization.Stringify(value); err != nil {
			return "", false
		} else {
			return val, true
		}
	}
}

func (s StringObjectMap) StringOrDefault(name string, defaultValue string) string {
	if value, ok := s.String(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) StringOrEmpty(name string) string {
	return s.StringOrDefault(name, "")
}

type NewObjectFunc func() interface{}

func (s StringObjectMap) Object(name string, obj interface{}) (interface{}, bool) {
	switch value := s[name].(type) {
	case string:
		if err := json.Unmarshal([]byte(value), obj); err != nil {
			return nil, false
		} else {
			return obj, true
		}

	case []byte:
		if err := json.Unmarshal(value, obj); err != nil {
			return nil, false
		} else {
			return obj, true
		}

	default:
		if reflect.TypeOf(obj) == reflect.TypeOf(value) {
			return value, true
		} else {
			return nil, false
		}
	}
}

func (s StringObjectMap) ObjectOrDefault(name string, obj interface{}, defaultValue interface{}) interface{} {
	if value, ok := s.Object(name, obj); ok {
		return value
	}
	return defaultValue
}
