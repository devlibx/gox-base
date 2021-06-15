package serialization

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameterizedString(t *testing.T) {

	tests := []struct {
		TestName    string
		Input       string
		Output      string
		err         error
		expectError bool
	}{
		{TestName: "1", Input: "a", Output: "a", err: nil},
		{TestName: "2", Input: "true", Output: "true", err: nil},
		{TestName: "3", Input: "false", Output: "false", err: nil},
		{TestName: "4", Input: "env:string: prod=1 ; dev=2", Output: "1", err: nil},
		{TestName: "5", Input: "env:string: prod=11 ; dev=2", Output: "11", err: nil},
		{TestName: "6", Input: "env:string: prod= space string ; dev=2", Output: "space string", err: nil},
		{TestName: "7", Input: "env:string: default= default value ; dev=2", Output: "default value", err: nil},
		{TestName: "8", Input: "env:string: dev=2", Output: "", err: nil, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			var value = ParameterizedValue(tt.Input)
			out, err := value.Get("prod")
			if tt.expectError && err == nil {
				assert.Fail(t, fmt.Sprintf("test=%s expected error but did not get error", tt.TestName))
			} else if !tt.expectError && out != tt.Output {
				assert.Fail(t, fmt.Sprintf("test=%s got %q, want %q", tt.TestName, out, tt.Output))
			}
		})
	}
}
