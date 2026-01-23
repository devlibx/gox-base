package goxSql

import (
	"fmt"
	"net/url"
	"os"
)

// CreateMySqlConnectionStringWithTimezoneEnv creates a MySQL connection string by reading the timezone from the "TIMEZONE" environment variable.
func CreateMySqlConnectionStringWithTimezoneEnv(user, password, host, port, dbName string) string {
	timezone := os.Getenv("TIMEZONE")
	return createMySqlConnectionString(user, password, host, port, dbName, timezone)
}

// CreateMySqlConnectionStringWithTimezone creates a MySQL connection string with a specified timezone.
func CreateMySqlConnectionStringWithTimezone(user, password, host, port, dbName, timezone string) string {
	return createMySqlConnectionString(user, password, host, port, dbName, timezone)
}

// createMySqlConnectionString is a helper function to construct the MySQL connection string.
func createMySqlConnectionString(user, password, host, port, dbName, timezone string) string {
	connectionUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)
	if timezone != "" {
		connectionUrl = fmt.Sprintf("%s&loc=%s", connectionUrl, url.QueryEscape(timezone))
	}
	return connectionUrl
}
