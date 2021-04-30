package util

import "strings"

func IsStringEmpty(input string) bool {
	return len(strings.TrimSpace(input)) <= 0
}
