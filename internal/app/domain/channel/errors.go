package channel

import "errors"

var (
	ErrInvalidIntegration = errors.New("invalid integration")
	ErrNameIsRequired     = errors.New("name is required")
)
