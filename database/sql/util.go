// Package goxSql provides utility functions for working with SQL database operations,
// particularly for converting between Go types and sql.NullString types, with support
// for encryption/decryption and JSON serialization.
package goxSql

import (
	"database/sql"
	errors2 "github.com/devlibx/gox-base/v2/errors"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
)

// StringToSqlNullString converts a Go string to a sql.NullString.
// The resulting sql.NullString will always have Valid set to true, regardless
// of whether the input string is empty or not. This function handles empty
// strings by storing them as valid empty strings in the database.
//
// Parameters:
//   - s: The input string to convert
//
// Returns:
//   - sql.NullString with Valid=true and String=s
//
// Example:
//
//	nullStr := StringToSqlNullString("hello")     // {String: "hello", Valid: true}
//	nullStr := StringToSqlNullString("")          // {String: "", Valid: true}
func StringToSqlNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{
			String: "",
			Valid:  true,
		}
	} else {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}
}

// SqlNullStringToString converts a sql.NullString to a Go string.
// If the sql.NullString is valid, it returns the contained string value.
// If the sql.NullString is not valid (represents NULL in database), it returns an empty string.
//
// Parameters:
//   - s: The sql.NullString to convert
//
// Returns:
//   - string: The contained string value if valid, empty string if not valid
//
// Example:
//
//	str := SqlNullStringToString(sql.NullString{String: "hello", Valid: true})   // "hello"
//	str := SqlNullStringToString(sql.NullString{String: "", Valid: false})      // ""
func SqlNullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}

// StringToEncryptedSqlNullString encrypts a string and converts it to a sql.NullString.
// This generic function accepts any type T that implements the EncryptAndOutputBase64Ciphertext method.
// The input string is encrypted using the provided encryptor, and the resulting base64-encoded
// ciphertext is stored in a sql.NullString with Valid=true.
//
// Type Parameters:
//   - T: Any type that implements EncryptAndOutputBase64Ciphertext(data string) (string, error)
//
// Parameters:
//   - s: The input string to encrypt and convert
//   - encryptor: The encryption service implementing the required interface
//
// Returns:
//   - sql.NullString: Contains the encrypted base64 ciphertext with Valid=true
//   - error: Any encryption error, wrapped with context
//
// Example:
//
//	encryptor := myEncryptionService{}
//	nullStr, err := StringToEncryptedSqlNullString("sensitive data", encryptor)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// nullStr.String contains the base64-encoded encrypted data
func StringToEncryptedSqlNullString[
	T interface {
		EncryptAndOutputBase64Ciphertext(data string) (string, error)
	}](
	s string,
	encryptor T,
) (sql.NullString, error) {
	const errString = "unable to encrypt data while converting sql to sql.NullString"
	if len(s) == 0 {
		if e, err := encryptor.EncryptAndOutputBase64Ciphertext(""); err != nil {
			return sql.NullString{}, errors2.Wrap(err, errString)
		} else {
			return sql.NullString{
				String: e,
				Valid:  true,
			}, nil
		}
	} else {
		if e, err := encryptor.EncryptAndOutputBase64Ciphertext(s); err != nil {
			return sql.NullString{}, errors2.Wrap(err, errString)
		} else {
			return sql.NullString{
				String: e,
				Valid:  true,
			}, nil
		}
	}
}

// EncryptedSqlNullStringToString decrypts a sql.NullString containing encrypted data and returns the original string.
// This generic function accepts any type T that implements the DecryptFromBase64Ciphertext method.
// If the sql.NullString is valid, it decrypts the base64-encoded ciphertext to recover the original string.
// If the sql.NullString is not valid (represents NULL), it returns an empty string without error.
//
// Type Parameters:
//   - T: Any type that implements DecryptFromBase64Ciphertext(data string) (string, error)
//
// Parameters:
//   - s: The sql.NullString containing encrypted base64 ciphertext
//   - decrypter: The decryption service implementing the required interface
//
// Returns:
//   - string: The decrypted original string, or empty string if sql.NullString is not valid
//   - error: Any decryption error, wrapped with context
//
// Example:
//
//	decrypter := myDecryptionService{}
//	originalStr, err := EncryptedSqlNullStringToString(encryptedNullStr, decrypter)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// originalStr contains the decrypted data
func EncryptedSqlNullStringToString[
	T interface {
		DecryptFromBase64Ciphertext(data string) (string, error)
	}](
	s sql.NullString,
	decrypter T,
) (string, error) {
	if s.Valid {
		if e, err := decrypter.DecryptFromBase64Ciphertext(s.String); err != nil {
			return "", errors2.Wrap(err, "unable to decrypt data from sql.NullString to string")
		} else {
			return e, nil
		}
	} else {
		return "", nil
	}
}

// SqlNullStringToStruct deserializes a JSON string stored in sql.NullString to a Go struct.
// This generic function can convert to any type T. If the sql.NullString is valid,
// it attempts to deserialize the JSON string to the target type. If the sql.NullString
// is not valid (represents NULL), it returns the zero value of type T without error.
//
// Type Parameters:
//   - T: The target Go type to deserialize to (must be JSON deserializable)
//
// Parameters:
//   - value: The sql.NullString containing JSON data
//
// Returns:
//   - T: The deserialized struct of type T, or zero value if sql.NullString is not valid
//   - error: Any JSON deserialization error
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	var user User
//	user, err := SqlNullStringToStruct[User](jsonNullStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
func SqlNullStringToStruct[T any](value sql.NullString) (T, error) {
	var retValue T
	if value.Valid {
		return goxJsonUtils.StringToObject[T](value.String)
	} else {
		return retValue, nil
	}
}

