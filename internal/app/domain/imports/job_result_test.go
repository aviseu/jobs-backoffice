package imports_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestJobResult(t *testing.T) {
	suite.Run(t, new(JobResultSuite))
}

type JobResultSuite struct {
	suite.Suite
}

func (suite *JobResultSuite) Test_Success() {
	// Prepare
	id := uuid.New()
	importID := uuid.New()

	// Execute
	j := imports.NewResult(id, importID, base.ImportJobResultNew)

	// Assert
	suite.Equal(id, j.JobID())
	suite.Equal(importID, j.ImportID())
	suite.Equal(base.ImportJobResultNew, j.Result())
}
