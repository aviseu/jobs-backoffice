package importing

import (
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"log/slog"
	"sync"
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

func (g *Gateway) Import(ctx context.Context, i *Import) error {
	if err := g.is.SetStatus(ctx, i, aggregator.ImportStatusFetching); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", i.ID(), err)
	}

	jobs, err := g.p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", g.p.Channel().ID, err)
		if err2 := g.is.MarkAsFailed(ctx, i, err); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", i.ID(), err2, err)
		}
		return err
	}

	if err := g.is.SetStatus(ctx, i, aggregator.ImportStatusProcessing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// create workers
	var wg sync.WaitGroup
	results := make(chan *Result, g.resultBufferSize)
	for w := 1; w <= g.resultWorkers; w++ {
		wg.Add(1)
		go g.worker(ctx, &wg, i, results)
	}

	if err := g.js.Sync(ctx, g.p.Channel().ID, jobs, results); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", g.p.Channel().ID, err)
	}

	close(results)
	wg.Wait()

	if err := g.is.SetStatus(ctx, i, aggregator.ImportStatusPublishing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// Publish (eventually)

	if err := g.is.MarkAsCompleted(ctx, i); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}
