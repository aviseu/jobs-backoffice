package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/google/uuid"
)

type Provider interface {
	Channel() *channel.Channel
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
		jr := imports.NewResult(r.JobID(), i.ID(), imports.JobStatus(r.Type()))
		if err := g.is.SaveJobResult(ctx, jr); err != nil {
			g.log.Error(fmt.Errorf("failed to save job result %s for import %s: %w", jr.JobID(), jr.ImportID(), err).Error())
			continue
		}
	}
	wg.Done()
}

func (g *Gateway) ImportChannel(ctx context.Context) error {
	i, err := g.is.Start(ctx, uuid.New(), g.p.Channel().ID())
	if err != nil {
		return fmt.Errorf("failed to create import for channel %s: %w", g.p.Channel().ID(), err)
	}

	if err := g.is.SetStatus(ctx, i, imports.StatusFetching); err != nil {
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

	if err := g.is.SetStatus(ctx, i, imports.StatusProcessing); err != nil {
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

	if err := g.is.SetStatus(ctx, i, imports.StatusPublishing); err != nil {
		return fmt.Errorf("failed to set status processing for import %s: %w", i.ID(), err)
	}

	// Publish (eventually)

	if err := g.is.MarkAsCompleted(ctx, i); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}
