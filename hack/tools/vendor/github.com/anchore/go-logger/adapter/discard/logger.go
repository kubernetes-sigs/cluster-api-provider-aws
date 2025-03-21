package discard

import (
	"io"

	iface "github.com/anchore/go-logger"
)

var _ iface.Logger = (*logger)(nil)
var _ iface.Controller = (*logger)(nil)

type logger struct {
}

func New() iface.Logger {
	return &logger{}
}

func (l *logger) Tracef(_ string, _ ...interface{}) {
}

func (l *logger) Debugf(_ string, _ ...interface{}) {}

func (l *logger) Infof(_ string, _ ...interface{}) {}

func (l *logger) Warnf(_ string, _ ...interface{}) {}

func (l *logger) Errorf(_ string, _ ...interface{}) {}

func (l *logger) Trace(_ ...interface{}) {}

func (l *logger) Debug(_ ...interface{}) {}

func (l *logger) Info(_ ...interface{}) {}

func (l *logger) Warn(_ ...interface{}) {}

func (l *logger) Error(_ ...interface{}) {}

func (l *logger) WithFields(_ ...interface{}) iface.MessageLogger {
	return l
}

func (l *logger) Nested(_ ...interface{}) iface.Logger { return l }

func (l *logger) SetOutput(_ io.Writer) {}

func (l *logger) GetOutput() io.Writer { return nil }
