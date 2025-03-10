package job_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"sync"
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
	s := job.NewService(r, 10, 10)

	results := make(chan *job.Result, 10)
	resultMap := make(map[uuid.UUID]*job.Result)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(results <-chan *job.Result) {
		for r := range results {
			resultMap[r.JobID()] = r
		}
		wg.Done()
	}(results)

	chID := uuid.New()
	existingNoChange := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	existingChange := job.New(uuid.New(), chID, job.StatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	existingActiveMissing := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	existingInactiveMissing := job.New(uuid.New(), chID, job.StatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished))
	r.Add(existingNoChange)
	r.Add(existingChange)
	r.Add(existingActiveMissing)
	r.Add(existingInactiveMissing)

	incomingNoChange := job.New(existingNoChange.ID(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, existingNoChange.PostedAt(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incomingChange := job.New(existingChange.ID(), chID, job.StatusActive, "https://example.com/job/id", "Title Changed!", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incomingNew := job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusUnpublished))
	incoming := []*job.Job{incomingNoChange, incomingChange, incomingNew}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)
	close(results)
	wg.Wait()

	// Assert
	suite.NoError(err)
	suite.Len(r.Jobs, 5)

	suite.Equal(job.StatusActive, r.Jobs[existingNoChange.ID()].Status())
	suite.Equal(job.PublishStatusPublished, r.Jobs[existingNoChange.ID()].PublishStatus())

	suite.Equal(job.StatusActive, r.Jobs[existingChange.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[existingChange.ID()].PublishStatus())

	suite.Equal(job.StatusInactive, r.Jobs[existingActiveMissing.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[existingActiveMissing.ID()].PublishStatus())

	suite.Equal(job.StatusInactive, r.Jobs[existingInactiveMissing.ID()].Status())
	suite.Equal(job.PublishStatusPublished, r.Jobs[existingInactiveMissing.ID()].PublishStatus())

	suite.Equal(job.StatusActive, r.Jobs[incomingNew.ID()].Status())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[incomingNew.ID()].PublishStatus())

	// Assert results
	suite.Len(resultMap, 4)

	suite.Equal(existingNoChange.ID(), resultMap[existingNoChange.ID()].JobID())
	suite.Equal(job.ResultTypeNoChange, resultMap[existingNoChange.ID()].Type())

	suite.Equal(existingChange.ID(), resultMap[existingChange.ID()].JobID())
	suite.Equal(job.ResultTypeUpdated, resultMap[existingChange.ID()].Type())

	suite.Equal(existingActiveMissing.ID(), resultMap[existingActiveMissing.ID()].JobID())
	suite.Equal(job.ResultTypeMissing, resultMap[existingActiveMissing.ID()].Type())

	suite.Equal(incomingNew.ID(), resultMap[incomingNew.ID()].JobID())
	suite.Equal(job.ResultTypeNew, resultMap[incomingNew.ID()].Type())
}

func (suite *ServiceSuite) Test_Sync_RepositoryFail() {
	// Prepare
	r := testutils.NewJobRepository()
	r.FailWith(errors.New("boom!"))
	s := job.NewService(r, 10, 10)
	results := make(chan *job.Result)

	chID := uuid.New()
	incoming := []*job.Job{
		job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished)),
		job.New(uuid.New(), chID, job.StatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.WithPublishStatus(job.PublishStatusPublished)),
	}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)

	// Assert
	suite.Error(err)
	suite.Equal("failed to get existing jobs: boom!", err.Error())
	close(results)
}
