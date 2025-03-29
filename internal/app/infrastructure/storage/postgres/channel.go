package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"name"`
	Integration int       `db:"integration"`
	Status      int       `db:"status"`
	ID          uuid.UUID `db:"id"`
}
