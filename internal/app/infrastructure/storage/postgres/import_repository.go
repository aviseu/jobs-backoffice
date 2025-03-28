package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ImportRepository struct {
	db *sqlx.DB
}

func NewImportRepository(db *sqlx.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

func (r *ImportRepository) SaveImport(ctx context.Context, i *imports.Import) error {
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
		fromDomainImport(i),
	)
	if err != nil {
		return fmt.Errorf("failed to save import %s: %w", i.ID(), err)
	}

	return nil
}

func (r *ImportRepository) FindImport(ctx context.Context, id uuid.UUID) (*imports.Import, error) {
	var i Import
	err := r.db.GetContext(ctx, &i, "SELECT * FROM imports WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, imports.ErrImportNotFound
		}

		return nil, fmt.Errorf("failed to get import %s: %w", id, err)
	}

	return toDomainImport(&i), nil
}

func (r *ImportRepository) GetImports(ctx context.Context) ([]*imports.Import, error) {
	var ii []*Import
	err := r.db.SelectContext(ctx, &ii, "SELECT * FROM imports order by started_at desc")
	if err != nil {
		return nil, fmt.Errorf("failed to get imports: %w", err)
	}

	result := make([]*imports.Import, 0, len(ii))
	for _, i := range ii {
		result = append(result, toDomainImport(i))
	}

	return result, nil
}

func (r *ImportRepository) SaveImportJob(ctx context.Context, jr *imports.JobResult) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO import_job_results (import_id, job_id, result)
				VALUES (:import_id, :job_id, :result)
				ON CONFLICT (import_id, job_id) DO UPDATE SET
					result = EXCLUDED.result`,
		fromDomainImportJobResult(jr),
	)
	if err != nil {
		return fmt.Errorf("failed to save import job result %s: %w", jr.JobID(), err)
	}

	return nil
}

func (r *ImportRepository) GetJobsByImportID(ctx context.Context, importID uuid.UUID) ([]*imports.JobResult, error) {
	var jj []*ImportJobResult
	err := r.db.SelectContext(ctx, &jj, "SELECT * FROM import_job_results WHERE import_id = $1", importID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs for import %s: %w", importID, err)
	}

	result := make([]*imports.JobResult, 0, len(jj))
	for _, j := range jj {
		result = append(result, toDomainImportJobResult(j))
	}

	return result, nil
}
