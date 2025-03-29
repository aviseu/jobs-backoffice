package domain

import (
	"context"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/google/uuid"
)

type PubSubService interface {
	PublishImportCommand(ctx context.Context, importID uuid.UUID) error
}

type ScheduleImportsAction struct {
	ia  *ScheduleImportAction
	chs *configuring.Service
}

func NewScheduleImportsAction(chs *configuring.Service, ia *ScheduleImportAction) *ScheduleImportsAction {
	return &ScheduleImportsAction{
		chs: chs,
		ia:  ia,
	}
}

func (s *ScheduleImportsAction) Execute(ctx context.Context) error {
	channels, err := s.chs.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, ch := range channels {
		if _, err := s.ia.Execute(ctx, ch); err != nil {
			return fmt.Errorf("failed to schedule import for channel %s: %w", ch.ID(), err)
		}
	}

	return nil
}
