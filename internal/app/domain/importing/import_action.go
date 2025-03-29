package importing

import (
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*aggregator.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
}

type ImportAction struct {
	chr ChannelRepository
	is  *ImportService
	f   *Factory
}

func NewImportAction(chr ChannelRepository, is *ImportService, f *Factory) *ImportAction {
	return &ImportAction{
		chr: chr,
		is:  is,
		f:   f,
	}
}

func (s *ImportAction) Execute(ctx context.Context, iID uuid.UUID) error {
	i, err := s.is.FindImport(ctx, iID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", iID, err)
	}

	dto, err := s.chr.Find(ctx, i.ChannelID())
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.ChannelID(), err)
	}

	ch := configuring.NewChannelFromDTO(dto)

	g := s.f.Create(ch.ToDTO())
	if err := g.Import(ctx, i); err != nil {
		return fmt.Errorf("failed to import channel %s: %w", ch.ID(), err)
	}

	return nil
}
