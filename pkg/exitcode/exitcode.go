package exitcode

import "errors"

// Error wraps an error with an exit code
type Error interface {
	error
	ExitCode() int
}

// NewError returns an exitcodeError which sets the exit code of the specified error.
func NewError(err error, code int) error {
	return exitcodeError{err, code}
}

type exitcodeError struct {
	error
	code int
}

func (e exitcodeError) ExitCode() int {
	return e.code
}

func (e exitcodeError) Unwrap() error {
	return e.error
}

// Get returns the exit code of an error
func Get(err error) int {
	if err == nil {
		return 0
	}

	if exitcodeErr := Error(nil); errors.As(err, &exitcodeErr) {
		return exitcodeErr.ExitCode()
	}

	return 1
}
