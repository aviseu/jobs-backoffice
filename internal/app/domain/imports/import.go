package imports

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Status int

const (
	StatusPending Status = iota
	StatusFetching
	StatusProcessing
	StatusPublishing
	StatusCompleted
	StatusFailed
)

func (s Status) String() string {
	return [...]string{"pending", "fetching", "processing", "publishing", "completed", "failed"}[s]
}

type Import struct {
	startedAt    time.Time
	endedAt      null.Time
	error        null.String
	status       Status
	newJobs      int
	updatedJobs  int
	noChangeJobs int
	missingJobs  int
	failedJobs   int
	id           uuid.UUID
	channelID    uuid.UUID
}

type ImportOptional func(*Import)

func WithStatus(s Status) ImportOptional {
	return func(i *Import) {
		i.status = s
	}
}

func WithError(err string) ImportOptional {
	return func(i *Import) {
		i.error = null.StringFrom(err)
	}
}

func WithStartAt(s time.Time) ImportOptional {
	return func(i *Import) {
		i.startedAt = s
	}
}

func WithEndAt(e time.Time) ImportOptional {
	return func(i *Import) {
		i.endedAt = null.TimeFrom(e)
	}
}

func WithMetadata(newJobs, updated, noChange, missing, failed int) ImportOptional {
	return func(i *Import) {
		i.newJobs = newJobs
		i.updatedJobs = updated
		i.noChangeJobs = noChange
		i.missingJobs = missing
		i.failedJobs = failed
	}
}

func New(id, channelID uuid.UUID, opts ...ImportOptional) *Import {
	i := &Import{
		id:        id,
		channelID: channelID,
		status:    StatusPending,
		startedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *Import) ID() uuid.UUID {
	return i.id
}

func (i *Import) ChannelID() uuid.UUID {
	return i.channelID
}

func (i *Import) Status() Status {
	return i.status
}

func (i *Import) Error() null.String {
	return i.error
}

func (i *Import) NewJobs() int {
	return i.newJobs
}

func (i *Import) UpdatedJobs() int {
	return i.updatedJobs
}

func (i *Import) NoChangeJobs() int {
	return i.noChangeJobs
}

func (i *Import) MissingJobs() int {
	return i.missingJobs
}

func (i *Import) FailedJobs() int {
	return i.failedJobs
}

func (i *Import) TotalJobs() int {
	return i.newJobs + i.updatedJobs + i.noChangeJobs + i.missingJobs + i.failedJobs
}

func (i *Import) StartedAt() time.Time {
	return i.startedAt
}

func (i *Import) EndedAt() null.Time {
	return i.endedAt
}

func (i *Import) resetMetadata() {
	i.newJobs = 0
	i.updatedJobs = 0
	i.noChangeJobs = 0
	i.missingJobs = 0
	i.failedJobs = 0
}

func (i *Import) addJobResult(r JobStatus) {
	switch r {
	case JobStatusNew:
		i.newJobs++
	case JobStatusUpdated:
		i.updatedJobs++
	case JobStatusNoChange:
		i.noChangeJobs++
	case JobStatusMissing:
		i.missingJobs++
	case JobStatusFailed:
		i.failedJobs++
	}
}
