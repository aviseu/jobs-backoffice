package aggregator_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestImport(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImportSuite))
}

type ImportSuite struct {
	suite.Suite
}

func (suite *ImportSuite) Test_ImportStatus_Success() {
	suite.Equal("pending", aggregator.ImportStatusPending.String())
	suite.Equal("fetching", aggregator.ImportStatusFetching.String())
	suite.Equal("processing", aggregator.ImportStatusProcessing.String())
	suite.Equal("publishing", aggregator.ImportStatusPublishing.String())
	suite.Equal("completed", aggregator.ImportStatusCompleted.String())
	suite.Equal("failed", aggregator.ImportStatusFailed.String())
}

func (suite *ImportSuite) Test_ImportPublishStatus_Success() {
	suite.Equal("new", aggregator.ImportJobResultNew.String())
	suite.Equal("updated", aggregator.ImportJobResultUpdated.String())
	suite.Equal("no_change", aggregator.ImportJobResultNoChange.String())
	suite.Equal("missing", aggregator.ImportJobResultMissing.String())
	suite.Equal("failed", aggregator.ImportJobResultFailed.String())
}

func (suite *ImportSuite) Test_Import_Success() {
	// Prepare
	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: uuid.New(),
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
		Jobs: []*aggregator.ImportJob{
			{ID: uuid.New(), Result: aggregator.ImportJobResultNew},
			{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated},
			{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated},
			{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange},
			{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange},
			{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange},
			{ID: uuid.New(), Result: aggregator.ImportJobResultMissing},
			{ID: uuid.New(), Result: aggregator.ImportJobResultMissing},
			{ID: uuid.New(), Result: aggregator.ImportJobResultMissing},
			{ID: uuid.New(), Result: aggregator.ImportJobResultMissing},
			{ID: uuid.New(), Result: aggregator.ImportJobResultFailed},
			{ID: uuid.New(), Result: aggregator.ImportJobResultFailed},
			{ID: uuid.New(), Result: aggregator.ImportJobResultFailed},
			{ID: uuid.New(), Result: aggregator.ImportJobResultFailed},
			{ID: uuid.New(), Result: aggregator.ImportJobResultFailed},
		},
	}

	// Assert
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(5, i.FailedJobs())
	suite.Equal(15, i.TotalJobs())
}
