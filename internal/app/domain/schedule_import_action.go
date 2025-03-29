package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
)

type ScheduleImportAction struct {
	is  *imports.Service
	log *slog.Logger

	ps PubSubService
}

func NewScheduleImportAction(is *imports.Service, ps PubSubService, log *slog.Logger) *ScheduleImportAction {
	return &ScheduleImportAction{
		is:  is,
		ps:  ps,
		log: log,
	}
}

func (s *ScheduleImportAction) Execute(ctx context.Context, ch *configuring.Channel) (*imports.Import, error) {
	s.log.Info(fmt.Sprintf("scheduling import for channel %s [%s] [name: %s]", ch.ID(), ch.Integration().String(), ch.Name()))

	i, err := s.is.Start(ctx, uuid.New(), ch.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to start import for channel %s: %w", ch.ID(), err)
	}

	if err := s.ps.PublishImportCommand(ctx, i.ID()); err != nil {
		return nil, fmt.Errorf("failed to publish import command for import %s: %w", i.ID(), err)
	}

	return i, nil
}
