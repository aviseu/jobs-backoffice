package importing

import (
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
)

type provider interface {
	GetJobs() ([]*aggregator.Job, error)
}

type factory struct {
	c   HTTPClient
	cfg Config
}

func newFactory(c HTTPClient, cfg Config) *factory {
	return &factory{
		cfg: cfg,
		c:   c,
	}
}

func (f *factory) create(ch *aggregator.Channel) (provider, error) {
	if ch.Integration == aggregator.IntegrationArbeitnow {
		return arbeitnow.NewService(f.c, f.cfg.Arbeitnow, ch), nil
	}

	return nil, fmt.Errorf("unsupported integration: %s", ch.Integration)
}
