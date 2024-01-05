package redis

import (
	"context"
	"crypto/tls"
	"github.com/Shopify/toxiproxy/client"
	"github.com/devlibx/gox-base/ratelimit"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"sync"
	"testing"
)

var ratelimitMutexToRunOneTestAtATime = sync.Mutex{}

type rateLimitE2ETestSuite struct {
	suite.Suite
	redisUrl        string
	url             string
	toxiProxyClient *toxiproxy.Client
	tProxy          *toxiproxy.Proxy
	id              string
}

func (s *rateLimitE2ETestSuite) SetupTest() {
	var err error
	s.id = uuid.NewString()

	s.redisUrl = "localhost:6379"
	s.toxiProxyClient = toxiproxy.NewClient("localhost:8474")

	s.tProxy, err = s.toxiProxyClient.CreateProxy("tests-resdis-"+s.id, "localhost:26379", s.redisUrl)
	assert.NoError(s.T(), err)
	s.url = "localhost:26379"
}

func (s *rateLimitE2ETestSuite) AfterTest(suiteName, testName string) {
	if s.tProxy != nil {
		_ = s.tProxy.Delete()
	}
}

func (rt *rateLimitE2ETestSuite) getRedisClient(url string) *goRedis.ClusterClient {
	client := goRedis.NewClusterClient(&goRedis.ClusterOptions{
		Addrs:        strings.Split(url, ","),
		PoolSize:     10,
		MinIdleConns: 10,
		IdleTimeout:  10,
		Password:     os.Getenv("REDIS_PASS"),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	result := client.Ping(context.Background())
	assert.NoError(rt.T(), result.Err())
	return client
}

// TestRateLimitE2E to test rate limit e2e
func (rt *rateLimitE2ETestSuite) TestRateLimitE2E() {
	ratelimitMutexToRunOneTestAtATime.Lock()
	defer ratelimitMutexToRunOneTestAtATime.Unlock()

	client := rt.getRedisClient(rt.redisUrl)
	lg := NewLimitGroup(&ratelimit.Config{GroupName: "test_e2e", Enabled: true, LimitPerSec: 100, RetryCount: 2}, client)
	out, err := lg.Allow(context.Background(), func() (interface{}, error) { return "ok", nil })
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

	client := goRedis.NewClusterClient(&goRedis.ClusterOptions{
		Addrs:        strings.Split(rt.url, ","),
		PoolSize:     10,
		MinIdleConns: 10,
		IdleTimeout:  10,
		Password:     os.Getenv("REDIS_PASS"),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	result := client.Ping(context.Background())
	assert.NoError(rt.T(), result.Err())

	lg := NewLimitGroup(&ratelimit.Config{
		GroupName:   "test_e2e",
		Enabled:     true,
		LimitPerSec: 100,
		RetryCount:  2,
	}, client)

	out, err := lg.Allow(context.Background(), func() (interface{}, error) {
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

	client := goRedis.NewClusterClient(&goRedis.ClusterOptions{
		Addrs:        strings.Split(rt.url, ","),
		PoolSize:     10,
		MinIdleConns: 10,
		IdleTimeout:  10,
		Password:     os.Getenv("REDIS_PASS"),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	lg := NewLimitGroup(&ratelimit.Config{
		GroupName:   "test_e2e",
		Enabled:     true,
		LimitPerSec: 100,
		RetryCount:  2,
	}, client)

	out, err := lg.Allow(context.Background(), func() (interface{}, error) {
		return "ok", nil
	})

	// In this test we have connection error - here we will run our function
	assert.NoError(rt.T(), err)
	assert.Equal(rt.T(), "ok", out)
}

func TestRateLimitE2ETestSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS_ENABLED") == "" {
		t.SkipNow()
	}
	suite.Run(t, new(rateLimitE2ETestSuite))
}
