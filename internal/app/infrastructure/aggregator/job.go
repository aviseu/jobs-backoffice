package aggregator

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus int

const (
	JobStatusInactive JobStatus = iota
	JobStatusActive
)

func (s JobStatus) String() string {
	return [...]string{"inactive", "active"}[s]
}

type JobPublishStatus int

const (
	JobPublishStatusUnpublished JobPublishStatus = iota
	JobPublishStatusPublished
)

func (s JobPublishStatus) String() string {
	return [...]string{"unpublished", "published"}[s]
}

type Job struct {
	PostedAt      time.Time        `db:"posted_at"`
	CreatedAt     time.Time        `db:"created_at"`
	UpdatedAt     time.Time        `db:"updated_at"`
	URL           string           `db:"url"`
	Title         string           `db:"title"`
	Description   string           `db:"description"`
	Source        string           `db:"source"`
	Location      string           `db:"location"`
	ID            uuid.UUID        `db:"id"`
	ChannelID     uuid.UUID        `db:"channel_id"`
	Remote        bool             `db:"remote"`
	Status        JobStatus        `db:"status"`
	PublishStatus JobPublishStatus `db:"publish_status"`
}
