package httpHelper

import (
	"github.com/devlibx/gox-base/serialization"
	"io/ioutil"
	"net/http"
)

type HttpPlayloadDeserializationError struct {
	error
	ErrorMessage string
	ErrorStatus  int
}

func ReadJsonPayload(request *http.Request, object interface{}) error {

	// read body from http_helper request
	_body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return &HttpPlayloadDeserializationError{
			error:        err,
			ErrorMessage: "failed to read request body",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	// read object from request body
	if err := serialization.JsonBytesToObject(_body, object); err != nil {
		return &HttpPlayloadDeserializationError{
			error:        err,
			ErrorMessage: "failed to read body in json object",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	return nil
}

func ReadYamlPayload(request *http.Request, object interface{}) error {

	// read body from http_helper request
	_body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return &HttpPlayloadDeserializationError{
			error:        err,
			ErrorMessage: "failed to read request body",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	// read object from request body
	if err := serialization.ReadYamlFromString(string(_body), object); err != nil {
		return &HttpPlayloadDeserializationError{
			error:        err,
			ErrorMessage: "failed to read body in json object",
			ErrorStatus:  http.StatusBadRequest,
		}
	}

	return nil
}
