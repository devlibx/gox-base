package test

import (
	"flag"
	"github.com/devlibx/gox-base/metrics"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"time"
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

	cf.Metric().Counter("").Inc(1)

	s := cf.Metric().Timer("").Start()
	s.Stop()

	cf.Metric().Gauge("g").Update(10)

	histogram := cf.Metric().Histogram("g", metrics.DefaultBuckets)
	hsw := histogram.Start()
	time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
	hsw.Stop()
}
