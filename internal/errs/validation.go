package errs

import "errors"

type ValidationError struct {
	err error
}

func NewValidationError(err error) error {
	return &ValidationError{err: err}
}

func (e *ValidationError) Error() string {
	return e.err.Error()
}

func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}
