package postgres

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/google/uuid"
)

type Job struct {
	PostedAt      time.Time `db:"posted_at"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	URL           string    `db:"url"`
	Title         string    `db:"title"`
	Description   string    `db:"description"`
	Source        string    `db:"source"`
	Location      string    `db:"location"`
	ID            uuid.UUID `db:"id"`
	ChannelID     uuid.UUID `db:"channel_id"`
	Remote        bool      `db:"remote"`
	Status        int       `db:"status"`
	PublishStatus int       `db:"publish_status"`
}

func fromDomainJob(j *job.Job) *Job {
	return &Job{
		ID:            j.ID(),
		ChannelID:     j.ChannelID(),
		URL:           j.URL(),
		Title:         j.Title(),
		Description:   j.Description(),
		Source:        j.Source(),
		Location:      j.Location(),
		Remote:        j.Remote(),
		PostedAt:      j.PostedAt(),
		CreatedAt:     j.CreatedAt(),
		UpdatedAt:     j.UpdatedAt(),
		Status:        int(j.Status()),
		PublishStatus: int(j.PublishStatus()),
	}
}

func toDomainJob(j *Job) *job.Job {
	return job.New(
		j.ID,
		j.ChannelID,
		job.Status(j.Status),
		j.URL,
		j.Title,
		j.Description,
		j.Source,
		j.Location,
		j.Remote,
		j.PostedAt,
		job.WithTimestamps(j.CreatedAt, j.UpdatedAt),
		job.WithPublishStatus(job.PublishStatus(j.PublishStatus)),
	)
}
