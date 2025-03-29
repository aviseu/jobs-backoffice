package postgres

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
	Name        string             `db:"name"`
	Integration base.Integration   `db:"integration"`
	Status      base.ChannelStatus `db:"status"`
	ID          uuid.UUID          `db:"id"`
}
