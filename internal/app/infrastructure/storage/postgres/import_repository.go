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

type importMetric struct {
	aggregator.ImportMetric

	ImportID uuid.UUID `db:"import_id"`
}

func (m *importMetric) toAggregator() *aggregator.ImportMetric {
	return &aggregator.ImportMetric{
		ID:         m.ID,
		JobID:      m.JobID,
		MetricType: m.MetricType,
		Err:        m.Err,
		CreatedAt:  m.CreatedAt,
	}
}

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

	ids := make([]uuid.UUID, len(imports))
	for i, imp := range imports {
		ids[i] = imp.ID
	}

	query, args, err := sqlx.In("SELECT id, import_id, job_id, metric_type, error, created_at FROM import_metrics WHERE import_id IN (?) ORDER BY created_at DESC", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build query for import metrics: %w", err)
	}
	query = r.db.Rebind(query)

	var metrics []*importMetric
	if err := r.db.SelectContext(ctx, &metrics, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get import metrics: %w", err)
	}

	var indexes = make(map[uuid.UUID]int)
	for i, imp := range imports {
		indexes[imp.ID] = i
	}

	for _, m := range metrics {
		imports[indexes[m.ImportID]].Metrics = append(imports[indexes[m.ImportID]].Metrics, m.toAggregator())
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
