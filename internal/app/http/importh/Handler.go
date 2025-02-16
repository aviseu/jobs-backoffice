package importh

import (
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/gateway"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	chs *channel.Service
	is  *imports.Service
	f   *gateway.Factory
	log *slog.Logger
}

func NewHandler(chs *channel.Service, is *imports.Service, f *gateway.Factory, log *slog.Logger) *Handler {
	return &Handler{
		chs: chs,
		is:  is,
		f:   f,
		log: log,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.Import)

	return r
}

func (h *Handler) Import(w http.ResponseWriter, _ *http.Request) {
	h.log.Info("received message!")
	w.WriteHeader(http.StatusOK)
}
