package postgres

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Import struct {
	StartedAt    time.Time         `db:"started_at"`
	EndedAt      null.Time         `db:"ended_at"`
	Error        null.String       `db:"error"`
	Status       base.ImportStatus `db:"status"`
	NewJobs      int               `db:"new_jobs"`
	UpdatedJobs  int               `db:"updated_jobs"`
	NoChangeJobs int               `db:"no_change_jobs"`
	MissingJobs  int               `db:"missing_jobs"`
	FailedJobs   int               `db:"failed_jobs"`
	ID           uuid.UUID         `db:"id"`
	ChannelID    uuid.UUID         `db:"channel_id"`
}
