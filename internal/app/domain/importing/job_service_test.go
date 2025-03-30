package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
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
	s := importing.NewJobService(r, 10, 10)

	results := make(chan *importing.Result, 10)
	resultMap := make(map[uuid.UUID]*importing.Result)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(results <-chan *importing.Result) {
		for r := range results {
			resultMap[r.JobID()] = r
		}
		wg.Done()
	}(results)

	chID := uuid.New()
	existingNoChange := importing.NewJob(uuid.New(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished))
	existingChange := importing.NewJob(uuid.New(), chID, aggregator.JobStatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished))
	existingActiveMissing := importing.NewJob(uuid.New(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished))
	existingInactiveMissing := importing.NewJob(uuid.New(), chID, aggregator.JobStatusInactive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished))
	r.Add(existingNoChange.ToDTO())
	r.Add(existingChange.ToDTO())
	r.Add(existingActiveMissing.ToDTO())
	r.Add(existingInactiveMissing.ToDTO())

	incomingNoChange := importing.NewJob(existingNoChange.ID(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, existingNoChange.PostedAt(), importing.JobWithPublishStatus(aggregator.JobPublishStatusUnpublished))
	incomingChange := importing.NewJob(existingChange.ID(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Title Changed!", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusUnpublished))
	incomingNew := importing.NewJob(uuid.New(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusUnpublished))
	incoming := []*aggregator.Job{incomingNoChange.ToDTO(), incomingChange.ToDTO(), incomingNew.ToDTO()}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)
	close(results)
	wg.Wait()

	// Assert
	suite.NoError(err)
	suite.Len(r.Jobs, 5)

	suite.Equal(aggregator.JobStatusActive, r.Jobs[existingNoChange.ID()].Status)
	suite.Equal(aggregator.JobPublishStatusPublished, r.Jobs[existingNoChange.ID()].PublishStatus)

	suite.Equal(aggregator.JobStatusActive, r.Jobs[existingChange.ID()].Status)
	suite.Equal(aggregator.JobPublishStatusUnpublished, r.Jobs[existingChange.ID()].PublishStatus)

	suite.Equal(aggregator.JobStatusInactive, r.Jobs[existingActiveMissing.ID()].Status)
	suite.Equal(aggregator.JobPublishStatusUnpublished, r.Jobs[existingActiveMissing.ID()].PublishStatus)

	suite.Equal(aggregator.JobStatusInactive, r.Jobs[existingInactiveMissing.ID()].Status)
	suite.Equal(aggregator.JobPublishStatusPublished, r.Jobs[existingInactiveMissing.ID()].PublishStatus)

	suite.Equal(aggregator.JobStatusActive, r.Jobs[incomingNew.ID()].Status)
	suite.Equal(aggregator.JobPublishStatusUnpublished, r.Jobs[incomingNew.ID()].PublishStatus)

	// Assert results
	suite.Len(resultMap, 4)

	suite.Equal(existingNoChange.ID(), resultMap[existingNoChange.ID()].JobID())
	suite.Equal(importing.ResultTypeNoChange, resultMap[existingNoChange.ID()].Type())

	suite.Equal(existingChange.ID(), resultMap[existingChange.ID()].JobID())
	suite.Equal(importing.ResultTypeUpdated, resultMap[existingChange.ID()].Type())

	suite.Equal(existingActiveMissing.ID(), resultMap[existingActiveMissing.ID()].JobID())
	suite.Equal(importing.ResultTypeMissing, resultMap[existingActiveMissing.ID()].Type())

	suite.Equal(incomingNew.ID(), resultMap[incomingNew.ID()].JobID())
	suite.Equal(importing.ResultTypeNew, resultMap[incomingNew.ID()].Type())
}

func (suite *JobServiceSuite) Test_Sync_RepositoryFail() {
	// Prepare
	r := testutils.NewJobRepository()
	r.FailWith(errors.New("boom!"))
	s := importing.NewJobService(r, 10, 10)
	results := make(chan *importing.Result)

	chID := uuid.New()
	incoming := []*aggregator.Job{
		importing.NewJob(uuid.New(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished)).ToDTO(),
		importing.NewJob(uuid.New(), chID, aggregator.JobStatusActive, "https://example.com/job/id", "Software Engineer", "Job Description", "Indeed", "Amsterdam", true, time.Now(), importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished)).ToDTO(),
	}

	// Execute
	err := s.Sync(context.Background(), chID, incoming, results)

	// Assert
	suite.Error(err)
	suite.Equal("failed to get existing jobs: boom!", err.Error())
	close(results)
}
