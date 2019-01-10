// Package nest allows users to use a Logger interface to
// create another Logger interface.
package nest

import (
	"fmt"

	"github.com/go-log/log"
)

const (
	// PreNest is a marker placed between the parent and child values
	// when calling the wrapped Log method.  For example:
	//
	//   parent := SomeLogger()
	//	 child := New(parent, "a", "b")
	//	 child.Log("c", "d")
	//
	// will result in:
	//
	//    parent.Log("a", "b", PreLog, "c", "d")
	PreNest Marker = "pre-nest"
)

// Marker is a string synonym.  The type difference allows underlying
// log implementations to distinguish between the PreNest marker and a
// "pre-nest" string literal.
type Marker string

// String returns a single space (regardless of the underlying marker
// string).  This makes the output of parent loggers based on a
// fmt.Print style more readable, because fmt.Print only inserts space
// between two non-string operands.
func (m Marker) String() string {
	return " "
}

type logger struct {
	logger log.Logger
	values []interface{}
}

func (logger *logger) Log(v ...interface{}) {
	values := append(logger.values, PreNest)
	logger.logger.Log(append(values, v...)...)
}

func (logger *logger) Logf(format string, v ...interface{}) {
	logger.Log(fmt.Sprintf(format, v...))
}

func New(log log.Logger, v ...interface{}) *logger {
	return &logger{
		logger: log,
		values: v,
	}
}

func Newf(log log.Logger, format string, v ...interface{}) *logger {
	return New(log, fmt.Sprintf(format, v...))
}
