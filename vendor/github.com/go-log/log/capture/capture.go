// Package capture implements the Logger interface by capturing logged
// lines.  This is useful for log inspection during unit-testing,
// if you want to assert that a particular line has, or has not, been
// logged.
package capture

import (
	"fmt"
	"sync"
)

// Logger implements the log.Logger interface by capturing logged
// lines.
type Logger struct {
	mutex *sync.Mutex

	// Entries holds logged entries in submission order.
	Entries []string
}

func (logger *Logger) Log(v ...interface{}) {
	logger.log(fmt.Sprint(v...))
}

func (logger *Logger) Logf(format string, v ...interface{}) {
	logger.log(fmt.Sprintf(format, v...))
}

func (logger *Logger) log(entry string) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.Entries = append(logger.Entries, entry)
}

func New() *Logger {
	return &Logger{
		mutex:   &sync.Mutex{},
		Entries: []string{},
	}
}
