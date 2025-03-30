package configuring

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type Channel struct {
	createdAt   time.Time
	updatedAt   time.Time
	name        string
	integration aggregator.Integration
	status      aggregator.ChannelStatus
	id          uuid.UUID
}

type Optional func(*Channel)

func WithTimestamps(c, u time.Time) Optional {
	return func(ch *Channel) {
		ch.createdAt = c
		ch.updatedAt = u
	}
}

func NewChannel(id uuid.UUID, name string, i aggregator.Integration, s aggregator.ChannelStatus, opts ...Optional) *Channel {
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

func (ch *Channel) Integration() aggregator.Integration {
	return ch.integration
}

func (ch *Channel) Status() aggregator.ChannelStatus {
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
	ch.status = aggregator.ChannelStatusActive
	ch.updatedAt = time.Now()
}

func (ch *Channel) Deactivate() {
	ch.status = aggregator.ChannelStatusInactive
	ch.updatedAt = time.Now()
}

func (ch *Channel) ToAggregator() *aggregator.Channel {
	return &aggregator.Channel{
		ID:          ch.id,
		Name:        ch.name,
		Integration: ch.integration,
		Status:      ch.status,
		CreatedAt:   ch.createdAt,
		UpdatedAt:   ch.updatedAt,
	}
}

func NewChannelFromAggregator(ch *aggregator.Channel) *Channel {
	return NewChannel(
		ch.ID,
		ch.Name,
		ch.Integration,
		ch.Status,
		WithTimestamps(ch.CreatedAt, ch.UpdatedAt),
	)
}
