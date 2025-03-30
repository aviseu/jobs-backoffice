package importing

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Import struct {
	startedAt time.Time
	endedAt   null.Time
	error     null.String
	jobs      []*aggregator.ImportJob
	status    aggregator.ImportStatus
	id        uuid.UUID
	channelID uuid.UUID
}

type ImportOptional func(*Import)

func ImportWithStatus(s aggregator.ImportStatus) ImportOptional {
	return func(i *Import) {
		i.status = s
	}
}

func ImportWithError(err string) ImportOptional {
	return func(i *Import) {
		i.error = null.StringFrom(err)
	}
}

func ImportWithStartAt(s time.Time) ImportOptional {
	return func(i *Import) {
		i.startedAt = s
	}
}

func ImportWithEndAt(e time.Time) ImportOptional {
	return func(i *Import) {
		i.endedAt = null.TimeFrom(e)
	}
}

func ImportWithJobs(jobs []*aggregator.ImportJob) ImportOptional {
	return func(i *Import) {
		i.jobs = jobs
	}
}

func NewImport(id, channelID uuid.UUID, opts ...ImportOptional) *Import {
	i := &Import{
		id:        id,
		channelID: channelID,
		status:    aggregator.ImportStatusPending,
		startedAt: time.Now(),
		endedAt:   null.NewTime(time.Now(), false),
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

func (i *Import) Status() aggregator.ImportStatus {
	return i.status
}

func (i *Import) Error() null.String {
	return i.error
}

func (i *Import) NewJobs() int {
	return i.jobCount(aggregator.ImportJobResultNew)
}

func (i *Import) UpdatedJobs() int {
	return i.jobCount(aggregator.ImportJobResultUpdated)
}

func (i *Import) NoChangeJobs() int {
	return i.jobCount(aggregator.ImportJobResultNoChange)
}

func (i *Import) MissingJobs() int {
	return i.jobCount(aggregator.ImportJobResultMissing)
}

func (i *Import) FailedJobs() int {
	return i.jobCount(aggregator.ImportJobResultFailed)
}

func (i *Import) jobCount(result aggregator.ImportJobResult) int {
	var count int
	for _, j := range i.jobs {
		if j.Result == result {
			count++
		}
	}
	return count
}

func (i *Import) TotalJobs() int {
	return len(i.jobs)
}

func (i *Import) StartedAt() time.Time {
	return i.startedAt
}

func (i *Import) EndedAt() null.Time {
	return i.endedAt
}

func (i *Import) ToAggregate() *aggregator.Import {
	return &aggregator.Import{
		ID:        i.ID(),
		ChannelID: i.ChannelID(),
		StartedAt: i.StartedAt(),
		EndedAt:   i.EndedAt(),
		Error:     i.Error(),
		Status:    i.Status(),
	}
}

func NewImportFromDTO(i *aggregator.Import) *Import {
	opts := []ImportOptional{
		ImportWithStartAt(i.StartedAt),
		ImportWithStatus(i.Status),
		ImportWithJobs(i.Jobs),
	}
	if i.Error.Valid {
		opts = append(opts, ImportWithError(i.Error.String))
	}
	if i.EndedAt.Valid {
		opts = append(opts, ImportWithEndAt(i.EndedAt.Time))
	}

	return NewImport(
		i.ID,
		i.ChannelID,
		opts...,
	)
}
