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

type importMetadata struct {
	aggregator.ImportMetadata
	ImportID uuid.UUID `db:"import_id"`
}

type importEntry struct {
	aggregator.Import
	importMetadata
}

func (i *importEntry) toImport() *aggregator.Import {
	return &aggregator.Import{
		StartedAt: i.Import.StartedAt,
		EndedAt:   i.Import.EndedAt,
		Error:     i.Import.Error,
		Metrics:   i.Import.Metrics,
		Status:    i.Import.Status,
		ID:        i.Import.ID,
		ChannelID: i.Import.ChannelID,
		Metadata:  i.toImportMetadata(),
	}
}

func (i *importEntry) toImportMetadata() *aggregator.ImportMetadata {
	return &aggregator.ImportMetadata{
		New:              i.ImportMetadata.New,
		Updated:          i.ImportMetadata.Updated,
		NoChange:         i.ImportMetadata.NoChange,
		Missing:          i.ImportMetadata.Missing,
		Errors:           i.ImportMetadata.Errors,
		Published:        i.ImportMetadata.Published,
		LatePublished:    i.ImportMetadata.LatePublished,
		MissingPublished: i.ImportMetadata.MissingPublished,
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

	// Get import metrics
	var groups []struct {
		MetricType aggregator.ImportMetricType `db:"metric_type"`
		Count      int                         `db:"count"`
	}
	if err := r.db.SelectContext(ctx, &groups, "SELECT metric_type, count(*) FROM import_metrics WHERE import_id = $1 GROUP BY metric_type", i.ID); err != nil {
		return fmt.Errorf("failed to get import metrics for import %s: %w", i.ID, err)
	}

	// Create import metadata
	var metadata importMetadata
	metadata.ImportID = i.ID
	for _, group := range groups {
		switch group.MetricType {
		case aggregator.ImportMetricTypeNew:
			metadata.New = group.Count
		case aggregator.ImportMetricTypeUpdated:
			metadata.Updated = group.Count
		case aggregator.ImportMetricTypeNoChange:
			metadata.NoChange = group.Count
		case aggregator.ImportMetricTypeMissing:
			metadata.Missing = group.Count
		case aggregator.ImportMetricTypeError:
			metadata.Errors = group.Count
		case aggregator.ImportMetricTypePublish:
			metadata.Published = group.Count
		case aggregator.ImportMetricTypeLatePublish:
			metadata.LatePublished = group.Count
		case aggregator.ImportMetricTypeMissingPublish:
			metadata.MissingPublished = group.Count
		default:
			return fmt.Errorf("unknown metric type %s", group.MetricType)
		}
	}
	// Save import metadata
	_, err = r.db.NamedExecContext(
		ctx,
		`INSERT INTO import_metadata (import_id, new_jobs, updated_jobs, no_change_jobs, missing_jobs, errors, published, late_published, missing_published)
				VALUES (:import_id, :new_jobs, :updated_jobs, :no_change_jobs, :missing_jobs, :errors, :published, :late_published, :missing_published)
				ON CONFLICT (import_id) DO UPDATE SET
				   new_jobs = EXCLUDED.new_jobs,
				   updated_jobs = EXCLUDED.updated_jobs,
				   no_change_jobs = EXCLUDED.no_change_jobs,
				   missing_jobs = EXCLUDED.missing_jobs,
				   errors = EXCLUDED.errors,
				   published = EXCLUDED.published,
				   late_published = EXCLUDED.late_published,
				   missing_published = EXCLUDED.missing_published`,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to save import %s metadata: %w", i.ID, err)
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
	var imports []*importEntry
	err := r.db.SelectContext(ctx, &imports, `
		SELECT imports.*,
       		COALESCE(im.new_jobs, 0) as new_jobs,
       		COALESCE(im.updated_jobs, 0) as updated_jobs,
       		COALESCE(im.no_change_jobs, 0) as no_change_jobs,
       		COALESCE(im.missing_jobs, 0) as missing_jobs,
       		COALESCE(im.errors, 0) as errors,
       		COALESCE(im.published, 0) as published,
       		COALESCE(im.late_published, 0) as late_published,
       		COALESCE(im.missing_published, 0) as missing_published
       	FROM imports LEFT OUTER JOIN import_metadata AS im ON id = import_id order by started_at desc
   `)
	if err != nil {
		return nil, fmt.Errorf("failed to get imports: %w", err)
	}

	var agImports []*aggregator.Import
	for _, imp := range imports {
		agImports = append(agImports, imp.toImport())
	}

	return agImports, nil
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
