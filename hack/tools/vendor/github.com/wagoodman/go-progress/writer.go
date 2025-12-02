package progress

// Writer will consume a throw away bytes its given (not a passthrough). This is intended to be used with io.MultiWriter
type Writer struct {
	current int64
	size    int64
	err     error
}

func NewSizedWriter(size int64) *Writer {
	return &Writer{
		size: size,
	}
}

func NewWriter() *Writer {
	return &Writer{
		size: -1,
	}
}

func (w *Writer) SetComplete() {
	w.err = ErrCompleted
}

func (w *Writer) Write(p []byte) (int, error) {
	n := len(p)
	w.current += int64(n)
	return n, nil
}

func (w *Writer) Current() int64 {
	return w.current
}

func (w *Writer) Size() int64 {
	return w.size
}

func (w *Writer) Error() error {
	return w.err
}
