package ioSerialization

import (
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/devlibx/gox-base/v2/serialization"
	"io"
	"net/http"
)

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
		return t, &serialization.DeserializationError{
			Err:          err,
			ErrorMessage: "error in reading request body form http request",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	// read object from request body
	var retValue T
	if err := serialization.JsonBytesToObject(body, &retValue); err != nil {
		return t, &serialization.DeserializationError{
			Err:          err,
			ErrorMessage: "error in parse request body form http request to object",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	return retValue, nil
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
	if out, err := serialization.Stringify(data); err == nil {
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
	if out, err := serialization.Stringify(data); err == nil {
		_, _ = writer.Write([]byte(out))
	}
}
