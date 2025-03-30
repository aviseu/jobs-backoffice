package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/scheduling"
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/errs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ChannelRepository interface {
	All(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type ChannelHandler struct {
	chs *configuring.Service
	chr ChannelRepository
	is  *scheduling.Service
	log *slog.Logger
}

func NewChannelHandler(chs *configuring.Service, chr ChannelRepository, is *scheduling.Service, log *slog.Logger) *ChannelHandler {
	return &ChannelHandler{
		chs: chs,
		chr: chr,
		log: log,
		is:  is,
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

	r.Put("/{id}/schedule", h.ScheduleImport)

	return r
}

func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleFail(w, fmt.Errorf("failed to decode post: %w", err), http.StatusBadRequest)
		return
	}

	cmd := configuring.NewCreateChannelCommand(req.Name, req.Integration)
	ch, err := h.chs.Create(r.Context(), cmd)
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

	resp := NewChannelResponse(ch.ToAggregator())
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.chr.All(r.Context())
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

	p, err := h.chr.Find(r.Context(), id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrChannelNotFound) {
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

	cmd := configuring.NewUpdateChannelCommand(id, req.Name)
	ch, err := h.chs.Update(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, configuring.ErrChannelNotFound) {
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
	resp := NewChannelResponse(ch.ToAggregator())
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

	if err := h.chs.Activate(r.Context(), id); err != nil {
		if errors.Is(err, configuring.ErrChannelNotFound) {
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

	if err := h.chs.Deactivate(r.Context(), id); err != nil {
		if errors.Is(err, configuring.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to deactivate post %s: %w", idStr, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChannelHandler) ScheduleImport(w http.ResponseWriter, r *http.Request) {
	channelIDStr := chi.URLParam(r, "id")
	if channelIDStr == "" {
		h.handleFail(w, errors.New("missing channel id"), http.StatusBadRequest)
		return
	}

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("invalid channel id: %w", err), http.StatusBadRequest)
		return
	}

	ch, err := h.chr.Find(r.Context(), channelID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrChannelNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to find channel: %w", err))
		return
	}

	// temporary
	i, err := h.is.ScheduleImport(context.Background(), ch)
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to schedule import: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := NewImportResponse(i, ch)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ChannelHandler) handleFail(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := NewErrorResponse(err)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(err.Error(), slog.Any("Error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ChannelHandler) handleError(w http.ResponseWriter, err error) {
	h.log.Error(err.Error(), slog.Any("Error", err))

	h.handleFail(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
}
