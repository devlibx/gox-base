// Code generated by MockGen. DO NOT EDIT.
// Source: metrics/interface.go

// Package mockGox is a generated GoMock package.
package mockGox

import (
	http "net/http"
	reflect "reflect"
	time "time"

	metrics "github.com/devlibx/gox-base/v2/metrics"
	gomock "github.com/golang/mock/gomock"
)

// MockReporter is a mock of Reporter interface.
type MockReporter struct {
	ctrl     *gomock.Controller
	recorder *MockReporterMockRecorder
}

// MockReporterMockRecorder is the mock recorder for MockReporter.
type MockReporterMockRecorder struct {
	mock *MockReporter
}

// NewMockReporter creates a new mock instance.
func NewMockReporter(ctrl *gomock.Controller) *MockReporter {
	mock := &MockReporter{ctrl: ctrl}
	mock.recorder = &MockReporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReporter) EXPECT() *MockReporterMockRecorder {
	return m.recorder
}

// HTTPHandler mocks base method.
func (m *MockReporter) HTTPHandler() http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPHandler")
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// HTTPHandler indicates an expected call of HTTPHandler.
func (mr *MockReporterMockRecorder) HTTPHandler() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPHandler", reflect.TypeOf((*MockReporter)(nil).HTTPHandler))
}

// MockScope is a mock of Scope interface.
type MockScope struct {
	ctrl     *gomock.Controller
	recorder *MockScopeMockRecorder
}

// MockScopeMockRecorder is the mock recorder for MockScope.
type MockScopeMockRecorder struct {
	mock *MockScope
}

// NewMockScope creates a new mock instance.
func NewMockScope(ctrl *gomock.Controller) *MockScope {
	mock := &MockScope{ctrl: ctrl}
	mock.recorder = &MockScopeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScope) EXPECT() *MockScopeMockRecorder {
	return m.recorder
}

// Capabilities mocks base method.
func (m *MockScope) Capabilities() metrics.Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capabilities")
	ret0, _ := ret[0].(metrics.Capabilities)
	return ret0
}

// Capabilities indicates an expected call of Capabilities.
func (mr *MockScopeMockRecorder) Capabilities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capabilities", reflect.TypeOf((*MockScope)(nil).Capabilities))
}

// Counter mocks base method.
func (m *MockScope) Counter(name string) metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Counter", name)
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// Counter indicates an expected call of Counter.
func (mr *MockScopeMockRecorder) Counter(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Counter", reflect.TypeOf((*MockScope)(nil).Counter), name)
}

// Gauge mocks base method.
func (m *MockScope) Gauge(name string) metrics.Gauge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gauge", name)
	ret0, _ := ret[0].(metrics.Gauge)
	return ret0
}

// Gauge indicates an expected call of Gauge.
func (mr *MockScopeMockRecorder) Gauge(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gauge", reflect.TypeOf((*MockScope)(nil).Gauge), name)
}

// Histogram mocks base method.
func (m *MockScope) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Histogram", name, buckets)
	ret0, _ := ret[0].(metrics.Histogram)
	return ret0
}

// Histogram indicates an expected call of Histogram.
func (mr *MockScopeMockRecorder) Histogram(name, buckets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Histogram", reflect.TypeOf((*MockScope)(nil).Histogram), name, buckets)
}

// SubScope mocks base method.
func (m *MockScope) SubScope(name string) metrics.Scope {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubScope", name)
	ret0, _ := ret[0].(metrics.Scope)
	return ret0
}

// SubScope indicates an expected call of SubScope.
func (mr *MockScopeMockRecorder) SubScope(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubScope", reflect.TypeOf((*MockScope)(nil).SubScope), name)
}

// Tagged mocks base method.
func (m *MockScope) Tagged(tags map[string]string) metrics.Scope {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tagged", tags)
	ret0, _ := ret[0].(metrics.Scope)
	return ret0
}

// Tagged indicates an expected call of Tagged.
func (mr *MockScopeMockRecorder) Tagged(tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tagged", reflect.TypeOf((*MockScope)(nil).Tagged), tags)
}

// Timer mocks base method.
func (m *MockScope) Timer(name string) metrics.Timer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timer", name)
	ret0, _ := ret[0].(metrics.Timer)
	return ret0
}

