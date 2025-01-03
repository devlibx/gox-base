package goxJsonUtils

import (
	"github.com/devlibx/gox-base/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestBytesToObject(t *testing.T) {
	// Test with valid JSON payload
	validJson := []byte(`{"field1": "value1", "field2": 123}`)
	result, err := BytesToObject[TestStruct](validJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid JSON payload
	invalidJson := []byte(`{"field1": "value1", "field2": }`)
	result, err = BytesToObject[TestStruct](invalidJson)
	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestStringObjectMapToString(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"key1": "value1", "key2": 123}
	result, err := StringObjectMapToString(input)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"key1":"value1","key2":123}`, result)

	// Test with empty StringObjectMap
	input = gox.StringObjectMap{}
	result, err = StringObjectMapToString(input)
	assert.NoError(t, err)
	assert.JSONEq(t, `{}`, result)
}

func TestStringObjectMapToBytes(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"key1": "value1", "key2": 123}
	result, err := StringObjectMapToBytes(input)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"key1":"value1","key2":123}`, string(result))

	// Test with empty StringObjectMap
	input = gox.StringObjectMap{}
	result, err = StringObjectMapToBytes(input)
	assert.NoError(t, err)
	assert.JSONEq(t, `{}`, string(result))
}

func TestStringObjectMapToObject(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"field1": "value1", "field2": 123}
	result, err := StringObjectMapToObject[TestStruct](input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with empty StringObjectMap
	input = gox.StringObjectMap{}
	result, err = StringObjectMapToObject[TestStruct](input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Zero(t, result)
}

func TestStringToStringObjectMap(t *testing.T) {
	// Test with valid JSON string
	validJson := `{"key1": "value1", "key2": 123}`
	result, err := StringToStringObjectMap(validJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, float64(123), result["key2"])

	// Test with invalid JSON string
	invalidJson := `{"key1": "value1", "key2": }`
	result, err = StringToStringObjectMap(invalidJson)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBytesToStringObjectMap(t *testing.T) {
	// Test with valid JSON byte slice
	validJson := []byte(`{"key1": "value1", "key2": 123}`)
	result, err := BytesToStringObjectMap(validJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, float64(123), result["key2"])

	// Test with invalid JSON byte slice
	invalidJson := []byte(`{"key1": "value1", "key2": }`)
	result, err = BytesToStringObjectMap(invalidJson)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestObjectToStringObjectMap(t *testing.T) {
	// Test with valid object
	input := TestStruct{Field1: "value1", Field2: 123}
	result, err := ObjectToStringObjectMap(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["field1"])
	assert.Equal(t, float64(123), result["field2"])

	// Test with invalid object
	inputInvalid := make(chan int)
	result, err = ObjectToStringObjectMap(inputInvalid)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBytesToObjectSuppressError(t *testing.T) {
	// Test with valid JSON byte slice
	validJson := []byte(`{"field1": "value1", "field2": 123}`)
	result := BytesToObjectSuppressError[TestStruct](validJson)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid JSON byte slice
	invalidJson := []byte(`{"field1": "value1", "field2": }`)
	result = BytesToObjectSuppressError[TestStruct](invalidJson)
	assert.Zero(t, result)
}

func TestStringObjectMapToStringSuppressError(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"key1": "value1", "key2": 123}
	result := StringObjectMapToStringSuppressError(input)
	assert.JSONEq(t, `{"key1":"value1","key2":123}`, result)

	// Test with invalid StringObjectMap
	inputInvalid := gox.StringObjectMap{"key1": make(chan int)}
	result = StringObjectMapToStringSuppressError(inputInvalid)
	assert.Equal(t, "", result)
}

func TestStringObjectMapToBytesSuppressError(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"key1": "value1", "key2": 123}
	result := StringObjectMapToBytesSuppressError(input)
	assert.JSONEq(t, `{"key1":"value1","key2":123}`, string(result))

	// Test with invalid StringObjectMap
	inputInvalid := gox.StringObjectMap{"key1": make(chan int)}
	result = StringObjectMapToBytesSuppressError(inputInvalid)
	assert.Nil(t, result)
}

func TestStringObjectMapToObjectSuppressError(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"field1": "value1", "field2": 123}
	result := StringObjectMapToObjectSuppressError[TestStruct](input)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid StringObjectMap
	inputInvalid := gox.StringObjectMap{"field1": make(chan int)}
	result = StringObjectMapToObjectSuppressError[TestStruct](inputInvalid)
	assert.Zero(t, result)
}

func TestStringToStringObjectMapSuppressError(t *testing.T) {
	// Test with valid JSON string
	validJson := `{"key1": "value1", "key2": 123}`
	result := StringToStringObjectMapSuppressError(validJson)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, float64(123), result["key2"])

	// Test with invalid JSON string
	invalidJson := `{"key1": "value1", "key2": }`
	result = StringToStringObjectMapSuppressError(invalidJson)
	assert.Empty(t, result)
}

func TestBytesToStringObjectMapSuppressError(t *testing.T) {
	// Test with valid JSON byte slice
	validJson := []byte(`{"key1": "value1", "key2": 123}`)
	result := BytesToStringObjectMapSuppressError(validJson)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, float64(123), result["key2"])

	// Test with invalid JSON byte slice
	invalidJson := []byte(`{"key1": "value1", "key2": }`)
	result = BytesToStringObjectMapSuppressError(invalidJson)
	assert.Empty(t, result)
}

func TestObjectToStringObjectMapSuppressError(t *testing.T) {
	// Test with valid object
	input := TestStruct{Field1: "value1", Field2: 123}
	result := ObjectToStringObjectMapSuppressError(input)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["field1"])
	assert.Equal(t, float64(123), result["field2"])

	// Test with invalid object
	inputInvalid := make(chan int)
	result = ObjectToStringObjectMapSuppressError(inputInvalid)
	assert.Empty(t, result)
}

func TestStringObjectMapToObjectIncludingElseCase(t *testing.T) {
	// Test with valid StringObjectMap
	input := gox.StringObjectMap{"field1": "value1", "field2": 123}
	result, err := StringObjectMapToObject[TestStruct](input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with empty StringObjectMap
	input = gox.StringObjectMap{}
	result, err = StringObjectMapToObject[TestStruct](input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Zero(t, result)

	// Test with invalid StringObjectMap (to trigger else case)
	inputInvalid := gox.StringObjectMap{"field1": make(chan int)}
	result, err = StringObjectMapToObject[TestStruct](inputInvalid)
	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestStringToObject(t *testing.T) {
	// Test with valid JSON string
	validJson := `{"field1": "value1", "field2": 123}`
	result, err := StringToObject[TestStruct](validJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid JSON string
	invalidJson := `{"field1": "value1", "field2": }`
	result, err = StringToObject[TestStruct](invalidJson)
	assert.Error(t, err)
	assert.Zero(t, result)
}
