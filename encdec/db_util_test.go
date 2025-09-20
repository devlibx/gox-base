package encryption

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/devlibx/gox-base/v2"
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
	Name string   `json:"name"`
	Age  int      `json:"age"`
	Tags []string `json:"tags,omitempty"`
}

type TestNestedStruct struct {
	User   TestUser `json:"user"`
	Active bool     `json:"active"`
	Score  float64  `json:"score"`
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

func TestStructToEncryptedSqlNullString(t *testing.T) {
	t.Run("successful encryption of simple struct", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		user := TestUser{
			Name: "John Doe",
			Age:  30,
			Tags: []string{"admin", "developer"},
		}

		result, err := goxSql.StructToEncryptedSqlNullString(user, encryptor)

		assert.NoError(t, err)
		assert.True(t, result.Valid)
		// Should contain encrypted JSON (prefixed with "encrypted_")
		assert.Contains(t, result.String, "encrypted_")

		// Verify it contains the JSON structure after decryption (mock removes "encrypted_" prefix)
		decryptedData := result.String[10:] // Remove "encrypted_" prefix
		assert.Contains(t, decryptedData, `"name":"John Doe"`)
		assert.Contains(t, decryptedData, `"age":30`)
	})

	t.Run("successful encryption of nested struct", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		nested := TestNestedStruct{
			User:   TestUser{Name: "Jane", Age: 25, Tags: []string{"user"}},
			Active: true,
			Score:  95.7,
		}

		result, err := goxSql.StructToEncryptedSqlNullString(nested, encryptor)

		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, "encrypted_")

		// Verify nested structure in decrypted data
		decryptedData := result.String[10:]
		assert.Contains(t, decryptedData, `"user":`)
		assert.Contains(t, decryptedData, `"name":"Jane"`)
		assert.Contains(t, decryptedData, `"active":true`)
	})

	t.Run("successful encryption of empty struct", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		user := TestUser{}

		result, err := goxSql.StructToEncryptedSqlNullString(user, encryptor)

		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, "encrypted_")
	})

	t.Run("successful encryption of nil pointer", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		var user *TestUser = nil

		result, err := goxSql.StructToEncryptedSqlNullString(user, encryptor)

		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Contains(t, result.String, "encrypted_")
		// Should contain encrypted "null"
		decryptedData := result.String[10:]
		assert.Equal(t, "null", decryptedData)
	})

	t.Run("encryption error during struct encryption", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: true, errorMsg: "encryption failed"}
		user := TestUser{Name: "John", Age: 30}

		result, err := goxSql.StructToEncryptedSqlNullString(user, encryptor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to encrypt JSON data while converting struct to sql.NullString")
		assert.Contains(t, err.Error(), "encryption failed")
		assert.Equal(t, sql.NullString{}, result)
	})

	t.Run("encryption with primitive types", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
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
				result, err := goxSql.StructToEncryptedSqlNullString(tt.input, encryptor)

				assert.NoError(t, err)
				assert.True(t, result.Valid)
				assert.Contains(t, result.String, "encrypted_")

				decryptedData := result.String[10:]
				assert.Equal(t, tt.expected, decryptedData)
			})
		}
	})
}

func TestEncryptedSqlNullStringToStruct(t *testing.T) {
	t.Run("successful decryption to simple struct", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		// Create encrypted JSON data (mock adds "encrypted_" prefix)
		jsonStr := `{"name":"John Doe","age":30,"tags":["admin","developer"]}`
		encryptedData := "encrypted_" + jsonStr
		input := sql.NullString{String: encryptedData, Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.NoError(t, err)
		assert.Equal(t, "John Doe", result.Name)
		assert.Equal(t, 30, result.Age)
		assert.Equal(t, []string{"admin", "developer"}, result.Tags)
	})

	t.Run("successful decryption to nested struct", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		jsonStr := `{"user":{"name":"Jane","age":25,"tags":["user"]},"active":true,"score":95.7}`
		encryptedData := "encrypted_" + jsonStr
		input := sql.NullString{String: encryptedData, Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestNestedStruct](input, decryptor)

		assert.NoError(t, err)
		assert.Equal(t, "Jane", result.User.Name)
		assert.Equal(t, 25, result.User.Age)
		assert.Equal(t, []string{"user"}, result.User.Tags)
		assert.True(t, result.Active)
		assert.Equal(t, 95.7, result.Score)
	})

	t.Run("invalid sql.NullString returns zero value", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_data", Valid: false}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.NoError(t, err)
		assert.Equal(t, TestUser{}, result) // Zero value
	})

	t.Run("decryption error", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: true, errorMsg: "decryption failed"}
		input := sql.NullString{String: "encrypted_data", Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to decrypt data from sql.NullString before JSON deserialization")
		assert.Contains(t, err.Error(), "decryption failed")
		assert.Equal(t, TestUser{}, result)
	})

	t.Run("invalid JSON after decryption", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		// Mock will return this invalid JSON after removing "encrypted_" prefix
		input := sql.NullString{String: "encrypted_{invalid json}", Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to deserialize decrypted JSON to struct")
		assert.Equal(t, TestUser{}, result)
	})

	t.Run("empty encrypted data", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_", Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.Error(t, err) // Empty string should cause JSON parsing error
		assert.Equal(t, TestUser{}, result)
	})

	t.Run("valid empty JSON object after decryption", func(t *testing.T) {
		decryptor := &mockDecryptor{shouldError: false}
		input := sql.NullString{String: "encrypted_{}", Valid: true}

		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](input, decryptor)

		assert.NoError(t, err)
		assert.Equal(t, TestUser{}, result) // Should be zero value but no error
	})
}

