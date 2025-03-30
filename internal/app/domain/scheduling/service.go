package scheduling

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
}

type PubSubService interface {
	PublishImportCommand(ctx context.Context, importID uuid.UUID) error
}

type ImportRepository interface {
	SaveImport(ctx context.Context, i *aggregator.Import) error
}

type Service struct {
	chr ChannelRepository
	ir  ImportRepository
	ps  PubSubService
	log *slog.Logger
}

func NewService(ir ImportRepository, chr ChannelRepository, ps PubSubService, log *slog.Logger) *Service {
	return &Service{
		ir:  ir,
		chr: chr,
		ps:  ps,
		log: log,
	}
}

func (s *Service) ScheduleActiveChannels(ctx context.Context) error {
	channels, err := s.chr.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, ch := range channels {
		if _, err := s.ScheduleImport(ctx, ch); err != nil {
			return fmt.Errorf("failed to schedule import for channel %s: %w", ch.ID, err)
		}
	}

	return nil
}

func (s *Service) ScheduleImport(ctx context.Context, ch *aggregator.Channel) (*aggregator.Import, error) {
	s.log.Info(fmt.Sprintf("scheduling import for channel %s [%s] [name: %s]", ch.ID, ch.Integration.String(), ch.Name))

	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: ch.ID,
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Now(),
		EndedAt:   null.NewTime(time.Now(), false),
	}
	if err := s.ir.SaveImport(ctx, i); err != nil {
		return nil, fmt.Errorf("failed to save import for channel %s while starting: %w", ch.ID, err)
	}

	if err := s.ps.PublishImportCommand(ctx, i.ID); err != nil {
		return nil, fmt.Errorf("failed to publish import %s for channel %s: %w", i.ID, ch.ID, err)
	}

	return i, nil
}
