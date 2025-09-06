package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/aviseu/jobs-protobuf/build/gen/go/commands/imports"
	"github.com/google/uuid"
)

type ImportService struct {
	client *client
}

func NewImportService(topic *pubsub.Topic, cfg Config) *ImportService {
	return &ImportService{
		client: newClient(topic, cfg.Timeout),
	}
}

func (s *ImportService) PublishImportCommand(ctx context.Context, importID uuid.UUID) error {
	msg := imports.ExecuteImportChannel{
		ImportId: importID.String(),
	}

	return s.client.publish(ctx, &msg)
}
