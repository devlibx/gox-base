package encryption

import (
	"database/sql"
	"errors"
	"testing"

	goxSql "github.com/devlibx/gox-base/v2/database/sql"
	"github.com/stretchr/testify/assert"
)

// Mock encryptor for testing
type mockEncryptor struct {
	shouldError bool
	errorMsg    string
}

func (m *mockEncryptor) EncryptAndOutputBase64Ciphertext(data string) (string, error) {
	if m.shouldError {
		return "", errors.New(m.errorMsg)
	}
	// Simple mock encryption - just prepend "encrypted_" to simulate encryption
	return "encrypted_" + data, nil
}

// Mock decryptor for testing
type mockDecryptor struct {
	shouldError bool
	errorMsg    string
}

func (m *mockDecryptor) DecryptFromBase64Ciphertext(data string) (string, error) {
	if m.shouldError {
		return "", errors.New(m.errorMsg)
	}
	// Simple mock decryption - remove "encrypted_" prefix to simulate decryption
	if len(data) >= 10 && data[:10] == "encrypted_" {
		return data[10:], nil
	}
	return data, nil
}

func TestStringToSqlNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "non-empty string",
			input:    "hello world",
			expected: sql.NullString{String: "hello world", Valid: true},
		},
		{
			name:     "empty string",
			input:    "",
			expected: sql.NullString{String: "", Valid: true},
		},
		{
			name:     "unicode string",
			input:    "你好世界",
			expected: sql.NullString{String: "你好世界", Valid: true},
		},
		{
			name:     "string with special characters",
			input:    "hello\nworld\t!@#$%^&*()",
			expected: sql.NullString{String: "hello\nworld\t!@#$%^&*()", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goxSql.StringToSqlNullString(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.True(t, result.Valid)
		})
	}
}

func TestSqlNullStringToString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected string
	}{
		{
			name:     "valid non-empty string",
			input:    sql.NullString{String: "hello world", Valid: true},
			expected: "hello world",
		},
		{
			name:     "valid empty string",
			input:    sql.NullString{String: "", Valid: true},
			expected: "",
		},
		{
			name:     "invalid null string",
			input:    sql.NullString{String: "some value", Valid: false},
			expected: "",
		},
		{
			name:     "valid unicode string",
			input:    sql.NullString{String: "你好世界", Valid: true},
			expected: "你好世界",
		},
		{
			name:     "valid string with special characters",
			input:    sql.NullString{String: "hello\nworld\t!@#$%^&*()", Valid: true},
			expected: "hello\nworld\t!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goxSql.SqlNullStringToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringToEncryptedSqlNullString(t *testing.T) {
	t.Run("successful encryption with non-empty string", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		input := "sensitive data"
		
		result, err := goxSql.StringToEncryptedSqlNullString(input, encryptor)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Equal(t, "encrypted_sensitive data", result.String)
	})

	t.Run("successful encryption with empty string", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		input := ""
		
		result, err := goxSql.StringToEncryptedSqlNullString(input, encryptor)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Equal(t, "encrypted_", result.String)
	})

	t.Run("encryption error with non-empty string", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: true, errorMsg: "encryption failed"}
		input := "sensitive data"
		
		result, err := goxSql.StringToEncryptedSqlNullString(input, encryptor)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to encrypt data while converting sql to sql.NullString")
		assert.Contains(t, err.Error(), "encryption failed")
		assert.Equal(t, sql.NullString{}, result)
	})

	t.Run("encryption error with empty string", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: true, errorMsg: "encryption failed"}
		input := ""
		
		result, err := goxSql.StringToEncryptedSqlNullString(input, encryptor)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to encrypt data while converting sql to sql.NullString")
		assert.Contains(t, err.Error(), "encryption failed")
		assert.Equal(t, sql.NullString{}, result)
	})

	t.Run("encryption with unicode string", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		input := "你好世界"
		
		result, err := goxSql.StringToEncryptedSqlNullString(input, encryptor)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Equal(t, "encrypted_你好世界", result.String)
	})
}

