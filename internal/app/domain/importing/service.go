package importing

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type ImportRepository interface {
	SaveImport(ctx context.Context, i *postgres.Import) error
	SaveImportJob(ctx context.Context, j *postgres.ImportJobResult) error

	GetImports(ctx context.Context) ([]*postgres.Import, error)
	FindImport(ctx context.Context, id uuid.UUID) (*postgres.Import, error)
	GetJobsByImportID(ctx context.Context, importID uuid.UUID) ([]*postgres.ImportJobResult, error)
}

type ImportService struct {
	r ImportRepository
}

func NewImportService(r ImportRepository) *ImportService {
	return &ImportService{r: r}
}

func (s *ImportService) Start(ctx context.Context, id, channelID uuid.UUID) (*Import, error) {
	i := New(id, channelID)
	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return nil, fmt.Errorf("failed to save import for channel %s while starting: %w", channelID, err)
	}

	return i, nil
}

func (s *ImportService) SaveJobResult(ctx context.Context, r *JobResult) error {
	return s.r.SaveImportJob(ctx, r.ToDTO())
}

func (s *ImportService) SetStatus(ctx context.Context, i *Import, status base.ImportStatus) error {
	i.status = status
	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return fmt.Errorf("failed to set status %s for import %s: %w", status.String(), i.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsCompleted(ctx context.Context, i *Import) error {
	i.status = base.ImportStatusCompleted
	i.endedAt = null.TimeFrom(time.Now())

	if err := s.setMetadataFromJobs(i); err != nil {
		return fmt.Errorf("failed to fill metadata from jobs for import %s while marking as completed: %w", i.ID(), err)
	}

	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsFailed(ctx context.Context, i *Import, err error) error {
	i.status = base.ImportStatusFailed
	i.endedAt = null.TimeFrom(time.Now())
	i.error = null.StringFrom(err.Error())

	if err := s.setMetadataFromJobs(i); err != nil {
		return fmt.Errorf("failed to fill metadata from jobs for import %s while marking as failed: %w", i.ID(), err)
	}

	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return fmt.Errorf("failed to mark import %s as failed: %w", i.ID(), err)
	}

	return nil
}

func (s *ImportService) GetImports(ctx context.Context) ([]*Import, error) {
	imports, err := s.r.GetImports(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get imports: %w", err)
	}

	results := make([]*Import, 0, len(imports))
	for _, dto := range imports {
		results = append(results, NewImportFromDTO(dto))
	}

	return results, nil
}

func (s *ImportService) FindImport(ctx context.Context, id uuid.UUID) (*Import, error) {
	dto, err := s.r.FindImport(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrImportNotFound) {
			return nil, ErrImportNotFound
		}
		return nil, fmt.Errorf("failed to find import %s: %w", id, err)
	}

	return NewImportFromDTO(dto), nil
}

func (s *ImportService) FindImportWithForcedMetadata(ctx context.Context, id uuid.UUID) (*Import, error) {
	dto, err := s.r.FindImport(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrImportNotFound) {
			return nil, ErrImportNotFound
		}
		return nil, fmt.Errorf("failed to find import %s: %w", id, err)
	}

	i := NewImportFromDTO(dto)

	if i.TotalJobs() == 0 {
		jj, err := s.r.GetJobsByImportID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs for import %s: %w", id, err)
		}

		i.resetMetadata()
		for _, r := range jj {
			i.addJobResult(NewJobResultFromDTO(r).Result())
		}
	}

	return i, nil
}

func (s *ImportService) setMetadataFromJobs(i *Import) error {
	jobs, err := s.r.GetJobsByImportID(context.Background(), i.ID())
	if err != nil {
		return fmt.Errorf("failed to get jobs for import %s while filling metadata: %w", i.ID(), err)
	}

	i.resetMetadata()
	for _, r := range jobs {
		i.addJobResult(NewJobResultFromDTO(r).Result())
	}

	return nil
}
