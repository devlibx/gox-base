package gox

import (
	"github.com/devlibx/gox-base/logger"
	"github.com/devlibx/gox-base/metrics"
)

// A holder to keep all cross function objects e.g. logger etc
type CrossFunction interface {
	logger.Logger
	metrics.MetricService
	TimeService
}

// Implementation of cross function
type crossFunction struct {
	logger.Logger
	metrics.MetricService
	TimeService
}

// Create a new cross function object
func NewCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{}
	for _, arg := range args {
		switch o := arg.(type) {
		case logger.Logger:
			obj.Logger = o
		case metrics.MetricService:
			obj.MetricService = o
		}
	}

	// Set default time-service
	if obj.TimeService == nil {
		obj.TimeService = &DefaultTimeService{}
	}

	return &obj
}

// A No Op cross function
func NewNoOpCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{TimeService: &DefaultTimeService{}}
	obj.Logger = logger.NewNoopLogger()
	obj.TimeService = &DefaultTimeService{}
	obj.MetricService = metrics.NewNoOpMetrics()
	return &obj
}
