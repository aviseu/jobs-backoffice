package importing

import (
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type Service struct {
	ir  ImportRepository
	chr ChannelRepository
	is  *ImportService
	f   *Factory
}

func NewService(chr ChannelRepository, ir ImportRepository, is *ImportService, f *Factory) *Service {
	return &Service{
		chr: chr,
		is:  is,
		ir:  ir,
		f:   f,
	}
}

func (s *Service) Import(ctx context.Context, importID uuid.UUID) error {
	i, err := s.ir.FindImport(ctx, importID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", importID, err)
	}

	ch, err := s.chr.Find(ctx, i.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.ChannelID, err)
	}

	g := s.f.Create(ch)
	if err := g.Import(ctx, i); err != nil {
		return fmt.Errorf("failed to import channel %s: %w", ch.ID, err)
	}

	return nil
}
