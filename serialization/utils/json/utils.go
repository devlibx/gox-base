package goxJsonUtils

import (
	"encoding/json"
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-base/v2/serialization"
)

// BytesToObject converts a byte slice to an object of type T.
//
// Parameters:
// - data: byte slice containing the JSON data to be unmarshaled.
//
// Returns:
// - T: the unmarshaled object of type T.
// - error: an error if the unmarshaling fails.
func BytesToObject[T any](data []byte) (T, error) {
	var retValue T
	if err := json.Unmarshal(data, &retValue); err != nil {
		return retValue, &serialization.DeserializationError{
			Err:          err,
			ErrorMessage: "error in parsing input data to object T",
		}
	}
	return retValue, nil
}

// StringToObject converts a string to a object
//
// Parameters:
// - input: the string to be converted.
//
// Returns:
// - T: the T object from inout string
// - error: an error if the conversion fails.
func StringToObject[T any](input string) (T, error) {
	return BytesToObject[T]([]byte(input))
}

// StringObjectMapToString converts a StringObjectMap to a JSON string.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - string: the JSON string representation of the input map.
// - error: an error if the conversion fails.
func StringObjectMapToString(input gox.StringObjectMap) (string, error) {
	return serialization.Stringify(input)
}

// StringObjectMapToBytes converts a StringObjectMap to a byte slice.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - []byte: the byte slice representation of the input map.
// - error: an error if the conversion fails.
func StringObjectMapToBytes(input gox.StringObjectMap) ([]byte, error) {
	if b, err := serialization.Stringify(input); err != nil {
		return nil, err
	} else {
		return []byte(b), nil
	}
}

// StringObjectMapToObject converts a StringObjectMap to an object of type T.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - T: the converted object of type T.
// - error: an error if the conversion fails.
func StringObjectMapToObject[T any](input gox.StringObjectMap) (T, error) {
	var retValue T
	if b, err := StringObjectMapToBytes(input); err != nil {
		return retValue, err
	} else {
		return BytesToObject[T](b)
	}
}

// StringToStringObjectMap converts a JSON string to a StringObjectMap.
//
// Parameters:
// - data: the JSON string to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap.
// - error: an error if the conversion fails.
func StringToStringObjectMap(data string) (gox.StringObjectMap, error) {
	return BytesToObject[gox.StringObjectMap]([]byte(data))
}

// BytesToStringObjectMap converts a byte slice to a StringObjectMap.
//
// Parameters:
// - data: the byte slice to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap.
// - error: an error if the conversion fails.
func BytesToStringObjectMap(data []byte) (gox.StringObjectMap, error) {
	return BytesToObject[gox.StringObjectMap](data)
}

// ObjectToStringObjectMap converts an object to a StringObjectMap.
//
// Parameters:
// - data: the object to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap.
// - error: an error if the conversion fails.
func ObjectToStringObjectMap(data any) (gox.StringObjectMap, error) {
	if d, err := json.Marshal(data); err != nil {
		return nil, &serialization.DeserializationError{
			Err:          err,
			ErrorMessage: "failed to convert input data to map data",
		}
	} else {
		return BytesToObject[gox.StringObjectMap](d)
	}
}

// BytesToObjectSuppressError converts a byte slice to an object of type T.
// Suppresses any error and returns a zero value of T in case of failure.
//
// Parameters:
// - data: the byte slice to be converted.
//
// Returns:
// - T: the converted object of type T, or zero value of T if conversion fails.
func BytesToObjectSuppressError[T any](data []byte) T {
	var retValue T
	if err := json.Unmarshal(data, &retValue); err != nil {
		return retValue
	}
	return retValue
}

// StringObjectMapToStringSuppressError converts a StringObjectMap to a JSON string.
// Suppresses any error and returns an empty string in case of failure.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - string: the JSON string representation of the input map, or an empty string if conversion fails.
func StringObjectMapToStringSuppressError(input gox.StringObjectMap) string {
	r, _ := serialization.Stringify(input)
	return r
}

// StringObjectMapToBytesSuppressError converts a StringObjectMap to a byte slice.
// Suppresses any error and returns nil in case of failure.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - []byte: the byte slice representation of the input map, or nil if conversion fails.
func StringObjectMapToBytesSuppressError(input gox.StringObjectMap) []byte {
	if b, err := serialization.Stringify(input); err != nil {
		return nil
	} else {
		return []byte(b)
	}
}

// StringObjectMapToObjectSuppressError converts a StringObjectMap to an object of type T.
// Suppresses any error and returns a zero value of T in case of failure.
//
// Parameters:
// - input: the StringObjectMap to be converted.
//
// Returns:
// - T: the converted object of type T, or zero value of T if conversion fails.
func StringObjectMapToObjectSuppressError[T any](input gox.StringObjectMap) T {
	var retValue T
	if b, err := StringObjectMapToBytes(input); err != nil {
		return retValue
	} else {
		return BytesToObjectSuppressError[T](b)
	}
}

// StringToStringObjectMapSuppressError converts a JSON string to a StringObjectMap.
// Suppresses any error and returns an empty StringObjectMap in case of failure.
//
// Parameters:
// - data: the JSON string to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap, or an empty StringObjectMap if conversion fails.
func StringToStringObjectMapSuppressError(data string) gox.StringObjectMap {
	return BytesToObjectSuppressError[gox.StringObjectMap]([]byte(data))
}

// BytesToStringObjectMapSuppressError converts a byte slice to a StringObjectMap.
// Suppresses any error and returns an empty StringObjectMap in case of failure.
//
// Parameters:
// - data: the byte slice to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap, or an empty StringObjectMap if conversion fails.
func BytesToStringObjectMapSuppressError(data []byte) gox.StringObjectMap {
	return BytesToObjectSuppressError[gox.StringObjectMap](data)
}

// ObjectToStringObjectMapSuppressError converts an object to a StringObjectMap.
// Suppresses any error and returns an empty StringObjectMap in case of failure.
//
// Parameters:
// - data: the object to be converted.
//
// Returns:
// - gox.StringObjectMap: the converted StringObjectMap, or an empty StringObjectMap if conversion fails.
func ObjectToStringObjectMapSuppressError(data any) gox.StringObjectMap {
	if d, err := json.Marshal(data); err != nil {
		return nil
	} else {
		return BytesToObjectSuppressError[gox.StringObjectMap](d)
	}
}

func StringToObjectSuppressError[T any](input string) T {
	ret, _ := StringToObject[T](input)
	return ret
}
