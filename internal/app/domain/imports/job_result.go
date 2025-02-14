package imports

import (
	"github.com/google/uuid"
)

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

type JobResult struct {
	id       uuid.UUID
	importID uuid.UUID
	result   JobStatus
}

func NewResult(id, importID uuid.UUID, result JobStatus) *JobResult {
	return &JobResult{
		id:       id,
		importID: importID,
		result:   result,
	}
}

func (j *JobResult) ID() uuid.UUID {
	return j.id
}

func (j *JobResult) ImportID() uuid.UUID {
	return j.importID
}

func (j *JobResult) Result() JobStatus {
	return j.result
}
