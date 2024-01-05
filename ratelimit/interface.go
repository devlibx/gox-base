package ratelimit

import "context"

// RateLimiter is a function which will allow an action to be performed - this is used to
type RateLimiter interface {
	Allow(ctx context.Context, toRun RateLimitedFunc) (out interface{}, err error)
}

type Configs struct {
	Enabled bool               `json:"enabled"`
	Configs map[string]*Config `json:"groups"`
}

// Config is the configuration for the rate limiter - each group may have its own name
// and max retry count
type Config struct {
	Enabled     bool   `json:"enabled"`
	GroupName   string `json:"group_name"`
	LimitPerSec int    `json:"limit_per_sec"`
	RetryCount  int    `json:"retry_count"`
}

// RateLimitedFunc is the function which is rate limited - this impl will make sure
// rate limited function is tried for N times and then if not able to get through, it will
// return error
type RateLimitedFunc func() (interface{}, error)

// noOpRateLimiter is a mock impl
type noOpRateLimiter struct{}

func (n noOpRateLimiter) Allow(ctx context.Context, toRun RateLimitedFunc) (out interface{}, err error) {
	return toRun()
}

func NewNoOpRateLimiter() RateLimiter {
	return &noOpRateLimiter{}
}
