package metrics

import "net/http"

// Labels represents a kev-value mapping. It is used to add information on metric
// e.g. counter.With(Labels({"type": "request", "code": 123}
type Labels map[string]interface{}

// Configure is used to provide initialization data for metric
type Configuration struct {
}

// A operation to increment a metrics e.g. counter.Inc() will increment a counter metric
type IncOperation interface {
	Inc()
	IncWithLabels(labels Labels)
}

// A increment operation with labels
type MetricWithLabelValues interface {
	IncOperation
}

// A counter to log metrics
type Counter interface {
	IncOperation
	WithLabels(labels Labels) MetricWithLabelValues
}

// Service represents a parent interface for metrics. It allows to register and
// get counter, timer etc
type Service interface {
	Initialize(configuration Configuration) error
	RegisterCounter(name string, help string, labels []string) error
	Counter(counterName string) Counter
	HttpHandler() http.Handler
}

// Returns a service which does no-op
func NewNoOpMetrics() Service {
	return &dummyService{}
}

// Metric name and data
type LabeledMetric struct {
	Name   string
	Labels Labels
}

func (l *LabeledMetric) NameWithErrorPrefix() string {
	return l.Name + "_error"
}

func (l *LabeledMetric) NameWithSuccessPrefix() string {
	return l.Name + "_error"
}
