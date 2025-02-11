package channel

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusInactive Status = iota
	StatusActive
)

func (s Status) String() string {
	return [...]string{"inactive", "active"}[s]
}

type Integration int

const (
	IntegrationArbeitnow Integration = iota
)

var integrations = map[Integration]string{
	IntegrationArbeitnow: "arbeitnow",
}

func (i Integration) String() string {
	return integrations[i]
}

func ParseIntegration(s string) (Integration, bool) {
	for _, i := range integrations {
		if i == s {
			return IntegrationArbeitnow, true
		}
	}

	return -1, false
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

func (ch *Channel) Update(name string) error {
	if name == "" {
		return ErrNameIsRequired
	}

	ch.name = name
	ch.updatedAt = time.Now()

	return nil
}
