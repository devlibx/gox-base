# Gox-Base

[![Go Reference](https://pkg.go.dev/badge/github.com/devlibx/gox-base.svg)](https://pkg.go.dev/github.com/devlibx/gox-base)
[![Go Report Card](https://goreportcard.com/badge/github.com/devlibx/gox-base)](https://goreportcard.com/report/github.com/devlibx/gox-base)

A comprehensive Go utility library providing essential building blocks for enterprise applications.

## Features

- **Serialization Utilities**
  - JSON, YAML, and XML serialization/deserialization
  - Support for parameterized YAML configurations
  - Environment variable substitution in config files

- **Rate Limiting**
  - Redis-based distributed rate limiting
  - Support for multiple rate limit types (per second/minute/hour/day)
  - Fail-safe design with configurable fallback behavior
  - [Detailed Documentation](ratelimiter/README.md)

- **HTTP Server Setup**
  - Easy-to-use REST application setup
  - Built-in OpenTracing support
  - Configurable middleware

- **Configuration Management**
  - Environment-aware configuration
  - Type-safe configuration parsing
  - Extensible configuration structures

- **Logging**
  - Structured logging with Uber's zap
  - Configurable output formats and destinations
  - Module-level logging support

## Installation

```bash
go get github.com/devlibx/gox-base/v2
```

## Quick Start

### Setting up a REST Server

```go
import (
    "github.com/gorilla/mux"
    "github.com/devlibx/gox-base/config"
    goxServer "github.com/devlibx/gox-base/server"
)

func main() {
    // Create router
    router := mux.NewRouter()

    // Configure your application
    appConfig := config.App{
        AppName:     "my_app",
        HttpPort:    8080,
        Environment: "dev",
    }

    // Start server
    serverInstance, _ := goxServer.NewServer(cf)
    if err := serverInstance.Start(router, &appConfig); err != nil {
        log.Fatal(err)
    }
}
```

### Using the Rate Limiter

```go
import (
    "github.com/devlibx/gox-base/v2/ratelimiter"
    "github.com/devlibx/gox-base/v2/ratelimiter/redis"
)

// Create rate limiter configuration
configs := &ratelimiter.Configs{
    Enabled: true,
    Configs: map[string]*ratelimiter.Config{
        "api": {
            Enabled:     true,
            GroupName:   "api",
            LimitPerSec: 100,
            RetryCount:  3,
            Redis: &ratelimiter.RedisConfig{
                URLs: []string{"localhost:6379"},
            },
        },
    },
}

// Create factory and get rate limiter
factory := ratelimiterRedis.NewRateLimitFactory(configs)
defer factory.Close()

limiter := factory.GetRateLimiter("api")

// Use rate limiter
result, err := limiter.Allow(ctx, func() (interface{}, error) {
    return doSomething()
})
```

### Working with JSON

The library provides powerful JSON utility functions for common operations:

#### Converting Between Types

```go
import "github.com/devlibx/gox-base/v2/serialization/utils/json"

// Convert JSON string to struct
jsonStr := `{"field1": "value1", "field2": 123}`
result, err := goxJsonUtils.StringToObject[MyStruct](jsonStr)

// Convert struct to JSON string
myObj := MyStruct{Field1: "value1", Field2: 123}
jsonStr, err := goxJsonUtils.ObjectToString(myObj)

// Convert to pretty-printed JSON (useful for logging)
prettyJson, err := goxJsonUtils.PrettyString(myObj)

// Working with bytes
jsonBytes := []byte(`{"field1": "value1", "field2": 123}`)
result, err := goxJsonUtils.BytesToObject[MyStruct](jsonBytes)

// Convert object to bytes
bytes, err := goxJsonUtils.ObjectToBytes(myObj)
```

#### Type Conversion Utilities

The library handles various type conversions automatically:

```go
// Convert primitive types to string
intStr, _ := goxJsonUtils.ObjectToString(42)          // "42"
boolStr, _ := goxJsonUtils.ObjectToString(true)       // "true"
floatStr, _ := goxJsonUtils.ObjectToString(3.14)      // "3.14"

// Convert complex objects
type User struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
}
user := User{Name: "John", Age: 30}
jsonStr, _ := goxJsonUtils.ObjectToString(user)  // {"name":"John","age":30}
```

#### Working with StringObjectMap

```go
// Convert JSON to map
jsonStr := `{"key1": "value1", "key2": 123}`
mapObj, err := goxJsonUtils.StringToStringObjectMap(jsonStr)

// Convert map to struct
type MyStruct struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}
mapObj := gox.StringObjectMap{"field1": "value1", "field2": 123}
result, err := goxJsonUtils.StringObjectMapToObject[MyStruct](mapObj)
```

#### Error-Suppressing Variants

For cases where you want to handle errors gracefully:

```go
// Convert without error handling (returns zero value on error)
result := goxJsonUtils.StringToObjectSuppressError[MyStruct](jsonStr)
mapObj := goxJsonUtils.StringToStringObjectMapSuppressError(jsonStr)
```

### Working with Configuration

```go
type Config struct {
    App    config.App    `yaml:"app"`
    Logger config.Logger `yaml:"logger"`
}

// Read from YAML
conf := &Config{}
err := serialization.ReadYamlFromString(yamlConfig, conf)
```

## Documentation

- [Rate Limiter Documentation](ratelimiter/README.md)
- [Configuration Guide](docs/configuration.md)
- [Logging Guide](docs/logging.md)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
