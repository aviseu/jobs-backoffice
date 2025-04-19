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

type ImportMetadata struct {
	New              int `db:"new_jobs"`
	Updated          int `db:"updated_jobs"`
	NoChange         int `db:"no_change_jobs"`
	Missing          int `db:"missing_jobs"`
	Errors           int `db:"errors"`
	Published        int `db:"published"`
	LatePublished    int `db:"late_published"`
	MissingPublished int `db:"missing_published"`
}

type Import struct {
	StartedAt time.Time       `db:"started_at"`
	EndedAt   null.Time       `db:"ended_at"`
	Error     null.String     `db:"error"`
	Metrics   []*ImportMetric `db:"jobs"`
	Status    ImportStatus    `db:"status"`
	ID        uuid.UUID       `db:"id"`
	ChannelID uuid.UUID       `db:"channel_id"`
	Metadata  *ImportMetadata `db:"-"`
}

func (i *Import) NewJobs() int {
	if i.Metadata != nil {
		return i.Metadata.New
	}
	return i.jobCount(ImportMetricTypeNew)
}

func (i *Import) UpdatedJobs() int {
	if i.Metadata != nil {
		return i.Metadata.Updated
	}
	return i.jobCount(ImportMetricTypeUpdated)
}

func (i *Import) NoChangeJobs() int {
	if i.Metadata != nil {
		return i.Metadata.NoChange
	}
	return i.jobCount(ImportMetricTypeNoChange)
}

func (i *Import) MissingJobs() int {
	if i.Metadata != nil {
		return i.Metadata.Missing
	}
	return i.jobCount(ImportMetricTypeMissing)
}

func (i *Import) TotalJobs() int {
	if i.Metadata != nil {
		return i.Metadata.New + i.Metadata.Updated + i.Metadata.NoChange
	}
	return i.NewJobs() + i.UpdatedJobs() + i.NoChangeJobs()
}

func (i *Import) Errors() int {
	if i.Metadata != nil {
		return i.Metadata.Errors
	}
	return i.jobCount(ImportMetricTypeError)
}

func (i *Import) Published() int {
	if i.Metadata != nil {
		return i.Metadata.Published
	}
	return i.jobCount(ImportMetricTypePublish)
}

func (i *Import) LatePublished() int {
	if i.Metadata != nil {
		return i.Metadata.LatePublished
	}
	return i.jobCount(ImportMetricTypeLatePublish)
}

func (i *Import) MissingPublished() int {
	if i.Metadata != nil {
		return i.Metadata.MissingPublished
	}
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
