package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/gateway"
)

type ImportActiveAction struct {
	chs *channel.Service
	f   *gateway.Factory
	log *slog.Logger
}

func NewImportActiveAction(chs *channel.Service, f *gateway.Factory, log *slog.Logger) *ImportActiveAction {
	return &ImportActiveAction{
		chs: chs,
		f:   f,
		log: log,
	}
}

func (s *ImportActiveAction) Execute(ctx context.Context) error {
	channels, err := s.chs.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active channels: %w", err)
	}

	for _, ch := range channels {
		s.log.Info(fmt.Sprintf("importing channel %s [%s] [name: %s]", ch.ID(), ch.Integration().String(), ch.Name()))

		g := s.f.Create(ch)
		if err := g.ImportChannel(ctx); err != nil {
			return fmt.Errorf("failed to import channel %s: %w", ch.ID(), err)
		}
	}

	return nil
}