// Test roundtrip conversion with encryption: struct -> encrypted sql.NullString -> struct
func TestEncryptedStructRoundtrip(t *testing.T) {
	encryptor := &mockEncryptor{shouldError: false}
	decryptor := &mockDecryptor{shouldError: false}

	t.Run("simple struct encrypted roundtrip", func(t *testing.T) {
		original := TestUser{
			Name: "Alice",
			Age:  28,
			Tags: []string{"developer", "admin"},
		}

		// Encrypt and convert to sql.NullString
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(original, encryptor)
		assert.NoError(t, err)
		assert.True(t, encryptedNullStr.Valid)
		assert.Contains(t, encryptedNullStr.String, "encrypted_")

		// Decrypt and convert back to struct
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedNullStr, decryptor)
		assert.NoError(t, err)

		// Verify they match
		assert.Equal(t, original, recovered)
	})

	t.Run("nested struct encrypted roundtrip", func(t *testing.T) {
		original := TestNestedStruct{
			User:   TestUser{Name: "Bob", Age: 35, Tags: []string{"manager"}},
			Active: true,
			Score:  88.7,
		}

		// Encrypt and convert to sql.NullString
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(original, encryptor)
		assert.NoError(t, err)

		// Decrypt and convert back to struct
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestNestedStruct](encryptedNullStr, decryptor)
		assert.NoError(t, err)

		// Verify they match
		assert.Equal(t, original, recovered)
	})

	// Note: Primitive types don't roundtrip cleanly through JSON serialization
	// when used directly with encrypted struct methods. These methods are designed
	// for structured data (structs), not raw primitive values.
}

// Test error scenarios with mixed encryption/decryption failures
func TestEncryptedStructErrorScenarios(t *testing.T) {
	t.Run("encryption succeeds but decryption fails", func(t *testing.T) {
		encryptor := &mockEncryptor{shouldError: false}
		failingDecryptor := &mockDecryptor{shouldError: true, errorMsg: "decryption service unavailable"}

		user := TestUser{Name: "John", Age: 30}

		// Encryption should succeed
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(user, encryptor)
		assert.NoError(t, err)
		assert.True(t, encryptedNullStr.Valid)

		// Decryption should fail
		result, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedNullStr, failingDecryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to decrypt data from sql.NullString before JSON deserialization")
		assert.Contains(t, err.Error(), "decryption service unavailable")
		assert.Equal(t, TestUser{}, result)
	})
}

