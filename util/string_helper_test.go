package util

import (
	"testing"
)

func TestIsStringEmpty(t *testing.T) {
	tests := []struct {
		TestName string
		Input    string
		Output   bool
		err      error
	}{
		{TestName: "1", Input: "", Output: true, err: nil},
		{TestName: "2", Input: " ", Output: true, err: nil},
		{TestName: "3", Input: "	", Output: true, err: nil},
		{TestName: "4", Input: "\t", Output: true, err: nil},
		{TestName: "5", Input: "a", Output: false, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := IsStringEmpty(tt.Input)
			if out != tt.Output {
				t.Errorf("got %t, want %t", out, tt.Output)
			}
		})
	}
}
