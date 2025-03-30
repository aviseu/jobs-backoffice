package importing

import (
	"context"
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

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type Service struct {
	ir  ImportRepository
	chr ChannelRepository
	is  *ImportService
	f   *Factory
	js  *JobService
	log *slog.Logger
	cfg Config
}

func NewService(chr ChannelRepository, ir ImportRepository, is *ImportService, js *JobService, f *Factory, cfg Config, log *slog.Logger) *Service {
	return &Service{
		chr: chr,
		is:  is,
		ir:  ir,
		f:   f,
		js:  js,
		log: log,
		cfg: cfg,
	}
}

func (s *Service) worker(ctx context.Context, wg *sync.WaitGroup, i *Import, results <-chan *Result) {
	for r := range results {
		j := &aggregator.ImportJob{
			ID:     r.JobID(),
			Result: aggregator.ImportJobResult(r.Type()),
		}
		if err := s.is.SaveJobResult(ctx, i.ID(), j); err != nil {
			s.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", j.ID, i.ID(), err).Error())
			continue
		}
	}
	wg.Done()
}

func (s *Service) Import(ctx context.Context, importID uuid.UUID) error {
	i, err := s.ir.FindImport(ctx, importID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", importID, err)
	}

	ch, err := s.chr.Find(ctx, i.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.ChannelID, err)
	}

	p, err := s.f.Create(ch)
	if err != nil {
		return fmt.Errorf("failed to create provider for channel %s: %w", ch.ID, err)
	}
	ip := NewImportFromDTO(i)

	if err := s.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusFetching); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", ip.ID(), err)
	}

	jobs, err := p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", ch.ID, err)
		if err2 := s.is.MarkAsFailed(ctx, ip.ToAggregate(), err); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", ip.ID(), err2, err)
		}
		return err
	}

	if err := s.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusProcessing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", ip.ID(), err)
	}

	// create workers
	var wg sync.WaitGroup
	results := make(chan *Result, s.cfg.Import.ResultBufferSize)
	for w := 1; w <= s.cfg.Import.ResultWorkers; w++ {
		wg.Add(1)
		go s.worker(ctx, &wg, ip, results)
	}

	if err := s.js.Sync(ctx, ch.ID, jobs, results); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", ch.ID, err)
	}

	close(results)
	wg.Wait()

	if err := s.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusPublishing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", ip.ID(), err)
	}

	// Publish (eventually)

	if err := s.is.MarkAsCompleted(ctx, ip.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", ip.ID(), err)
	}

	return nil
}
