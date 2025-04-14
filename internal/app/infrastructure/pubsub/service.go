package pubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/aviseu/jobs-protobuf/build/gen/commands/jobs"
	"github.com/google/uuid"
)

type Config struct {
	Timeout time.Duration `default:"1s"`
}

type Service struct {
	importClient *client
}

func NewService(importTopic *pubsub.Topic, cfg Config) *Service {
	return &Service{
		importClient: newClient(importTopic, cfg.Timeout),
	}
}

func (s *Service) PublishImportCommand(ctx context.Context, importID uuid.UUID) error {
	msg := jobs.ExecuteImportChannel{
		ImportId: importID.String(),
	}

	return s.importClient.publish(ctx, &msg)
}
