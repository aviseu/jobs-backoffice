package importing

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/google/uuid"
)

type ImportJobResult struct {
	jobID    uuid.UUID
	importID uuid.UUID
	result   base.ImportJobResult
}

func NewImportJobResult(jobID, importID uuid.UUID, result base.ImportJobResult) *ImportJobResult {
	return &ImportJobResult{
		jobID:    jobID,
		importID: importID,
		result:   result,
	}
}

func (j *ImportJobResult) JobID() uuid.UUID {
	return j.jobID
}

func (j *ImportJobResult) ImportID() uuid.UUID {
	return j.importID
}

func (j *ImportJobResult) Result() base.ImportJobResult {
	return j.result
}

func (j *ImportJobResult) ToDTO() *postgres.ImportJobResult {
	return &postgres.ImportJobResult{
		ID:       j.JobID(),
		ImportID: j.ImportID(),
		Result:   j.Result(),
	}
}

func NewJobResultFromDTO(j *postgres.ImportJobResult) *ImportJobResult {
	return NewImportJobResult(
		j.ID,
		j.ImportID,
		j.Result,
	)
}