// Timer indicates an expected call of Timer.
func (mr *MockScopeMockRecorder) Timer(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timer", reflect.TypeOf((*MockScope)(nil).Timer), name)
}

// MockClosableScope is a mock of ClosableScope interface.
type MockClosableScope struct {
	ctrl     *gomock.Controller
	recorder *MockClosableScopeMockRecorder
}

// MockClosableScopeMockRecorder is the mock recorder for MockClosableScope.
type MockClosableScopeMockRecorder struct {
	mock *MockClosableScope
}

// NewMockClosableScope creates a new mock instance.
func NewMockClosableScope(ctrl *gomock.Controller) *MockClosableScope {
	mock := &MockClosableScope{ctrl: ctrl}
	mock.recorder = &MockClosableScopeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClosableScope) EXPECT() *MockClosableScopeMockRecorder {
	return m.recorder
}

// Capabilities mocks base method.
func (m *MockClosableScope) Capabilities() metrics.Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capabilities")
	ret0, _ := ret[0].(metrics.Capabilities)
	return ret0
}

// Capabilities indicates an expected call of Capabilities.
func (mr *MockClosableScopeMockRecorder) Capabilities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capabilities", reflect.TypeOf((*MockClosableScope)(nil).Capabilities))
}

// Counter mocks base method.
func (m *MockClosableScope) Counter(name string) metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Counter", name)
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// Counter indicates an expected call of Counter.
func (mr *MockClosableScopeMockRecorder) Counter(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Counter", reflect.TypeOf((*MockClosableScope)(nil).Counter), name)
}

// Gauge mocks base method.
func (m *MockClosableScope) Gauge(name string) metrics.Gauge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gauge", name)
	ret0, _ := ret[0].(metrics.Gauge)
	return ret0
}

// Gauge indicates an expected call of Gauge.
func (mr *MockClosableScopeMockRecorder) Gauge(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gauge", reflect.TypeOf((*MockClosableScope)(nil).Gauge), name)
}

// Histogram mocks base method.
func (m *MockClosableScope) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Histogram", name, buckets)
	ret0, _ := ret[0].(metrics.Histogram)
	return ret0
}

// Histogram indicates an expected call of Histogram.
func (mr *MockClosableScopeMockRecorder) Histogram(name, buckets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Histogram", reflect.TypeOf((*MockClosableScope)(nil).Histogram), name, buckets)
}

// Stop mocks base method.
func (m *MockClosableScope) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockClosableScopeMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockClosableScope)(nil).Stop))
}

// SubScope mocks base method.
func (m *MockClosableScope) SubScope(name string) metrics.Scope {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubScope", name)
	ret0, _ := ret[0].(metrics.Scope)
	return ret0
}

// SubScope indicates an expected call of SubScope.
func (mr *MockClosableScopeMockRecorder) SubScope(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubScope", reflect.TypeOf((*MockClosableScope)(nil).SubScope), name)
}

// Tagged mocks base method.
func (m *MockClosableScope) Tagged(tags map[string]string) metrics.Scope {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tagged", tags)
	ret0, _ := ret[0].(metrics.Scope)
	return ret0
}

// Tagged indicates an expected call of Tagged.
func (mr *MockClosableScopeMockRecorder) Tagged(tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tagged", reflect.TypeOf((*MockClosableScope)(nil).Tagged), tags)
}

// Timer mocks base method.
func (m *MockClosableScope) Timer(name string) metrics.Timer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timer", name)
	ret0, _ := ret[0].(metrics.Timer)
	return ret0
}

// Timer indicates an expected call of Timer.
func (mr *MockClosableScopeMockRecorder) Timer(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timer", reflect.TypeOf((*MockClosableScope)(nil).Timer), name)
}

// MockCounter is a mock of Counter interface.
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter.
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance.
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Inc mocks base method.
func (m *MockCounter) Inc(delta int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc", delta)
}

// Inc indicates an expected call of Inc.
func (mr *MockCounterMockRecorder) Inc(delta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc), delta)
}

// MockGauge is a mock of Gauge interface.
type MockGauge struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeMockRecorder
}

// MockGaugeMockRecorder is the mock recorder for MockGauge.
type MockGaugeMockRecorder struct {
	mock *MockGauge
}

