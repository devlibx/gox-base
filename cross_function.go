package gox

import (
	"github.com/harishb2k/gox-base/logger"
	"github.com/harishb2k/gox-base/metrics"
)

type Metrics metrics.Service

// A holder to keep all cross function objects e.g. logger etc
type CrossFunction interface {
	logger.Logger
	Metrics
	TimeService
}

// Implementation of cross function
type crossFunction struct {
	logger.Logger
	Metrics
	TimeService
}

// Create a new cross function object
func NewCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{}
	for _, arg := range args {
		switch o := arg.(type) {
		case logger.Logger:
			obj.Logger = o
		case metrics.Service:
			obj.Metrics = o
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
	obj.Logger = logger.NoOpLogger(logger.Configuration{})
	obj.TimeService = &DefaultTimeService{}
	obj.Metrics = metrics.NewNoOpMetrics()
	return &obj
}
