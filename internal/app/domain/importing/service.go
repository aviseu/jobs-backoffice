package importing

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"sync"
)

type PubSubService interface {
	PublishJobInformation(ctx context.Context, job *aggregator.Job) error
}

type ConfigWorker struct {
	BufferSize int `split_words:"true" default:"10"`
	Workers    int `split_words:"true" default:"10"`
}

type Config struct {
	Arbeitnow arbeitnow.Config

	Import struct {
		Metric  ConfigWorker
		Job     ConfigWorker
		Publish ConfigWorker
	}
}

type ImportRepository interface {
	SaveImport(ctx context.Context, i *aggregator.Import) error
	SaveImportMetric(ctx context.Context, importID uuid.UUID, m *aggregator.ImportMetric) error

	FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error)
}

type JobRepository interface {
	Save(ctx context.Context, j *aggregator.Job) error
	GetByChannelID(ctx context.Context, chID uuid.UUID) ([]*aggregator.Job, error)
	GetActiveUnpublishedByChannelID(ctx context.Context, chID uuid.UUID) ([]*aggregator.Job, error)
}

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	jr  JobRepository
	ir  ImportRepository
	chr ChannelRepository
	f   *factory
	log *slog.Logger
	cfg Config
	pjs PubSubService
}

func NewService(chr ChannelRepository, ir ImportRepository, jr JobRepository, c HTTPClient, cfg Config, pjs PubSubService, log *slog.Logger) *Service {
	return &Service{
		chr: chr,
		jr:  jr,
		ir:  ir,
		f:   newFactory(c, cfg),
		pjs: pjs,
		log: log,
		cfg: cfg,
	}
}

func (s *Service) metricWorker(ctx context.Context, wg *sync.WaitGroup, i *importEntry, results <-chan *aggregator.ImportMetric, errs chan<- error) {
	for j := range results {
		if err := s.ir.SaveImportMetric(ctx, i.id, j); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.ID, err)
			s.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", j.ID, i.id, err).Error())
			continue
		}
	}
	wg.Done()
}

func (*Service) jobWorker(ctx context.Context, wg *sync.WaitGroup, r JobRepository, jobs <-chan *job, metrics chan<- *aggregator.ImportMetric, publish chan<- *job, errs chan<- error) {
	for j := range jobs {
		if err := r.Save(ctx, j.toAggregator()); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.id, err)
			metrics <- &aggregator.ImportMetric{ID: j.id, MetricType: aggregator.ImportMetricTypeError}
		}
		if j.needsPublishing() {
			publish <- j
		}
	}
	wg.Done()
}

func (s *Service) publishWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan *job, metrics chan<- *aggregator.ImportMetric, errs chan<- error) {
	for j := range jobs {
		// Publish
		err := s.pjs.PublishJobInformation(ctx, j.toAggregator())
		if err != nil {
			errs <- fmt.Errorf("failed to publish job %s: %w", j.id, err)
			metrics <- &aggregator.ImportMetric{ID: j.id, MetricType: aggregator.ImportMetricTypeError}
			continue
		}

		// Mark as published
		j.markAsPublished()

		// Save in db
		if err := s.jr.Save(ctx, j.toAggregator()); err != nil {
			errs <- fmt.Errorf("failed to save job %s: %w", j.id, err)
			metrics <- &aggregator.ImportMetric{ID: j.id, MetricType: aggregator.ImportMetricTypeError}
			continue
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
	i := newImportFromAggregator(importAggr)

	// Find related channel
	ch, err := s.chr.Find(ctx, i.channelID)
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.channelID, err)
	}

	// Create provider that will fetch jobs from external API
	p, err := s.f.create(ch)
	if err != nil {
		return fmt.Errorf("failed to create provider for channel %s: %w", ch.ID, err)
	}

	// *******************************************************
	// Import status: fetching
	// *******************************************************
	i.markAsFetching()
	if err := s.ir.SaveImport(ctx, i.toAggregate()); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", i.id, err)
	}

	// Fetch jobs from external API
	pJobs, err := p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", ch.ID, err)
		i.markAsFailed(err)
		if err2 := s.ir.SaveImport(ctx, i.toAggregate()); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", i.id, err2, err)
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
	if err := s.ir.SaveImport(ctx, i.toAggregate()); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.id, err)
	}

	errs := make(chan error, s.cfg.Import.Job.BufferSize)

	// Metric workers
	var metricsWG sync.WaitGroup
	metrics := make(chan *aggregator.ImportMetric, s.cfg.Import.Metric.BufferSize)
	for w := 1; w <= s.cfg.Import.Metric.Workers; w++ {
		metricsWG.Add(1)
		go s.metricWorker(ctx, &metricsWG, i, metrics, errs)
	}

	// Publish workers
	var publishWG sync.WaitGroup
	jobsToPublish := make(chan *job, s.cfg.Import.Publish.BufferSize)
	for w := 1; w <= s.cfg.Import.Publish.Workers; w++ {
		publishWG.Add(1)
		go s.publishWorker(ctx, &publishWG, jobsToPublish, metrics, errs)
	}

	// Job workers
	var jobsWG sync.WaitGroup
	jobsToSave := make(chan *job, s.cfg.Import.Job.BufferSize)
	for w := 1; w <= s.cfg.Import.Job.Workers; w++ {
		jobsWG.Add(1)
		go s.jobWorker(ctx, &jobsWG, s.jr, jobsToSave, metrics, jobsToPublish, errs)
	}

	// Error workers
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
					metrics <- &aggregator.ImportMetric{ID: incoming.id, MetricType: aggregator.ImportMetricTypeNoChange}
					goto next
				}
			}
		}

		incoming.markAsChanged()
		jobsToSave <- incoming
		if found {
			metrics <- &aggregator.ImportMetric{ID: incoming.id, MetricType: aggregator.ImportMetricTypeUpdated}
		} else {
			metrics <- &aggregator.ImportMetric{ID: incoming.id, MetricType: aggregator.ImportMetricTypeNew}
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
		metrics <- &aggregator.ImportMetric{ID: existing.id, MetricType: aggregator.ImportMetricTypeMissing}

	skip:
	}

	// Close channels and wait for workers to finish
	close(jobsToSave)
	jobsWG.Wait()
	close(metrics)
	metricsWG.Wait()

	// *******************************************************
	// Import status: publishing
	// *******************************************************
	i.markAsPublishing()
	if err := s.ir.SaveImport(ctx, i.toAggregate()); err != nil {
		return fmt.Errorf("failed to set status publishing for import %s: %w", i.id, err)
	}

	// Get all jobs needing publishing
	jj, err := s.jr.GetActiveUnpublishedByChannelID(ctx, ch.ID)
	if err != nil {
		return fmt.Errorf("failed to get unpublished jobs for channel %s: %w", ch.ID, err)
	}

	for _, j := range jj {
		jobsToPublish <- newJobFromAggregator(j)
	}

	close(jobsToPublish)
	publishWG.Wait()
	close(errs)
	errorWG.Wait()

	// *******************************************************
	// Import status: completed
	// *******************************************************
	i.markAsCompleted()
	if err := s.ir.SaveImport(ctx, i.toAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.id, err)
	}

	return nil
}
