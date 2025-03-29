package configuring

import (
	"errors"

	"github.com/aviseu/jobs-backoffice/internal/errs"
)

var (
	ErrInvalidIntegration = errs.NewValidationError(errors.New("invalid integration"))
	ErrNameIsRequired     = errs.NewValidationError(errors.New("name is required"))
	ErrChannelNotFound    = errs.NewValidationError(errors.New("channel not found"))
)
