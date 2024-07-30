package log

import (
	"fmt"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
)

// assert interface compliance.
var _ Interface = (*Entry)(nil)

// Now returns the current time.
var Now = time.Now

// Entry represents a single log entry.
type Entry struct {
	Logger  *Logger
	Level   Level
	Message string
	Padding int
	Fields  *orderedmap.OrderedMap[string, any]
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *Logger) *Entry {
	return &Entry{
		Logger:  log,
		Padding: log.Padding,
		Fields:  orderedmap.NewOrderedMap[string, any](),
	}
}

// ResetPadding resets the padding to default.
func (e *Entry) ResetPadding() {
	e.Logger.ResetPadding()
}

// IncreasePadding increases the padding 1 times.
func (e *Entry) IncreasePadding() {
	e.Logger.IncreasePadding()
}

// DecreasePadding decreases the padding 1 times.
func (e *Entry) DecreasePadding() {
	e.Logger.DecreasePadding()
}

// WithField returns a new entry with the `key` and `value` set.
func (e *Entry) WithField(key string, value interface{}) *Entry {
	f := e.Fields.Copy()
	f.Set(key, value)
	return &Entry{
		Logger:  e.Logger,
		Padding: e.Padding,
		Fields:  f,
	}
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}
	ctx := e.WithField("error", err.Error())
	return ctx
}

// WithoutPadding returns a new entry with padding set to default.
func (e *Entry) WithoutPadding() *Entry {
	return &Entry{
		Logger:  e.Logger,
		Padding: defaultPadding,
		Fields:  e.Fields,
	}
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Logger.log(DebugLevel, e, msg)
}

// Info level message.
func (e *Entry) Info(msg string) {
	e.Logger.log(InfoLevel, e, msg)
}

// Warn level message.
func (e *Entry) Warn(msg string) {
	e.Logger.log(WarnLevel, e, msg)
}

// Error level message.
func (e *Entry) Error(msg string) {
	e.Logger.log(ErrorLevel, e, msg)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string) {
	e.Logger.log(FatalLevel, e, msg)
	os.Exit(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Info(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Warn(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Error(fmt.Sprintf(msg, v...))
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Fatal(fmt.Sprintf(msg, v...))
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string) *Entry {
	return &Entry{
		Logger:  e.Logger,
		Padding: e.Padding,
		Fields:  e.Fields,
		Level:   level,
		Message: msg,
	}
}
