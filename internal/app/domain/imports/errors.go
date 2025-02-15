package imports

import (
	"errors"

	"github.com/aviseu/jobs-backoffice/internal/app/errs"
)

var ErrImportNotFound = errs.NewValidationError(errors.New("import not found"))
