package postgres

import (
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Integration int       `db:"integration"`
	Status      int       `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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
