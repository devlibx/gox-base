package util

import (
	"fmt"
	"sync"
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
	name   string
	mu     *sync.Mutex
}

func (t *captureImpl) Record(msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
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
	t.Record("end")

	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.enable {
		return ""
	}
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
	name   string
}

func (t timeTrackerImpl) Capture() Capture {
	if t.enable {
		s := &captureImpl{enable: true, times: make([]TimeTrack, 0), name: t.name, mu: &sync.Mutex{}}
		s.Record("start")
		return s
	} else {
		s := &captureImpl{enable: false, name: t.name, mu: &sync.Mutex{}}
		s.Record("start")
		return s
	}
}

func NewTimeTracker(enable bool) TimeTracker {
	t := &timeTrackerImpl{enable: enable}
	return t
}

func NewTimeTrackerWithName(enable bool, name string) TimeTracker {
	t := &timeTrackerImpl{enable: enable, name: name}
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
