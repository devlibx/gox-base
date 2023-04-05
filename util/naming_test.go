package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkGetCurrentMethodName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCurrentMethodName()
	}
}

func TestGetCurrentMethodName(t *testing.T) {
	assert.Equal(t, "util.TestGetCurrentMethodName", GetCurrentMethodName())
}

func TestGetCallingMethodName(t *testing.T) {
	assert.Equal(t, "util.TestGetCallingMethodName", caller())
}

func caller() string {
	return GetCallingMethodName()
}
