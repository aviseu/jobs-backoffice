package postgres

import (
	"context"
	"fmt"

	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/jmoiron/sqlx"
)

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Save(ctx context.Context, job *job.Job) error {
	j := fromDomainJob(job)
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO jobs (id, url, title, description, source, location, remote, posted_at, created_at, updated_at)
				VALUES (:id, :url, :title, :description, :source, :location, :remote, :posted_at, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET
					url = EXCLUDED.url,
					title = EXCLUDED.title,
					description = EXCLUDED.description,
					source = EXCLUDED.source,
					location = EXCLUDED.location,
					remote = EXCLUDED.remote,
					posted_at = EXCLUDED.posted_at,
					updated_at = EXCLUDED.updated_at`,
		j,
	)
	if err != nil {
		return fmt.Errorf("failed to save job %s: %w", job.ID(), err)
	}

	return nil
}
