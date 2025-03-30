package importing

import (
	"fmt"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
)

type Provider interface {
	GetJobs() ([]*aggregator.Job, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Factory struct {
	c   HTTPClient
	cfg Config
}

func NewFactory(c HTTPClient, cfg Config) *Factory {
	return &Factory{
		cfg: cfg,
		c:   c,
	}
}

func (f *Factory) Create(ch *aggregator.Channel) (Provider, error) {
	if ch.Integration == aggregator.IntegrationArbeitnow {
		return arbeitnow.NewService(f.c, f.cfg.Arbeitnow, ch), nil
	}

	return nil, fmt.Errorf("unsupported integration: %s", ch.Integration)
}
