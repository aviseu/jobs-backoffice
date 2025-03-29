package base

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
