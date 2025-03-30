package aggregator

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type ImportStatus int

const (
	ImportStatusPending ImportStatus = iota
	ImportStatusFetching
	ImportStatusProcessing
	ImportStatusPublishing
	ImportStatusCompleted
	ImportStatusFailed
)

func (s ImportStatus) String() string {
	return [...]string{"pending", "fetching", "processing", "publishing", "completed", "failed"}[s]
}

type ImportJobResult int

const (
	ImportJobResultNew ImportJobResult = iota
	ImportJobResultUpdated
	ImportJobResultNoChange
	ImportJobResultMissing
	ImportJobResultFailed
)

func (s ImportJobResult) String() string {
	return [...]string{"new", "updated", "no_change", "missing", "failed"}[s]
}

type ImportJob struct {
	ID     uuid.UUID       `db:"job_id"`
	Result ImportJobResult `db:"result"`
}

type Import struct {
	StartedAt time.Time    `db:"started_at"`
	EndedAt   null.Time    `db:"ended_at"`
	Error     null.String  `db:"error"`
	Status    ImportStatus `db:"status"`
	ID        uuid.UUID    `db:"id"`
	ChannelID uuid.UUID    `db:"channel_id"`
	Jobs      []*ImportJob `db:"jobs"`
}

func (i *Import) NewJobs() int {
	return i.jobCount(ImportJobResultNew)
}

func (i *Import) UpdatedJobs() int {
	return i.jobCount(ImportJobResultUpdated)
}

func (i *Import) NoChangeJobs() int {
	return i.jobCount(ImportJobResultNoChange)
}

func (i *Import) MissingJobs() int {
	return i.jobCount(ImportJobResultMissing)
}

func (i *Import) FailedJobs() int {
	return i.jobCount(ImportJobResultFailed)
}

func (i *Import) jobCount(result ImportJobResult) int {
	var count int
	for _, j := range i.Jobs {
		if j.Result == result {
			count++
		}
	}
	return count
}
