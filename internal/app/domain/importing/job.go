package importing

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	postedAt      time.Time
	createdAt     time.Time
	updatedAt     time.Time
	url           string
	title         string
	description   string
	source        string
	location      string
	id            uuid.UUID
	remote        bool
	channelID     uuid.UUID
	status        base.JobStatus
	publishStatus base.JobPublishStatus
}

type JobOptional func(*Job)

func JobWithTimestamps(c, u time.Time) JobOptional {
	return func(j *Job) {
		j.createdAt = c
		j.updatedAt = u
	}
}

func JobWithPublishStatus(s base.JobPublishStatus) JobOptional {
	return func(j *Job) {
		j.publishStatus = s
	}
}

func NewJob(id, channelID uuid.UUID, s base.JobStatus, url, title, description, source, location string, remote bool, postedAt time.Time, opts ...JobOptional) *Job {
	j := &Job{
		id:            id,
		channelID:     channelID,
		status:        s,
		publishStatus: base.JobPublishStatusUnpublished,
		url:           url,
		title:         title,
		description:   description,
		source:        source,
		location:      location,
		remote:        remote,
		postedAt:      postedAt,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *Job) ID() uuid.UUID {
	return j.id
}

func (j *Job) ChannelID() uuid.UUID {
	return j.channelID
}

func (j *Job) URL() string {
	return j.url
}

func (j *Job) Title() string {
	return j.title
}

func (j *Job) Description() string {
	return j.description
}

func (j *Job) Source() string {
	return j.source
}

func (j *Job) Location() string {
	return j.location
}

func (j *Job) Remote() bool {
	return j.remote
}

func (j *Job) PostedAt() time.Time {
	return j.postedAt
}

func (j *Job) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Job) UpdatedAt() time.Time {
	return j.updatedAt
}

func (j *Job) Status() base.JobStatus {
	return j.status
}

func (j *Job) PublishStatus() base.JobPublishStatus {
	return j.publishStatus
}

func (j *Job) MarkAsMissing() {
	j.status = base.JobStatusInactive
	j.publishStatus = base.JobPublishStatusUnpublished
	j.updatedAt = time.Now()
}

func (j *Job) MarkAsChanged() {
	j.status = base.JobStatusActive
	j.publishStatus = base.JobPublishStatusUnpublished
	j.updatedAt = time.Now()
}

func (j *Job) IsEqual(other *Job) bool {
	// ignore publish status
	return j.status == other.status &&
		j.id == other.id &&
		j.channelID == other.channelID &&
		j.url == other.url &&
		j.title == other.title &&
		j.description == other.description &&
		j.source == other.source &&
		j.location == other.location &&
		j.remote == other.remote &&
		j.postedAt.Equal(other.postedAt)
}

func (j *Job) ToDTO() *postgres.Job {
	return &postgres.Job{
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
		Status:        j.Status(),
		PublishStatus: j.PublishStatus(),
	}
}

func NewJobFromDTO(j *postgres.Job) *Job {
	return NewJob(
		j.ID,
		j.ChannelID,
		j.Status,
		j.URL,
		j.Title,
		j.Description,
		j.Source,
		j.Location,
		j.Remote,
		j.PostedAt,
		JobWithTimestamps(j.CreatedAt, j.UpdatedAt),
		JobWithPublishStatus(j.PublishStatus),
	)
}
