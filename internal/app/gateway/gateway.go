package gateway

import (
	"context"
	"fmt"

	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/domain/job"
)

type Provider interface {
	Channel() *channel.Channel
	GetJobs() ([]*job.Job, error)
}

type Gateway struct {
	p Provider
	s *job.Service
}

func NewGateway(p Provider, s *job.Service) *Gateway {
	return &Gateway{
		p: p,
		s: s,
	}
}

func (g *Gateway) ImportChannel(ctx context.Context) error {
	jobs, err := g.p.GetJobs()
	if err != nil {
		return fmt.Errorf("failed to import channel %s: %w", g.p.Channel().ID(), err)
	}

	if err := g.s.Sync(ctx, g.p.Channel().ID(), jobs); err != nil {
		return fmt.Errorf("failed to sync jobs for channel %s: %w", g.p.Channel().ID(), err)
	}

	return nil
}
