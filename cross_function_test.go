package gox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNoOpCrossFunction(t *testing.T) {
	cf := NewNoOpCrossFunction()

	assert.NotNil(t, cf)
}
