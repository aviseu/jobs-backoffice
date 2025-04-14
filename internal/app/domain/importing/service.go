package importing

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/google/uuid"
	"log/slog"
	"sync"
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

func (s *Service) metricsWorker(ctx context.Context, wg *sync.WaitGroup, i *Import, results <-chan *Result) {
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

func (*Service) jobWorker(ctx context.Context, wg *sync.WaitGroup, r JobRepository, jobs <-chan *job, importJobs chan<- *Result, errs chan<- error) {
	for j := range jobs {
		if err := r.Save(ctx, j.toAggregator()); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.id, err)
			importJobs <- NewImportMetric(j.id, ResultTypeFailed, WithError(err.Error()))
		}
	}
	wg.Done()
}

func (s *Service) Import(ctx context.Context, importID uuid.UUID) error {
	// *******************************************************
	// Setup for importing
	// *******************************************************

	// Find import
	importAggr, err := s.ir.FindImport(ctx, importID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", importID, err)
	}
	i := NewImportFromAggregator(importAggr)

	// Find related channel
	ch, err := s.chr.Find(ctx, i.ChannelID())
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.ChannelID(), err)
	}

	// Create provider that will fetch jobs from external API
	p, err := s.f.Create(ch)
	if err != nil {
		return fmt.Errorf("failed to create provider for channel %s: %w", ch.ID, err)
	}

	// *******************************************************
	// Import status: fetching
	// *******************************************************
	i.markAsFetching()
	if err := s.ir.SaveImport(ctx, i.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", i.ID(), err)
	}

	// Fetch jobs from external API
	pJobs, err := p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", ch.ID, err)
		i.markAsFailed(err)
		if err2 := s.ir.SaveImport(ctx, i.ToAggregate()); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", i.ID(), err2, err)
		}

		return err
	}

	// Convert aggregator jobs into domain jobs
	incomingJobs := make([]*job, len(pJobs))
	for i, job := range pJobs {
		incomingJobs[i] = newJobFromAggregator(job)
	}

	// *******************************************************
	// Import status: processing
	// *******************************************************
	i.markAsProcessing()
	if err := s.ir.SaveImport(ctx, i.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// Create worker group for saving metrics (ImportJobs)
	var metricsWG sync.WaitGroup
	metrics := make(chan *Result, s.cfg.Import.ResultBufferSize)
	for w := 1; w <= s.cfg.Import.ResultWorkers; w++ {
		metricsWG.Add(1)
		go s.metricsWorker(ctx, &metricsWG, i, metrics)
	}

	// Create worker group for saving Jobs in the database
	var jobsWG sync.WaitGroup
	jobsToSave := make(chan *job, s.workerBuffer)
	errs := make(chan error, s.workerBuffer)
	for w := 1; w <= s.workerCount; w++ {
		jobsWG.Add(1)
		go s.jobWorker(ctx, &jobsWG, s.jr, jobsToSave, metrics, errs)
	}

	// Create worker for handling errors
	var syncErrs error
	var errorWG sync.WaitGroup
	errorWG.Add(1)
	go func(errs <-chan error) {
		for err := range errs {
			syncErrs = errors.Join(syncErrs, err)
		}
		errorWG.Done()
	}(errs)

	// Get existing jobs from the database
	dbJobs, err := s.jr.GetByChannelID(ctx, ch.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing jobs: %w", err)
	}

	// Convert aggregator jobs into domain jobs
	existingJobs := make([]*job, len(dbJobs))
	for i, job := range dbJobs {
		existingJobs[i] = newJobFromAggregator(job)
	}

	// Save incoming job if different or new
	for _, incoming := range incomingJobs {
		found := false
		for _, existing := range existingJobs {
			if incoming.id == existing.id {
				found = true
				if incoming.IsEqual(existing) {
					metrics <- NewImportMetric(incoming.id, ResultTypeNoChange)
					goto next
				}
			}
		}

		incoming.markAsChanged()
		jobsToSave <- incoming
		if found {
			metrics <- NewImportMetric(incoming.id, ResultTypeUpdated)
		} else {
			metrics <- NewImportMetric(incoming.id, ResultTypeNew)
		}

	next:
	}

	// Mark as missing if exists but didn't income
	for _, existing := range existingJobs {
		if existing.status == aggregator.JobStatusInactive {
			continue
		}
		for _, incoming := range incomingJobs {
			if existing.id == incoming.id {
				goto skip
			}
		}

		existing.markAsMissing()
		jobsToSave <- existing
		metrics <- NewImportMetric(existing.id, ResultTypeMissing)

	skip:
	}

	// Close channels and wait for workers to finish
	close(jobsToSave)
	jobsWG.Wait()
	close(errs)
	errorWG.Wait()
	close(metrics)
	metricsWG.Wait()

	// *******************************************************
	// Import status: publishing
	// *******************************************************
	i.markAsPublishing()
	if err := s.ir.SaveImport(ctx, i.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status publishing for import %s: %w", i.ID(), err)
	}

	// Publish (eventually)

	// *******************************************************
	// Import status: completed
	// *******************************************************
	i.markAsCompleted()
	if err := s.ir.SaveImport(ctx, i.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}
