package importing

import (
	"context"
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

	GetImports(ctx context.Context) ([]*aggregator.Import, error)
	FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error)
}

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type Service struct {
	ir  ImportRepository
	chr ChannelRepository
	f   *Factory
	js  *JobService
	log *slog.Logger
	cfg Config
}

func NewService(chr ChannelRepository, ir ImportRepository, js *JobService, f *Factory, cfg Config, log *slog.Logger) *Service {
	return &Service{
		chr: chr,
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
		if err := s.ir.SaveImportJob(ctx, i.ID(), j); err != nil {
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

	if err := s.js.Sync(ctx, ch.ID, jobs, results); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", ch.ID, err)
	}

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
