package serialization

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func TestStringToObject(t *testing.T) {
	// Test with valid JSON payload
	validJson := `{"field1": "value1", "field2": 123}`
	result, err := StringToObject[TestStruct](validJson)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Equal(t, 123, result.Field2)

	// Test with invalid JSON payload
	invalidJson := `{"field1": "value1", "field2": }`
	result, err = StringToObject[TestStruct](invalidJson)
	assert.Error(t, err)
	assert.Zero(t, result)
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
	invalidJson := `{"field1": "value1", "field2": 123 }`
	reader = bytes.NewReader([]byte(invalidJson))
	result, err = ReadPayload[TestStruct](reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)
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
		_, _ = fmt.Fprintln(w, StringifyOrEmptyJsonOnError(ts))
	}))
	defer server.Close()

	// Make a Post request to the test server
	reqBody := TestStruct{Field1: "value1", Field2: 123}
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewReader([]byte(StringifyOrEmptyJsonOnError(reqBody))))
	assert.NoError(t, err)
	defer resp.Body.Close()

	ts, err := ReadPayload[TestStruct](resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "resp_value1", ts.Field1)
}

func TestHttpRequestResponse_WritePayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check payload
		ts, err := ReadPayload[TestStruct](r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "value1", ts.Field1)

		// Send modified response
		ts.Field1 = "resp_" + ts.Field1
		w.WriteHeader(http.StatusOK)
		assert.NoError(t, WritePayload(w, ts))
	}))
	defer server.Close()

	// Make a Post request to the test server
	reqBody := TestStruct{Field1: "value1", Field2: 123}
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewReader([]byte(StringifyOrEmptyJsonOnError(reqBody))))
	assert.NoError(t, err)
	defer resp.Body.Close()

	ts, err := ReadPayload[TestStruct](resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "resp_value1", ts.Field1)
}
