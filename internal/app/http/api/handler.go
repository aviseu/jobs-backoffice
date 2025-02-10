package api

import (
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	s   *channel.Service
	log *slog.Logger
}

func NewHandler(s *channel.Service, log *slog.Logger) *Handler {
	return &Handler{
		s:   s,
		log: log,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/channels", h.CreateChannel)

	return r
}

func (*Handler) CreateChannel(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("CreateChannel"))
	if err != nil {
		return
	}
}
