package utils

import "fmt"

// SafeError wraps an error which may contain sensitive information, allowing
// it to be matched without exposing sensitive details when logging the error.
type SafeError struct {
	Err error
	Msg string
}

// NewSafeError returns a new SafeError.
func NewSafeError(err error, msg string) error {
	return &SafeError{err, msg}
}

// NewSafeErrorf returns a new SafeError with formatting.
func NewSafeErrorf(err error, format string, args ...interface{}) error {
	return &SafeError{err, fmt.Sprintf(format, args...)}
}

// Error returns the safe error message.
func (s *SafeError) Error() string {
	return s.Msg
}

// Unwrap returns the wrapped error, which may contain sensitive information.
func (s *SafeError) Unwrap() error {
	return s.Err
}

// UnsafeError returns the full error message, which may contain sensitive
// information.
func (s *SafeError) UnsafeError() string {
	return fmt.Sprintf("%s: %s", s.Msg, s.Err)
}