// StructToSqlNullString serializes a Go struct to JSON and stores it in a sql.NullString.
// This generic function can accept any type as input. The struct is serialized to JSON
// and stored in a sql.NullString with Valid=true. The function will return an error
// if the JSON serialization fails.
//
// Type Parameters:
//   - T: The type constraint (any) - accepts any Go type for serialization
//
// Parameters:
//   - value: The Go struct/value to serialize to JSON
//
// Returns:
//   - sql.NullString: Contains the JSON string with Valid=true
//   - error: Any JSON serialization error, wrapped with context
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	user := User{Name: "John", Age: 30}
//	nullStr, err := StructToSqlNullString(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// nullStr.String contains: {"name":"John","age":30}
func StructToSqlNullString[T any](value any) (sql.NullString, error) {
	val, err := goxJsonUtils.ObjectToString(value)
	if err != nil {
		return sql.NullString{}, errors2.Wrap(err, "unable to serialize data to JSON")
	}
	return sql.NullString{Valid: true, String: val}, nil
}

// StructToEncryptedSqlNullString serializes a Go struct to JSON, encrypts it, and stores it in a sql.NullString.
// This generic function combines JSON serialization with encryption for secure storage of complex data structures.
// The struct is first serialized to JSON, then the JSON string is encrypted using the provided encryptor,
// and finally stored in a sql.NullString with Valid=true.
//
// Type Parameters:
//   - T: Any type that implements EncryptAndOutputBase64Ciphertext(data string) (string, error)
//
// Parameters:
//   - value: The Go struct/value to serialize and encrypt
//   - encryptor: The encryption service implementing the required interface
//
// Returns:
//   - sql.NullString: Contains the encrypted JSON string with Valid=true
//   - error: Any JSON serialization or encryption error, wrapped with context
//
// Example:
//
//	type UserProfile struct {
//	    Name   string `json:"name"`
//	    Email  string `json:"email"`
//	    Secret string `json:"secret"`
//	}
//
//	profile := UserProfile{
//	    Name:   "John Doe",
//	    Email:  "john@example.com",
//	    Secret: "sensitive-data",
//	}
//
//	encryptor := &MyEncryptor{}
//	nullStr, err := StructToEncryptedSqlNullString(profile, encryptor)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// nullStr.String contains encrypted JSON: base64-encoded ciphertext
//	// Store nullStr in database for secure persistence
func StructToEncryptedSqlNullString[
	T interface {
		EncryptAndOutputBase64Ciphertext(data string) (string, error)
	}](
	value any,
	encryptor T,
) (sql.NullString, error) {
	// First, serialize the struct to JSON
	jsonStr, err := goxJsonUtils.ObjectToString(value)
	if err != nil {
		return sql.NullString{}, errors2.Wrap(err, "unable to serialize struct to JSON before encryption")
	}

	// Then encrypt the JSON string
	encryptedData, err := encryptor.EncryptAndOutputBase64Ciphertext(jsonStr)
	if err != nil {
		return sql.NullString{}, errors2.Wrap(err, "unable to encrypt JSON data while converting struct to sql.NullString")
	}

	return sql.NullString{Valid: true, String: encryptedData}, nil
}

// EncryptedSqlNullStringToStruct decrypts a sql.NullString containing encrypted JSON data and deserializes it to a Go struct.
// This generic function combines decryption with JSON deserialization for secure retrieval of complex data structures.
// If the sql.NullString is valid, it decrypts the encrypted JSON string, then deserializes it to the target struct type.
// If the sql.NullString is not valid (represents NULL), it returns the zero value of type T without error.
//
// Type Parameters:
//   - T: The target Go type to deserialize to (must be JSON deserializable)
//   - U: Any type that implements DecryptFromBase64Ciphertext(data string) (string, error)
//
// Parameters:
//   - value: The sql.NullString containing encrypted JSON data
//   - decryptor: The decryption service implementing the required interface
//
// Returns:
//   - T: The deserialized struct of type T, or zero value if sql.NullString is not valid
//   - error: Any decryption or JSON deserialization error, wrapped with context
//
// Example:
//
//	type UserProfile struct {
//	    Name   string `json:"name"`
//	    Email  string `json:"email"`
//	    Secret string `json:"secret"`
//	}
//
//	// Retrieve encrypted data from database
//	var encryptedNullStr sql.NullString
//	err := db.QueryRow("SELECT encrypted_profile FROM users WHERE id = ?", userID).Scan(&encryptedNullStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	decryptor := &MyDecryptor{}
//	profile, err := EncryptedSqlNullStringToStruct[UserProfile](encryptedNullStr, decryptor)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// profile now contains the decrypted and deserialized UserProfile struct
func EncryptedSqlNullStringToStruct[
	T any,
	U interface {
		DecryptFromBase64Ciphertext(data string) (string, error)
	}](
	value sql.NullString,
	decryptor U,
) (T, error) {
	var retValue T

	if !value.Valid {
		// Return zero value if sql.NullString represents NULL
		return retValue, nil
	}

	// First, decrypt the encrypted JSON string
	jsonStr, err := decryptor.DecryptFromBase64Ciphertext(value.String)
	if err != nil {
		return retValue, errors2.Wrap(err, "unable to decrypt data from sql.NullString before JSON deserialization")
	}

	// Then deserialize the JSON to the target struct
	result, err := goxJsonUtils.StringToObject[T](jsonStr)
	if err != nil {
		return retValue, errors2.Wrap(err, "unable to deserialize decrypted JSON to struct")
	}

	return result, nil
}
