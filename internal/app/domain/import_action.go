package domain

import (
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/google/uuid"
)

type ChannelRepository interface {
	GetActive(ctx context.Context) ([]*postgres.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*postgres.Channel, error)
}

type ImportAction struct {
	chr ChannelRepository
	is  *importing.ImportService
	f   *importing.Factory
}

func NewImportAction(chr ChannelRepository, is *importing.ImportService, f *importing.Factory) *ImportAction {
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
