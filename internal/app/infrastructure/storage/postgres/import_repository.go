package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type ImportRepository struct {
	db *sqlx.DB
}

func NewImportRepository(db *sqlx.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

func (r *ImportRepository) SaveImport(ctx context.Context, i *aggregator.Import) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO imports (id, channel_id, status, started_at, ended_at, error)
				VALUES (:id, :channel_id, :status, :started_at, :ended_at, :error)
				ON CONFLICT (id) DO UPDATE SET
					channel_id = EXCLUDED.channel_id,
					status = EXCLUDED.status,
					started_at = EXCLUDED.started_at,
					ended_at = EXCLUDED.ended_at,
					error = EXCLUDED.error`,
		i,
	)
	if err != nil {
		return fmt.Errorf("failed to save import %s: %w", i.ID, err)
	}

	// Save import jobs
	var eg errgroup.Group
	for _, job := range i.Metrics {
		eg.Go(func() error {
			return r.SaveImportMetric(ctx, i.ID, job)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to save import jobs for import %s: %w", i.ID, err)
	}

	return nil
}

func (r *ImportRepository) FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error) {
	var i aggregator.Import
	var metrics []*aggregator.ImportMetric
	var eg errgroup.Group

	eg.Go(func() error {
		err := r.db.GetContext(ctx, &i, "SELECT * FROM imports WHERE id = $1", id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return infrastructure.ErrImportNotFound
			}

			return fmt.Errorf("failed to get import %s: %w", id, err)
		}
		return nil
	})

	eg.Go(func() error {
		err := r.db.SelectContext(ctx, &metrics, "SELECT id, job_id, metric_type, error, created_at FROM import_metrics WHERE import_id = $1 ORDER BY created_at DESC", id)
		if err != nil {
			return fmt.Errorf("failed to get metrics for import %s: %w", id, err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	i.Metrics = metrics

	return &i, nil
}

func (r *ImportRepository) GetImports(ctx context.Context) ([]*aggregator.Import, error) {
	var imports []*aggregator.Import
	err := r.db.SelectContext(ctx, &imports, "SELECT * FROM imports order by started_at desc")
	if err != nil {
		return nil, fmt.Errorf("failed to get imports: %w", err)
	}

	var eg errgroup.Group
	for _, i := range imports {
		eg.Go(func() error {
			var metrics []*aggregator.ImportMetric
			err := r.db.SelectContext(ctx, &metrics, "SELECT id, job_id, metric_type, error, created_at FROM import_metrics WHERE import_id = $1 ORDER BY created_at DESC", i.ID)
			if err != nil {
				return fmt.Errorf("failed to get metrics for import %s: %w", i.ID, err)
			}
			i.Metrics = metrics
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return imports, nil
}

func (r *ImportRepository) SaveImportMetric(ctx context.Context, importID uuid.UUID, m *aggregator.ImportMetric) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO import_metrics (id, import_id, job_id, metric_type, error, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO UPDATE SET
					metric_type = EXCLUDED.metric_type,
					error = EXCLUDED.error`,
		m.ID,
		importID,
		m.JobID,
		m.MetricType,
		m.Err,
		m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save import job metric %s: %w", m.ID, err)
	}

	return nil
}
