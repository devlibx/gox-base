package gox

import (
	"github.com/devlibx/gox-base/metrics"
	"github.com/devlibx/gox-base/util"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mockGox
//go:generate mockgen -source=metrics/interface.go -destination=mocks/mock_metrics_interface.go -package=mockGox

type TimeService interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// A holder to keep all cross function objects e.g. logger etc
type CrossFunction interface {
	Logger() *zap.Logger
	Metric() metrics.Scope
	TimeService
	Config() StringObjectMap
	TimeTracker() util.TimeTracker
}
