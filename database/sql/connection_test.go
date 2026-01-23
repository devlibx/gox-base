package goxSql

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
)

func TestCreateMySqlConnectionStringWithTimezone(t *testing.T) {
	user := "testUser"
	password := "testPass"
	host := "localhost"
	port := "3306"
	dbName := "testDb"
	timezone := "Asia/Kolkata"

	expected := "testUser:testPass@tcp(localhost:3306)/testDb?parseTime=true&loc=" + url.QueryEscape(timezone)
	actual := CreateMySqlConnectionStringWithTimezone(user, password, host, port, dbName, timezone)
	assert.Equal(t, expected, actual)

	// Test with empty timezone
	expectedEmpty := "testUser:testPass@tcp(localhost:3306)/testDb?parseTime=true"
	actualEmpty := CreateMySqlConnectionStringWithTimezone(user, password, host, port, dbName, "")
	assert.Equal(t, expectedEmpty, actualEmpty)
}

func TestCreateMySqlConnectionStringWithTimezoneEnv(t *testing.T) {
	user := "testUser"
	password := "testPass"
	host := "localhost"
	port := "3306"
	dbName := "testDb"
	timezone := "America/New_York"

	// Save current env and defer restore
	oldTimezone := os.Getenv("TIMEZONE")
	defer func() {
		os.Setenv("TIMEZONE", oldTimezone)
	}()

	// Test with env var set
	os.Setenv("TIMEZONE", timezone)
	expected := "testUser:testPass@tcp(localhost:3306)/testDb?parseTime=true&loc=" + url.QueryEscape(timezone)
	actual := CreateMySqlConnectionStringWithTimezoneEnv(user, password, host, port, dbName)
	assert.Equal(t, expected, actual)

	// Test with env var unset (or empty)
	os.Setenv("TIMEZONE", "")
	expectedEmpty := "testUser:testPass@tcp(localhost:3306)/testDb?parseTime=true"
	actualEmpty := CreateMySqlConnectionStringWithTimezoneEnv(user, password, host, port, dbName)
	assert.Equal(t, expectedEmpty, actualEmpty)
}
