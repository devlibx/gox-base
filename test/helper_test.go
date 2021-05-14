package test

import (
	"flag"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func init() {
	var ignore bool
	flag.BoolVar(&ignore, "dynamoTest", false, "run all database tests for dynamo")
	flag.BoolVar(&ignore, "test.real.schema", false, "run all database tests for dynamo")
}

func TestBuildMockCf(t *testing.T) {
	ctrl := gomock.NewController(t)
	cf := BuildMockCf(t, ctrl)
	cf.Logger().Debug("make sure I do not crash.")
}

func TestMockCf(t *testing.T) {
	cf, _ := MockCf(t, zap.DebugLevel)
	cf.Logger().Debug("make sure I do not crash.")
}
