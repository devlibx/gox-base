package logger

import "io"

// Maintains all fields to put with logs
type Fields map[string]interface{}

type Level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Internal Entry for logging
type Entry interface {
	// Output a debug log
	Debug(args ...interface{})

	// Output a debug log
	Info(args ...interface{})

	// Output a debug log
	Warn(args ...interface{})

	// Output a debug log
	Error(args ...interface{})

	// Output a debug log
	Trace(args ...interface{})
}

// Config for logger
type Configuration struct {
	LogLevel Level
	Out      io.Writer
}

type Logger interface {

	// Initialize log configuration
	InitLogger(configuration Configuration)

	// Output a debug log
	Trace(args ...interface{})

	// Output a debug log
	Debug(args ...interface{})

	// Output a debug log
	Info(args ...interface{})

	// Output a debug log
	Warn(args ...interface{})

	// Output a debug log
	Error(args ...interface{})

	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool

	// Add a single field in the log
	WithField(key string, value interface{}) Entry

	// Add one or more field in logs
	WithFields(fields Fields) Entry

	// Add a single field in the log
	WithError(err error) Entry
}

// Create a new logger
func NewLogger(configuration Configuration) Logger {
	logger := &internalImpl{}
	logger.InitLogger(configuration)
	return logger
}

// Create a new logger
func NewNoOpLogger(configuration Configuration) Logger {
	return &noOpLogger{}
}
