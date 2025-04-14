package importing

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
)

type job struct {
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
	status        aggregator.JobStatus
	publishStatus aggregator.JobPublishStatus
}

func newJob(id, channelID uuid.UUID, s aggregator.JobStatus, url, title, description, source, location string, remote bool, postedAt time.Time, publishStatus aggregator.JobPublishStatus, createdAt, updatedAt time.Time) *job {
	return &job{
		id:            id,
		channelID:     channelID,
		status:        s,
		publishStatus: publishStatus,
		url:           url,
		title:         title,
		description:   description,
		source:        source,
		location:      location,
		remote:        remote,
		postedAt:      postedAt,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (j *job) markAsMissing() {
	j.status = aggregator.JobStatusInactive
	j.publishStatus = aggregator.JobPublishStatusUnpublished
	j.updatedAt = time.Now()
}

func (j *job) markAsChanged() {
	j.status = aggregator.JobStatusActive
	j.publishStatus = aggregator.JobPublishStatusUnpublished
	j.updatedAt = time.Now()
}

func (j *job) IsEqual(other *job) bool {
	// ignore publish status, created at and updated at
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

func (j *job) toAggregator() *aggregator.Job {
	return &aggregator.Job{
		ID:            j.id,
		ChannelID:     j.channelID,
		URL:           j.url,
		Title:         j.title,
		Description:   j.description,
		Source:        j.source,
		Location:      j.location,
		Remote:        j.remote,
		PostedAt:      j.postedAt,
		CreatedAt:     j.createdAt,
		UpdatedAt:     j.updatedAt,
		Status:        j.status,
		PublishStatus: j.publishStatus,
	}
}

func newJobFromAggregator(j *aggregator.Job) *job {
	return newJob(
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
		j.PublishStatus,
		j.CreatedAt,
		j.UpdatedAt,
	)
}
