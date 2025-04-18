package ratelimiterRedis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/devlibx/gox-base/v2/ratelimiter"
	goRedisRate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

const (
	// Default values for Redis configuration
	DefaultPoolSize     = 10
	DefaultMinIdleConns = 10
	DefaultTimeoutMs    = 100
)

// redisClientKey creates a unique key for caching Redis clients
func redisClientKey(config *ratelimiter.RedisConfig) string {
	return strings.Join(config.URLs, ",") + ":" + config.Password + ":" +
		strings.Join([]string{
			strings.Join(config.URLs, ","),
			config.Password,
			boolToString(config.UseCluster),
			boolToString(config.UseTLS),
		}, ":")
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// RedisRateLimiterClient wraps a Redis client with rate limiting capabilities
type RedisRateLimiterClient interface {
	// Ping checks if the Redis server is available
	Ping(ctx context.Context) *redis.StatusCmd

	// CreateLimiter creates a new rate limiter using the Redis client
	CreateLimiter() *goRedisRate.Limiter

	// Close closes the underlying Redis client
	Close() error
}

// NewRedisClient creates a new Redis client based on the configuration
func NewRedisClient(config *ratelimiter.RedisConfig) (RedisRateLimiterClient, error) {
	if config == nil {
		return nil, nil
	}

	if config.UseCluster {
		return NewClusterRedisClient(config)
	}
	return NewStandaloneRedisClient(config)
}

// ClusterRedisClient implements RedisRateLimiterClient for Redis Cluster
type ClusterRedisClient struct {
	client *redis.ClusterClient
}

// applyDefaultValues applies default values to Redis configuration
func applyDefaultValues(config *ratelimiter.RedisConfig) {
	// Apply default pool size
	if config.PoolSize <= 0 {
		config.PoolSize = DefaultPoolSize
	}

	// Apply default min idle connections
	if config.MinIdleConns <= 0 {
		config.MinIdleConns = DefaultMinIdleConns
	}

	// Apply default timeouts
	if config.ReadTimeoutMs <= 0 {
		config.ReadTimeoutMs = DefaultTimeoutMs
	}

	if config.WriteTimeoutMs <= 0 {
		config.WriteTimeoutMs = DefaultTimeoutMs
	}

	if config.PutTimeoutMs <= 0 {
		config.PutTimeoutMs = DefaultTimeoutMs
	}

	if config.GetTimeoutMs <= 0 {
		config.GetTimeoutMs = DefaultTimeoutMs
	}

	// Ensure we have at least one URL
	if len(config.URLs) == 0 {
		config.URLs = getRedisURLsFromEnv()
	}
}

// NewClusterRedisClient creates a new ClusterRedisClient from configuration
func NewClusterRedisClient(config *ratelimiter.RedisConfig) (*ClusterRedisClient, error) {
	// Apply default values
	applyDefaultValues(config)

	// Ensure we have at least one URL
	if len(config.URLs) == 0 {
		return nil, fmt.Errorf("no Redis URLs provided")
	}

	opts := &redis.ClusterOptions{
		Addrs:        config.URLs,
		Password:     config.Password,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		ReadTimeout:  time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeoutMs) * time.Millisecond,
	}

	if config.UseTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := redis.NewClusterClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return &ClusterRedisClient{client: client}, nil
}

func (c *ClusterRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *ClusterRedisClient) CreateLimiter() *goRedisRate.Limiter {
	return goRedisRate.NewLimiter(c.client)
}

func (c *ClusterRedisClient) Close() error {
	return c.client.Close()
}

// StandaloneRedisClient implements RedisRateLimiterClient for standalone Redis
type StandaloneRedisClient struct {
	client *redis.Client
}

// NewStandaloneRedisClient creates a new StandaloneRedisClient from configuration
func NewStandaloneRedisClient(config *ratelimiter.RedisConfig) (*StandaloneRedisClient, error) {
	// Apply default values
	applyDefaultValues(config)

	// Ensure we have at least one URL
	if len(config.URLs) == 0 {
		return nil, fmt.Errorf("no Redis URLs provided")
	}

	// For standalone Redis, we use the first URL
	opts := &redis.Options{
		Addr:         config.URLs[0],
		Password:     config.Password,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		ReadTimeout:  time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeoutMs) * time.Millisecond,
	}

	if config.UseTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return &StandaloneRedisClient{client: client}, nil
}

func (c *StandaloneRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *StandaloneRedisClient) CreateLimiter() *goRedisRate.Limiter {
	return goRedisRate.NewLimiter(c.client)
}

func (c *StandaloneRedisClient) Close() error {
	return c.client.Close()
}
