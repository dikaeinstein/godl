package exitcode

import "errors"

// NewError returns an exitcode Error which sets the exit code of the specified error.
func NewError(err error, code int) error {
	if err == nil {
		return nil
	}

	return Error{err, code}
}

type Error struct {
	error
	code int
}

func (e Error) ExitCode() int {
	return e.code
}

func (e Error) Unwrap() error {
	return e.error
}

// Get returns the exit code of an error if it is an exitcode Error,
// otherwise it returns 1 for non-nil errors and 0 for nil errors.
func Get(err error) int {
	if err == nil {
		return 0
	}

	var exitcodeErr Error
	if errors.As(err, &exitcodeErr) {
		return exitcodeErr.ExitCode()
	}

	return 1
}
