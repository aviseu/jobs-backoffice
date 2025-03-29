package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Save(ctx context.Context, j *Job) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO jobs (id, channel_id, status, publish_status, url, title, description, source, location, remote, posted_at, created_at, updated_at)
				VALUES (:id, :channel_id, :status, :publish_status, :url, :title, :description, :source, :location, :remote, :posted_at, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET
					channel_id = EXCLUDED.channel_id,
					status = EXCLUDED.status,
					publish_status = EXCLUDED.publish_status,
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
		return fmt.Errorf("failed to save job %s: %w", j.ID, err)
	}

	return nil
}

func (r *JobRepository) GetByChannelID(ctx context.Context, chID uuid.UUID) ([]*Job, error) {
	var results []*Job
	err := r.db.SelectContext(ctx, &results, "SELECT * FROM jobs WHERE channel_id = $1 ORDER BY posted_at DESC", chID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by channel id %s: %w", chID, err)
	}

	return results, nil
}
