package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
)

type PubSubService interface {
	PublishImportCommand(ctx context.Context, importID uuid.UUID) error
}

type ScheduleImportsAction struct {
	chs *channel.Service
	is  *imports.Service
	log *slog.Logger

	ps PubSubService
}

func NewScheduleImportsAction(chs *channel.Service, is *imports.Service, ps PubSubService, log *slog.Logger) *ScheduleImportsAction {
	return &ScheduleImportsAction{
		chs: chs,
		is:  is,
		ps:  ps,
		log: log,
	}
}

func (s *ScheduleImportsAction) Execute(ctx context.Context) error {
	channels, err := s.chs.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, ch := range channels {
		s.log.Info(fmt.Sprintf("scheduling import for channel %s [%s] [name: %s]", ch.ID(), ch.Integration().String(), ch.Name()))

		i, err := s.is.Start(ctx, uuid.New(), ch.ID())
		if err != nil {
			return fmt.Errorf("failed to start import for channel %s: %w", ch.ID(), err)
		}

		if err := s.ps.PublishImportCommand(ctx, i.ID()); err != nil {
			return fmt.Errorf("failed to publish import command for import %s: %w", i.ID(), err)
		}
	}

	return nil
}
