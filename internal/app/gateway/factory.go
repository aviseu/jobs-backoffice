package gateway

import (
	"net/http"

	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/aviseu/jobs/internal/app/gateway/arbeitnow"
)

type Config struct {
	Arbeitnow arbeitnow.Config
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Factory struct {
	c   HTTPClient
	s   *job.Service
	cfg Config
}

func NewFactory(s *job.Service, c HTTPClient, cfg Config) *Factory {
	return &Factory{
		cfg: cfg,
		s:   s,
		c:   c,
	}
}

func (f *Factory) Create(ch *channel.Channel) *Gateway {
	var p Provider
	if ch.Integration() == channel.IntegrationArbeitnow {
		p = arbeitnow.NewService(f.c, f.cfg.Arbeitnow, ch)
	}

	return NewGateway(p, f.s)
}
