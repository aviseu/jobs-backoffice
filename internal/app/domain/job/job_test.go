package job_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	suite.Run(t, new(JobSuite))
}

type JobSuite struct {
	suite.Suite
}

func (suite *JobSuite) Test_Success() {
	// Prepare
	id := uuid.New()
	chID := uuid.New()

	// Execute
	j := job.New(
		id,
		chID,
		job.StatusActive,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
	)

	// Assert
	suite.Equal(id, j.ID())
	suite.Equal(chID, j.ChannelID())
	suite.Equal(job.StatusActive, j.Status())
	suite.Equal(job.PublishStatusUnpublished, j.PublishStatus())
	suite.Equal("https://example.com/job/id", j.URL())
	suite.Equal("Software Engineer", j.Title())
	suite.Equal("Job Description", j.Description())
	suite.Equal("Indeed", j.Source())
	suite.Equal("Amsterdam", j.Location())
	suite.True(j.Remote())
	suite.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), j.PostedAt())
	suite.True(j.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(j.UpdatedAt().After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobSuite) Test_WithTimestamps__WithPublishStatus_Success() {
	// Prepare
	id := uuid.New()
	chID := uuid.New()

	// Execute
	j := job.New(
		id,
		chID,
		job.StatusActive,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		job.WithTimestamps(
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
		),
		job.WithPublishStatus(job.PublishStatusPublished),
	)

	// Assert
	suite.Equal(id, j.ID())
	suite.Equal(chID, j.ChannelID())
	suite.Equal(job.StatusActive, j.Status())
	suite.Equal(job.PublishStatusPublished, j.PublishStatus())
	suite.Equal("https://example.com/job/id", j.URL())
	suite.Equal("Software Engineer", j.Title())
	suite.Equal("Job Description", j.Description())
	suite.Equal("Indeed", j.Source())
	suite.Equal("Amsterdam", j.Location())
	suite.True(j.Remote())
	suite.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), j.PostedAt())
	suite.Equal(time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC), j.CreatedAt())
	suite.Equal(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC), j.UpdatedAt())
}
