package redis

import (
	"context"
	"errors"
	"github.com/devlibx/gox-base/ratelimit"
	"github.com/go-redis/redis_rate/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRateLimit_WithAllRetries_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionAllower := NewMockactionAllower(ctrl)
	mockActionAllower.EXPECT().
		Allow(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&redis_rate.Result{RetryAfter: 100 * time.Millisecond}, errors.New("not allowed")).
		AnyTimes()

	rl := &limitGroup{
		limiter: mockActionAllower,
		cfg: &ratelimit.Config{
			GroupName:   "test",
			Enabled:     true,
			LimitPerSec: 3,
			RetryCount:  3,
		}}

	out, err := rl.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})
	var rateLimitError *Error
	assert.ErrorAs(t, err, &rateLimitError)
	assert.Nil(t, out)
}

func TestRateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionAllower := NewMockactionAllower(ctrl)
	mockActionAllower.EXPECT().
		Allow(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&redis_rate.Result{RetryAfter: 100 * time.Millisecond}, nil).
		AnyTimes()

	rl := &limitGroup{
		limiter: mockActionAllower,
		cfg: &ratelimit.Config{
			GroupName:   "test",
			Enabled:     true,
			LimitPerSec: 3,
			RetryCount:  3,
		}}

	out, err := rl.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok", out)
}

func TestRateLimit_WithRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	count := 0
	mockActionAllower := NewMockactionAllower(ctrl)
	mockActionAllower.EXPECT().
		Allow(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error) {
			defer func() {
				count++
			}()
			if count == 0 {
				return &redis_rate.Result{RetryAfter: 100}, errors.New("error")
			}
			return &redis_rate.Result{RetryAfter: 100}, nil
		}).AnyTimes()

	rl := &limitGroup{
		limiter: mockActionAllower,
		cfg: &ratelimit.Config{
			GroupName:   "test",
			Enabled:     true,
			LimitPerSec: 3,
			RetryCount:  3,
		}}

	out, err := rl.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok", out)
	assert.Equal(t, 2, count)
}
