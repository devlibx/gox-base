package gox

import (
	"github.com/devlibx/gox-base/metrics"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestNewNoOpCrossFunction(t *testing.T) {
	cf := NewNoOpCrossFunction()

	assert.NotNil(t, cf)
}

func TestTypes_Mock(t *testing.T) {
	cf := NewNoOpCrossFunction()
	cf.Metric().Counter("").Inc(1)
	s := cf.Metric().Timer("").Start()
	s.Stop()
	cf.Metric().Gauge("g").Update(10)
	histogram := cf.Metric().Histogram("g", metrics.DefaultBuckets)
	hsw := histogram.Start()
	time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
	hsw.Stop()
}

func TestCrossFunctionWithNoOp(t *testing.T) {
	cf := NewCrossFunction()
	cf.Logger().Debug("nothing")
	cf.Metric().Counter("no").Inc(1)
}
