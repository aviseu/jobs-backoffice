package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/application/http/api"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http/importh"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/scheduling"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Config struct {
	Addr            string        `env:"ADDR" envDefault:":8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	Cors            bool          `env:"CORS" envDefault:"false"`
	MaxConnections  int           `env:"MAX_CONNECTIONS" envDefault:"0"`
}

func SetupServer(ctx context.Context, cfg Config, h http.Handler) http.Server {
	return http.Server{
		Addr:    cfg.Addr,
		Handler: h,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
}

func APIRootHandler(chs *configuring.Service, chr api.ChannelRepository, ir api.ImportRepository, is *scheduling.Service, cfg Config, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	if cfg.Cors {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}

	r.Mount("/api/channels", api.NewChannelHandler(chs, chr, is, log).Routes())
	r.Mount("/api/integrations", api.NewIntegrationHandler(chs, log).Routes())
	r.Mount("/api/imports", api.NewImportHandler(chr, ir, log).Routes())

	return r
}

func ImportRootHandler(is *importing.Service, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Mount("/import", importh.NewHandler(is, log).Routes())

	return r
}
