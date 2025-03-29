package postgres

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/google/uuid"
)

type ImportJobResult struct {
	ID       uuid.UUID      `db:"job_id"`
	ImportID uuid.UUID      `db:"import_id"`
	Result   base.JobStatus `db:"result"`
}
