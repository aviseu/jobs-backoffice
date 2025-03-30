package importing

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Import struct {
	startedAt    time.Time
	endedAt      null.Time
	error        null.String
	status       base.ImportStatus
	newJobs      int
	updatedJobs  int
	noChangeJobs int
	missingJobs  int
	failedJobs   int
	id           uuid.UUID
	channelID    uuid.UUID
}

type ImportOptional func(*Import)

func ImportWithStatus(s base.ImportStatus) ImportOptional {
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

func ImportWithMetadata(newJobs, updated, noChange, missing, failed int) ImportOptional {
	return func(i *Import) {
		i.newJobs = newJobs
		i.updatedJobs = updated
		i.noChangeJobs = noChange
		i.missingJobs = missing
		i.failedJobs = failed
	}
}

func NewImport(id, channelID uuid.UUID, opts ...ImportOptional) *Import {
	i := &Import{
		id:        id,
		channelID: channelID,
		status:    base.ImportStatusPending,
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

func (i *Import) Status() base.ImportStatus {
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

func (i *Import) addJobResult(r base.ImportJobResult) {
	switch r {
	case base.ImportJobResultNew:
		i.newJobs++
	case base.ImportJobResultUpdated:
		i.updatedJobs++
	case base.ImportJobResultNoChange:
		i.noChangeJobs++
	case base.ImportJobResultMissing:
		i.missingJobs++
	case base.ImportJobResultFailed:
		i.failedJobs++
	}
}

func (i *Import) ToDTO() *aggregator.Import {
	return &aggregator.Import{
		ID:           i.ID(),
		ChannelID:    i.ChannelID(),
		StartedAt:    i.StartedAt(),
		EndedAt:      i.EndedAt(),
		Error:        i.Error(),
		Status:       i.Status(),
		NewJobs:      i.NewJobs(),
		UpdatedJobs:  i.UpdatedJobs(),
		NoChangeJobs: i.NoChangeJobs(),
		MissingJobs:  i.MissingJobs(),
		FailedJobs:   i.FailedJobs(),
	}
}

func NewImportFromDTO(i *aggregator.Import) *Import {
	opts := []ImportOptional{
		ImportWithStartAt(i.StartedAt),
		ImportWithStatus(i.Status),
		ImportWithMetadata(i.NewJobs, i.UpdatedJobs, i.NoChangeJobs, i.MissingJobs, i.FailedJobs),
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
