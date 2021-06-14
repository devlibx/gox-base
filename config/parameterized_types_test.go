package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameterizedString(t *testing.T) {

	tests := []struct {
		TestName string
		Input    string
		Output   string
		err      error
	}{
		{TestName: "1", Input: "a", Output: "a", err: nil},
		{TestName: "2", Input: "true", Output: "true", err: nil},
		{TestName: "3", Input: "false", Output: "false", err: nil},
		{TestName: "4", Input: "env: prod=1 ; dev=2", Output: "1", err: nil},
		{TestName: "5", Input: "env: prod=11 ; dev=2", Output: "11", err: nil},
		{TestName: "6", Input: "env: prod= space string ; dev=2", Output: "space string", err: nil},
		{TestName: "7", Input: "env: default= default value ; dev=2", Output: "default value", err: nil},
		{TestName: "8", Input: "env:  dev=2", Output: "", err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			var value = ParameterizedString(tt.Input)
			out, _ := value.Get("prod")
			if out != tt.Output {
				assert.Fail(t, fmt.Sprintf("test=%s got %q, want %q", tt.TestName, out, tt.Output))
			}
		})
	}
}
