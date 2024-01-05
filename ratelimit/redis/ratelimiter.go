package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/devlibx/gox-base/ratelimit"
	goRedis "github.com/go-redis/redis/v8"
	goRedisRate "github.com/go-redis/redis_rate/v9"
	"net"
	"time"
)

//go:generate mockgen -destination=rate_limiter_mock.go -package=redis -source=ratelimiter.go

var _ ratelimit.RateLimiter = &limitGroup{}
var _ actionAllower = &goRedisRate.Limiter{}

// ActionAllower is the interface to allow an action to be performed - this is used to
// abstract the rate limiter implementation
type actionAllower interface {
	Allow(ctx context.Context, key string, limit goRedisRate.Limit) (*goRedisRate.Result, error)
}

type limitGroup struct {
	cfg     *ratelimit.Config
	limiter actionAllower
}

const (
	RetryNeeded    = true
	RetryNotNeeded = false
)

func (l *limitGroup) checkRetryNeeded(result *goRedisRate.Result, err error) (time.Duration, bool) {

	if err == nil {
		return time.Millisecond, RetryNotNeeded
	}

	// We have error and result is not nil - this means we need to wait for some time (given by result.RetryAfter)
	if result != nil && err != nil {
		return result.RetryAfter, RetryNeeded
	}

	// Some connection error - we do not want to fail due to redis connection errors
	// So we will have safe default and will run the action
	var connError *net.OpError
	if errors.As(err, &connError) {
		return time.Duration(1 * time.Millisecond), RetryNotNeeded
	}

	// Some unknown issue - we will retry with 10 sec delay
	if result == nil && err != nil {
		return time.Duration(10 * time.Millisecond), RetryNeeded
	}

	// No need to retry
	return time.Millisecond, RetryNotNeeded
}

// Allow will allow an action to be performed - this is used to
// It take a function to run - which returns a result and error.
// If this function is not able to run due to rate limit exceeded then it will return error
// If it is able to run then it will return result which is returned by the function itself
func (l *limitGroup) Allow(ctx context.Context, toRun ratelimit.RateLimitedFunc) (out interface{}, err error) {

	// Try for N times (n = retry count) to get to allow=ok from rate limiter
	for i := 0; i < l.cfg.RetryCount; i++ {

		var result *goRedisRate.Result
		result, err = l.limiter.Allow(ctx, l.cfg.GroupName, goRedisRate.PerSecond(l.cfg.LimitPerSec))

		// Check if we should retry or now
		delay, retryNeeded := l.checkRetryNeeded(result, err)
		if retryNeeded == RetryNeeded {
			time.Sleep(delay)
		} else {
			return toRun()
		}
	}

	// Result based on the success status
	return nil, &Error{
		Err:       err,
		Message:   fmt.Sprintf("rate limit exceeded: group=%s", l.cfg.GroupName),
		ErrorCode: "failed",
	}
}

func NewLimitGroup(cfg *ratelimit.Config, redisClient *goRedis.ClusterClient) ratelimit.RateLimiter {
	if !cfg.Enabled {
		return ratelimit.NewNoOpRateLimiter()
	}

	limiter := goRedisRate.NewLimiter(redisClient)
	return &limitGroup{
		cfg:     cfg,
		limiter: limiter,
	}
}

type Error struct {
	Err       error
	Message   string
	ErrorCode string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error_message=%s, error_code=%s, error=%s", e.Message, e.ErrorCode, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
