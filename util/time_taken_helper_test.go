package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeTracker(t *testing.T) {
	capture := NewTimeTracker(true).Capture()
	time.Sleep(4 * time.Millisecond)
	capture.Record("first")
	time.Sleep(14 * time.Millisecond)
	capture.Record("second")
	time.Sleep(11 * time.Millisecond)
	capture.Record("third")
	time.Sleep(100 * time.Millisecond)
	capture.Record("forth")
	time.Sleep(1 * time.Nanosecond)
	fmt.Println(capture.DumpMillis())
}

func TestTimeTracker_NoRecord(t *testing.T) {
	capture := NewTimeTracker(true).Capture()
	fmt.Println(capture.DumpMillis())
}

func TestTimeTracker_Only_One(t *testing.T) {
	capture := NewTimeTracker(true).Capture()
	time.Sleep(4 * time.Millisecond)
	capture.Record("first")
	fmt.Println(capture.DumpMillis())
}

func TestTimeTracker_Only_One_NoWait(t *testing.T) {
	capture := NewTimeTracker(true).Capture()
	time.Sleep(4 * time.Millisecond)
	capture.Record("first")
	fmt.Println(capture.DumpMillis())
}
