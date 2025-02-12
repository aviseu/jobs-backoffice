package postgres

import (
	"time"

	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/google/uuid"
)

type Job struct {
	PostedAt    time.Time `db:"posted_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Source      string    `db:"source"`
	Location    string    `db:"location"`
	ID          uuid.UUID `db:"id"`
	Remote      bool      `db:"remote"`
	ChannelID   uuid.UUID `db:"channel_id"`
}

func fromDomainJob(j *job.Job) *Job {
	return &Job{
		ID:          j.ID(),
		URL:         j.URL(),
		Title:       j.Title(),
		Description: j.Description(),
		Source:      j.Source(),
		Location:    j.Location(),
		Remote:      j.Remote(),
		PostedAt:    j.PostedAt(),
		CreatedAt:   j.CreatedAt(),
		UpdatedAt:   j.UpdatedAt(),
		ChannelID:   j.ChannelID(),
	}
}
