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
	chr ChannelRepository
}

func NewScheduleImportsAction(chr ChannelRepository, ia *ScheduleImportAction) *ScheduleImportsAction {
	return &ScheduleImportsAction{
		chr: chr,
		ia:  ia,
	}
}

func (s *ScheduleImportsAction) Execute(ctx context.Context) error {
	channels, err := s.chr.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, dto := range channels {
		ch := configuring.NewChannelFromDTO(dto)
		if _, err := s.ia.Execute(ctx, ch); err != nil {
			return fmt.Errorf("failed to schedule import for channel %s: %w", ch.ID(), err)
		}
	}

	return nil
}
