package importing

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type ImportRepository interface {
	SaveImport(ctx context.Context, i *aggregator.Import) error
	SaveImportJob(ctx context.Context, importID uuid.UUID, j *aggregator.ImportJob) error

	GetImports(ctx context.Context) ([]*aggregator.Import, error)
	FindImport(ctx context.Context, id uuid.UUID) (*aggregator.Import, error)
}

type ImportService struct {
	r ImportRepository
}

func NewImportService(r ImportRepository) *ImportService {
	return &ImportService{r: r}
}

func (s *ImportService) Start(ctx context.Context, id, channelID uuid.UUID) (*Import, error) {
	i := NewImport(id, channelID)
	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return nil, fmt.Errorf("failed to save import for channel %s while starting: %w", channelID, err)
	}

	return i, nil
}

func (s *ImportService) SaveJobResult(ctx context.Context, importID uuid.UUID, j *aggregator.ImportJob) error {
	return s.r.SaveImportJob(ctx, importID, j)
}

func (s *ImportService) SetStatus(ctx context.Context, i *Import, status aggregator.ImportStatus) error {
	i.status = status
	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return fmt.Errorf("failed to set status %s for import %s: %w", status.String(), i.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsCompleted(ctx context.Context, i *Import) error {
	i.status = aggregator.ImportStatusCompleted
	i.endedAt = null.TimeFrom(time.Now())

	if err := s.r.SaveImport(ctx, i.ToDTO()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsFailed(ctx context.Context, i *Import, err error) error {
	i.status = aggregator.ImportStatusFailed
	i.endedAt = null.TimeFrom(time.Now())
	i.error = null.StringFrom(err.Error())

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
		if errors.Is(err, infrastructure.ErrImportNotFound) {
			return nil, ErrImportNotFound
		}
		return nil, fmt.Errorf("failed to find import %s: %w", id, err)
	}

	return NewImportFromDTO(dto), nil
}
