package gateway

import (
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/app/gateway/arbeitnow"
)

type Config struct {
	Arbeitnow arbeitnow.Config

	Import struct {
		ResultBufferSize int `split_words:"true" default:"10"`
		ResultWorkers    int `split_words:"true" default:"10"`
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Factory struct {
	c   HTTPClient
	js  *job.Service
	is  *imports.Service
	log *slog.Logger
	cfg Config
}

func NewFactory(js *job.Service, is *imports.Service, c HTTPClient, cfg Config, log *slog.Logger) *Factory {
	return &Factory{
		cfg: cfg,
		js:  js,
		is:  is,
		c:   c,
		log: log,
	}
}

func (f *Factory) Create(ch *channel.Channel) *Gateway {
	var p Provider
	if ch.Integration() == channel.IntegrationArbeitnow {
		p = arbeitnow.NewService(f.c, f.cfg.Arbeitnow, ch)
	}

	return NewGateway(p, f.js, f.is, f.log, f.cfg.Import.ResultBufferSize, f.cfg.Import.ResultWorkers)
}
