package errs_test

import (
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/errs"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestValidation(t *testing.T) {
	suite.Run(t, new(ValidationSuite))
}

type ValidationSuite struct {
	suite.Suite
}

func (suite *ValidationSuite) Test_Create_Success() {
	// Prepare
	err := errs.NewValidationError(errors.New("Name is required"))

	// Assert
	suite.ErrorContains(err, "Name is required")
	suite.True(errs.IsValidationError(err))
}
