package metrics

import (
	"time"
)

type x interface {
	Scope
	Counter
	Gauge
	Timer
	Histogram
}
type noOpScope struct {
}

func (n *noOpScope) Reporting() bool {
	return false
}

func (n *noOpScope) Tagging() bool {
	return false
}

func (n *noOpScope) RecordValue(value float64) {
}

func (n *noOpScope) RecordDuration(value time.Duration) {
}

func (n *noOpScope) RecordStopwatch(stopwatchStart time.Time) {
}

func (n *noOpScope) Record(value time.Duration) {
}

func (n *noOpScope) Start() Stopwatch {
	return NewStopwatch(time.Now(), n)
}

func (n *noOpScope) Update(value float64) {
}

func (n *noOpScope) Inc(delta int64) {
}

func (n *noOpScope) Counter(name string) Counter {
	return n
}

func (n *noOpScope) Gauge(name string) Gauge {
	return n
}

func (n *noOpScope) Timer(name string) Timer {
	return n
}

func (n *noOpScope) Histogram(name string, buckets Buckets) Histogram {
	return n
}

func (n *noOpScope) Tagged(tags map[string]string) Scope {
	return n
}

func (n *noOpScope) SubScope(name string) Scope {
	return n
}

func (n *noOpScope) Capabilities() Capabilities {
	return n
}

func NoOpMetric() Scope {
	return &noOpScope{}
}
