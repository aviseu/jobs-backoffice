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

type ImportMetricType int

const (
	ImportMetricTypeNew ImportMetricType = iota
	ImportMetricTypeUpdated
	ImportMetricTypeNoChange
	ImportMetricTypeMissing
	ImportMetricTypeError
	ImportMetricTypePublish
	ImportMetricTypeLatePublish
	ImportMetricTypeMissingPublish
)

func (s ImportMetricType) String() string {
	return [...]string{"new", "updated", "no_change", "missing", "error", "publish", "late_publish", "missing_publish"}[s]
}

type ImportMetric struct {
	ID         uuid.UUID        `db:"id"`
	JobID      uuid.UUID        `db:"job_id"`
	MetricType ImportMetricType `db:"metric_type"`
	Err        null.String      `db:"error"`
	CreatedAt  time.Time        `db:"created_at"`
}

type Import struct {
	StartedAt time.Time       `db:"started_at"`
	EndedAt   null.Time       `db:"ended_at"`
	Error     null.String     `db:"error"`
	Metrics   []*ImportMetric `db:"jobs"`
	Status    ImportStatus    `db:"status"`
	ID        uuid.UUID       `db:"id"`
	ChannelID uuid.UUID       `db:"channel_id"`
}

func (i *Import) NewJobs() int {
	return i.jobCount(ImportMetricTypeNew)
}

func (i *Import) UpdatedJobs() int {
	return i.jobCount(ImportMetricTypeUpdated)
}

func (i *Import) NoChangeJobs() int {
	return i.jobCount(ImportMetricTypeNoChange)
}

func (i *Import) MissingJobs() int {
	return i.jobCount(ImportMetricTypeMissing)
}

func (i *Import) TotalJobs() int {
	return i.NewJobs() + i.UpdatedJobs() + i.NoChangeJobs()
}

func (i *Import) Errors() int {
	return i.jobCount(ImportMetricTypeError)
}

func (i *Import) Published() int {
	return i.jobCount(ImportMetricTypePublish)
}

func (i *Import) LatePublished() int {
	return i.jobCount(ImportMetricTypeLatePublish)
}

func (i *Import) MissingPublished() int {
	return i.jobCount(ImportMetricTypeMissingPublish)
}

func (i *Import) jobCount(metricType ImportMetricType) int {
	var count int
	for _, j := range i.Metrics {
		if j.MetricType == metricType {
			count++
		}
	}
	return count
}
