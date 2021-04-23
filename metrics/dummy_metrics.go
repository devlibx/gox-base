package metrics

import "net/http"

type dummyIncOperations struct {
}

func (d *dummyIncOperations) Inc() {
}

func (d *dummyIncOperations) IncWithLabels(labels Labels) {
}

type dummyMetricWithLabelValues struct {
	dummyIncOperations
}

type dummyCounter struct {
	dummyIncOperations
}

func (d *dummyCounter) WithLabels(labels Labels) MetricWithLabelValues {
	return &dummyMetricWithLabelValues{dummyIncOperations{}}
}

type dummyService struct {
}

func (d *dummyService) HttpHandler() http.Handler {
	return nil
}

func (d *dummyService) Initialize(configuration Configuration) error {
	return nil
}

func (d *dummyService) RegisterCounter(name string, help string, labels []string) error {
	return nil
}

func (d *dummyService) Counter(counterName string) Counter {
	return &dummyCounter{dummyIncOperations{}}
}

func NewNoOpCounter() Counter {
	return &dummyCounter{dummyIncOperations{}}
}
