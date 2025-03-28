package imports

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type Repository interface {
	SaveImport(ctx context.Context, i *Import) error
	SaveImportJob(ctx context.Context, j *JobResult) error

	GetImports(ctx context.Context) ([]*Import, error)
	FindImport(ctx context.Context, id uuid.UUID) (*Import, error)
	GetJobsByImportID(ctx context.Context, importID uuid.UUID) ([]*JobResult, error)
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Start(ctx context.Context, id, channelID uuid.UUID) (*Import, error) {
	i := New(id, channelID)
	if err := s.r.SaveImport(ctx, i); err != nil {
		return nil, fmt.Errorf("failed to save import for channel %s while starting: %w", channelID, err)
	}

	return i, nil
}

func (s *Service) SaveJobResult(ctx context.Context, r *JobResult) error {
	return s.r.SaveImportJob(ctx, r)
}

func (s *Service) SetStatus(ctx context.Context, i *Import, status Status) error {
	i.status = status
	if err := s.r.SaveImport(ctx, i); err != nil {
		return fmt.Errorf("failed to set status %s for import %s: %w", status.String(), i.ID(), err)
	}

	return nil
}

func (s *Service) MarkAsCompleted(ctx context.Context, i *Import) error {
	i.status = StatusCompleted
	i.endedAt = null.TimeFrom(time.Now())

	if err := s.setMetadataFromJobs(i); err != nil {
		return fmt.Errorf("failed to fill metadata from jobs for import %s while marking as completed: %w", i.ID(), err)
	}

	if err := s.r.SaveImport(ctx, i); err != nil {
		return fmt.Errorf("failed to mark import %s as completed: %w", i.ID(), err)
	}

	return nil
}

func (s *Service) MarkAsFailed(ctx context.Context, i *Import, err error) error {
	i.status = StatusFailed
	i.endedAt = null.TimeFrom(time.Now())
	i.error = null.StringFrom(err.Error())

	if err := s.setMetadataFromJobs(i); err != nil {
		return fmt.Errorf("failed to fill metadata from jobs for import %s while marking as failed: %w", i.ID(), err)
	}

	if err := s.r.SaveImport(ctx, i); err != nil {
		return fmt.Errorf("failed to mark import %s as failed: %w", i.ID(), err)
	}

	return nil
}

func (s *Service) GetImports(ctx context.Context) ([]*Import, error) {
	return s.r.GetImports(ctx)
}

func (s *Service) FindImport(ctx context.Context, id uuid.UUID) (*Import, error) {
	return s.r.FindImport(ctx, id)
}

func (s *Service) FindImportWithForcedMetadata(ctx context.Context, id uuid.UUID) (*Import, error) {
	i, err := s.r.FindImport(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find import %s: %w", id, err)
	}

	if i.TotalJobs() == 0 {
		jj, err := s.r.GetJobsByImportID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs for import %s: %w", id, err)
		}

		i.resetMetadata()
		for _, r := range jj {
			i.addJobResult(r.Result())
		}
	}

	return i, nil
}

func (s *Service) setMetadataFromJobs(i *Import) error {
	jobs, err := s.r.GetJobsByImportID(context.Background(), i.ID())
	if err != nil {
		return fmt.Errorf("failed to get jobs for import %s while filling metadata: %w", i.ID(), err)
	}

	i.resetMetadata()
	for _, r := range jobs {
		i.addJobResult(r.Result())
	}

	return nil
}