// NewMockGauge creates a new mock instance.
func NewMockGauge(ctrl *gomock.Controller) *MockGauge {
	mock := &MockGauge{ctrl: ctrl}
	mock.recorder = &MockGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGauge) EXPECT() *MockGaugeMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockGauge) Update(value float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", value)
}

// Update indicates an expected call of Update.
func (mr *MockGaugeMockRecorder) Update(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGauge)(nil).Update), value)
}

// MockTimer is a mock of Timer interface.
type MockTimer struct {
	ctrl     *gomock.Controller
	recorder *MockTimerMockRecorder
}

// MockTimerMockRecorder is the mock recorder for MockTimer.
type MockTimerMockRecorder struct {
	mock *MockTimer
}

// NewMockTimer creates a new mock instance.
func NewMockTimer(ctrl *gomock.Controller) *MockTimer {
	mock := &MockTimer{ctrl: ctrl}
	mock.recorder = &MockTimerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimer) EXPECT() *MockTimerMockRecorder {
	return m.recorder
}

// Record mocks base method.
func (m *MockTimer) Record(value time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Record", value)
}

// Record indicates an expected call of Record.
func (mr *MockTimerMockRecorder) Record(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Record", reflect.TypeOf((*MockTimer)(nil).Record), value)
}

// Start mocks base method.
func (m *MockTimer) Start() metrics.Stopwatch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(metrics.Stopwatch)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockTimerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTimer)(nil).Start))
}

// MockHistogram is a mock of Histogram interface.
type MockHistogram struct {
	ctrl     *gomock.Controller
	recorder *MockHistogramMockRecorder
}

// MockHistogramMockRecorder is the mock recorder for MockHistogram.
type MockHistogramMockRecorder struct {
	mock *MockHistogram
}

// NewMockHistogram creates a new mock instance.
func NewMockHistogram(ctrl *gomock.Controller) *MockHistogram {
	mock := &MockHistogram{ctrl: ctrl}
	mock.recorder = &MockHistogramMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistogram) EXPECT() *MockHistogramMockRecorder {
	return m.recorder
}

// RecordDuration mocks base method.
func (m *MockHistogram) RecordDuration(value time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecordDuration", value)
}

// RecordDuration indicates an expected call of RecordDuration.
func (mr *MockHistogramMockRecorder) RecordDuration(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordDuration", reflect.TypeOf((*MockHistogram)(nil).RecordDuration), value)
}

// RecordValue mocks base method.
func (m *MockHistogram) RecordValue(value float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecordValue", value)
}

// RecordValue indicates an expected call of RecordValue.
func (mr *MockHistogramMockRecorder) RecordValue(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordValue", reflect.TypeOf((*MockHistogram)(nil).RecordValue), value)
}

// Start mocks base method.
func (m *MockHistogram) Start() metrics.Stopwatch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(metrics.Stopwatch)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockHistogramMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockHistogram)(nil).Start))
}

// MockStopwatch is a mock of Stopwatch interface.
type MockStopwatch struct {
	ctrl     *gomock.Controller
	recorder *MockStopwatchMockRecorder
}

// MockStopwatchMockRecorder is the mock recorder for MockStopwatch.
type MockStopwatchMockRecorder struct {
	mock *MockStopwatch
}

// NewMockStopwatch creates a new mock instance.
func NewMockStopwatch(ctrl *gomock.Controller) *MockStopwatch {
	mock := &MockStopwatch{ctrl: ctrl}
	mock.recorder = &MockStopwatchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStopwatch) EXPECT() *MockStopwatchMockRecorder {
	return m.recorder
}

// Stop mocks base method.
func (m *MockStopwatch) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockStopwatchMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockStopwatch)(nil).Stop))
}

// MockStopwatchRecorder is a mock of StopwatchRecorder interface.
type MockStopwatchRecorder struct {
	ctrl     *gomock.Controller
	recorder *MockStopwatchRecorderMockRecorder
}

// MockStopwatchRecorderMockRecorder is the mock recorder for MockStopwatchRecorder.
type MockStopwatchRecorderMockRecorder struct {
	mock *MockStopwatchRecorder
}

