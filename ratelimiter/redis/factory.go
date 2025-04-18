package ratelimiterRedis

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/devlibx/gox-base/v2/ratelimiter"
)

const (
	// DefaultRedisURL is the default Redis URL if none is specified
	DefaultRedisURL = "localhost:6379"

	// DefaultRedisPort is the default Redis port
	DefaultRedisPort = "6379"

	// Environment variables for Redis configuration
	EnvRedisHost         = "REDIS_HOST"
	EnvRedisPassword     = "REDIS_PASSWORD"
	EnvRedisReadTimeout  = "REDIS_READ_TIMEOUT"
	EnvRedisWriteTimeout = "REDIS_WRITE_TIMEOUT"
	EnvRedisPutTimeout   = "REDIS_PUT_TIMEOUT"
	EnvRedisGetTimeout   = "REDIS_GET_TIMEOUT"
)

// getEnvOrDefault gets an environment variable or returns the default value
func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets an environment variable as int or returns the default value
func getEnvIntOrDefault(envVar string, defaultValue int) int {
	if value := os.Getenv(envVar); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// ensureRedisPort ensures that the Redis URL has a port
func ensureRedisPort(url string) string {
	// If URL already contains a port, return as is
	if strings.Contains(url, ":") {
		return url
	}
	// Add default port
	return url + ":" + DefaultRedisPort
}

// getRedisURLsFromEnv gets Redis URLs from environment variables
func getRedisURLsFromEnv() []string {
	hostEnv := os.Getenv(EnvRedisHost)
	if hostEnv == "" {
		return []string{DefaultRedisURL}
	}

	// Split URLs and ensure each has a port
	urls := strings.Split(hostEnv, ",")
	for i, url := range urls {
		urls[i] = ensureRedisPort(strings.TrimSpace(url))
	}
	return urls
}

// RateLimitFactory provides a factory interface to create and manage rate limiters.
// It supports creating rate limiters based on configuration and caches them for reuse.
//
// The factory implements a fail-safe approach to rate limiting:
// 1. If rate limiting is globally disabled (Configs.Enabled = false), all operations are allowed
// 2. If a group's configuration is not found, a no-op rate limiter is returned that allows all operations
// 3. If a specific group's rate limiting is disabled (Config.Enabled = false), a no-op rate limiter is returned
// 4. Only when both global rate limiting is enabled AND group configuration exists AND is enabled, rate limiting is applied
type RateLimitFactory interface {
	// GetRateLimiter returns a rate limiter for the given group name.
	//
	// The method implements a fail-safe approach:
	// 1. If the group exists in configuration and is enabled, returns a properly configured rate limiter
	// 2. If the group doesn't exist in configuration or is disabled, returns a no-op rate limiter that allows all operations
	// 3. If rate limiting is globally disabled, always returns a no-op rate limiter
	//
	// This behavior ensures that missing configurations or disabled rate limiting won't block operations.
	//
	// Parameters:
	//   - groupName: The name of the rate limit group to get/create a limiter for
	//
	// Returns:
	//   - A RateLimiter instance that can be used to rate limit operations
	GetRateLimiter(groupName string) ratelimiter.RateLimiter

	// Close closes all resources managed by the factory.
	// This should be called when the factory is no longer needed to release resources.
	// For the no-op factory, this is a no-op.
	Close() error
}

type rateLimitFactory struct {
	configs      *ratelimiter.Configs
	limiters     map[string]ratelimiter.RateLimiter
	redisClients map[string]RedisRateLimiterClient
	mu           sync.RWMutex
}

// NewRateLimitFactory creates a new rate limit factory that manages rate limiters based on configuration.
// The factory creates and manages Redis clients internally based on the configuration.
//
// The factory maintains a thread-safe cache of rate limiters and Redis clients to avoid recreating them.
// If rate limiting is globally disabled (Configs.Enabled = false), it returns a no-op factory
// that allows all operations without rate limiting.
//
// Parameters:
//   - configs: Rate limiting configuration containing group-specific settings and Redis configurations
//
// Returns:
//   - A RateLimitFactory implementation that can create and manage rate limiters
func NewRateLimitFactory(configs *ratelimiter.Configs) RateLimitFactory {
	if configs == nil || !configs.Enabled {
		return NewNoOpRateLimitFactory()
	}

	return &rateLimitFactory{
		configs:      configs,
		limiters:     make(map[string]ratelimiter.RateLimiter),
		redisClients: make(map[string]RedisRateLimiterClient),
	}
}

func (f *rateLimitFactory) GetRateLimiter(groupName string) ratelimiter.RateLimiter {
	// First try to get from cache
	f.mu.RLock()
	if limiter, exists := f.limiters[groupName]; exists {
		f.mu.RUnlock()
		return limiter
	}
	f.mu.RUnlock()

	// If not in cache, create new with write lock
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double check after acquiring write lock
	if limiter, exists := f.limiters[groupName]; exists {
		return limiter
	}

	// Get config for the group
	config, exists := f.configs.Configs[groupName]
	if !exists || !config.Enabled {
		// If no config found or rate limiter is disabled, return a no-op rate limiter
		return ratelimiter.NewNoOpRateLimiter()
	}

	// Get or create Redis client for this group
	redisClient, err := f.getOrCreateRedisClient(config)
	if err != nil {
		// If Redis client creation fails, return a no-op rate limiter
		// This ensures the application continues to function even if Redis is unavailable
		return ratelimiter.NewNoOpRateLimiter()
	}

	// Create new limiter
	limiter := NewLimitGroup(config, redisClient)
	f.limiters[groupName] = limiter
	return limiter
}

// getOrCreateRedisClient gets an existing Redis client or creates a new one
// Note: This method assumes the caller already holds the mutex lock
func (f *rateLimitFactory) getOrCreateRedisClient(config *ratelimiter.Config) (RedisRateLimiterClient, error) {
	// Get Redis config from the group or use environment variables
	var redisConfig *ratelimiter.RedisConfig
	if config.Redis != nil {
		redisConfig = config.Redis
		// Ensure all URLs have ports
		for i, url := range redisConfig.URLs {
			redisConfig.URLs[i] = ensureRedisPort(url)
		}
	} else {
		// Create Redis config from environment variables
		redisConfig = &ratelimiter.RedisConfig{
			URLs:     getRedisURLsFromEnv(),
			Password: os.Getenv(EnvRedisPassword),
		}
	}

	// Generate a key for this Redis configuration
	key := redisClientKey(redisConfig)

	// Check if we already have a client for this configuration
	if client, exists := f.redisClients[key]; exists {
		return client, nil
	}

	// Create new client
	client, err := NewRedisClient(redisConfig)
	if err != nil {
		return nil, err
	}

	// Cache the client
	f.redisClients[key] = client
	return client, nil
}

// NoOpRateLimitFactory provides a rate limiter factory that always returns no-op rate limiters.
// This is useful in several scenarios:
// 1. When rate limiting is globally disabled
// 2. In test environments where you want to bypass rate limiting
// 3. When you want to explicitly disable rate limiting for certain components
type NoOpRateLimitFactory struct{}

// NewNoOpRateLimitFactory creates a new factory that always returns no-op rate limiters
func NewNoOpRateLimitFactory() RateLimitFactory {
	return &NoOpRateLimitFactory{}
}

func (f *NoOpRateLimitFactory) GetRateLimiter(groupName string) ratelimiter.RateLimiter {
	return ratelimiter.NewNoOpRateLimiter()
}

// Close is a no-op for NoOpRateLimitFactory
func (f *NoOpRateLimitFactory) Close() error {
	return nil
}

// Close closes all Redis clients managed by the factory
func (f *rateLimitFactory) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var lastErr error
	for _, client := range f.redisClients {
		if err := client.Close(); err != nil {
			lastErr = err
		}
	}

	// Clear the maps
	f.redisClients = make(map[string]RedisRateLimiterClient)
	f.limiters = make(map[string]ratelimiter.RateLimiter)

	return lastErr
}
