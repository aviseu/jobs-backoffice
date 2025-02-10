package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aviseu/jobs/internal/app/http/api"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr            string        `default:":8080"`
	ShutdownTimeout time.Duration `default:"5s"`
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

func APIRootHandler() http.Handler {
	r := chi.NewRouter()

	h := api.NewHandler()
	r.Mount("/api", h.Routes())

	return r
}
