package importing

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"sync"

	"github.com/google/uuid"
)

type JobRepository interface {
	Save(ctx context.Context, j *postgres.Job) error
	GetByChannelID(ctx context.Context, chID uuid.UUID) ([]*postgres.Job, error)
}

type JobService struct {
	r JobRepository

	workerBuffer int
	workerCount  int
}

func NewJobService(r JobRepository, workerBuffer, workerCount int) *JobService {
	return &JobService{
		r:            r,
		workerBuffer: workerBuffer,
		workerCount:  workerCount,
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, r JobRepository, jobs <-chan *Job, results chan<- *Result, errs chan<- error) {
	for j := range jobs {
		if err := r.Save(ctx, j.ToDTO()); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.ID(), err)
			results <- NewResult(j.ID(), ResultTypeFailed, WithError(err.Error()))
		}
	}
	wg.Done()
}

func (s *JobService) Sync(ctx context.Context, chID uuid.UUID, incoming []*postgres.Job, results chan<- *Result) error {
	// get existing jobs
	existing, err := s.r.GetByChannelID(ctx, chID)
	if err != nil {
		return fmt.Errorf("failed to get existing jobs: %w", err)
	}

	// create job workers
	var wgWorkers sync.WaitGroup
	jobs := make(chan *Job, s.workerBuffer)
	errs := make(chan error, s.workerBuffer)
	for w := 1; w <= s.workerCount; w++ {
		wgWorkers.Add(1)
		go worker(ctx, &wgWorkers, s.r, jobs, results, errs)
	}

	// create error worker
	var syncErrs error
	var wgError sync.WaitGroup
	wgError.Add(1)
	go func(errs <-chan error) {
		for err := range errs {
			syncErrs = errors.Join(syncErrs, err)
		}
		wgError.Done()
	}(errs)

	incomingJobs := make([]*Job, len(incoming))
	for i, in := range incoming {
		inj := NewJobFromDTO(in)
		incomingJobs[i] = inj
	}

	// save if incoming does not exist or is different
	for _, in := range incomingJobs {
		found := false
		for _, ex := range existing {
			if ex.ID == in.ID() {
				found = true
				if in.IsEqual(NewJobFromDTO(ex)) {
					results <- NewResult(in.ID(), ResultTypeNoChange)
					goto next
				}
			}
		}

		in.MarkAsChanged()
		jobs <- in
		if found {
			results <- NewResult(in.ID(), ResultTypeUpdated)
		} else {
			results <- NewResult(in.ID(), ResultTypeNew)
		}

	next:
	}

	// save if existing does not exist in incoming
	for _, ex := range existing {
		exj := NewJobFromDTO(ex)
		if ex.Status == base.JobStatusInactive {
			continue
		}
		for _, in := range incoming {
			if exj.ID() == in.ID {
				goto skip
			}
		}

		exj.MarkAsMissing()
		jobs <- exj
		results <- NewResult(exj.ID(), ResultTypeMissing)

	skip:
	}

	close(jobs)
	wgWorkers.Wait()
	close(errs)
	wgError.Wait()

	return syncErrs
}
