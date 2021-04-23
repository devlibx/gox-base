package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// A constant error - to test errors.Is() function
var errTestConstant = errors.New("error_constant")

// A custom error - to test errors.As() function
type customError struct {
}

// Implementing error interface
func (e customError) Error() string {
	return "CustomError"
}

// A test object to verify different error methods
type errorTestInterfaceImpl struct {
}

func (e *errorTestInterfaceImpl) MethodWithErrorCreatedUsingNewErrorMethod() error {
	return NewError(
		"test_error_name",
		"test_error_description",
		errors.New("errors"),
		10,
	)
}

// Test a error created using "NewError" method is valid
func TestNewErrorMethod(t *testing.T) {
	var eImpl errorTestInterfaceImpl
	e := eImpl.MethodWithErrorCreatedUsingNewErrorMethod()
	var errorObj Error
	if errors.As(e, &errorObj) {
		assert.Equal(t, "test_error_name", errorObj.GetCode())
		assert.Equal(t, "test_error_description", errorObj.GetMessage())
		assert.Equal(t, "errors", errorObj.GetError().Error())
		assert.Equal(t, 10, errorObj.GetData())
	} else {
		assert.Fail(t, "Expected it to be DetailedError", e)
	}
}

func (e *errorTestInterfaceImpl) MethodDetailedError() error {
	return &DetailedError{
		Code:    "test_error_name",
		Message: "test_error_description",
		Data:    10,
		Err:     errors.New("errors"),
	}
}

// Test a error created using "DetailedError"
func TestDetailedError(t *testing.T) {
	var eImpl errorTestInterfaceImpl
	e := eImpl.MethodDetailedError()
	var errorObj Error
	if errors.As(e, &errorObj) {
		assert.Equal(t, "test_error_name", errorObj.GetCode())
		assert.Equal(t, "test_error_description", errorObj.GetMessage())
		assert.Equal(t, "errors", errorObj.GetError().Error())
		assert.Equal(t, 10, errorObj.GetData())
	} else {
		assert.Fail(t, "Expected it to be DetailedError", e)
	}
}

func (e *errorTestInterfaceImpl) MethodToCreateErrorConstant() error {
	return &DetailedError{
		Code:    "test_error_name",
		Message: "test_error_description",
		Data:    10,
		Err:     errTestConstant,
	}
}

// Test errors.Is() method
func TestErrorConstant(t *testing.T) {
	var eImpl errorTestInterfaceImpl
	e := eImpl.MethodToCreateErrorConstant()
	if !errors.Is(e, errTestConstant) {
		assert.Fail(t, "We expected error to be errTestConstant", e)
	}
}

func (e *errorTestInterfaceImpl) MethodToCreateCustomError() error {
	return &DetailedError{
		Code:    "test_error_name",
		Message: "test_error_description",
		Data:    10,
		Err:     customError{},
	}
}

// Test errors.As() method
func TestErrorAsFunction(t *testing.T) {
	var eImpl errorTestInterfaceImpl
	e := eImpl.MethodToCreateCustomError()
	var ce customError
	if !errors.As(e, &ce) {
		assert.Fail(t, "We expected error to be customError", e)
	} else {
		assert.Equal(t, "CustomError", ce.Error())
	}
}

func TestUsage(t *testing.T) {

	// Use case 1 - When we want to generate a detailed error with a error constant
	var errorGenerateFunction = func() error {
		return &DetailedError{
			Code:    "1",
			Message: "",
			Data:    nil,
			Err:     os.ErrNoDeadline,
		}
	}

	// Client may want to check for specific error
	var err = errorGenerateFunction()
	if !errors.Is(err, os.ErrNoDeadline) {
		assert.Fail(t, "Expected it to be ErrNoDeadline")
	}

	// Use case 2 - Wrap DetailedError inside DetailedError
	var wrappedErrorGenerateFunction = func() error {
		return &DetailedError{
			Code:    "2",
			Message: "",
			Data:    nil,
			Err:     err,
		}
	}

	// Let's see if we can recover a wrapped error or not
	var wrappedErr = wrappedErrorGenerateFunction()
	if !errors.Is(wrappedErr, err) {
		assert.Fail(t, "Expected it to be err (DetailedError)")
	}

	// Let's see of errors.As() works
	var errorAs Error
	if !errors.As(wrappedErr, &errorAs) {
		assert.Fail(t, "Expected Error from As method")
	} else {
		assert.Equal(t, "2", errorAs.GetCode())
	}
}
