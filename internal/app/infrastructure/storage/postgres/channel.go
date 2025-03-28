package postgres

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/google/uuid"
)

type Channel struct {
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"name"`
	Integration int       `db:"integration"`
	Status      int       `db:"status"`
	ID          uuid.UUID `db:"id"`
}

func fromDomainChannel(ch *channel.Channel) *Channel {
	return &Channel{
		ID:          ch.ID(),
		Name:        ch.Name(),
		Integration: int(ch.Integration()),
		Status:      int(ch.Status()),
		CreatedAt:   ch.CreatedAt(),
		UpdatedAt:   ch.UpdatedAt(),
	}
}

func toDomainChannel(ch *Channel) *channel.Channel {
	return channel.New(
		ch.ID,
		ch.Name,
		channel.Integration(ch.Integration),
		channel.Status(ch.Status),
		channel.WithTimestamps(ch.CreatedAt, ch.UpdatedAt),
	)
}
