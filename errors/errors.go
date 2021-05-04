package errors

import (
	"errors"
	"fmt"
	errors1 "github.com/pkg/errors"
)

// Extended errors object
type Error interface {
	error
	GetCode() string
	GetMessage() string
	GetData() interface{}
	GetError() error
}

// Extended errors object
type DetailedError struct {
	Code    string
	Message string
	Data    interface{}
	Err     error
}

func (e *DetailedError) GetCode() string {
	return e.Code
}

func (e *DetailedError) GetMessage() string {
	return e.Message
}

func (e *DetailedError) GetData() interface{} {
	return e.Data
}

func (e *DetailedError) GetError() error {
	return e.Err
}

// Build string representation
func (e *DetailedError) Error() string {
	return fmt.Sprintf("Code=%s, Message=[%s] Error=[%v] Data=[%v]", e.Code, e.Message, e.Err, e.Data)
}

// Build string representation
func (e *DetailedError) Unwrap() error {
	return e.Err
}

// Create a new errors with more information
func NewError(code string, message string, err error, object interface{}) error {
	return &DetailedError{
		Code:    code,
		Message: message,
		Data:    object,
		Err:     err,
	}
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Wrap(err error, message string, obj ...interface{}) error {
	if obj == nil || len(obj) == 0 {
		return errors1.Wrap(err, message)
	} else {
		return errors1.Wrap(err, fmt.Sprintf(message, obj...))
	}
}
