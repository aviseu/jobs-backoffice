package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/errs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ChannelHandler struct {
	s   *channel.Service
	log *slog.Logger
}

func NewChannelHandler(s *channel.Service, log *slog.Logger) *ChannelHandler {
	return &ChannelHandler{
		s:   s,
		log: log,
	}
}

func (h *ChannelHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.ListChannels)
	r.Get("/{id}", h.FindChannel)
	r.Post("/", h.CreateChannel)
	r.Patch("/{id}", h.UpdateChannel)
	r.Put("/{id}/activate", h.ActivateChannel)
	r.Put("/{id}/deactivate", h.DeactivateChannel)

	r.Get("/integrations", h.ListIntegrations)

	return r
}

func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
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

func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.s.All(r.Context())
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to get channels: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := NewListChannelsResponse(channels)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ChannelHandler) FindChannel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("failed to parse post uuid %s: %w", idStr, err), http.StatusBadRequest)
		return
	}

	p, err := h.s.Find(r.Context(), id)
	if err != nil {
		if errors.Is(err, channel.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to find post %s: %w", idStr, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := NewChannelResponse(p)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode post %s: %w", idStr, err))
		return
	}
}

func (h *ChannelHandler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("failed to parse post uuid %s: %w", idStr, err), http.StatusBadRequest)
		return
	}

	var req UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleFail(w, fmt.Errorf("failed to decode post: %w", err), http.StatusBadRequest)
		return
	}

	cmd := channel.NewUpdateCommand(id, req.Name)
	ch, err := h.s.Update(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, channel.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		if errs.IsValidationError(err) {
			h.handleFail(w, err, http.StatusBadRequest)
			return
		}

		h.handleError(w, fmt.Errorf("failed to update post %s: %w", idStr, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := NewChannelResponse(ch)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode post %s: %w", idStr, err))
		return
	}
}

func (h *ChannelHandler) ActivateChannel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("failed to parse post uuid %s: %w", idStr, err), http.StatusBadRequest)
		return
	}

	if err := h.s.Activate(r.Context(), id); err != nil {
		if errors.Is(err, channel.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to activate post %s: %w", idStr, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChannelHandler) DeactivateChannel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("failed to parse post uuid %s: %w", idStr, err), http.StatusBadRequest)
		return
	}

	if err := h.s.Deactivate(r.Context(), id); err != nil {
		if errors.Is(err, channel.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to deactivate post %s: %w", idStr, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChannelHandler) ListIntegrations(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := NewListIntegrationsResponse(h.s.Integrations())
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ChannelHandler) handleFail(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := NewErrorResponse(err)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(err.Error(), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ChannelHandler) handleError(w http.ResponseWriter, err error) {
	h.log.Error(err.Error(), slog.Any("error", err))

	h.handleFail(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
}
