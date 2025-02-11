package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aviseu/jobs/internal/app/errs"
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

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleFail(w, fmt.Errorf("failed to decode post: %w", err), http.StatusBadRequest)
		return
	}

	cmd := channel.NewCreateCommand(req.Name, req.Integration)
	ch, err := h.s.Create(r.Context(), cmd)
	if err != nil {
		if errs.IsValidationError(err) {
			h.handleFail(w, err, http.StatusBadRequest)
			return
		}

		h.handleError(w, fmt.Errorf("failed to create channel: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := NewChannelResponse(ch)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *Handler) handleFail(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := NewErrorResponse(err)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(err.Error(), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	h.log.Error(err.Error(), slog.Any("error", err))

	h.handleFail(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
}
