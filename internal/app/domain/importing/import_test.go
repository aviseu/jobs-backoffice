package importing_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestImport(t *testing.T) {
	suite.Run(t, new(ImportSuite))
}

type ImportSuite struct {
	suite.Suite
}

func (suite *ImportSuite) Test_Success() {
	// Execute
	i := importing.NewImport(uuid.New(), uuid.New(), importing.ImportWithMetadata(1, 2, 3, 4, 5))

	// Assert
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(5, i.FailedJobs())
	suite.Equal(15, i.TotalJobs())
}
