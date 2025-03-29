package importing

import (
	"errors"

	"github.com/aviseu/jobs-backoffice/internal/errs"
)

var ErrImportNotFound = errs.NewValidationError(errors.New("import not found"))
