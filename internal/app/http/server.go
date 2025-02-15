package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/http/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Config struct {
	Addr            string        `default:":8080"`
	ShutdownTimeout time.Duration `default:"5s"`
	Cors            bool          `default:"false"`
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

func APIRootHandler(chs *channel.Service, is *imports.Service, cfg Config, log *slog.Logger) http.Handler {
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

	r.Mount("/api/channels", api.NewChannelHandler(chs, log).Routes())
	r.Mount("/api/integrations", api.NewIntegrationHandler(chs, log).Routes())
	r.Mount("/api/imports", api.NewImportHandler(chs, is, log).Routes())

	return r
}
