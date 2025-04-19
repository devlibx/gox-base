# Rate Limiter

A Redis-based distributed rate limiter implementation with support for multiple Redis configurations and fail-safe behavior.

> **⚠️ IMPORTANT NOTES:**
> 1. **Redis Dependency**: This rate limiter requires a Redis server/cluster. Without Redis, it falls back to a no-op implementation.
> 2. **Fail-Safe Design**: By default, if Redis is unavailable, all operations are allowed to prevent application failures.
> 3. **Performance Impact**: Configure appropriate timeouts and retry counts based on your latency requirements.
> 4. **Resource Management**: Always call `Close()` on the factory to properly release Redis connections.
> 5. **Cluster Mode**: When using Redis Cluster, ensure all nodes are properly configured and reachable.

## Features

- Redis-based distributed rate limiting
- Support for both standalone Redis and Redis Cluster
- Multiple rate limit types:
  - Per second
  - Per minute
  - Per hour
  - Per day
- Fail-safe behavior:
  - Falls back to no-op rate limiter if Redis is unavailable
  - Configurable retry behavior
  - Fast-fail option with NoRetryToAcquire
- Thread-safe client management
- Support for multiple Redis configurations within the same application

## Usage

### Basic Usage

```go
import (
    "github.com/devlibx/gox-base/v2/ratelimiter"
    "github.com/devlibx/gox-base/v2/ratelimiter/redis"
)

// Create configuration
configs := &ratelimiter.Configs{
    Enabled: true,
    Configs: map[string]*ratelimiter.Config{
        "api": {
            Enabled:     true,
            GroupName:   "api",
            LimitPerSec: 100,
            RetryCount:  3,
            Redis: &ratelimiter.RedisConfig{
                URLs:     []string{"redis-1:6379", "redis-2:6379"},
                Password: "password",
            },
        },
    },
}

// Create factory
factory := ratelimiterRedis.NewRateLimitFactory(configs)
defer factory.Close()

// Get rate limiter
limiter := factory.GetRateLimiter("api")

// Use rate limiter
result, err := limiter.Allow(ctx, func() (interface{}, error) {
    return doSomething()
})
```

### Configuration

#### Redis Configuration

```go
type RedisConfig struct {
    URLs           []string // Redis URLs (e.g., ["redis-1:6379", "redis-2:6379"])
    Password       string   // Redis password
    UseCluster     bool     // Use Redis Cluster mode
    UseTLS         bool     // Enable TLS
    PoolSize       int      // Connection pool size (default: 10)
    MinIdleConns   int      // Minimum idle connections (default: 10)
    ReadTimeoutMs  int      // Read timeout in milliseconds (default: 100)
    WriteTimeoutMs int      // Write timeout in milliseconds (default: 100)
    PutTimeoutMs   int      // Put timeout in milliseconds (default: 100)
    GetTimeoutMs   int      // Get timeout in milliseconds (default: 100)
}
```

> **⚠️ Note:** If port is not specified in URLs, default port 6379 will be automatically added.

#### Rate Limit Configuration

```go
type Config struct {
    Enabled          bool         // Enable/disable rate limiting for this group
    LimitPerSec      int         // Rate limit per second
    LimitPerMin      int         // Rate limit per minute
    LimitPerHour     int         // Rate limit per hour
    LimitPerDay      int         // Rate limit per day
    RetryCount       int         // Number of retries when rate limit is exceeded
    NoRetryToAcquire bool        // Fast-fail without retrying
    Redis            *RedisConfig // Redis configuration for this group
}
```

> **⚠️ Note:** Only one of LimitPerSec, LimitPerMin, LimitPerHour, or LimitPerDay should be set.

### Environment Variables

- `REDIS_HOST` - Comma-separated list of Redis URLs (default: "localhost:6379")
- `REDIS_PASSWORD` - Redis password

> **⚠️ Note:** These environment variables are primarily used for tests. In production code, you should provide Redis configuration directly in your Config objects.
> **⚠️ Note:** For cluster mode, provide multiple Redis URLs separated by commas in `REDIS_HOST`.

### YAML Configuration

You can also define your rate limiter configuration in YAML:

