package job

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

type PublishStatus int

const (
	PublishStatusUnpublished PublishStatus = iota
	PublishStatusPublished
)

func (s PublishStatus) String() string {
	return [...]string{"unpublished", "published"}[s]
}

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
	status        Status
	publishStatus PublishStatus
}

type Optional func(*Job)

func WithTimestamps(c, u time.Time) Optional {
	return func(j *Job) {
		j.createdAt = c
		j.updatedAt = u
	}
}

func WithPublishStatus(s PublishStatus) Optional {
	return func(j *Job) {
		j.publishStatus = s
	}
}

func New(id, channelID uuid.UUID, s Status, url, title, description, source, location string, remote bool, postedAt time.Time, opts ...Optional) *Job {
	j := &Job{
		id:            id,
		channelID:     channelID,
		status:        s,
		publishStatus: PublishStatusUnpublished,
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

func (j *Job) Status() Status {
	return j.status
}

func (j *Job) PublishStatus() PublishStatus {
	return j.publishStatus
}
