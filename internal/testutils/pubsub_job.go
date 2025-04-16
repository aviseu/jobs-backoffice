package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"sync"
)

type PubSubJobService struct {
	Jobs []*aggregator.Job
	err  error
	m    sync.Mutex
}

func NewPubSubJobService() *PubSubJobService {
	return &PubSubJobService{}
}

func (p *PubSubJobService) FailWith(err error) {
	p.err = err
}

func (p *PubSubJobService) PublishJobInformation(_ context.Context, job *aggregator.Job) error {
	if p.err != nil {
		return p.err
	}

	p.m.Lock()
	defer p.m.Unlock()
	p.Jobs = append(p.Jobs, job)
	return nil
}