// Tests using real EncryptorDecryptService instead of mocks
func TestEncryptedStructMethodsWithRealEncryption(t *testing.T) {
	// Setup real encryption service
	key, err := GenerateAESKey(32)
	assert.NoError(t, err)

	encConfig := &EncryptDecryptConfigs{
		Group: map[string]*EncryptDecryptConfig{
			"test": {
				Algo: "aes_32",
				AesConfig: &AesConfig{
					Base64CodedKey: base64.StdEncoding.EncodeToString(key),
				},
			},
		},
	}

	factory, err := NewServiceFactory(gox.NewNoOpCrossFunction(), encConfig)
	assert.NoError(t, err)

	encryptorDecryptor, err := factory.GetEncryptorDecryptService("test")
	assert.NoError(t, err)

	t.Run("struct to encrypted sql.NullString with real encryption", func(t *testing.T) {
		user := TestUser{
			Name: "John Doe",
			Age:  30,
			Tags: []string{"admin", "developer"},
		}

		// Test encryption
		result, err := goxSql.StructToEncryptedSqlNullString(user, encryptorDecryptor)
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.NotEmpty(t, result.String)

		// Verify it's actually encrypted (not readable JSON)
		assert.NotContains(t, result.String, "John Doe")
		assert.NotContains(t, result.String, "admin")
	})

	t.Run("encrypted sql.NullString to struct with real decryption", func(t *testing.T) {
		// First encrypt a struct
		original := TestUser{
			Name: "Jane Smith",
			Age:  25,
			Tags: []string{"user", "tester"},
		}

		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(original, encryptorDecryptor)
		assert.NoError(t, err)

		// Then decrypt it back
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)

		// Verify they match
		assert.Equal(t, original, recovered)
	})

	t.Run("nested struct roundtrip with real encryption", func(t *testing.T) {
		original := TestNestedStruct{
			User:   TestUser{Name: "Alice", Age: 28, Tags: []string{"manager", "lead"}},
			Active: true,
			Score:  95.5,
		}

		// Encrypt
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(original, encryptorDecryptor)
		assert.NoError(t, err)
		assert.True(t, encryptedNullStr.Valid)

		// Verify it's encrypted
		assert.NotContains(t, encryptedNullStr.String, "Alice")
		assert.NotContains(t, encryptedNullStr.String, "manager")

		// Decrypt
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestNestedStruct](encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)

		// Verify perfect roundtrip
		assert.Equal(t, original, recovered)
	})

	t.Run("encrypted string methods with real encryption service", func(t *testing.T) {
		originalData := "sensitive personal information"

		// Test string encryption
		encryptedNullStr, err := goxSql.StringToEncryptedSqlNullString(originalData, encryptorDecryptor)
		assert.NoError(t, err)
		assert.True(t, encryptedNullStr.Valid)
		assert.NotEmpty(t, encryptedNullStr.String)

		// Verify it's actually encrypted
		assert.NotContains(t, encryptedNullStr.String, "sensitive")
		assert.NotContains(t, encryptedNullStr.String, "personal")

		// Test string decryption
		decryptedData, err := goxSql.EncryptedSqlNullStringToString(encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)
		assert.Equal(t, originalData, decryptedData)
	})

	t.Run("empty struct with real encryption", func(t *testing.T) {
		empty := TestUser{}

		// Encrypt empty struct
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(empty, encryptorDecryptor)
		assert.NoError(t, err)

		// Decrypt back
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)

		// Should match
		assert.Equal(t, empty, recovered)
	})

	t.Run("nil pointer with real encryption", func(t *testing.T) {
		var nilPtr *TestUser = nil

		// Encrypt nil pointer
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(nilPtr, encryptorDecryptor)
		assert.NoError(t, err)

		// Decrypt back as pointer
		var recovered *TestUser
		recovered, err = goxSql.EncryptedSqlNullStringToStruct[*TestUser](encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)

		// Should be nil
		assert.Nil(t, recovered)
	})

	t.Run("unicode data with real encryption", func(t *testing.T) {
		unicodeUser := TestUser{
			Name: "测试用户",
			Age:  30,
			Tags: []string{"用户", "测试员"},
		}

		// Encrypt unicode data
		encryptedNullStr, err := goxSql.StructToEncryptedSqlNullString(unicodeUser, encryptorDecryptor)
		assert.NoError(t, err)

		// Decrypt back
		recovered, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedNullStr, encryptorDecryptor)
		assert.NoError(t, err)

		// Should preserve unicode
		assert.Equal(t, unicodeUser, recovered)
	})
}

// Test compatibility between different encryption methods
func TestEncryptionCompatibility(t *testing.T) {
	// Setup real encryption service
	key, err := GenerateAESKey(32)
	assert.NoError(t, err)

	encConfig := &EncryptDecryptConfigs{
		Group: map[string]*EncryptDecryptConfig{
			"test": {
				Algo: "aes_32",
				AesConfig: &AesConfig{
					Base64CodedKey: base64.StdEncoding.EncodeToString(key),
				},
			},
		},
	}

	factory, err := NewServiceFactory(gox.NewNoOpCrossFunction(), encConfig)
	assert.NoError(t, err)

	service, err := factory.GetEncryptorDecryptService("test")
	assert.NoError(t, err)

	t.Run("same service works for both string and struct methods", func(t *testing.T) {
		// Test string methods
		testString := "test data"
		encryptedStr, err := goxSql.StringToEncryptedSqlNullString(testString, service)
		assert.NoError(t, err)

		decryptedStr, err := goxSql.EncryptedSqlNullStringToString(encryptedStr, service)
		assert.NoError(t, err)
		assert.Equal(t, testString, decryptedStr)

		// Test struct methods
		testUser := TestUser{Name: "Test User", Age: 25}
		encryptedStruct, err := goxSql.StructToEncryptedSqlNullString(testUser, service)
		assert.NoError(t, err)

		decryptedStruct, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encryptedStruct, service)
		assert.NoError(t, err)
		assert.Equal(t, testUser, decryptedStruct)
	})

	t.Run("verify different data produces different ciphertext", func(t *testing.T) {
		user1 := TestUser{Name: "User One", Age: 25}
		user2 := TestUser{Name: "User Two", Age: 30}

		encrypted1, err := goxSql.StructToEncryptedSqlNullString(user1, service)
		assert.NoError(t, err)

		encrypted2, err := goxSql.StructToEncryptedSqlNullString(user2, service)
		assert.NoError(t, err)

		// Different data should produce different ciphertext
		assert.NotEqual(t, encrypted1.String, encrypted2.String)

		// But both should decrypt correctly
		decrypted1, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encrypted1, service)
		assert.NoError(t, err)
		assert.Equal(t, user1, decrypted1)

		decrypted2, err := goxSql.EncryptedSqlNullStringToStruct[TestUser](encrypted2, service)
		assert.NoError(t, err)
		assert.Equal(t, user2, decrypted2)
	})
}
