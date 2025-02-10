package postgres

import (
	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/google/uuid"
	"time"
)

type Job struct {
	ID          uuid.UUID `db:"id"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Source      string    `db:"source"`
	Location    string    `db:"location"`
	Remote      bool      `db:"remote"`
	PostedAt    time.Time `db:"posted_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func FromDomainJob(job *job.Job) *Job {
	return &Job{
		ID:          job.ID(),
		URL:         job.URL(),
		Title:       job.Title(),
		Description: job.Description(),
		Source:      job.Source(),
		Location:    job.Location(),
		Remote:      job.Remote(),
		PostedAt:    job.PostedAt(),
		CreatedAt:   job.CreatedAt(),
		UpdatedAt:   job.UpdatedAt(),
	}
}
