package gox

import (
	"github.com/devlibx/gox-base/metrics"
	"go.uber.org/zap"
)

// Implementation of cross function
type crossFunction struct {
	logger *zap.Logger
	metrics.Scope
	TimeService
}

func (c *crossFunction) Metric() metrics.Scope {
	return c.Scope
}

func (c *crossFunction) Logger() *zap.Logger {
	return c.logger
}

// Create a new cross function object
func NewCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{}
	for _, arg := range args {
		switch o := arg.(type) {
		case *zap.Logger:
			obj.logger = o
		case metrics.Scope:
			obj.Scope = o
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
	obj.logger = zap.NewNop()
	obj.TimeService = &DefaultTimeService{}
	obj.Scope = metrics.NoOpMetric()
	return &obj
}
