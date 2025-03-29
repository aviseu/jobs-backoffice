package configuring

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"

	"github.com/google/uuid"
)

type Repository interface {
	All(ctx context.Context) ([]*postgres.Channel, error)
	GetActive(ctx context.Context) ([]*postgres.Channel, error)
	Find(ctx context.Context, id uuid.UUID) (*postgres.Channel, error)
	Save(context.Context, *postgres.Channel) error
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, cmd *CreateCommand) (*Channel, error) {
	var errs error

	i, ok := base.ParseIntegration(cmd.Integration)
	if !ok {
		errs = errors.Join(errs, fmt.Errorf("failed to find integration %s: %w", cmd.Integration, ErrInvalidIntegration))
	}

	if cmd.Name == "" {
		errs = errors.Join(errs, ErrNameIsRequired)
	}

	if errs != nil {
		return nil, errs
	}

	ch := New(
		uuid.New(),
		cmd.Name,
		i,
		base.ChannelStatusInactive,
	)

	if err := s.r.Save(ctx, ch.DTO()); err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return ch, nil
}

func (s *Service) All(ctx context.Context) ([]*Channel, error) {
	dtos, err := s.r.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}

	channels := make([]*Channel, 0, len(dtos))
	for _, dto := range dtos {
		ch := fromDTO(dto)
		channels = append(channels, ch)
	}

	return channels, nil
}

func (s *Service) GetActive(ctx context.Context) ([]*Channel, error) {
	dtos, err := s.r.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active channels: %w", err)
	}

	channels := make([]*Channel, 0, len(dtos))
	for _, dto := range dtos {
		ch := fromDTO(dto)
		channels = append(channels, ch)
	}

	return channels, nil
}

func (s *Service) Find(ctx context.Context, id uuid.UUID) (*Channel, error) {
	dto, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrChannelNotFound) {
			return nil, ErrChannelNotFound
		}
		return nil, fmt.Errorf("failed to find channel: %w", err)
	}

	return fromDTO(dto), nil
}

func (s *Service) Update(ctx context.Context, cmd *UpdateCommand) (*Channel, error) {
	dto, err := s.r.Find(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrChannelNotFound) {
			return nil, ErrChannelNotFound
		}
		return nil, fmt.Errorf("failed to find channel: %w", err)
	}

	ch := fromDTO(dto)

	if err := ch.Update(cmd.Name); err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	if err := s.r.Save(ctx, ch.DTO()); err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	return ch, nil
}

func (s *Service) Activate(ctx context.Context, id uuid.UUID) error {
	dto, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrChannelNotFound) {
			return ErrChannelNotFound
		}
		return fmt.Errorf("failed to find channel: %w", err)
	}

	ch := fromDTO(dto)

	ch.Activate()

	if err := s.r.Save(ctx, ch.DTO()); err != nil {
		return fmt.Errorf("failed to activate channel: %w", err)
	}

	return nil
}

func (s *Service) Deactivate(ctx context.Context, id uuid.UUID) error {
	dto, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrChannelNotFound) {
			return ErrChannelNotFound
		}
		return fmt.Errorf("failed to find channel: %w", err)
	}

	ch := fromDTO(dto)

	ch.Deactivate()

	if err := s.r.Save(ctx, ch.DTO()); err != nil {
		return fmt.Errorf("failed to activate channel: %w", err)
	}

	return nil
}

func (*Service) Integrations() []base.Integration {
	ii := make([]base.Integration, 0, len(base.Integrations))
	for i := range base.Integrations {
		ii = append(ii, i)
	}

	return ii
}
