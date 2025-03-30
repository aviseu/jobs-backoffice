package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"slices"
	"sync"
)

type ImportRepository struct {
	Imports map[uuid.UUID]*aggregator.Import
	err     error
	im      sync.Mutex
	jrm     sync.Mutex
}

func NewImportRepository() *ImportRepository {
	return &ImportRepository{
		Imports: make(map[uuid.UUID]*aggregator.Import),
	}
}

func (r *ImportRepository) First() *aggregator.Import {
	for _, i := range r.Imports {
		return i
	}

	return nil
}

func (r *ImportRepository) ImportJobs() map[uuid.UUID]*aggregator.ImportJob {
	jobs := make(map[uuid.UUID]*aggregator.ImportJob)
	for _, i := range r.Imports {
		for _, j := range i.Jobs {
			jobs[j.ID] = j
		}
	}
	return jobs
}

func (r *ImportRepository) AddImport(i *aggregator.Import) {
	r.Imports[i.ID] = i
}

func (r *ImportRepository) AddImportJob(importID uuid.UUID, j *aggregator.ImportJob) {
	i, ok := r.Imports[importID]
	if !ok {
		panic("import not found")
	}
	for idx, job := range i.Jobs {
		if job.ID == j.ID {
			i.Jobs[idx] = j
			return
		}
	}
	i.Jobs = append(i.Jobs, j)
}

func (r *ImportRepository) FailWith(err error) {
	r.err = err
}

func (r *ImportRepository) SaveImport(_ context.Context, i *aggregator.Import) error {
	if r.err != nil {
		return r.err
	}
	r.im.Lock()
	old, ok := r.Imports[i.ID]

	if ok {
		for _, job := range old.Jobs {
			for _, newJob := range i.Jobs {
				if job.ID == newJob.ID {
					continue
				}
			}
			i.Jobs = append(i.Jobs, job)
		}
	}
	r.Imports[i.ID] = i
	r.im.Unlock()
	return nil
}

func (r *ImportRepository) FindImport(_ context.Context, id uuid.UUID) (*aggregator.Import, error) {
	if r.err != nil {
		return nil, r.err
	}

	i, ok := r.Imports[id]
	if !ok {
		return nil, infrastructure.ErrImportNotFound
	}

	return i, nil
}

func (r *ImportRepository) GetImports(_ context.Context) ([]*aggregator.Import, error) {
	if r.err != nil {
		return nil, r.err
	}

	var ii []*aggregator.Import
	for _, i := range r.Imports {
		ii = append(ii, i)
	}

	slices.SortFunc(ii, func(a, b *aggregator.Import) int {
		return b.StartedAt.Compare(a.StartedAt)
	})

	return ii, nil
}

func (r *ImportRepository) SaveImportJob(_ context.Context, importID uuid.UUID, j *aggregator.ImportJob) error {
	if r.err != nil {
		return r.err
	}

	r.jrm.Lock()
	defer r.jrm.Unlock()

	i, ok := r.Imports[importID]
	if !ok {
		return infrastructure.ErrImportNotFound
	}
	for idx, job := range i.Jobs {
		if job.ID == j.ID {
			i.Jobs[idx] = j
			return nil
		}
	}
	i.Jobs = append(i.Jobs, j)

	return nil
}
