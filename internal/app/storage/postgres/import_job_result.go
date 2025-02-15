package postgres

import (
	"github.com/aviseu/jobs/internal/app/domain/imports"
	"github.com/google/uuid"
)

type ImportJobResult struct {
	ID       uuid.UUID `db:"job_id"`
	ImportID uuid.UUID `db:"import_id"`
	Result   int       `db:"result"`
}

func fromDomainImportJobResult(j *imports.JobResult) *ImportJobResult {
	return &ImportJobResult{
		ID:       j.JobID(),
		ImportID: j.ImportID(),
		Result:   int(j.Result()),
	}
}

func toDomainImportJobResult(j *ImportJobResult) *imports.JobResult {
	return imports.NewResult(
		j.ID,
		j.ImportID,
		imports.JobStatus(j.Result),
	)
}
