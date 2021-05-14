package test

import (
	"github.com/devlibx/gox-base"
	mockGox "github.com/devlibx/gox-base/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

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
