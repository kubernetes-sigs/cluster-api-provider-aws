package fmt

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type fmtLogger struct {
	writer io.Writer
}

func (t *fmtLogger) Log(v ...interface{}) {
	t.output(fmt.Sprint(v...))
}

func (t *fmtLogger) Logf(format string, v ...interface{}) {
	t.output(fmt.Sprintf(format, v...))
}

func (logger *fmtLogger) output(line string) (n int, err error) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	return logger.writer.Write([]byte(line))
}

// New creates a new fmt logger which writes to stdout.
func New() *fmtLogger {
	return &fmtLogger{writer: os.Stdout}
}

// NewFromWriter creates a new fmt logger which writes to the given
// writer.
func NewFromWriter(writer io.Writer) *fmtLogger {
	return &fmtLogger{writer: writer}
}
