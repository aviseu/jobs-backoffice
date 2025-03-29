package channel

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	createdAt   time.Time
	updatedAt   time.Time
	name        string
	integration base.Integration
	status      base.ChannelStatus
	id          uuid.UUID
}

type Optional func(*Channel)

func WithTimestamps(c, u time.Time) Optional {
	return func(ch *Channel) {
		ch.createdAt = c
		ch.updatedAt = u
	}
}

func New(id uuid.UUID, name string, i base.Integration, s base.ChannelStatus, opts ...Optional) *Channel {
	ch := &Channel{
		id:          id,
		name:        name,
		integration: i,
		status:      s,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	for _, opt := range opts {
		opt(ch)
	}

	return ch
}

func (ch *Channel) ID() uuid.UUID {
	return ch.id
}

func (ch *Channel) Name() string {
	return ch.name
}

func (ch *Channel) Integration() base.Integration {
	return ch.integration
}

func (ch *Channel) Status() base.ChannelStatus {
	return ch.status
}

func (ch *Channel) CreatedAt() time.Time {
	return ch.createdAt
}

func (ch *Channel) UpdatedAt() time.Time {
	return ch.updatedAt
}

func (ch *Channel) Update(name string) error {
	if name == "" {
		return ErrNameIsRequired
	}

	ch.name = name
	ch.updatedAt = time.Now()

	return nil
}

func (ch *Channel) Activate() {
	ch.status = base.ChannelStatusActive
	ch.updatedAt = time.Now()
}

func (ch *Channel) Deactivate() {
	ch.status = base.ChannelStatusInactive
	ch.updatedAt = time.Now()
}

func (ch *Channel) DTO() *postgres.Channel {
	return &postgres.Channel{
		ID:          ch.id,
		Name:        ch.name,
		Integration: int(ch.integration),
		Status:      int(ch.status),
		CreatedAt:   ch.createdAt,
		UpdatedAt:   ch.updatedAt,
	}
}

func fromDTO(dto *postgres.Channel) *Channel {
	return New(
		dto.ID,
		dto.Name,
		base.Integration(dto.Integration),
		base.ChannelStatus(dto.Status),
		WithTimestamps(dto.CreatedAt, dto.UpdatedAt),
	)
}
