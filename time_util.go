package gox

import "time"

type DefaultTimeService struct {
	TimeService
}

func (t *DefaultTimeService) Now() time.Time {
	return time.Now()
}

func (t *DefaultTimeService) Sleep(d time.Duration) {
	time.Sleep(d)
}
