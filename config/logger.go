package config

type Logger struct {
	LogLevel         string `json:"level" yaml:"level"`
	EnableConsoleLog bool   `json:"enable_console_log" yaml:"enable_console_log"`
}
