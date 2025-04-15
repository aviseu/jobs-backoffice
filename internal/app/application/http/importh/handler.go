package importh

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-protobuf/build/gen/commands/jobs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	is  *importing.Service
	log *slog.Logger
}

func NewHandler(is *importing.Service, log *slog.Logger) *Handler {
	return &Handler{
		is:  is,
		log: log,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.Import)

	return r
}

type pubSubMessage struct {
	Subscription string `json:"subscription"`
	Message      struct {
		ID   string `json:"id"`
		Data []byte `json:"data,omitempty"`
	} `json:"message"`
}

func (h *Handler) Import(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received message!")

	var msg pubSubMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		h.log.Error(fmt.Errorf("failed to json decode request body: %w", err).Error())
		http.Error(w, "skipped message", http.StatusOK) // 200 will ack message
		return
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			h.log.Error(fmt.Errorf("failed to close body: %w", err).Error())
		}
	}(r.Body)

	var data jobs.ExecuteImportChannel
	if err := proto.Unmarshal(msg.Message.Data, &data); err != nil {
		h.log.Error(fmt.Errorf("failed to unmarshal pubsub message: %w", err).Error())
		http.Error(w, "skipped message", http.StatusOK) // 200 will ack message
		return
	}

	importID, err := uuid.Parse(data.ImportId)
	if err != nil {
		h.log.Error(fmt.Errorf("failed to convert import id %s to uuid: %w", data.ImportId, err).Error())
		http.Error(w, "skipped message", http.StatusOK) // 200 will ack message
		return
	}

	h.log.Info("processing import " + importID.String())
	if err := h.is.Import(r.Context(), importID); err != nil {
		h.log.Error(fmt.Errorf("failed to execute import %s: %w", importID, err).Error())
		http.Error(w, "skipped message", http.StatusOK) // 200 will ack message
		return
	}
	h.log.Info("completed import " + importID.String())

	w.WriteHeader(http.StatusOK)
}
