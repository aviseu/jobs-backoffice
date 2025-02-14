package job

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, j *Job) error
	GetByChannelID(ctx context.Context, chID uuid.UUID) ([]*Job, error)
}

type Service struct {
	r Repository

	workerBuffer int
	workerCount  int
}

func NewService(r Repository, workerBuffer, workerCount int) *Service {
	return &Service{
		r:            r,
		workerBuffer: workerBuffer,
		workerCount:  workerCount,
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, r Repository, jobs <-chan *Job, results chan<- *Result, errs chan<- error) {
	for j := range jobs {
		if err := r.Save(ctx, j); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.ID(), err)
			results <- NewResult(j.ID(), ResultTypeFailed, WithError(err.Error()))
		}
	}
	wg.Done()
}

func (s *Service) Sync(ctx context.Context, chID uuid.UUID, incoming []*Job, results chan<- *Result) error {
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

	// save if incoming does not exist or is different
	for _, in := range incoming {
		found := false
		for _, ex := range existing {
			if ex.ID() == in.ID() {
				found = true
				if in.IsEqual(ex) {
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
		for _, in := range incoming {
			if ex.ID() == in.ID() {
				goto skip
			}
		}

		ex.MarkAsMissing()
		jobs <- ex
		results <- NewResult(ex.ID(), ResultTypeMissing)

	skip:
	}

	close(jobs)
	wgWorkers.Wait()
	close(errs)
	wgError.Wait()

	return syncErrs
}
