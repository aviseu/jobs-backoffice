package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"sync"
)

type JobRepository struct {
	Jobs map[uuid.UUID]*aggregator.Job
	err  error
	m    sync.Mutex
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		Jobs: make(map[uuid.UUID]*aggregator.Job),
	}
}

func (r *JobRepository) First() *aggregator.Job {
	for _, j := range r.Jobs {
		return j
	}

	return nil
}

func (r *JobRepository) Add(j *aggregator.Job) {
	r.Jobs[j.ID] = j
}

func (r *JobRepository) FailWith(err error) {
	r.err = err
}

func (r *JobRepository) Save(_ context.Context, j *aggregator.Job) error {
	if r.err != nil {
		return r.err
	}

	r.m.Lock()
	r.Jobs[j.ID] = j
	r.m.Unlock()
	return nil
}

func (r *JobRepository) GetByChannelID(_ context.Context, chID uuid.UUID) ([]*aggregator.Job, error) {
	if r.err != nil {
		return nil, r.err
	}

	var jobs []*aggregator.Job
	for _, j := range r.Jobs {
		if j.ChannelID == chID {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

func (r *JobRepository) GetActiveUnpublishedByChannelID(_ context.Context, chID uuid.UUID) ([]*aggregator.Job, error) {
	if r.err != nil {
		return nil, r.err
	}

	var jobs []*aggregator.Job
	for _, j := range r.Jobs {
		if j.ChannelID == chID && j.PublishStatus == aggregator.JobPublishStatusUnpublished && j.Status == aggregator.JobStatusActive {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}
