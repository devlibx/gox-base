package util

import (
	"fmt"
)

// SafeRun will make sure that the function which is passed will not crash the app. It will print the error in case
// of error on console
func SafeRun(toRun func(), errorMessage string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Got panic which was recovered: errorMessage=%s, info=%v\n", errorMessage, r)
		}
	}()
	toRun()
}

// PanicErrorWrapper will provide more detail of the error which is recovered from panic
type PanicErrorWrapper struct {
	ValueFromRecover interface{}
	errorMessage     string
}

func (receiver PanicErrorWrapper) Error() string {
	return fmt.Sprintf("Got panic which was recovered: errorMessage=%s, info=%v", receiver.errorMessage, receiver.ValueFromRecover)
}

// SafeRunWithReturn will make sure that the function which is passed will not crash the app. It will print the error in case
// of error on console

func SafeRunWithReturn(toRun func() (interface{}, error), errorMessage string) (val interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &PanicErrorWrapper{
				ValueFromRecover: r,
				errorMessage:     errorMessage,
			}
		}
	}()
	return toRun()
}
