package util

import (
	"hash/fnv"
	"math"
	"strings"
)

func IsStringEmpty(input string) bool {
	return len(strings.TrimSpace(input)) <= 0
}

func StringToHashMod(item string, count int) int {
	h := fnv.New64()
	_, _ = h.Write([]byte(item))
	sha := h.Sum64() % uint64(count)
	return int(math.Abs(float64(sha)))
}
