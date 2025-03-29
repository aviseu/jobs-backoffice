package testutils

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/google/uuid"
	"slices"
	"sync"
)

type ImportRepository struct {
	Imports    map[uuid.UUID]*postgres.Import
	JobResults map[uuid.UUID]*postgres.ImportJobResult
	err        error
	im         sync.Mutex
	jrm        sync.Mutex
}

func NewImportRepository() *ImportRepository {
	return &ImportRepository{
		Imports:    make(map[uuid.UUID]*postgres.Import),
		JobResults: make(map[uuid.UUID]*postgres.ImportJobResult),
	}
}

func (r *ImportRepository) First() *postgres.Import {
	for _, i := range r.Imports {
		return i
	}

	return nil
}

func (r *ImportRepository) Add(i *postgres.Import) {
	r.Imports[i.ID] = i
}

func (r *ImportRepository) AddResult(j *postgres.ImportJobResult) {
	r.JobResults[j.ID] = j
}

func (r *ImportRepository) FailWith(err error) {
	r.err = err
}

func (r *ImportRepository) SaveImport(_ context.Context, i *postgres.Import) error {
	if r.err != nil {
		return r.err
	}
	r.im.Lock()
	r.Imports[i.ID] = i
	r.im.Unlock()
	return nil
}

func (r *ImportRepository) FindImport(_ context.Context, id uuid.UUID) (*postgres.Import, error) {
	if r.err != nil {
		return nil, r.err
	}

	i, ok := r.Imports[id]
	if !ok {
		return nil, postgres.ErrImportNotFound
	}

	return i, nil
}

func (r *ImportRepository) GetImports(_ context.Context) ([]*postgres.Import, error) {
	if r.err != nil {
		return nil, r.err
	}

	var ii []*postgres.Import
	for _, i := range r.Imports {
		ii = append(ii, i)
	}

	slices.SortFunc(ii, func(a, b *postgres.Import) int {
		return b.StartedAt.Compare(a.StartedAt)
	})

	return ii, nil
}

func (r *ImportRepository) SaveImportJob(_ context.Context, jr *postgres.ImportJobResult) error {
	if r.err != nil {
		return r.err
	}

	r.jrm.Lock()
	r.JobResults[jr.ID] = jr
	r.jrm.Unlock()

	return nil
}

func (r *ImportRepository) GetJobsByImportID(_ context.Context, importID uuid.UUID) ([]*postgres.ImportJobResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	var results []*postgres.ImportJobResult
	for _, jr := range r.JobResults {
		if jr.ImportID == importID {
			results = append(results, jr)
		}
	}

	return results, nil
}
