package aggregator

import (
	"time"

	"github.com/google/uuid"
)

type ChannelStatus int

const (
	ChannelStatusInactive ChannelStatus = iota
	ChannelStatusActive
)

func (s ChannelStatus) String() string {
	return [...]string{"inactive", "active"}[s]
}

type Channel struct {
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	Name        string        `db:"name"`
	Integration Integration   `db:"integration"`
	Status      ChannelStatus `db:"status"`
	ID          uuid.UUID     `db:"id"`
}
