package gateway

import (
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"log/slog"
	"sync"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
)

type Provider interface {
	Channel() *configuring.Channel
	GetJobs() ([]*job.Job, error)
}

type Gateway struct {
	p   Provider
	js  *job.Service
	is  *imports.Service
	log *slog.Logger

	resultBufferSize int
	resultWorkers    int
}

func NewGateway(p Provider, js *job.Service, is *imports.Service, log *slog.Logger, resultBufferSize, resultWorkers int) *Gateway {
	return &Gateway{
		p:   p,
		js:  js,
		is:  is,
		log: log,

		resultBufferSize: resultBufferSize,
		resultWorkers:    resultWorkers,
	}
}

func (g *Gateway) worker(ctx context.Context, wg *sync.WaitGroup, i *imports.Import, results <-chan *job.Result) {
	for r := range results {
		jr := imports.NewResult(r.JobID(), i.ID(), base.ImportJobResult(r.Type()))
		if err := g.is.SaveJobResult(ctx, jr); err != nil {
			g.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", jr.JobID(), jr.ImportID(), err).Error())
			continue
		}
	}
	wg.Done()
}

func (g *Gateway) Import(ctx context.Context, i *imports.Import) error {
	if err := g.is.SetStatus(ctx, i, base.ImportStatusFetching); err != nil {
		return fmt.Errorf("failed to set status fetching for import %s: %w", i.ID(), err)
	}

	jobs, err := g.p.GetJobs()
	if err != nil {
		err := fmt.Errorf("failed to import channel %s: %w", g.p.Channel().ID(), err)
		if err2 := g.is.MarkAsFailed(ctx, i, err); err2 != nil {
			return fmt.Errorf("failed to mark import %s as failed: %w: %w", i.ID(), err2, err)
		}
		return err
	}

	if err := g.is.SetStatus(ctx, i, base.ImportStatusProcessing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// create workers
	var wg sync.WaitGroup
	results := make(chan *job.Result, g.resultBufferSize)
	for w := 1; w <= g.resultWorkers; w++ {
		wg.Add(1)
		go g.worker(ctx, &wg, i, results)
	}

	if err := g.js.Sync(ctx, g.p.Channel().ID(), jobs, results); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", g.p.Channel().ID(), err)
	}

	close(results)
	wg.Wait()

	if err := g.is.SetStatus(ctx, i, base.ImportStatusPublishing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// Publish (eventually)

	if err := g.is.MarkAsCompleted(ctx, i); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}
