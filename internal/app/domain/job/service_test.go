package job_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

func TestJobService(t *testing.T) {
	suite.Run(t, new(JobServiceSuite))
}

type JobServiceSuite struct {
	suite.Suite
}

func (suite *JobServiceSuite) Test_Sync_Success() {
	// Prepare
	r := testutils.NewJobRepository()
	s := job.NewJobService(r, 10, 10)

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
	existingNoChange := job.NewJob(uuid.New(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished))
	existingChange := job.NewJob(uuid.New(), chID, base.JobStatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished))
	existingActiveMissing := job.NewJob(uuid.New(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished))
	existingInactiveMissing := job.NewJob(uuid.New(), chID, base.JobStatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished))
	r.Add(existingNoChange.ToDTO())
	r.Add(existingChange.ToDTO())
	r.Add(existingActiveMissing.ToDTO())
	r.Add(existingInactiveMissing.ToDTO())

	incomingNoChange := job.NewJob(existingNoChange.ID(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, existingNoChange.PostedAt(), job.JobWithPublishStatus(base.JobPublishStatusUnpublished))
	incomingChange := job.NewJob(existingChange.ID(), chID, base.JobStatusActive, "https://example.com/job/id", "Title Changed!", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusUnpublished))
	incomingNew := job.NewJob(uuid.New(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusUnpublished))
	incoming := []*job.Job{incomingNoChange, incomingChange, incomingNew}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)
	close(results)
	wg.Wait()

	// Assert
	suite.NoError(err)
	suite.Len(r.Jobs, 5)

	suite.Equal(base.JobStatusActive, r.Jobs[existingNoChange.ID()].Status)
	suite.Equal(base.JobPublishStatusPublished, r.Jobs[existingNoChange.ID()].PublishStatus)

	suite.Equal(base.JobStatusActive, r.Jobs[existingChange.ID()].Status)
	suite.Equal(base.JobPublishStatusUnpublished, r.Jobs[existingChange.ID()].PublishStatus)

	suite.Equal(base.JobStatusInactive, r.Jobs[existingActiveMissing.ID()].Status)
	suite.Equal(base.JobPublishStatusUnpublished, r.Jobs[existingActiveMissing.ID()].PublishStatus)

	suite.Equal(base.JobStatusInactive, r.Jobs[existingInactiveMissing.ID()].Status)
	suite.Equal(base.JobPublishStatusPublished, r.Jobs[existingInactiveMissing.ID()].PublishStatus)

	suite.Equal(base.JobStatusActive, r.Jobs[incomingNew.ID()].Status)
	suite.Equal(base.JobPublishStatusUnpublished, r.Jobs[incomingNew.ID()].PublishStatus)

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

func (suite *JobServiceSuite) Test_Sync_RepositoryFail() {
	// Prepare
	r := testutils.NewJobRepository()
	r.FailWith(errors.New("boom!"))
	s := job.NewJobService(r, 10, 10)
	results := make(chan *job.Result)

	chID := uuid.New()
	incoming := []*job.Job{
		job.NewJob(uuid.New(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished)),
		job.NewJob(uuid.New(), chID, base.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), job.JobWithPublishStatus(base.JobPublishStatusPublished)),
	}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)

	// Assert
	suite.Error(err)
	suite.Equal("failed to get existing jobs: boom!", err.Error())
	close(results)
}
