package imports

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/google/uuid"
)

type JobResult struct {
	jobID    uuid.UUID
	importID uuid.UUID
	result   base.JobStatus
}

func NewResult(jobID, importID uuid.UUID, result base.JobStatus) *JobResult {
	return &JobResult{
		jobID:    jobID,
		importID: importID,
		result:   result,
	}
}

func (j *JobResult) JobID() uuid.UUID {
	return j.jobID
}

func (j *JobResult) ImportID() uuid.UUID {
	return j.importID
}

func (j *JobResult) Result() base.JobStatus {
	return j.result
}

func (j *JobResult) ToDTO() *postgres.ImportJobResult {
	return &postgres.ImportJobResult{
		ID:       j.JobID(),
		ImportID: j.ImportID(),
		Result:   int(j.Result()),
	}
}

func NewJobResultFromDTO(j *postgres.ImportJobResult) *JobResult {
	return NewResult(
		j.ID,
		j.ImportID,
		base.JobStatus(j.Result),
	)
}
