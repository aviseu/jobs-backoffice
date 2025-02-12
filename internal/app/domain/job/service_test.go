package job_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) Test_Sync_Success() {
	// Prepare
	r := testutils.NewJobRepository()
	s := job.NewService(r)

	chID := uuid.New()
	existingNoChange := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	existingChange := job.New(uuid.New(), chID, job.StatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	existingMissing := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	r.Add(existingNoChange)
	r.Add(existingChange)
	r.Add(existingMissing)

	incomingNoChange := job.New(existingNoChange.ID(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incomingChange := job.New(existingChange.ID(), chID, job.StatusActive, "https://example.com/job/id", "Title Changed!", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incomingNew := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incoming := []*job.Job{incomingNoChange, incomingChange, incomingNew}

	// Execute
	err := s.Sync(context.Background(), chID, incoming)

	// Assert
	suite.NoError(err)
	suite.Len(r.Jobs, 4)

	suite.Equal(job.StatusActive, r.Jobs[existingNoChange.ID()].Status())
	suite.Equal(job.PublishStatusPublished, r.Jobs[existingNoChange.ID()].PublishStatus())

	suite.Equal(job.StatusActive, r.Jobs[existingChange.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[existingChange.ID()].PublishStatus())

	suite.Equal(job.StatusInactive, r.Jobs[existingMissing.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[existingMissing.ID()].PublishStatus())

	suite.Equal(job.StatusActive, r.Jobs[incomingNew.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[incomingNew.ID()].PublishStatus())
}

func (suite *ServiceSuite) Test_Sync_RepositoryFail() {
	// Prepare
	r := testutils.NewJobRepository()
	r.FailWith(errors.New("boom!"))
	s := job.NewService(r)

	chID := uuid.New()
	incoming := []*job.Job{
		job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished)),
		job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished)),
	}

	// Execute
	err := s.Sync(context.Background(), chID, incoming)

	// Assert
	suite.Error(err)
	suite.Equal("failed to get existing jobs: boom!", err.Error())
}
