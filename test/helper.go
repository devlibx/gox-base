package test

import (
	"github.com/devlibx/gox-base"
	mockGox "github.com/devlibx/gox-base/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"testing"
)

func MockCf(args ...interface{}) (gox.CrossFunction, *gomock.Controller) {
	var controller *gomock.Controller
	var logger *zap.Logger
	logLevel := zap.ErrorLevel

	// See if we have got a controller from client
	for _, arg := range args {
		switch o := arg.(type) {
		case *gomock.Controller:
			controller = o
		case zapcore.Level:
			logLevel = o
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

	// Build cross function and return
	cf := mockGox.NewMockCrossFunction(controller)
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	return cf, controller
}

func BuildMockCf(t *testing.T, controller *gomock.Controller) gox.CrossFunction {
	cf := mockGox.NewMockCrossFunction(controller)
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	return cf
}

func BuildMockCfB(b *testing.B, controller *gomock.Controller) gox.CrossFunction {
	cf := mockGox.NewMockCrossFunction(controller)
	logger := zaptest.NewLogger(b, zaptest.Level(zap.ErrorLevel))
	cf.EXPECT().Logger().Return(logger).AnyTimes()
	return cf
}
