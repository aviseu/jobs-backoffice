package channel

import (
	"context"
	"errors"
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
	var errs error

	i, ok := ParseIntegration(cmd.Integration)
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
		StatusInactive,
	)

	if err := s.r.Save(ctx, ch); err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return ch, nil
}
