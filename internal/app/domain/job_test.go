package domain_test

import (
	"github.com/aviseu/jobs/internal/app/domain"
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

	// Execute
	job := domain.NewJob(
		id,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
	)

	// Assert
	suite.Equal(id, job.ID())
	suite.Equal("https://example.com/job/id", job.URL())
	suite.Equal("Software Engineer", job.Title())
	suite.Equal("Job Description", job.Description())
	suite.Equal("Indeed", job.Source())
	suite.Equal("Amsterdam", job.Location())
	suite.True(job.Remote())
	suite.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), job.PostedAt())
	suite.True(job.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(job.UpdatedAt().After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobSuite) Test_WithTimestamps_Success() {
	// Prepare
	id := uuid.New()

	// Execute
	job := domain.NewJob(
		id,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		domain.WithTimestamps(
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
		),
	)

	// Assert
	suite.Equal(id, job.ID())
	suite.Equal("https://example.com/job/id", job.URL())
	suite.Equal("Software Engineer", job.Title())
	suite.Equal("Job Description", job.Description())
	suite.Equal("Indeed", job.Source())
	suite.Equal("Amsterdam", job.Location())
	suite.True(job.Remote())
	suite.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), job.PostedAt())
	suite.Equal(time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC), job.CreatedAt())
	suite.Equal(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC), job.UpdatedAt())
}
