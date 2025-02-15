package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/google/uuid"
	"sync"
)

type ImportRepository struct {
	Imports    map[uuid.UUID]*imports.Import
	JobResults map[uuid.UUID]*imports.JobResult
	err        error
	im         sync.Mutex
	jrm        sync.Mutex
}

func NewImportRepository() *ImportRepository {
	return &ImportRepository{
		Imports:    make(map[uuid.UUID]*imports.Import),
		JobResults: make(map[uuid.UUID]*imports.JobResult),
	}
}

func (r *ImportRepository) First() *imports.Import {
	for _, i := range r.Imports {
		return i
	}

	return nil
}

func (r *ImportRepository) FailWith(err error) {
	r.err = err
}

func (r *ImportRepository) SaveImport(_ context.Context, i *imports.Import) error {
	if r.err != nil {
		return r.err
	}
	r.im.Lock()
	r.Imports[i.ID()] = i
	r.im.Unlock()
	return nil
}

func (r *ImportRepository) FindImport(_ context.Context, id uuid.UUID) (*imports.Import, error) {
	if r.err != nil {
		return nil, r.err
	}

	i, ok := r.Imports[id]
	if !ok {
		return nil, imports.ErrImportNotFound
	}

	return i, nil
}

func (r *ImportRepository) SaveImportJob(_ context.Context, jr *imports.JobResult) error {
	if r.err != nil {
		return r.err
	}

	r.jrm.Lock()
	r.JobResults[jr.JobID()] = jr
	r.jrm.Unlock()

	return nil
}

func (r *ImportRepository) GetJobsByImportID(_ context.Context, importID uuid.UUID) ([]*imports.JobResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	var results []*imports.JobResult
	for _, jr := range r.JobResults {
		if jr.ImportID() == importID {
			results = append(results, jr)
		}
	}

	return results, nil
}
