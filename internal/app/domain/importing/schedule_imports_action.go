package importing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PubSubService interface {
	PublishImportCommand(ctx context.Context, importID uuid.UUID) error
}

type ScheduleImportsAction struct {
	s   *Service
	chr ChannelRepository
}

func NewScheduleImportsAction(chr ChannelRepository, s *Service) *ScheduleImportsAction {
	return &ScheduleImportsAction{
		chr: chr,
		s:   s,
	}
}

func (s *ScheduleImportsAction) Execute(ctx context.Context) error {
	channels, err := s.chr.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, ch := range channels {
		if _, err := s.s.ScheduleImport(ctx, ch); err != nil {
			return fmt.Errorf("failed to schedule import for channel %s: %w", ch.ID, err)
		}
	}

	return nil
}
