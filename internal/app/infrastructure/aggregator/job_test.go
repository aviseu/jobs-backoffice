package aggregator_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestJob(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobSuite))
}

type JobSuite struct {
	suite.Suite
}

func (suite *IntegrationSuite) Test_JobStatus_Success() {
	suite.Equal("inactive", aggregator.JobStatusInactive.String())
	suite.Equal("active", aggregator.JobStatusActive.String())
}

func (suite *IntegrationSuite) Test_JobPublishStatus_Success() {
	suite.Equal("unpublished", aggregator.JobPublishStatusUnpublished.String())
	suite.Equal("published", aggregator.JobPublishStatusPublished.String())
}
