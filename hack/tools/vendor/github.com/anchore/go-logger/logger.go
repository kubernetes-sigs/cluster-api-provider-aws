package logger

import (
	"fmt"
	"io"
	"strings"
)

type Level string

const (
	DisabledLevel Level = ""
	ErrorLevel    Level = "error"
	WarnLevel     Level = "warn"
	InfoLevel     Level = "info"
	DebugLevel    Level = "debug"
	TraceLevel    Level = "trace"
)

func Levels() []Level {
	return []Level{
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel,
	}
}

type Logger interface {
	MessageLogger
	FieldLogger
	NestedLogger
}

type Controller interface {
	SetOutput(io.Writer)
	GetOutput() io.Writer
}

type NestedLogger interface {
	Nested(fields ...interface{}) Logger
}

type FieldLogger interface {
	WithFields(fields ...interface{}) MessageLogger
}

type Fields map[string]interface{}

type MessageLogger interface {
	ErrorMessageLogger
	WarnMessageLogger
	InfoMessageLogger
	DebugMessageLogger
	TraceMessageLogger
}

// type MessageLogger interface {
//	Logf(level Level, format string, args ...interface{})
//	Log(level Level, args ...interface{})
//}

type ErrorMessageLogger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
}

type WarnMessageLogger interface {
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
}

type InfoMessageLogger interface {
	Infof(format string, args ...interface{})
	Info(args ...interface{})
}

type DebugMessageLogger interface {
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

type TraceMessageLogger interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
}

func LevelFromString(l string) (Level, error) {
	switch strings.ToLower(l) {
	case "":
		return DisabledLevel, nil
	case "error", "err", "e":
		return ErrorLevel, nil
	case "warn", "warning", "w":
		return WarnLevel, nil
	case "info", "information", "informational", "i":
		return InfoLevel, nil
	case "debug", "debugging", "d":
		return DebugLevel, nil
	case "trace", "t":
		return TraceLevel, nil
	}

	return Level(l), fmt.Errorf("not a valid log level: %q", l)
}

func LevelFromVerbosity(v int, levels ...Level) Level {
	if len(levels) == 0 {
		return DisabledLevel
	}
	if v >= len(levels) {
		return levels[len(levels)-1]
	}
	if v <= 0 {
		return levels[0]
	}
	return levels[v]
}

func IsLevel(l Level, levels ...Level) bool {
	for _, level := range levels {
		if l == level {
			return true
		}
	}
	return false
}

func IsVerbose(level Level) bool {
	return IsLevel(level, InfoLevel, DebugLevel)
}
