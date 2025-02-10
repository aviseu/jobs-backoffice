package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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
