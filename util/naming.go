package util

import (
	"path/filepath"
	"runtime"
)

func GetMethodName(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	_, file := filepath.Split(fn.Name())
	return file
}

func GetCurrentMethodName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	_, file := filepath.Split(fn.Name())
	return file
}

func GetCallingMethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	_, file := filepath.Split(fn.Name())
	return file
}

func GetMethodNameName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	_, file := filepath.Split(fn.Name())
	return file
}