// NewMockStopwatchRecorder creates a new mock instance.
func NewMockStopwatchRecorder(ctrl *gomock.Controller) *MockStopwatchRecorder {
	mock := &MockStopwatchRecorder{ctrl: ctrl}
	mock.recorder = &MockStopwatchRecorderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStopwatchRecorder) EXPECT() *MockStopwatchRecorderMockRecorder {
	return m.recorder
}

// RecordStopwatch mocks base method.
func (m *MockStopwatchRecorder) RecordStopwatch(stopwatchStart time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecordStopwatch", stopwatchStart)
}

// RecordStopwatch indicates an expected call of RecordStopwatch.
func (mr *MockStopwatchRecorderMockRecorder) RecordStopwatch(stopwatchStart interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordStopwatch", reflect.TypeOf((*MockStopwatchRecorder)(nil).RecordStopwatch), stopwatchStart)
}

// MockBuckets is a mock of Buckets interface.
type MockBuckets struct {
	ctrl     *gomock.Controller
	recorder *MockBucketsMockRecorder
}

// MockBucketsMockRecorder is the mock recorder for MockBuckets.
type MockBucketsMockRecorder struct {
	mock *MockBuckets
}

// NewMockBuckets creates a new mock instance.
func NewMockBuckets(ctrl *gomock.Controller) *MockBuckets {
	mock := &MockBuckets{ctrl: ctrl}
	mock.recorder = &MockBucketsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuckets) EXPECT() *MockBucketsMockRecorder {
	return m.recorder
}

// AsDurations mocks base method.
func (m *MockBuckets) AsDurations() []time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AsDurations")
	ret0, _ := ret[0].([]time.Duration)
	return ret0
}

// AsDurations indicates an expected call of AsDurations.
func (mr *MockBucketsMockRecorder) AsDurations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsDurations", reflect.TypeOf((*MockBuckets)(nil).AsDurations))
}

// AsValues mocks base method.
func (m *MockBuckets) AsValues() []float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AsValues")
	ret0, _ := ret[0].([]float64)
	return ret0
}

// AsValues indicates an expected call of AsValues.
func (mr *MockBucketsMockRecorder) AsValues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsValues", reflect.TypeOf((*MockBuckets)(nil).AsValues))
}

// Len mocks base method.
func (m *MockBuckets) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockBucketsMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockBuckets)(nil).Len))
}

// Less mocks base method.
func (m *MockBuckets) Less(i, j int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Less", i, j)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Less indicates an expected call of Less.
func (mr *MockBucketsMockRecorder) Less(i, j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Less", reflect.TypeOf((*MockBuckets)(nil).Less), i, j)
}

// String mocks base method.
func (m *MockBuckets) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockBucketsMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockBuckets)(nil).String))
}

// Swap mocks base method.
func (m *MockBuckets) Swap(i, j int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Swap", i, j)
}

// Swap indicates an expected call of Swap.
func (mr *MockBucketsMockRecorder) Swap(i, j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Swap", reflect.TypeOf((*MockBuckets)(nil).Swap), i, j)
}

// MockCapabilities is a mock of Capabilities interface.
type MockCapabilities struct {
	ctrl     *gomock.Controller
	recorder *MockCapabilitiesMockRecorder
}

// MockCapabilitiesMockRecorder is the mock recorder for MockCapabilities.
type MockCapabilitiesMockRecorder struct {
	mock *MockCapabilities
}

// NewMockCapabilities creates a new mock instance.
func NewMockCapabilities(ctrl *gomock.Controller) *MockCapabilities {
	mock := &MockCapabilities{ctrl: ctrl}
	mock.recorder = &MockCapabilitiesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCapabilities) EXPECT() *MockCapabilitiesMockRecorder {
	return m.recorder
}

// Reporting mocks base method.
func (m *MockCapabilities) Reporting() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reporting")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Reporting indicates an expected call of Reporting.
func (mr *MockCapabilitiesMockRecorder) Reporting() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reporting", reflect.TypeOf((*MockCapabilities)(nil).Reporting))
}

// Tagging mocks base method.
func (m *MockCapabilities) Tagging() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tagging")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Tagging indicates an expected call of Tagging.
func (mr *MockCapabilitiesMockRecorder) Tagging() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tagging", reflect.TypeOf((*MockCapabilities)(nil).Tagging))
}
