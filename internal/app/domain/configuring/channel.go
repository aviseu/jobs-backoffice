package configuring

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type channel struct {
	createdAt   time.Time
	updatedAt   time.Time
	name        string
	integration aggregator.Integration
	status      aggregator.ChannelStatus
	id          uuid.UUID
}

type optional func(*channel)

func withTimestamps(c, u time.Time) optional {
	return func(ch *channel) {
		ch.createdAt = c
		ch.updatedAt = u
	}
}

func newChannel(id uuid.UUID, name string, i aggregator.Integration, s aggregator.ChannelStatus, opts ...optional) *channel {
	ch := &channel{
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

func (ch *channel) update(name string) error {
	if name == "" {
		return ErrNameIsRequired
	}

	ch.name = name
	ch.updatedAt = time.Now()

	return nil
}

func (ch *channel) activate() {
	ch.status = aggregator.ChannelStatusActive
	ch.updatedAt = time.Now()
}

func (ch *channel) deactivate() {
	ch.status = aggregator.ChannelStatusInactive
	ch.updatedAt = time.Now()
}

func (ch *channel) toAggregator() *aggregator.Channel {
	return &aggregator.Channel{
		ID:          ch.id,
		Name:        ch.name,
		Integration: ch.integration,
		Status:      ch.status,
		CreatedAt:   ch.createdAt,
		UpdatedAt:   ch.updatedAt,
	}
}

func newChannelFromAggregator(ch *aggregator.Channel) *channel {
	return newChannel(
		ch.ID,
		ch.Name,
		ch.Integration,
		ch.Status,
		withTimestamps(ch.CreatedAt, ch.UpdatedAt),
	)
}
