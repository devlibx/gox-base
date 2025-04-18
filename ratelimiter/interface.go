package ratelimiter

import "context"

// RateLimiter is a function which will allow an action to be performed - this is used to
type RateLimiter interface {
	Allow(ctx context.Context, toRun RateLimitedFunc) (out interface{}, err error)
}

// RedisConfig holds the configuration for connecting to Redis
type RedisConfig struct {
	// Redis connection URLs (e.g., ["redis-1:6379", "redis-2:6379"])
	URLs []string `json:"urls" yaml:"urls"`

	// Redis password (optional)
	Password string `json:"password" yaml:"password"`

	// Whether to use Redis Cluster (if false, uses standalone Redis)
	UseCluster bool `json:"use_cluster" yaml:"use_cluster"`

	// TLS configuration (optional)
	UseTLS bool `json:"use_tls" yaml:"use_tls"`

	// Connection pool size (default: 10)
	PoolSize int `json:"pool_size" yaml:"pool_size"`

	// Minimum idle connections (default: 10)
	MinIdleConns int `json:"min_idle_conns" yaml:"min_idle_conns"`

	// Read timeout in milliseconds (default: 100)
	ReadTimeoutMs int `json:"read_timeout_ms" yaml:"read_timeout_ms"`

	// Write timeout in milliseconds (default: 100)
	WriteTimeoutMs int `json:"write_timeout_ms" yaml:"write_timeout_ms"`

	// Put timeout in milliseconds (default: 100)
	PutTimeoutMs int `json:"put_timeout_ms" yaml:"put_timeout_ms"`

	// Get timeout in milliseconds (default: 100)
	GetTimeoutMs int `json:"get_timeout_ms" yaml:"get_timeout_ms"`
}

// Configs holds the global rate limiting configuration
type Configs struct {
	// Whether rate limiting is globally enabled
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Group-specific configurations
	Configs map[string]*Config `json:"groups" yaml:"groups"`
}

// Config is the configuration for a specific rate limiter group
type Config struct {
	// Whether this rate limiter is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Unique name for this rate limit group
	GroupName string `json:"group_name" yaml:"group_name"`

	// Rate limits (only one should be set)
	LimitPerSec  int `json:"limit_per_sec" yaml:"limit_per_sec"`
	LimitPerMin  int `json:"limit_per_min" yaml:"limit_per_min"`
	LimitPerHour int `json:"limit_per_hour" yaml:"limit_per_hour"`
	LimitPerDay  int `json:"limit_per_day" yaml:"limit_per_day"`

	// Number of retry attempts when rate limited
	RetryCount int `json:"retry_count" yaml:"retry_count"`

	// Whether to fail immediately without retrying when rate limited
	NoRetryToAcquire bool `json:"no_retry_to_acquire" yaml:"no_retry_to_acquire"`

	// Redis configuration specific to this group (overrides default)
	Redis *RedisConfig `json:"redis" yaml:"redis"`
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
