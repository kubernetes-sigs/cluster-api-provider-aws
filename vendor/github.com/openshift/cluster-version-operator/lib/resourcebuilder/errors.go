package resourcebuilder

// RetryLaterError instructs the resource can't be reconciled right now, so retry later.
type RetryLaterError struct {
	Message string
}

func (e *RetryLaterError) Error() string {
	return e.Message
}
