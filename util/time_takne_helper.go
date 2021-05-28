package util

import (
	"fmt"
	"time"
)

//go:generate mockgen -source=time_takne_helper.go -destination=../mocks/util/mock_time_takne_helper.go -package=mockUtil

type TimeTracker interface {
	Capture() Capture
}
type Capture interface {
	Record(msg string)
	DumpMillis() string
	DumpMicros() string
	DumpNanos() string
}

type TimeTrack struct {
	Message string
	Time    time.Time
}

type captureImpl struct {
	times  []TimeTrack
	enable bool
}

func (t *captureImpl) Record(msg string) {
	if t.enable {
		t.times = append(t.times, TimeTrack{
			Message: msg,
			Time:    time.Now(),
		})
	}
}

func (t *captureImpl) DumpMillis() string {
	return t.dump("ms")
}

func (t *captureImpl) DumpMicros() string {
	return t.dump("micro")
}

func (t *captureImpl) DumpNanos() string {
	return t.dump("ns")
}

func (t *captureImpl) dump(unit string) string {
	if !t.enable {
		return ""
	}
	t.Record("end")
	result := ""
	length := len(t.times)
	for i := 1; i < length; i++ {
		f := t.times[i-1]
		s := t.times[i]

		switch unit {
		case "ms":
			u := unit
			timeTaken := s.Time.Sub(f.Time).Milliseconds()
			if timeTaken == 0 {
				timeTaken = s.Time.Sub(f.Time).Microseconds()
				u = "micro"
			}
			result += fmt.Sprintf("[%s %s]=%d %s ", f.Message, s.Message, timeTaken, u)
		case "micro":
			timeTaken := s.Time.Sub(f.Time).Microseconds()
			result += fmt.Sprintf("[%s %s]=%d %s ", f.Message, s.Message, timeTaken, unit)
		case "ns":
			timeTaken := s.Time.Sub(f.Time).Nanoseconds()
			result += fmt.Sprintf("[%s %s]=%d %s ", f.Message, s.Message, timeTaken, unit)
		default:
			timeTaken := s.Time.Sub(f.Time).Milliseconds()
			result += fmt.Sprintf("[%s %s]=%d %s ", f.Message, s.Message, timeTaken, unit)
		}
	}
	return result
}

type timeTrackerImpl struct {
	enable bool
}

func (t timeTrackerImpl) Capture() Capture {
	if t.enable {
		s := &captureImpl{enable: true, times: make([]TimeTrack, 0)}
		s.Record("start")
		return s
	} else {
		s := &captureImpl{enable: false}
		s.Record("start")
		return s
	}
}

func NewTimeTracker(enable bool) TimeTracker {
	t := &timeTrackerImpl{enable: enable}
	return t
}

type noOpImpl struct {
}

func (t *noOpImpl) Capture() Capture {
	return t
}

func (t *noOpImpl) Record(msg string) {
}

func (t *noOpImpl) DumpMillis() string {
	return ""
}

func (t *noOpImpl) DumpMicros() string {
	return ""
}

func (t *noOpImpl) DumpNanos() string {
	return ""
}

func (t *noOpImpl) Active() bool {
	return false
}

func NewNoOpTimeTracker() TimeTracker {
	return &noOpImpl{}
}
