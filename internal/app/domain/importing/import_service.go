package importing

import (
	"context"
	"fmt"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
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

func (s *ImportService) SaveJobResult(ctx context.Context, importID uuid.UUID, j *aggregator.ImportJob) error {
	return s.r.SaveImportJob(ctx, importID, j)
}

func (s *ImportService) SetStatus(ctx context.Context, i *aggregator.Import, status aggregator.ImportStatus) error {
	ip := NewImportFromDTO(i)
	ip.status = status
	if err := s.r.SaveImport(ctx, ip.ToAggregate()); err != nil {
		return fmt.Errorf("failed to set status %s for import %s: %w", status.String(), ip.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsCompleted(ctx context.Context, i *aggregator.Import) error {
	ip := NewImportFromDTO(i)

	ip.status = aggregator.ImportStatusCompleted
	ip.endedAt = null.TimeFrom(time.Now())

	if err := s.r.SaveImport(ctx, ip.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", ip.ID(), err)
	}

	return nil
}

func (s *ImportService) MarkAsFailed(ctx context.Context, i *aggregator.Import, err error) error {
	ip := NewImportFromDTO(i)

	ip.status = aggregator.ImportStatusFailed
	ip.endedAt = null.TimeFrom(time.Now())
	ip.error = null.StringFrom(err.Error())

	if err := s.r.SaveImport(ctx, ip.ToAggregate()); err != nil {
		return fmt.Errorf("failed to mark import %s as failed: %w", i.ID, err)
	}

	return nil
}
