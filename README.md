# Gox-Base

[![Go Reference](https://pkg.go.dev/badge/github.com/devlibx/gox-base.svg)](https://pkg.go.dev/github.com/devlibx/gox-base)
[![Go Report Card](https://goreportcard.com/badge/github.com/devlibx/gox-base)](https://goreportcard.com/report/github.com/devlibx/gox-base)

A comprehensive Go utility library providing essential building blocks for enterprise applications.

## Features

- **Serialization Utilities**
  - JSON, YAML, and XML serialization/deserialization
  - Support for parameterized YAML configurations
  - Environment variable substitution in config files

- **SQL Database Utilities**
  - Type-safe conversions between Go types and `sql.NullString`
  - Built-in encryption/decryption support for sensitive data
  - JSON serialization for complex structs in database fields
  - Comprehensive error handling with contextual messages

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

### SQL Database Utilities

The library provides comprehensive utilities for working with SQL database operations, particularly for converting between Go types and `sql.NullString` types, with support for encryption/decryption and JSON serialization.

#### Basic String Conversions

```go
import goxSql "github.com/devlibx/gox-base/v2/database/sql"

// Convert string to sql.NullString (always Valid=true)
nullStr := goxSql.StringToSqlNullString("hello world")
// Result: sql.NullString{String: "hello world", Valid: true}

// Convert sql.NullString to string
str := goxSql.SqlNullStringToString(nullStr)
// Result: "hello world"

// Handle invalid NullString (represents NULL in database)
invalidNull := sql.NullString{String: "some value", Valid: false}
str := goxSql.SqlNullStringToString(invalidNull)
// Result: "" (empty string)
```

#### Encrypted Data Storage

For secure data storage, use the encryption-enabled functions:

```go
// Define your encryption service (must implement the interface)
type MyEncryptor struct{}
func (e *MyEncryptor) EncryptAndOutputBase64Ciphertext(data string) (string, error) {
    // Your encryption logic here
    return encryptedData, nil
}

type MyDecryptor struct{}
func (d *MyDecryptor) DecryptFromBase64Ciphertext(data string) (string, error) {
    // Your decryption logic here
    return decryptedData, nil
}

// Encrypt and store sensitive data
encryptor := &MyEncryptor{}
encryptedNull, err := goxSql.StringToEncryptedSqlNullString("sensitive data", encryptor)
if err != nil {
    log.Fatal(err)
}
// Store encryptedNull.String in database

// Decrypt data when reading from database
decryptor := &MyDecryptor{}
originalData, err := goxSql.EncryptedSqlNullStringToString(encryptedNull, decryptor)
if err != nil {
    log.Fatal(err)
}
// originalData contains "sensitive data"
```

#### Encrypted Struct Serialization

For maximum security, combine struct serialization with encryption for storing sensitive complex data:

```go
// Define your sensitive data structures
type UserProfile struct {
    Name         string   `json:"name"`
    Email        string   `json:"email"`
    PersonalData string   `json:"personal_data"`
    Preferences  []string `json:"preferences"`
}

type FinancialData struct {
    AccountNumber string  `json:"account_number"`
    Balance       float64 `json:"balance"`
    Transactions  []Transaction `json:"transactions"`
}

// Setup encryption service (using the built-in AES encryption)
key, err := encryption.GenerateAESKey(32)
if err != nil {
    log.Fatal(err)
}

encConfig := &encryption.EncryptDecryptConfigs{
    Group: map[string]*encryption.EncryptDecryptConfig{
        "user_data": {
            Algo: "aes_32",
            AesConfig: &encryption.AesConfig{
                Base64CodedKey: base64.StdEncoding.EncodeToString(key),
            },
        },
    },
}

factory, err := encryption.NewServiceFactory(gox.NewNoOpCrossFunction(), encConfig)
if err != nil {
    log.Fatal(err)
}

encryptorDecryptor, err := factory.GetEncryptorDecryptService("user_data")
if err != nil {
    log.Fatal(err)
}

// Encrypt and store sensitive struct data
profile := UserProfile{
    Name:         "John Doe",
    Email:        "john@example.com",
    PersonalData: "SSN: 123-45-6789",
    Preferences:  []string{"privacy", "security"},
}

encryptedProfileData, err := goxSql.StructToEncryptedSqlNullString(profile, encryptorDecryptor)
if err != nil {
    log.Fatal(err)
}
// encryptedProfileData.String contains fully encrypted JSON - not readable

// Store in database
query := `INSERT INTO users (id, encrypted_profile) VALUES (?, ?)`
_, err = db.ExecContext(ctx, query, userID, encryptedProfileData)
if err != nil {
    log.Fatal(err)
}

// Retrieve and decrypt struct data
var encryptedData sql.NullString
query = `SELECT encrypted_profile FROM users WHERE id = ?`
err = db.QueryRowContext(ctx, query, userID).Scan(&encryptedData)
if err != nil {
    log.Fatal(err)
}

// Decrypt back to original struct
var retrievedProfile UserProfile
retrievedProfile, err = goxSql.EncryptedSqlNullStringToStruct[UserProfile](encryptedData, encryptorDecryptor)
if err != nil {
    log.Fatal(err)
}
// retrievedProfile now contains the original decrypted UserProfile data
```

#### JSON Serialization for Complex Types

Store and retrieve complex Go structs as JSON in database fields:

```go
// Define your data structures
type User struct {
    Name string   `json:"name"`
    Age  int      `json:"age"`
    Tags []string `json:"tags,omitempty"`
}

type UserProfile struct {
    User     User    `json:"user"`
    Active   bool    `json:"active"`
    Score    float64 `json:"score"`
}

// Serialize struct to sql.NullString
user := User{
    Name: "John Doe",
    Age:  30,
    Tags: []string{"admin", "developer"},
}

nullStr, err := goxSql.StructToSqlNullString[User](user)
if err != nil {
    log.Fatal(err)
}
// nullStr.String contains: {"name":"John Doe","age":30,"tags":["admin","developer"]}

// Deserialize sql.NullString back to struct
var retrievedUser User
retrievedUser, err = goxSql.SqlNullStringToStruct[User](nullStr)
if err != nil {
    log.Fatal(err)
}
// retrievedUser now contains the original user data
```

#### Complete Database Integration Example

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    goxSql "github.com/devlibx/gox-base/v2/database/sql"
)

type UserRepository struct {
    db *sql.DB
}

// Store user with encrypted email and JSON metadata
func (r *UserRepository) CreateUser(ctx context.Context, user User, email string, metadata UserProfile) error {
    // Convert struct to JSON for storage
    metadataJSON, err := goxSql.StructToSqlNullString[UserProfile](metadata)
    if err != nil {
        return err
    }
    
    // Encrypt sensitive email data
    encryptor := &MyEncryptor{}
    encryptedEmail, err := goxSql.StringToEncryptedSqlNullString(email, encryptor)
    if err != nil {
        return err
    }
    
    // Convert regular fields
    userName := goxSql.StringToSqlNullString(user.Name)
    
    query := `INSERT INTO users (name, encrypted_email, metadata_json) VALUES (?, ?, ?)`
    _, err = r.db.ExecContext(ctx, query, userName, encryptedEmail, metadataJSON)
    return err
}

// Retrieve user with decryption and deserialization
func (r *UserRepository) GetUser(ctx context.Context, userID int) (*User, string, *UserProfile, error) {
    var userName, encryptedEmail, metadataJSON sql.NullString
    
    query := `SELECT name, encrypted_email, metadata_json FROM users WHERE id = ?`
    err := r.db.QueryRowContext(ctx, query, userID).Scan(&userName, &encryptedEmail, &metadataJSON)
    if err != nil {
        return nil, "", nil, err
    }
    
    // Convert basic field
    name := goxSql.SqlNullStringToString(userName)
    
    // Decrypt email
    decryptor := &MyDecryptor{}
    email, err := goxSql.EncryptedSqlNullStringToString(encryptedEmail, decryptor)
    if err != nil {
        return nil, "", nil, err
    }
    
    // Deserialize JSON metadata
    var metadata UserProfile
    metadata, err = goxSql.SqlNullStringToStruct[UserProfile](metadataJSON)
    if err != nil {
        return nil, "", nil, err
    }
    
    user := &User{Name: name}
    return user, email, &metadata, nil
}
```

#### Error Handling

All utility functions provide comprehensive error handling:

```go
// Encryption errors are wrapped with context
_, err := goxSql.StringToEncryptedSqlNullString("data", failingEncryptor)
if err != nil {
    // Error contains: "unable to encrypt data while converting sql to sql.NullString: <original error>"
}

// JSON serialization errors are wrapped
_, err = goxSql.StructToSqlNullString[MyStruct](invalidStruct)
if err != nil {
    // Error contains: "unable to serialize data to JSON: <original error>"
}

// Invalid JSON causes deserialization errors
invalidJSON := sql.NullString{String: "{invalid json}", Valid: true}
_, err = goxSql.SqlNullStringToStruct[MyStruct](invalidJSON)
// err will contain JSON parsing error details
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

# Dummy PR for testing
