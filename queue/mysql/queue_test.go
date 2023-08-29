package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEndOfWeek(t *testing.T) {

	inputTime, _ := time.Parse("2006-01-02 15:04:05", "2023-08-29 15:04:05")
	outTime := endOfWeek(inputTime)
	fmt.Println(inputTime, outTime)
	assert.Equal(t, time.Month(9), outTime.Month())
	assert.Equal(t, 3, outTime.Day())

	inputTime = outTime
	outTime = endOfWeek(inputTime)
	fmt.Println(inputTime, outTime)
	assert.Equal(t, time.Month(9), outTime.Month())
	assert.Equal(t, 3, outTime.Day())

	inputTime = outTime.Add(time.Hour)
	outTime = endOfWeek(inputTime)
	fmt.Println(inputTime, outTime)
	assert.Equal(t, time.Month(9), outTime.Month())
	assert.Equal(t, 10, outTime.Day())
}
