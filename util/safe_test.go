package util

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafeRun(t *testing.T) {
	SafeRun(func() {
		willPanic()
	}, "got error in function")
}

func TestSafeRunWithReturn_NoPanic(t *testing.T) {
	value, err := SafeRunWithReturn(func() (interface{}, error) {
		return willPanicWithReturnType(false)
	}, "got error in function")
	assert.Equal(t, "good", value)
	assert.NoError(t, err)
}

func TestSafeRunWithReturn_WithPanic(t *testing.T) {
	value, err := SafeRunWithReturn(func() (interface{}, error) {
		return willPanicWithReturnType(true)
	}, "got error in function")
	_ = value
	assert.Error(t, err)
	fmt.Println(err.Error())
}

func willPanic() {
	panic(errors.New("err"))
}

func willPanicWithReturnType(fail bool) (string, error) {
	if fail {
		panic(errors.New("err"))
	}
	return "good", nil
}
