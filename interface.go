package gox

import (
	"github.com/devlibx/gox-base/metrics"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mock_gox

type TimeService interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// A holder to keep all cross function objects e.g. logger etc
type CrossFunction interface {
	Logger() *zap.Logger
	metrics.MetricService
	TimeService
}
