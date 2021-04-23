package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type entry struct {
	fields Fields
}

func (e *entry) buildFields() log.Fields {
	f := log.Fields{}
	for k, v := range e.fields {
		f[k] = v
	}
	return f
}

// Output debug logs
func (e *entry) Trace(args ...interface{}) {
	f := e.buildFields()
	if len(f) == 0 {
		log.Trace(args...)
	} else {
		log.WithFields(f).Trace(args...)
	}
}

// Output debug logs
func (e *entry) Debug(args ...interface{}) {
	f := e.buildFields()
	if len(f) == 0 {
		log.Debug(args...)
	} else {
		log.WithFields(f).Debug(args...)
	}
}

// Output debug logs
func (e *entry) Info(args ...interface{}) {
	f := e.buildFields()
	if len(f) == 0 {
		log.Info(args...)
	} else {
		log.WithFields(f).Info(args...)
	}
}

// Output debug logs
func (e *entry) Warn(args ...interface{}) {
	f := e.buildFields()
	if len(f) == 0 {
		log.Warn(args...)
	} else {
		log.WithFields(f).Warn(args...)
	}
}

// Output debug logs
func (e *entry) Error(args ...interface{}) {
	f := e.buildFields()
	if len(f) == 0 {
		log.Error(args...)
	} else {
		log.WithFields(f).Error(args...)
	}
}

type internalImpl struct {
	configuration Configuration
}

func (i *internalImpl) InitLogger(configuration Configuration) {
	i.configuration = configuration
	switch configuration.LogLevel {
	case DebugLevel:
		log.SetLevel(log.DebugLevel)
	case InfoLevel:
		log.SetLevel(log.InfoLevel)
	case WarnLevel:
		log.SetLevel(log.WarnLevel)
	case ErrorLevel:
		log.SetLevel(log.ErrorLevel)
	case FatalLevel:
		log.SetLevel(log.FatalLevel)
	case TraceLevel:
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	if configuration.Out != nil {
		log.SetOutput(configuration.Out)
	}
}

func (i *internalImpl) Trace(args ...interface{}) {
	log.Trace(args...)
}

func (i *internalImpl) Debug(args ...interface{}) {
	log.Debug(args...)
}

func (i *internalImpl) Info(args ...interface{}) {
	log.Info(args...)
}

func (i *internalImpl) Warn(args ...interface{}) {
	log.Warn(args...)
}

func (i *internalImpl) Error(args ...interface{}) {
	log.Error(args...)
}

func (i *internalImpl) WithField(key string, value interface{}) Entry {
	f := Fields{}
	f[key] = value
	e := &entry{fields: f}
	return e
}

func (i *internalImpl) WithError(err error) Entry {
	f := Fields{}
	f["error"] = fmt.Sprintf("%v", err)
	e := &entry{fields: f}
	return e
}

func (i *internalImpl) WithFields(fields Fields) Entry {
	e := &entry{fields: fields}
	return e
}

func (i *internalImpl) IsTraceEnabled() bool {
	return log.IsLevelEnabled(log.TraceLevel)
}

func (i *internalImpl) IsDebugEnabled() bool {
	return log.IsLevelEnabled(log.DebugLevel)
}

func (i *internalImpl) IsInfoEnabled() bool {
	return log.IsLevelEnabled(log.InfoLevel)
}
