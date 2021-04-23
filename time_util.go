package gox

import "time"

//go:generate mockgen -source=time_util.go -destination=mocks/mock_time_util.go -package=mock_gox
type TimeService interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type DefaultTimeService struct {
	TimeService
}

func (t *DefaultTimeService) Now() time.Time {
	return time.Now()
}

func (t *DefaultTimeService) Sleep(d time.Duration) {
	time.Sleep(d)
}
