package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImportRepository interface {
	GetImports(ctx context.Context) ([]*aggregator.Import, error)
	FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error)
}

type ImportHandler struct {
	ir  ImportRepository
	chr ChannelRepository
	log *slog.Logger
}

func NewImportHandler(chr ChannelRepository, ir ImportRepository, log *slog.Logger) *ImportHandler {
	return &ImportHandler{
		chr: chr,
		ir:  ir,
		log: log,
	}
}

func (h *ImportHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.ListImports)
	r.Get("/{id}", h.FindImport)

	return r
}

func (h *ImportHandler) ListImports(w http.ResponseWriter, r *http.Request) {
	ii, err := h.ir.GetImports(r.Context())
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to get imports: %w", err))
		return
	}

	cc, err := h.chr.All(r.Context())
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to get channels: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := NewImportsResponse(ii, cc)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ImportHandler) FindImport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.handleFail(w, errors.New("missing import id"), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(w, fmt.Errorf("invalid import id: %w", err), http.StatusBadRequest)
		return
	}

	i, err := h.ir.FindImport(r.Context(), id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrImportNotFound) {
			h.handleFail(w, err, http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to find import %s: %w", idStr, err))
		return
	}

	ch, err := h.chr.Find(r.Context(), i.ChannelID)
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to find channel: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := NewImportResponse(i, ch)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(w, fmt.Errorf("failed to encode response: %w", err))
	}
}

func (h *ImportHandler) handleFail(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := NewErrorResponse(err)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(err.Error(), slog.Any("Error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ImportHandler) handleError(w http.ResponseWriter, err error) {
	h.log.Error(err.Error(), slog.Any("Error", err))

	h.handleFail(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
}
