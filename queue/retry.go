package queue

import "time"

type FixedDelayRetryBackoffAlgo struct {
	fixedDelay time.Duration
}

func (d *FixedDelayRetryBackoffAlgo) NextRetryAfter(attempt int, maxExecution int) (time.Duration, error) {
	if attempt > maxExecution {
		return time.Hour, ErrNoMoreRetry
	}
	return d.fixedDelay, nil
}

func NewDefaultRetryBackoffAlgo(fixedDelay time.Duration) RetryBackoffAlgo {
	return &FixedDelayRetryBackoffAlgo{fixedDelay: fixedDelay}
}
