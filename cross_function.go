package gox

import (
	"context"
	"time"

	"github.com/devlibx/gox-base/v2/metrics"
	"github.com/devlibx/gox-base/v2/util"
	"go.uber.org/zap"
)

// Implementation of cross function
type crossFunction struct {
	logger *zap.Logger
	metrics.Scope
	TimeService
	config      StringObjectMap
	timeTracker util.TimeTracker
	publisher   metrics.Publisher
}

func (c *crossFunction) Publisher() metrics.Publisher {
	return c.publisher
}

func (c *crossFunction) TimeTracker() util.TimeTracker {
	return c.timeTracker
}

func (c *crossFunction) Metric() metrics.Scope {
	return c.Scope
}

func (c *crossFunction) Logger() *zap.Logger {
	return c.logger
}

func (c *crossFunction) Config() StringObjectMap {
	return c.config
}

func (c *crossFunction) CurrentTime(ctx context.Context, input any) time.Time {
	if EnableInterceptorInDefaultTimeService {
		if ts, ok := c.TimeService.(InterceptableTimeService); ok {
			return ts.CurrentTime(ctx, input)
		}
	}
	return c.Now()
}

// NewCrossFunction a no-op cross function object which does not have a side effect
func NewCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{}
	for _, arg := range args {
		switch o := arg.(type) {
		case *zap.Logger:
			obj.logger = o
		case metrics.Scope:
			obj.Scope = o
		case StringObjectMap:
			obj.config = o
		case util.TimeTracker:
			obj.timeTracker = o
		case metrics.Publisher:
			obj.publisher = o
		case TimeService:
			obj.TimeService = o
		}
	}

	// Set default time-service
	if obj.TimeService == nil {
		obj.TimeService = &DefaultTimeService{}
	}

	// Setup no-op logger if it is not passed
	if obj.logger == nil {
		obj.logger = zap.NewNop()
	}

	// Setup no-op metrics
	if obj.Scope == nil {
		obj.Scope = metrics.NoOpMetric()
	}

	// Set default config
	if obj.config == nil {
		obj.config = StringObjectMap{}
	}

	// Set dummy trim tracker if not provided
	if obj.timeTracker == nil {
		obj.timeTracker = util.NewNoOpTimeTracker()
	}

	// Setup no-op publisher
	if obj.publisher == nil {
		obj.publisher = metrics.NewNoOpPublisher()
	}

	return &obj
}

// A No Op cross function
func NewNoOpCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{TimeService: &DefaultTimeService{}}
	obj.logger = zap.NewNop()
	obj.TimeService = &DefaultTimeService{}
	obj.Scope = metrics.NoOpMetric()
	obj.config = StringObjectMap{}
	obj.timeTracker = util.NewNoOpTimeTracker()
	obj.publisher = metrics.NewNoOpPublisher()
	return &obj
}

// InterceptableCurrentTime this is a helper on top of CF to give current time
func InterceptableCurrentTime(cf CrossFunction, input any) time.Time {
	if !EnableInterceptorInDefaultTimeService {
		return cf.Now()
	}

	var timeToReturn time.Time
	if tf, ok := cf.(InterceptableTimeService); ok {
		timeToReturn = tf.CurrentTime(context.Background(), input)
	} else {
		timeToReturn = cf.Now()
	}
	return timeToReturn
}
