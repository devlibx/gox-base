package gox

import (
	"encoding/json"
	"github.com/devlibx/gox-base/errors"
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

func (s StringObjectMap) IntOrZero(name string) int {
	return s.IntOrDefault(name, 0)
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

func (s StringObjectMap) Float64OrZero(name string) float64 {
	return s.Float64OrDefault(name, 0)
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
	case bool:
		return value, true
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

func (s StringObjectMap) BoolOrFalse(name string) bool {
	return s.BoolOrDefault(name, false)
}

func (s StringObjectMap) BoolOrTrue(name string) bool {
	return s.BoolOrDefault(name, true)
}

func (s StringObjectMap) BoolOrDefault(name string, defaultValue bool) bool {
	if value, ok := s.Bool(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) MapOrDefault(name string, defaultValue map[string]interface{}) map[string]interface{} {
	if value, ok := s[name].(map[string]interface{}); ok {
		return value
	} else if value, ok := s[name].(StringObjectMap); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) MapOrEmpty(name string) map[string]interface{} {
	if value, ok := s[name].(map[string]interface{}); ok {
		return value
	} else if value, ok := s[name].(StringObjectMap); ok {
		return value
	} else {
		return map[string]interface{}{}
	}
}

func (s StringObjectMap) StringObjectMapOrDefault(name string, defaultValue map[string]interface{}) StringObjectMap {
	if value, ok := s[name].(map[string]interface{}); ok {
		return value
	} else if value, ok := s[name].(StringObjectMap); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) StringObjectMapOrEmpty(name string) StringObjectMap {
	if value, ok := s[name].(map[string]interface{}); ok {
		return value
	} else if value, ok := s[name].(StringObjectMap); ok {
		return value
	} else {
		return map[string]interface{}{}
	}
}

func (s StringObjectMap) String(name string) (string, bool) {
	switch value := s[name].(type) {
	case string:
		return value, true
	default:
		if value == nil {
			return "", false
		} else if val, err := serialization.Stringify(value); err != nil {
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

func (s StringObjectMap) Bytes(name string) ([]byte, bool) {
	switch value := s[name].(type) {
	case []byte:
		return value, true
	case string:
		return []byte(value), true
	default:
		if value == nil {
			return nil, false
		} else if val, err := serialization.Stringify(value); err != nil {
			return nil, false
		} else {
			return []byte(val), true
		}
	}
}

func (s StringObjectMap) BytesOrDefault(name string, defaultValue []byte) []byte {
	if value, ok := s.Bytes(name); ok {
		return value
	} else {
		return defaultValue
	}
}

func (s StringObjectMap) BytesOrEmpty(name string) []byte {
	if value, ok := s.Bytes(name); ok {
		return value
	} else {
		return []byte{}
	}
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

func StringObjectMapFromString(input string) (StringObjectMap, error) {
	out := StringObjectMap{}
	if err := serialization.JsonToObject(input, &out); err != nil {
		return nil, errors.Wrap(err, "failed to read map from input string")
	}
	return out, nil
}

// Convert input json string to StringObjectMap
func StringObjectMapFromJson(input string) (StringObjectMap, error) {
	out := StringObjectMap{}
	if err := serialization.JsonToObject(input, &out); err != nil {
		return nil, errors.Wrap(err, "failed to read map from input string")
	}
	return out, nil
}

// Convert input json string to StringObjectMap, and give empty map in case of error
func StringObjectMapFromJsonOrEmpty(input string) StringObjectMap {
	out := StringObjectMap{}
	if err := serialization.JsonToObject(input, &out); err != nil {
		return out
	}
	return out
}

// Convert StringObjectMap to a Json string
func (s StringObjectMap) JsonString() (string, error) {
	return serialization.Stringify(s)
}

// Convert StringObjectMap to a Json string, or give empty json "{}" on error
func (s StringObjectMap) JsonStringOrEmptyJson() string {
	str, err := serialization.Stringify(s)
	if err != nil {
		return "{}"
	} else {
		return str
	}
}

// ------------------------------------------ Helper to get sub key ----------------------------------------------------
func (s StringObjectMap) BoolOrFalse2(key1 string, key2 string) bool {
	return s.StringObjectMapOrEmpty(key1).BoolOrFalse(key2)
}

func (s StringObjectMap) BoolOrFalse3(key1 string, key2 string, key3 string) bool {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).BoolOrFalse(key3)
}

func (s StringObjectMap) BoolOrFalse4(key1 string, key2 string, key3 string, key4 string) bool {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).StringObjectMapOrEmpty(key3).BoolOrFalse(key4)
}

func (s StringObjectMap) BoolOrTrue2(key1 string, key2 string) bool {
	return s.StringObjectMapOrEmpty(key1).BoolOrTrue(key2)
}
func (s StringObjectMap) BoolOrTrue3(key1 string, key2 string, key3 string) bool {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).BoolOrTrue(key3)
}

func (s StringObjectMap) BoolOrTrue4(key1 string, key2 string, key3 string, key4 string) bool {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).StringObjectMapOrEmpty(key3).BoolOrTrue(key4)
}

func (s StringObjectMap) StringOrEmpty2(key1 string, key2 string) string {
	return s.StringObjectMapOrEmpty(key1).StringOrEmpty(key2)
}
func (s StringObjectMap) StringOrEmpty3(key1 string, key2 string, key3 string) string {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).StringOrEmpty(key3)
}

func (s StringObjectMap) StringOrEmpty4(key1 string, key2 string, key3 string, key4 string) string {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).StringObjectMapOrEmpty(key3).StringOrEmpty(key4)
}

func (s StringObjectMap) IntOrZero2(key1 string, key2 string) int {
	return s.StringObjectMapOrEmpty(key1).IntOrZero(key2)
}
func (s StringObjectMap) IntOrZero3(key1 string, key2 string, key3 string) int {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).IntOrZero(key3)
}

func (s StringObjectMap) IntOrZero4(key1 string, key2 string, key3 string, key4 string) int {
	return s.StringObjectMapOrEmpty(key1).StringObjectMapOrEmpty(key2).StringObjectMapOrEmpty(key3).IntOrZero(key4)
}

// Convert map to String Object Map
func ConvertStringObjectMapToMap(in StringObjectMap, out map[string]interface{}) {
	for k, v := range in {
		switch val := v.(type) {
		case StringObjectMap:
			{
				m := map[string]interface{}{}
				Map(val, m)
				out[k] = m
			}

		case []StringObjectMap:
			l := make([]interface{}, 0)
			for _, _v := range val {
				m := map[string]interface{}{}
				Map(_v, m)
				l = append(l, m)
			}
			out[k] = l

		case []interface{}:
			l := make([]interface{}, 0)
			for _, _v := range val {
				switch _val := _v.(type) {
				case StringObjectMap:
					m := map[string]interface{}{}
					Map(_val, m)
					l = append(l, m)
				default:
					l = append(l, _val)
				}
			}
			out[k] = l
		default:
			out[k] = v
		}
	}
}
