package postgres

import (
	"time"

	"github.com/aviseu/jobs/internal/app/domain/imports"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Import struct {
	StartedAt    time.Time   `db:"started_at"`
	EndedAt      null.Time   `db:"ended_at"`
	Error        null.String `db:"error"`
	Status       int         `db:"status"`
	NewJobs      int         `db:"new_jobs"`
	UpdatedJobs  int         `db:"updated_jobs"`
	NoChangeJobs int         `db:"no_change_jobs"`
	MissingJobs  int         `db:"missing_jobs"`
	FailedJobs   int         `db:"failed_jobs"`
	ID           uuid.UUID   `db:"id"`
	ChannelID    uuid.UUID   `db:"channel_id"`
}

func fromDomainImport(i *imports.Import) *Import {
	return &Import{
		ID:           i.ID(),
		ChannelID:    i.ChannelID(),
		StartedAt:    i.StartedAt(),
		EndedAt:      i.EndedAt(),
		Error:        i.Error(),
		Status:       int(i.Status()),
		NewJobs:      i.NewJobs(),
		UpdatedJobs:  i.UpdatedJobs(),
		NoChangeJobs: i.NoChangeJobs(),
		MissingJobs:  i.MissingJobs(),
		FailedJobs:   i.FailedJobs(),
	}
}

func toDomainImport(i *Import) *imports.Import {
	return imports.New(
		i.ID,
		i.ChannelID,
		imports.WithError(i.Error.String),
		imports.WithStartAt(i.StartedAt),
		imports.WithEndAt(i.EndedAt.Time),
		imports.WithStatus(imports.Status(i.Status)),
		imports.WithMetadata(i.NewJobs, i.UpdatedJobs, i.NoChangeJobs, i.MissingJobs, i.FailedJobs),
	)
}
