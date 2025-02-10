package channel

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusInactive = iota
	StatusActive
)

func (s Status) String() string {
	return [...]string{"inactive", "active"}[s]
}

type Integration int

const (
	IntegrationArbeitnow = iota
)

func (i Integration) String() string {
	return [...]string{"arbeitnow"}[i]
}

type Channel struct {
	createdAt   time.Time
	updatedAt   time.Time
	name        string
	integration Integration
	status      Status
	id          uuid.UUID
}

type Optional func(*Channel)

func WithTimestamps(c, u time.Time) Optional {
	return func(ch *Channel) {
		ch.createdAt = c
		ch.updatedAt = u
	}
}

func New(id uuid.UUID, name string, i Integration, s Status, opts ...Optional) *Channel {
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

func (ch *Channel) Integration() Integration {
	return ch.integration
}

func (ch *Channel) Status() Status {
	return ch.status
}

func (ch *Channel) CreatedAt() time.Time {
	return ch.createdAt
}

func (ch *Channel) UpdatedAt() time.Time {
	return ch.updatedAt
}
