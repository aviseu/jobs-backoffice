package base

type JobStatus int

const (
	JobStatusNew JobStatus = iota
	JobStatusUpdated
	JobStatusNoChange
	JobStatusMissing
	JobStatusFailed
)

func (s JobStatus) String() string {
	return [...]string{"new", "updated", "no_change", "missing", "failed"}[s]
}
