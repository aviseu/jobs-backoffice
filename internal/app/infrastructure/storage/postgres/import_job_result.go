package postgres

import (
	"github.com/google/uuid"
)

type ImportJobResult struct {
	ID       uuid.UUID `db:"job_id"`
	ImportID uuid.UUID `db:"import_id"`
	Result   int       `db:"result"`
}
