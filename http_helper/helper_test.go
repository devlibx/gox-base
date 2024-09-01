package httpHelper

import (
	"bytes"
	"fmt"
	"github.com/devlibx/gox-base/v2/serialization"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestReadPayload(t *testing.T) {
	// Test with valid JSON payload
	validJson := `{"field1": "value1", "field2": 123}`
	reader := bytes.NewReader([]byte(validJson))
	result, err := ReadPayload[TestStruct](reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid JSON payload
	invalidJson := `{"field1": "value1", "field2": }`
	reader = bytes.NewReader([]byte(invalidJson))
	result, err = ReadPayload[TestStruct](reader)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReadStringObjectMap(t *testing.T) {
	// Test with valid JSON payload
	validJson := `{"key1": "value1", "key2": 123}`
	reader := bytes.NewReader([]byte(validJson))
	result, err := ReadStringObjectMap(reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, 123, result.IntOrZero("key2"))

	// Test with invalid JSON payload
	invalidJson := `{"key1": "value1", "key2": }`
	reader = bytes.NewReader([]byte(invalidJson))
	result, err = ReadStringObjectMap(reader)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test with empty JSON payload
	emptyJson := `{}`
	reader = bytes.NewReader([]byte(emptyJson))
	result, err = ReadStringObjectMap(reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))

	// Test with malformed JSON payload
	malformedJson := `{"key1": "value1", "key2": 123`
	reader = bytes.NewReader([]byte(malformedJson))
	result, err = ReadStringObjectMap(reader)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestHttpRequestResponse_Pojo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check payload
		ts, err := ReadPayload[TestStruct](r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "value1", ts.Field1)

		// Send modified response
		ts.Field1 = "resp_" + ts.Field1
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, serialization.StringifyOrEmptyJsonOnError(ts))
	}))
	defer server.Close()

	// Make a Post request to the test server
	reqBody := TestStruct{Field1: "value1", Field2: 123}
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewReader([]byte(serialization.StringifyOrEmptyJsonOnError(reqBody))))
	assert.NoError(t, err)
	defer resp.Body.Close()

	ts, err := ReadPayload[TestStruct](resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "resp_value1", ts.Field1)
}

func TestHttpRequestResponse_StringObjectMap(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check payload
		ts, err := ReadStringObjectMap(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "value1", ts.StringOrEmpty("field1"))

		// Send modified response
		ts["field1"] = "resp_" + ts.StringOrEmpty("field1")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, serialization.StringifyOrEmptyJsonOnError(ts))
	}))
	defer server.Close()

	// Make a Post request to the test server
	reqBody := TestStruct{Field1: "value1", Field2: 123}
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewReader([]byte(serialization.StringifyOrEmptyJsonOnError(reqBody))))
	assert.NoError(t, err)
	defer resp.Body.Close()

	ts, err := ReadStringObjectMap(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "resp_value1", ts.StringOrEmpty("field1"))
}
