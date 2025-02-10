package channel

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Repository interface {
	Save(context.Context, *Channel) error
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, cmd *CreateCommand) (*Channel, error) {
	i, ok := ParseIntegration(cmd.Integration)
	if !ok {
		return nil, fmt.Errorf("failed to find integration %s: %w", cmd.Integration, ErrInvalidIntegration)
	}

	if cmd.Name == "" {
		return nil, ErrNameIsRequired
	}

	ch := New(
		uuid.New(),
		cmd.Name,
		i,
		StatusInactive,
	)

	if err := s.r.Save(ctx, ch); err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return ch, nil
}
