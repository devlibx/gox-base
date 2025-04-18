package ratelimiterRedis

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/Shopify/toxiproxy/client"
	"github.com/devlibx/gox-base/v2/ratelimiter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var ratelimitMutexToRunOneTestAtATime = sync.Mutex{}

type rateLimitE2ETestSuite struct {
	suite.Suite
	redisUrl        string
	url             string
	toxiProxyClient *toxiproxy.Client
	tProxy          *toxiproxy.Proxy
	id              string
	factory         RateLimitFactory
}

func (s *rateLimitE2ETestSuite) SetupTest() {
	var err error
	s.id = uuid.NewString()

	s.redisUrl = getEnvOrDefault(EnvRedisHost, DefaultRedisURL)
	s.toxiProxyClient = toxiproxy.NewClient("localhost:8474")

	s.tProxy, err = s.toxiProxyClient.CreateProxy("tests-resdis-"+s.id, "localhost:26379", s.redisUrl)
	assert.NoError(s.T(), err)
	s.url = "localhost:26379"

	// Create factory with default Redis config
	s.factory = NewRateLimitFactory(&ratelimiter.Configs{
		Enabled: true,
		Configs: map[string]*ratelimiter.Config{
			"test_e2e": {
				Enabled:     true,
				GroupName:   "test_e2e",
				LimitPerSec: 100,
				RetryCount:  2,
				Redis: &ratelimiter.RedisConfig{
					URLs:     []string{s.redisUrl},
					Password: os.Getenv(EnvRedisPassword),
				},
			},
		},
	})
}

func (s *rateLimitE2ETestSuite) AfterTest(suiteName, testName string) {
	if s.tProxy != nil {
		_ = s.tProxy.Delete()
	}

	// Close factory
	if f, ok := s.factory.(*rateLimitFactory); ok {
		f.Close()
	}
}

// TestRateLimitE2E to test rate limit e2e
func (rt *rateLimitE2ETestSuite) TestRateLimitE2E() {
	ratelimitMutexToRunOneTestAtATime.Lock()
	defer ratelimitMutexToRunOneTestAtATime.Unlock()

	// Create factory with specific config for this test
	factory := NewRateLimitFactory(&ratelimiter.Configs{
		Enabled: true,
		Configs: map[string]*ratelimiter.Config{
			"test_e2e": {
				Enabled:     true,
				GroupName:   "test_e2e",
				LimitPerSec: 100,
				RetryCount:  2,
				Redis: &ratelimiter.RedisConfig{
					URLs:     []string{rt.redisUrl},
					Password: os.Getenv(EnvRedisPassword),
				},
			},
		},
	})
	defer func() {
		if f, ok := factory.(*rateLimitFactory); ok {
			f.Close()
		}
	}()

	// Get rate limiter
	limiter := factory.GetRateLimiter("test_e2e")

	// Test rate limiter
	out, err := limiter.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})
	assert.NoError(rt.T(), err)
	assert.Equal(rt.T(), "ok", out)
}

// TestRateLimitE2E to test rate limit e2e
// NOTE - must install and start toxiProxy before running this test
func (rt *rateLimitE2ETestSuite) TestRateLimitE2E_WithLatency_1sec_latency() {
	ratelimitMutexToRunOneTestAtATime.Lock()
	defer ratelimitMutexToRunOneTestAtATime.Unlock()

	id := uuid.NewString()
	_, err := rt.tProxy.AddToxic("tests_redis_latency_down_"+id, "latency", "downstream", 1.0, toxiproxy.Attributes{
		"latency": 1000,
	})
	assert.NoError(rt.T(), err)
	defer rt.tProxy.RemoveToxic("tests_redis_latency_down_" + id)

	// Create factory with specific config for this test
	factory := NewRateLimitFactory(&ratelimiter.Configs{
		Enabled: true,
		Configs: map[string]*ratelimiter.Config{
			"test_e2e": {
				Enabled:     true,
				GroupName:   "test_e2e",
				LimitPerSec: 100,
				RetryCount:  2,
				Redis: &ratelimiter.RedisConfig{
					URLs:     []string{rt.url}, // Use toxiproxy URL
					Password: os.Getenv(EnvRedisPassword),
				},
			},
		},
	})
	defer func() {
		if f, ok := factory.(*rateLimitFactory); ok {
			f.Close()
		}
	}()

	// Get rate limiter
	limiter := factory.GetRateLimiter("test_e2e")

	// Test rate limiter
	out, err := limiter.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})
	assert.NoError(rt.T(), err)
	assert.Equal(rt.T(), "ok", out)
}

// TestRateLimitE2E_WithErrorFromRedis to test rate limit e2e
// NOTE - must install and start toxiProxy before running this test
func (rt *rateLimitE2ETestSuite) TestRateLimitE2E_RedisIsDown() {
	ratelimitMutexToRunOneTestAtATime.Lock()
	defer ratelimitMutexToRunOneTestAtATime.Unlock()

	id := uuid.NewString()
	rt.tProxy.Disable()
	_ = id

	// Create factory with specific config for this test
	factory := NewRateLimitFactory(&ratelimiter.Configs{
		Enabled: true,
		Configs: map[string]*ratelimiter.Config{
			"test_e2e": {
				Enabled:     true,
				GroupName:   "test_e2e",
				LimitPerSec: 100,
				RetryCount:  2,
				Redis: &ratelimiter.RedisConfig{
					URLs:     []string{rt.url}, // Use toxiproxy URL (which is disabled)
					Password: os.Getenv(EnvRedisPassword),
				},
			},
		},
	})
	defer func() {
		if f, ok := factory.(*rateLimitFactory); ok {
			f.Close()
		}
	}()

	// Get rate limiter
	limiter := factory.GetRateLimiter("test_e2e")

	// Test rate limiter
	out, err := limiter.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})

	// In this test we have connection error - here we will run our function
	// The factory should return a no-op limiter when Redis is down
	assert.NoError(rt.T(), err)
	assert.Equal(rt.T(), "ok", out)
}

func TestRateLimitE2ETestSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS_ENABLED") == "" {
		t.SkipNow()
	}
	suite.Run(t, new(rateLimitE2ETestSuite))
}
