package logger

type noOpLogger struct {
}

func (n noOpLogger) InitLogger(configuration Configuration) {
}

func (n noOpLogger) Trace(args ...interface{}) {
}

func (n noOpLogger) Debug(args ...interface{}) {
}

func (n noOpLogger) Info(args ...interface{}) {
}

func (n noOpLogger) Warn(args ...interface{}) {
}

func (n noOpLogger) Error(args ...interface{}) {
}

func (n noOpLogger) WithField(key string, value interface{}) Entry {
	return &noOpEntry{}
}

func (n noOpLogger) WithError(err error) Entry {
	return &noOpEntry{}
}

func (n noOpLogger) WithFields(fields Fields) Entry {
	return &noOpEntry{}
}

func (n noOpLogger) IsTraceEnabled() bool {
	return false
}

func (n noOpLogger) IsDebugEnabled() bool {
	return false
}

func (n noOpLogger) IsInfoEnabled() bool {
	return false
}

func NoOpLogger(configuration Configuration) Logger {
	return noOpLogger{}
}

type noOpEntry struct {
}

func (e *noOpEntry) Info(args ...interface{}) {
}

func (e *noOpEntry) Warn(args ...interface{}) {
}

func (e *noOpEntry) Error(args ...interface{}) {
}

func (e *noOpEntry) Debug(args ...interface{}) {

}

func (e *noOpEntry) Trace(args ...interface{}) {

}