```yaml
# Rate limiting configuration
enabled: true  # Enable/disable rate limiting globally

# Group-specific configurations
groups:
  api-group:
    enabled: true
    limit_per_sec: 100  # Rate limit per second
    retry_count: 3      # Number of retries when rate limited
    redis:
      urls:
        - localhost:6379  # Redis URLs (can be multiple for cluster)
      password: ""        # Redis password (optional)
      use_cluster: false  # Use Redis cluster mode
      use_tls: false     # Use TLS for Redis connection
      pool_size: 10      # Connection pool size
      min_idle_conns: 5  # Minimum idle connections

  high-throughput:
    enabled: true
    limit_per_min: 1000  # Rate limit per minute
    retry_count: 5
    no_retry_to_acquire: true  # Fail immediately without retrying
    redis:
      urls:
        - redis-1:6379
        - redis-2:6379
      password: "your-password"
      use_cluster: true
      use_tls: true
```

Load the YAML configuration:

```go
import (
    "os"
    "gopkg.in/yaml.v3"
    "github.com/devlibx/gox-base/v2/ratelimiter"
    ratelimiterRedis "github.com/devlibx/gox-base/v2/ratelimiter/redis"
)

// Read and parse YAML config
data, err := os.ReadFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

var configs ratelimiter.Configs
if err := yaml.Unmarshal(data, &configs); err != nil {
    log.Fatal(err)
}

// Create factory with YAML config
factory := ratelimiterRedis.NewRateLimitFactory(&configs)
defer factory.Close()

// Example 1: Using the api-group rate limiter
apiLimiter := factory.GetRateLimiter("api-group")
result, err := apiLimiter.Allow(context.Background(), func() (interface{}, error) {
    // Your API operation here
    return "API operation successful", nil
})
if err != nil {
    log.Printf("API operation failed: %v", err)
} else {
    log.Printf("API result: %v", result)
}

// Example 2: Using the high-throughput rate limiter with fast-fail
htLimiter := factory.GetRateLimiter("high-throughput")
result, err = htLimiter.Allow(context.Background(), func() (interface{}, error) {
    // Your high-throughput operation here
    return "High throughput operation completed", nil
})
if err != nil {
    // This might fail fast due to no_retry_to_acquire: true
    log.Printf("High throughput operation failed: %v", err)
} else {
    log.Printf("High throughput result: %v", result)
}
```

### Multiple Redis Configurations

You can use different Redis configurations for different rate limit groups:

```go
configs := &ratelimiter.Configs{
    Enabled: true,
    Configs: map[string]*ratelimiter.Config{
        "api": {
            Enabled:     true,
            LimitPerSec: 100,
            Redis: &ratelimiter.RedisConfig{
                URLs: []string{"redis-1:6379"},
            },
        },
        "background": {
            Enabled:     true,
            LimitPerMin: 1000,
            Redis: &ratelimiter.RedisConfig{
                URLs:       []string{"redis-cluster:6379"},
                UseCluster: true,
            },
        },
    },
}
```

> **⚠️ Note:** Each group can have its own Redis configuration, allowing for different Redis servers/clusters per group.

### Fast-Fail Behavior

To fail immediately without retrying when rate limit is exceeded:

```go
config := &ratelimiter.Config{
    Enabled:          true,
    GroupName:        "fast-fail",
    LimitPerSec:      100,
    RetryCount:       3,
    NoRetryToAcquire: true, // Enable fast-fail
}
```

> **⚠️ Note:** When NoRetryToAcquire is true, the rate limiter will fail immediately if the first attempt to acquire the rate limit fails.

### Fail-Safe Behavior

The rate limiter implements fail-safe behavior:

1. If rate limiting is globally disabled (Configs.Enabled = false), all operations are allowed
2. If a group's configuration is not found, a no-op rate limiter is returned
3. If Redis is unavailable, a no-op rate limiter is returned
4. If a specific group's rate limiting is disabled, a no-op rate limiter is returned

> **⚠️ IMPORTANT:** This fail-safe behavior means your application will continue to function even if Redis is unavailable. However, this also means that rate limiting will not be enforced in these scenarios. Monitor your Redis health and application metrics to ensure rate limiting is working as expected.

This ensures that your application continues to function even if Redis is unavailable or misconfigured.
