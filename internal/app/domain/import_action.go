package domain

import (
	"context"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/gateway"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
)

type ImportAction struct {
	chs *channel.Service
	is  *imports.Service
	f   *gateway.Factory
}

func NewImportAction(chs *channel.Service, is *imports.Service, f *gateway.Factory) *ImportAction {
	return &ImportAction{
		chs: chs,
		is:  is,
		f:   f,
	}
}

func (s *ImportAction) Execute(ctx context.Context, iID uuid.UUID) error {
	i, err := s.is.FindImport(ctx, iID)
	if err != nil {
		return fmt.Errorf("failed to find import %s: %w", iID, err)
	}

	ch, err := s.chs.Find(ctx, i.ChannelID())
	if err != nil {
		return fmt.Errorf("failed to find channel %s: %w", i.ChannelID(), err)
	}

	g := s.f.Create(ch)
	if err := g.Import(ctx, i); err != nil {
		return fmt.Errorf("failed to import channel %s: %w", ch.ID(), err)
	}

	return nil
}
