package httpHelper

import (
	"fmt"
	"github.com/devlibx/gox-base/serialization"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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

// PortHelper is a helper to get random port
type PortHelper struct {
	Listener net.Listener
	Port     int
}

func (p *PortHelper) Dump() string {
	return fmt.Sprintf("PortHelper{Port=%d}", p.Port)
}

func NewPortHelper() (*PortHelper, func(), error) {
	p := &PortHelper{}
	if listener, err := net.Listen("tcp", ":0"); err != nil {
		return nil, func() {}, err
	} else {
		p.Listener = listener
	}
	p.Port = p.Listener.Addr().(*net.TCPAddr).Port
	return p, func() { _ = p.Listener.Close() }, nil
}

// AllocateFreePortsAndAssignToEnvironmentVariables will allocate free ports and assign to environment variables
//
// For example, you need to allocate 2 random ports and assign to environment variables (1) HTTP_PORT, (2) STUB_PORT
// This function will allocate 2 random ports and assign to environment variables
func AllocateFreePortsAndAssignToEnvironmentVariables(envVariables ...string) (map[string]int, error) {
	result := map[string]int{}
	for _, env := range envVariables {
		if portHelper, closeFunc, err := NewPortHelper(); err == nil {
			_ = os.Setenv(env, fmt.Sprintf("%d", portHelper.Port))
			closeFunc()
			result[env] = portHelper.Port
		} else {
			return result, err
		}
	}
	return result, nil
}
