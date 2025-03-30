package importing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Config struct {
	Arbeitnow arbeitnow.Config

	Import struct {
		ResultBufferSize int `split_words:"true" default:"10"`
		ResultWorkers    int `split_words:"true" default:"10"`
	}
}

type ImportRepository interface {
	SaveImport(ctx context.Context, i *aggregator.Import) error
	SaveImportJob(ctx context.Context, importID uuid.UUID, j *aggregator.ImportJob) error

	FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error)
}

type JobRepository interface {
	Save(ctx context.Context, j *aggregator.Job) error
	GetByChannelID(ctx context.Context, chID uuid.UUID) ([]*aggregator.Job, error)
}

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type Service struct {
	jr           JobRepository
	ir           ImportRepository
	chr          ChannelRepository
	f            *Factory
	log          *slog.Logger
	cfg          Config
	workerBuffer int
	workerCount  int
}

func NewService(chr ChannelRepository, ir ImportRepository, jr JobRepository, f *Factory, cfg Config, buffer, workers int, log *slog.Logger) *Service {
	return &Service{
		chr:          chr,
		jr:           jr,
		ir:           ir,
		f:            f,
		log:          log,
		cfg:          cfg,
		workerBuffer: buffer,
		workerCount:  workers,
	}
}

func (s *Service) worker(ctx context.Context, wg *sync.WaitGroup, i *Import, results <-chan *Result) {
	for r := range results {
		j := &aggregator.ImportJob{
			ID:     r.JobID(),
			Result: aggregator.ImportJobResult(r.Type()),
		}
		if err := s.ir.SaveImportJob(ctx, i.ID(), j); err != nil {
			s.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", j.ID, i.ID(), err).Error())
			continue
		}
	}
	wg.Done()
}

func jsworker(ctx context.Context, wg *sync.WaitGroup, r JobRepository, jobs <-chan *Job, results chan<- *Result, errs chan<- error) {
	for j := range jobs {
		if err := r.Save(ctx, j.ToDTO()); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.ID(), err)
			results <- NewResult(j.ID(), ResultTypeFailed, WithError(err.Error()))
		}
	}
	wg.Done()
}

func (s *Service) Import(ctx context.Context, importID uuid.UUID) error {
	i, err := s.ir.FindImport(ctx, importID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", importID, err)
	}

	idm := NewImportFromDTO(i)

	ch, err := s.chr.Find(ctx, idm.ChannelID())
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", idm.ChannelID(), err)
	}

	p, err := s.f.Create(ch)
	if err != nil {
		return fmt.Errorf("failed to create provider for channel %s: %w", ch.ID, err)
	}

	idm.status = aggregator.ImportStatusFetching
	if err := s.ir.SaveImport(ctx, idm.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", idm.ID(), err)
	}

	jobs, err := p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", ch.ID, err)

		idm.status = aggregator.ImportStatusFailed
		idm.endedAt = null.TimeFrom(time.Now())
		idm.error = null.StringFrom(err.Error())

		if err2 := s.ir.SaveImport(ctx, idm.ToAggregate()); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", idm.ID(), err2, err)
		}

		return err
	}

	idm.status = aggregator.ImportStatusProcessing
	if err := s.ir.SaveImport(ctx, idm.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", idm.ID(), err)
	}

	// create workers
	var wg sync.WaitGroup
	results := make(chan *Result, s.cfg.Import.ResultBufferSize)
	for w := 1; w <= s.cfg.Import.ResultWorkers; w++ {
		wg.Add(1)
		go s.worker(ctx, &wg, idm, results)
	}

	// get existing jobs
	existing, err := s.jr.GetByChannelID(ctx, ch.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing jobs: %w", err)
	}

	// create job workers
	var wgWorkers sync.WaitGroup
	jobChan := make(chan *Job, s.workerBuffer)
	errs := make(chan error, s.workerBuffer)
	for w := 1; w <= s.workerCount; w++ {
		wgWorkers.Add(1)
		go jsworker(ctx, &wgWorkers, s.jr, jobChan, results, errs)
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

	incomingJobs := make([]*Job, len(jobs))
	for i, in := range jobs {
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
		jobChan <- in
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
		if ex.Status == aggregator.JobStatusInactive {
			continue
		}
		for _, in := range jobs {
			if exj.ID() == in.ID {
				goto skip
			}
		}

		exj.MarkAsMissing()
		jobChan <- exj
		results <- NewResult(exj.ID(), ResultTypeMissing)

	skip:
	}

	close(jobChan)
	wgWorkers.Wait()
	close(errs)
	wgError.Wait()

	close(results)
	wg.Wait()

	idm.status = aggregator.ImportStatusPublishing
	if err := s.ir.SaveImport(ctx, idm.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status publishing for import %s: %w", idm.ID(), err)
	}

	// Publish (eventually)

	idm.status = aggregator.ImportStatusCompleted
	idm.endedAt = null.TimeFrom(time.Now())
	if err := s.ir.SaveImport(ctx, idm.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", idm.ID(), err)
	}

	return nil
}
