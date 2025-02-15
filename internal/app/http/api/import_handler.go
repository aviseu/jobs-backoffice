package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImportHandler struct {
	is  *imports.Service
	chs *channel.Service
	log *slog.Logger
}

func NewImportHandler(chs *channel.Service, is *imports.Service, log *slog.Logger) *ImportHandler {
	return &ImportHandler{
		chs: chs,
		is:  is,
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
	ii, err := h.is.GetImports(r.Context())
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to get imports: %w", err))
		return
	}

	cc, err := h.chs.All(r.Context())
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

	i, err := h.is.FindImport(r.Context(), id)
	if err != nil {
		if errors.Is(err, imports.ErrImportNotFound) {
			h.handleFail(w, fmt.Errorf("import not found: %w", err), http.StatusNotFound)
			return
		}

		h.handleError(w, fmt.Errorf("failed to get import: %w", err))
		return
	}

	ch, err := h.chs.Find(r.Context(), i.ChannelID())
	if err != nil {
		h.handleError(w, fmt.Errorf("failed to get channels: %w", err))
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
