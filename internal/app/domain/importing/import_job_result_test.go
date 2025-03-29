package importing_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestImportJobResult(t *testing.T) {
	suite.Run(t, new(ImportJobResultSuite))
}

type ImportJobResultSuite struct {
	suite.Suite
}

func (suite *ImportJobResultSuite) Test_Success() {
	// Prepare
	id := uuid.New()
	importID := uuid.New()

	// Execute
	j := importing.NewImportJobResult(id, importID, base.ImportJobResultNew)

	// Assert
	suite.Equal(id, j.JobID())
	suite.Equal(importID, j.ImportID())
	suite.Equal(base.ImportJobResultNew, j.Result())
}
