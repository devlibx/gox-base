package serialization

import (
	"bytes"
	"encoding/json"
	"github.com/devlibx/gox-base/v2/errors"
	"io"
	"net/http"
)

// BytesToObject is a helper to read payload from http request
//
// Parameter:
// - data: data in bytes
//
// Returns:
// - *T: type of object from data from reader
// - error: error
func BytesToObject[T any](data []byte) (T, error) {
	var retValue T
	if err := JsonBytesToObject(data, &retValue); err != nil {
		return retValue, &DeserializationError{
			Err:          err,
			ErrorMessage: "error in parse request body form http request to object",
			ErrorStatus:  http.StatusBadRequest,
		}
	}
	return retValue, nil
}

// BytesToObjectSuppressError is a helper to read payload from http request. Suppress any error
//
// Parameter:
// - data: data in bytes
//
// Returns:
// - *T: type of object from data from reader
func BytesToObjectSuppressError[T any](data []byte) T {
	retValue, _ := BytesToObject[T](data)
	return retValue
}

// StringToObject is a helper to read payload from http request
//
// Parameter:
// - data: data in string
//
// Returns:
// - *T: type of object from data from reader
// - error: error
func StringToObject[T any](data string) (T, error) {
	return BytesToObject[T]([]byte(data))
}

// StringToObjectSuppressError is a helper to read payload from http request. Suppress any error
//
// Parameter:
// - data: data in string
//
// Returns:
// - *T: type of object from data from reader
func StringToObjectSuppressError[T any](data string) T {
	retValue, _ := BytesToObject[T]([]byte(data))
	return retValue
}

// ReadPayload is a helper to read payload from http request
//
// Parameter:
// - request: io.Reader
//
// Returns:
// - *T: type of object from data from reader
// - error: error
func ReadPayload[T any](request io.Reader) (T, error) {
	var t T

	// Read body from request
	body, err := io.ReadAll(request)
	if err != nil {
		return t, &DeserializationError{
			Err:          err,
			ErrorMessage: "error in reading request body form http request",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	return BytesToObject[T](body)
}

// ReadPayloadSuppressError is a helper to read payload from http request. Suppress any error.
//
// Parameter:
// - request: io.Reader
//
// Returns:
// - *T: type of object from data from reader
func ReadPayloadSuppressError[T any](request io.Reader) T {
	retValue, _ := ReadPayload[T](request)
	return retValue
}

// WritePayload is a helper to write payload to writer
//
// Parameter:
// - writer: io.Writer
// - data: any
//
// Returns:
// - error: if we fail to write data to writer
func WritePayload(writer io.Writer, data any) error {
	if out, err := Stringify(data); err == nil {
		if _, err = writer.Write([]byte(out)); err != nil {
			return errors.Wrap(err, "error in writing data to writer")
		}
	} else {
		return errors.Wrap(err, "error in writing data to writer - failed to stringify")
	}
	return nil
}

// WritePayloadSuppressError is a helper to write payload to writer
//
// Parameter:
// - writer: io.Writer
// - data: any
func WritePayloadSuppressError(writer io.Writer, data any) {
	if out, err := Stringify(data); err == nil {
		_, _ = writer.Write([]byte(out))
	}
}

// WritePayloadJsonPrettySuppressError is a helper to write payload to writer - send formated Json for easy debuging
//
// Parameter:
// - writer: io.Writer
// - data: any
func WritePayloadJsonPrettySuppressError(writer io.Writer, data any) {
	if out, err := Stringify(data); err == nil {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(out), "", "\t"); err == nil {
			_, _ = writer.Write([]byte(prettyJSON.String()))
		} else {
			_, _ = writer.Write([]byte(out))
		}
	}
}
