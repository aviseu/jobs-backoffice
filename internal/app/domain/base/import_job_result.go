package base

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
