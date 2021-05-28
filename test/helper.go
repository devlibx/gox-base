package test

import (
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/metrics"
	mockGox "github.com/devlibx/gox-base/mocks"
	"github.com/devlibx/gox-base/util"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"testing"
)

func MockCf(args ...interface{}) (gox.CrossFunction, *gomock.Controller) {
	var controller *gomock.Controller
	var config gox.StringObjectMap
	var scope metrics.Scope
	var timeTracker util.TimeTracker
	var logger *zap.Logger
	logLevel := zap.ErrorLevel

	// See if we have got a controller from client
	for _, arg := range args {
		switch o := arg.(type) {
		case *gomock.Controller:
			controller = o
		case zapcore.Level:
			logLevel = o
		case metrics.Scope:
			scope = o
		case gox.StringObjectMap:
			config = o
		case util.TimeTracker:
			timeTracker = o
		}
	}

	// build dummy cf
	for _, arg := range args {
		switch o := arg.(type) {
		case *testing.T:
			if controller == nil {
				controller = gomock.NewController(o)
			}
			logger = zaptest.NewLogger(o, zaptest.Level(logLevel))
		case *testing.B:
			if controller == nil {
				controller = gomock.NewController(o)
			}
			logger = zaptest.NewLogger(o, zaptest.Level(logLevel))
		}
	}

	if scope == nil {
		scope = buildNoOfMetricsScope(controller)
	}

	if config == nil {
		config = gox.StringObjectMap{}
	}

	if timeTracker == nil {
		timeTracker = util.NewNoOpTimeTracker()
	}

	// Build cross function and return
	cf := mockGox.NewMockCrossFunction(controller)
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	cf.EXPECT().Metric().Return(scope).AnyTimes()
	cf.EXPECT().Config().Return(config).AnyTimes()
	cf.EXPECT().TimeTracker().Return(timeTracker).AnyTimes()
	return cf, controller
}

func BuildMockCf(t *testing.T, controller *gomock.Controller) gox.CrossFunction {
	cf := mockGox.NewMockCrossFunction(controller)
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	cf.EXPECT().Metric().Return(metrics.NoOpMetric()).AnyTimes()
	cf.EXPECT().Config().Return(gox.StringObjectMap{}).AnyTimes()
	cf.EXPECT().TimeTracker().Return(util.NewNoOpTimeTracker()).AnyTimes()
	return cf
}

func BuildMockCfB(b *testing.B, controller *gomock.Controller) gox.CrossFunction {
	cf := mockGox.NewMockCrossFunction(controller)
	logger := zaptest.NewLogger(b, zaptest.Level(zap.ErrorLevel))
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	cf.EXPECT().Metric().Return(metrics.NoOpMetric()).AnyTimes()
	cf.EXPECT().Config().Return(gox.StringObjectMap{}).AnyTimes()
	cf.EXPECT().Config().Return(util.NewNoOpTimeTracker()).AnyTimes()
	return cf
}

func buildNoOfMetricsScope(controller *gomock.Controller) metrics.Scope {
	mockScope := mockGox.NewMockScope(controller)
	mockCounter := mockGox.NewMockCounter(controller)
	mockGauge := mockGox.NewMockGauge(controller)
	mockTimer := mockGox.NewMockTimer(controller)
	mockHistogram := mockGox.NewMockHistogram(controller)

	mockScope.EXPECT().Counter(gomock.Any()).Return(mockCounter).AnyTimes()
	mockScope.EXPECT().Gauge(gomock.Any()).Return(mockGauge).AnyTimes()
	mockScope.EXPECT().Timer(gomock.Any()).Return(mockTimer).AnyTimes()
	mockScope.EXPECT().Histogram(gomock.Any(), gomock.Any()).Return(mockHistogram).AnyTimes()

	mockCounter.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockStopWatch := mockGox.NewMockStopwatch(controller)
	mockStopWatch.EXPECT().Stop().AnyTimes()
	mockTimer.EXPECT().Start().Return(mockStopWatch).AnyTimes()
	mockGauge.EXPECT().Update(gomock.Any()).AnyTimes()
	mockHistogram.EXPECT().Start().Return(mockStopWatch).AnyTimes()

	return mockScope
}
