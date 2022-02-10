package logger

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	// INFO human-readable log level.
	INFO = "INFO"
	// WARNING human-readable log level.
	WARNING = "WARNING"
	// ERROR human-readable log level.
	ERROR = "ERROR"
	// FATAL human-readable log level.
	FATAL = "FATAL"
	// DEBUG human-readable log level.
	DEBUG = "DEBUG"
	// TRACE human-readable log level.
	TRACE = "TRACE"

	// LogLevelInfo for V(k) based numbering.
	LogLevelInfo = iota
	// LogLevelWarning for V(k) based numbering.
	LogLevelWarning
	// LogLevelError for V(k) based numbering.
	LogLevelError
	// LogLevelFatal for V(k) based numbering.
	LogLevelFatal
	// LogLevelDebug for V(k) based numbering.
	LogLevelDebug
	// LogLevelTrace for V(k) based numbering.
	LogLevelTrace
)

type Logger struct {
	logger logr.Logger
}

func (l *Logger) Enabled() bool {
	return l.logger.Enabled()
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelInfo)
	l.logger.WithValues("level", INFO).Info(msg, keysAndValues...)
}

func (l *Logger) Warning(msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelWarning)
	l.logger.WithValues("level", WARNING).Info(msg, keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelError)
	l.logger.WithValues("level", ERROR).Error(err, msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelFatal)
	l.logger.WithValues("level", FATAL).Info(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelDebug)
	l.logger.WithValues("level", DEBUG).Info(msg, keysAndValues...)
}

func (l *Logger) Trace(err error, msg string, keysAndValues ...interface{}) {
	l.logger = l.logger.V(LogLevelTrace)
	l.logger.WithValues("level", TRACE).Error(err, msg, keysAndValues...)
}

func (l *Logger) WithValues(keysAndValues ...interface{}) *Logger {
	l.logger = l.logger.WithValues(keysAndValues...)
	return l
}

func (l *Logger) WithName(name string) *Logger {
	l.logger = l.logger.WithName(name)
	return l
}

func (l *Logger) Logger() logr.Logger {
	return l.logger
}

// FromContext gets a logr.Logger from context and wraps it into a verbose logger
// with grep-able log levels displayed as `level=INFO`.
func FromContext(ctx context.Context) *Logger {
	log := ctrl.LoggerFrom(ctx)
	return &Logger{
		logger: log,
	}
}

// NewWithLogger creates a new logger using a passed in logr.Logger.
func NewWithLogger(log logr.Logger) *Logger {
	return &Logger{
		logger: log,
	}
}
