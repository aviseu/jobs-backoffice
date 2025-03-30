package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrImportNotFound = errors.New("import not found")

type ImportRepository struct {
	db *sqlx.DB
}

func NewImportRepository(db *sqlx.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

func (r *ImportRepository) SaveImport(ctx context.Context, i *aggregator.Import) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO imports (id, channel_id, status, started_at, ended_at, error, new_jobs, updated_jobs, no_change_jobs, missing_jobs, failed_jobs)
				VALUES (:id, :channel_id, :status, :started_at, :ended_at, :error, :new_jobs, :updated_jobs, :no_change_jobs, :missing_jobs, :failed_jobs)
				ON CONFLICT (id) DO UPDATE SET
					channel_id = EXCLUDED.channel_id,
					status = EXCLUDED.status,
					started_at = EXCLUDED.started_at,
					ended_at = EXCLUDED.ended_at,
					error = EXCLUDED.error,
					new_jobs = EXCLUDED.new_jobs,
					updated_jobs = EXCLUDED.updated_jobs,
					no_change_jobs = EXCLUDED.no_change_jobs,
					missing_jobs = EXCLUDED.missing_jobs,
					failed_jobs = EXCLUDED.failed_jobs`,
		i,
	)
	if err != nil {
		return fmt.Errorf("failed to save import %s: %w", i.ID, err)
	}

	return nil
}

func (r *ImportRepository) FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error) {
	var i aggregator.Import
	err := r.db.GetContext(ctx, &i, "SELECT * FROM imports WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrImportNotFound
		}

		return nil, fmt.Errorf("failed to get import %s: %w", id, err)
	}

	return &i, nil
}

func (r *ImportRepository) GetImports(ctx context.Context) ([]*aggregator.Import, error) {
	var result []*aggregator.Import
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM imports order by started_at desc")
	if err != nil {
		return nil, fmt.Errorf("failed to get imports: %w", err)
	}

	return result, nil
}

func (r *ImportRepository) SaveImportJob(ctx context.Context, jr *ImportJobResult) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO import_job_results (import_id, job_id, result)
				VALUES (:import_id, :job_id, :result)
				ON CONFLICT (import_id, job_id) DO UPDATE SET
					result = EXCLUDED.result`,
		jr,
	)
	if err != nil {
		return fmt.Errorf("failed to save import job result %s: %w", jr.ID, err)
	}

	return nil
}

func (r *ImportRepository) GetJobsByImportID(ctx context.Context, importID uuid.UUID) ([]*ImportJobResult, error) {
	var result []*ImportJobResult
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM import_job_results WHERE import_id = $1", importID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs for import %s: %w", importID, err)
	}

	return result, nil
}
