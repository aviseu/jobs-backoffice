package importing

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
)

type Provider interface {
	Channel() *aggregator.Channel
	GetJobs() ([]*aggregator.Job, error)
}

type Gateway struct {
	p   Provider
	js  *JobService
	is  *ImportService
	log *slog.Logger

	resultBufferSize int
	resultWorkers    int
}

func NewGateway(p Provider, js *JobService, is *ImportService, log *slog.Logger, resultBufferSize, resultWorkers int) *Gateway {
	return &Gateway{
		p:   p,
		js:  js,
		is:  is,
		log: log,

		resultBufferSize: resultBufferSize,
		resultWorkers:    resultWorkers,
	}
}

func (g *Gateway) worker(ctx context.Context, wg *sync.WaitGroup, i *Import, results <-chan *Result) {
	for r := range results {
		j := &aggregator.ImportJob{
			ID:     r.JobID(),
			Result: aggregator.ImportJobResult(r.Type()),
		}
		if err := g.is.SaveJobResult(ctx, i.ID(), j); err != nil {
			g.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", j.ID, i.ID(), err).Error())
			continue
		}
	}
	wg.Done()
}

func (g *Gateway) Import(ctx context.Context, i *aggregator.Import) error {
	ip := NewImportFromDTO(i)

	if err := g.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusFetching); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", ip.ID(), err)
	}

	jobs, err := g.p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", g.p.Channel().ID, err)
		if err2 := g.is.MarkAsFailed(ctx, ip.ToAggregate(), err); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", ip.ID(), err2, err)
		}
		return err
	}

	if err := g.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusProcessing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", ip.ID(), err)
	}

	// create workers
	var wg sync.WaitGroup
	results := make(chan *Result, g.resultBufferSize)
	for w := 1; w <= g.resultWorkers; w++ {
		wg.Add(1)
		go g.worker(ctx, &wg, ip, results)
	}

	if err := g.js.Sync(ctx, g.p.Channel().ID, jobs, results); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", g.p.Channel().ID, err)
	}

	close(results)
	wg.Wait()

	if err := g.is.SetStatus(ctx, ip.ToAggregate(), aggregator.ImportStatusPublishing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", ip.ID(), err)
	}

	// Publish (eventually)

	if err := g.is.MarkAsCompleted(ctx, ip.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", ip.ID(), err)
	}

	return nil
}
