package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/google/uuid"
	"sync"
)

type JobRepository struct {
	Jobs map[uuid.UUID]*postgres.Job
	err  error
	m    sync.Mutex
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		Jobs: make(map[uuid.UUID]*postgres.Job),
	}
}

func (r *JobRepository) First() *postgres.Job {
	for _, j := range r.Jobs {
		return j
	}

	return nil
}

func (r *JobRepository) Add(j *postgres.Job) {
	r.Jobs[j.ID] = j
}

func (r *JobRepository) FailWith(err error) {
	r.err = err
}

func (r *JobRepository) Save(_ context.Context, j *postgres.Job) error {
	if r.err != nil {
		return r.err
	}

	r.m.Lock()
	r.Jobs[j.ID] = j
	r.m.Unlock()
	return nil
}

func (r *JobRepository) GetByChannelID(_ context.Context, chID uuid.UUID) ([]*postgres.Job, error) {
	if r.err != nil {
		return nil, r.err
	}

	var jobs []*postgres.Job
	for _, j := range r.Jobs {
		if j.ChannelID == chID {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}
