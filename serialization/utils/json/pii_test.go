package goxJsonUtils_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPrettyStringLoggingSuppressError(t *testing.T) {

	// Ensure the registry is clear after this test suite runs

	t.Cleanup(goxJsonUtils.ClearPiiMaskersForTesting)



	t.Run("simple object", func(t *testing.T) {

		goxJsonUtils.RegisterPiiMasker("email", func(email string) string {

			return "masked-email"

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"name":  "John Doe",

			"email": "john.doe@example.com",

			"age":   30,

		}

		expected := `{

	"age": 30,

	"email": "masked-email",

	"name": "John Doe"

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("nested object", func(t *testing.T) {

		goxJsonUtils.RegisterPiiMasker("email", func(email string) string {

			return "masked-email"

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"user": map[string]interface{}{

				"name":  "Jane Doe",

				"email": "jane.doe@example.com",

			},

			"active": true,

		}

		expected := `{

	"active": true,

	"user": {

		"email": "masked-email",

		"name": "Jane Doe"

	}

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("array of objects", func(t *testing.T) {

		goxJsonUtils.RegisterPiiMasker("email", func(email string) string {

			return "masked-email"

		})

		goxJsonUtils.RegisterPiiMasker("ssn", func(ssn string) string {

			hasher := sha256.New()

			hasher.Write([]byte(ssn))

			return hex.EncodeToString(hasher.Sum(nil))

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"users": []interface{}{

				map[string]interface{}{

					"name":  "User 1",

					"email": "user1@example.com",

				},

				map[string]interface{}{

					"name":  "User 2",

					"ssn":   "123-45-678",

				},

			},

		}



		h := sha256.Sum256([]byte("123-45-678"))

		hashedSsn := hex.EncodeToString(h[:])



		expected := fmt.Sprintf(`{

	"users": [

		{

			"email": "masked-email",

			"name": "User 1"

		},

		{

			"name": "User 2",

			"ssn": "%s"

		}

	]

}`, hashedSsn)

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("panic handling", func(t *testing.T) {

		goxJsonUtils.RegisterPiiMasker("panicking_field", func(s string) string {

			panic("test panic")

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"field":           "some value",

			"panicking_field": "this will panic",

		}

		expected := `{

	"field": "some value",

	"panicking_field": "***************"

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("no masking needed", func(t *testing.T) {

		goxJsonUtils.ClearPiiMaskersForTesting() // Ensure no maskers are present

		input := map[string]interface{}{

			"name": "No PII here",

			"info": "nothing to see",

		}

		expected := `{

	"info": "nothing to see",

	"name": "No PII here"

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("non-string value for masked field", func(t *testing.T) {

		goxJsonUtils.RegisterPiiMasker("user_id", func(id string) string {

			return "masked-id"

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"user_id": 12345,

			"data":    "some data",

		}

		expected := `{

	"data": "some data",

	"user_id": "masked-id"

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("masking with a custom helper", func(t *testing.T) {

		hashEmail := func(email string) string {

			parts := strings.Split(email, "@")

			if len(parts) != 2 {

				return "invalid-email"

			}

			return "hashed-" + parts[0] + "@" + parts[1]

		}

		goxJsonUtils.RegisterPiiMasker("custom_email", hashEmail)

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := map[string]interface{}{

			"user":         "test_user",

			"custom_email": "test@example.com",

		}

		expected := `{

	"custom_email": "hashed-test@example.com",

	"user": "test_user"

}`

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)

	})



	t.Run("masking with struct object", func(t *testing.T) {

		type requestObj struct {

			Type   string `json:"type"`

			UserId string `json:"user_id"`

			MailId string `json:"mail_id"`

		}



		goxJsonUtils.RegisterPiiMasker("mail_id", func(mailId string) string {

			return "masked-mail-id"

		})

		defer goxJsonUtils.ClearPiiMaskersForTesting()



		input := requestObj{

			Type:   "test_request",

			UserId: "user-123",

			MailId: "test.user@example.com",

		}



		expected := `{

	"type": "test_request",

	"user_id": "user-123",

	"mail_id": "masked-mail-id"

}`



		// Test with struct value

		result := goxJsonUtils.PrettyStringLoggingSuppressError(input)

		assert.JSONEq(t, expected, result)



		// Test with pointer to struct

		result = goxJsonUtils.PrettyStringLoggingSuppressError(&input)

		assert.JSONEq(t, expected, result)

	})

}
