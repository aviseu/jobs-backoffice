package postgres

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	PostedAt      time.Time             `db:"posted_at"`
	CreatedAt     time.Time             `db:"created_at"`
	UpdatedAt     time.Time             `db:"updated_at"`
	URL           string                `db:"url"`
	Title         string                `db:"title"`
	Description   string                `db:"description"`
	Source        string                `db:"source"`
	Location      string                `db:"location"`
	ID            uuid.UUID             `db:"id"`
	ChannelID     uuid.UUID             `db:"channel_id"`
	Remote        bool                  `db:"remote"`
	Status        base.JobStatus        `db:"status"`
	PublishStatus base.JobPublishStatus `db:"publish_status"`
}
