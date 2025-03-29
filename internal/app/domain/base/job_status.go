package base

type JobStatus int

const (
	JobStatusInactive JobStatus = iota
	JobStatusActive
)

func (s JobStatus) String() string {
	return [...]string{"inactive", "active"}[s]
}
