package configuring

import (
	"context"
	"errors"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type Repository interface {
	Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error)
	Save(context.Context, *aggregator.Channel) error
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, cmd *CreateChannelCommand) (*aggregator.Channel, error) {
	var errs error

	i, ok := aggregator.ParseIntegration(cmd.Integration)
	if !ok {
		errs = errors.Join(errs, fmt.Errorf("failed to find integration %s: %w", cmd.Integration, ErrInvalidIntegration))
	}

	if cmd.Name == "" {
		errs = errors.Join(errs, ErrNameIsRequired)
	}

	if errs != nil {
		return nil, errs
	}

	ch := newChannel(
		uuid.New(),
		cmd.Name,
		i,
		aggregator.ChannelStatusInactive,
	)

	if err := s.r.Save(ctx, ch.toAggregator()); err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return ch.toAggregator(), nil
}

func (s *Service) Update(ctx context.Context, cmd *UpdateChannelCommand) (*aggregator.Channel, error) {
	aggr, err := s.r.Find(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrChannelNotFound) {
			return nil, ErrChannelNotFound
		}
		return nil, fmt.Errorf("failed to find channel: %w", err)
	}

	ch := newChannelFromAggregator(aggr)

	if err := ch.update(cmd.Name); err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	if err := s.r.Save(ctx, ch.toAggregator()); err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	return ch.toAggregator(), nil
}

func (s *Service) Activate(ctx context.Context, id uuid.UUID) error {
	aggr, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrChannelNotFound) {
			return ErrChannelNotFound
		}
		return fmt.Errorf("failed to find channel: %w", err)
	}

	ch := newChannelFromAggregator(aggr)

	ch.activate()

	if err := s.r.Save(ctx, ch.toAggregator()); err != nil {
		return fmt.Errorf("failed to activate channel: %w", err)
	}

	return nil
}

func (s *Service) Deactivate(ctx context.Context, id uuid.UUID) error {
	aggr, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrChannelNotFound) {
			return ErrChannelNotFound
		}
		return fmt.Errorf("failed to find channel: %w", err)
	}

	ch := newChannelFromAggregator(aggr)

	ch.deactivate()

	if err := s.r.Save(ctx, ch.toAggregator()); err != nil {
		return fmt.Errorf("failed to activate channel: %w", err)
	}

	return nil
}