func TestEncryptedSqlNullStringToString(t *testing.T) {
	t.Run("successful decryption with valid data", func(t *testing.T) {
		decrypter := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_sensitive data", Valid: true}
		
		result, err := goxSql.EncryptedSqlNullStringToString(input, decrypter)
		
		assert.NoError(t, err)
		assert.Equal(t, "sensitive data", result)
	})

	t.Run("successful decryption with empty encrypted data", func(t *testing.T) {
		decrypter := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_", Valid: true}
		
		result, err := goxSql.EncryptedSqlNullStringToString(input, decrypter)
		
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("invalid sql.NullString returns empty string", func(t *testing.T) {
		decrypter := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_some data", Valid: false}
		
		result, err := goxSql.EncryptedSqlNullStringToString(input, decrypter)
		
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("decryption error", func(t *testing.T) {
		decrypter := &mockDecryptor{shouldError: true, errorMsg: "decryption failed"}
		input := sql.NullString{String: "encrypted_sensitive data", Valid: true}
		
		result, err := goxSql.EncryptedSqlNullStringToString(input, decrypter)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to decrypt data from sql.NullString to string")
		assert.Contains(t, err.Error(), "decryption failed")
		assert.Equal(t, "", result)
	})

	t.Run("decryption with unicode data", func(t *testing.T) {
		decrypter := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_你好世界", Valid: true}
		
		result, err := goxSql.EncryptedSqlNullStringToString(input, decrypter)
		
		assert.NoError(t, err)
		assert.Equal(t, "你好世界", result)
	})
}

// Test struct for JSON serialization tests
type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Tags []string `json:"tags,omitempty"`
}

type TestNestedStruct struct {
	User    TestUser `json:"user"`
	Active  bool     `json:"active"`
	Score   float64  `json:"score"`
}

func TestSqlNullStringToStruct(t *testing.T) {
	t.Run("valid JSON to simple struct", func(t *testing.T) {
		jsonStr := `{"name":"John Doe","age":30,"tags":["admin","user"]}`
		input := sql.NullString{String: jsonStr, Valid: true}
		
		result, err := goxSql.SqlNullStringToStruct[TestUser](input)
		
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", result.Name)
		assert.Equal(t, 30, result.Age)
		assert.Equal(t, []string{"admin", "user"}, result.Tags)
	})

	t.Run("valid JSON to nested struct", func(t *testing.T) {
		jsonStr := `{"user":{"name":"Jane","age":25},"active":true,"score":95.5}`
		input := sql.NullString{String: jsonStr, Valid: true}
		
		result, err := goxSql.SqlNullStringToStruct[TestNestedStruct](input)
		
		assert.NoError(t, err)
		assert.Equal(t, "Jane", result.User.Name)
		assert.Equal(t, 25, result.User.Age)
		assert.True(t, result.Active)
		assert.Equal(t, 95.5, result.Score)
	})

	t.Run("invalid sql.NullString returns zero value", func(t *testing.T) {
		input := sql.NullString{String: `{"name":"test"}`, Valid: false}
		
		result, err := goxSql.SqlNullStringToStruct[TestUser](input)
		
		assert.NoError(t, err)
		assert.Equal(t, TestUser{}, result) // Zero value
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		input := sql.NullString{String: `{"name":"test",invalid}`, Valid: true}
		
		result, err := goxSql.SqlNullStringToStruct[TestUser](input)
		
		assert.Error(t, err)
		assert.Equal(t, TestUser{}, result)
	})

	t.Run("empty JSON string", func(t *testing.T) {
		input := sql.NullString{String: "", Valid: true}
		
		result, err := goxSql.SqlNullStringToStruct[TestUser](input)
		
		assert.Error(t, err) // Empty string should cause JSON parsing error
		assert.Equal(t, TestUser{}, result)
	})

	t.Run("valid empty JSON object", func(t *testing.T) {
		input := sql.NullString{String: "{}", Valid: true}
		
		result, err := goxSql.SqlNullStringToStruct[TestUser](input)
		
		assert.NoError(t, err)
		assert.Equal(t, TestUser{}, result) // Should be zero value but no error
	})
}

func TestStructToSqlNullString(t *testing.T) {
	t.Run("simple struct to JSON", func(t *testing.T) {
		user := TestUser{
			Name: "John Doe",
			Age:  30,
			Tags: []string{"admin", "user"},
		}
		
		result, err := goxSql.StructToSqlNullString[TestUser](user)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, `"name":"John Doe"`)
		assert.Contains(t, result.String, `"age":30`)
		assert.Contains(t, result.String, `"tags":["admin","user"]`)
	})

	t.Run("nested struct to JSON", func(t *testing.T) {
		nested := TestNestedStruct{
			User:   TestUser{Name: "Jane", Age: 25},
			Active: true,
			Score:  95.5,
		}
		
		result, err := goxSql.StructToSqlNullString[TestNestedStruct](nested)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, `"user":`)
		assert.Contains(t, result.String, `"name":"Jane"`)
		assert.Contains(t, result.String, `"active":true`)
		assert.Contains(t, result.String, `"score":95.5`)
	})

	t.Run("empty struct to JSON", func(t *testing.T) {
		user := TestUser{}
		
		result, err := goxSql.StructToSqlNullString[TestUser](user)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, `"name":""`)
		assert.Contains(t, result.String, `"age":0`)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var user *TestUser = nil
		
		result, err := goxSql.StructToSqlNullString[*TestUser](user)
		
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Equal(t, "null", result.String)
	})

	t.Run("primitive types", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected string
		}{
			{"string", "hello", "hello"},
			{"int", 42, "42"},
			{"bool", true, "true"},
			{"float", 3.14, "3.14"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := goxSql.StructToSqlNullString[interface{}](tt.input)
				
				assert.NoError(t, err)
				assert.True(t, result.Valid)
				assert.Equal(t, tt.expected, result.String)
			})
		}
	})
}

// Test roundtrip conversion: struct -> sql.NullString -> struct
func TestStructRoundtrip(t *testing.T) {
	t.Run("simple struct roundtrip", func(t *testing.T) {
		original := TestUser{
			Name: "Alice",
			Age:  28,
			Tags: []string{"developer", "admin"},
		}

		// Convert to sql.NullString
		nullStr, err := goxSql.StructToSqlNullString[TestUser](original)
		assert.NoError(t, err)

		// Convert back to struct
		recovered, err := goxSql.SqlNullStringToStruct[TestUser](nullStr)
		assert.NoError(t, err)

		// Verify they match
		assert.Equal(t, original, recovered)
	})

	t.Run("nested struct roundtrip", func(t *testing.T) {
		original := TestNestedStruct{
			User:   TestUser{Name: "Bob", Age: 35, Tags: []string{"manager"}},
			Active: true,
			Score:  88.7,
		}

		// Convert to sql.NullString
		nullStr, err := goxSql.StructToSqlNullString[TestNestedStruct](original)
		assert.NoError(t, err)

		// Convert back to struct
		recovered, err := goxSql.SqlNullStringToStruct[TestNestedStruct](nullStr)
		assert.NoError(t, err)

		// Verify they match
		assert.Equal(t, original, recovered)
	})
}