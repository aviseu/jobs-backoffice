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
	m       sync.Mutex
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

func (r *ImportRepository) ImportMetrics() map[uuid.UUID]*aggregator.ImportMetric {
	metrics := make(map[uuid.UUID]*aggregator.ImportMetric)
	for _, i := range r.Imports {
		for _, m := range i.Metrics {
			metrics[m.ID] = m
		}
	}
	return metrics
}

func (r *ImportRepository) AddImport(i *aggregator.Import) {
	r.Imports[i.ID] = i
}

func (r *ImportRepository) AddImportMetric(importID uuid.UUID, m *aggregator.ImportMetric) {
	i, ok := r.Imports[importID]
	if !ok {
		panic("import not found")
	}
	for idx, metric := range i.Metrics {
		if metric.ID == m.ID {
			i.Metrics[idx] = m
			return
		}
	}
	i.Metrics = append(i.Metrics, m)
}

func (r *ImportRepository) FailWith(err error) {
	r.err = err
}

func (r *ImportRepository) SaveImport(_ context.Context, i *aggregator.Import) error {
	if r.err != nil {
		return r.err
	}
	r.m.Lock()
	defer r.m.Unlock()
	old, ok := r.Imports[i.ID]

	if ok {
		for _, metric := range old.Metrics {
			for _, newMetric := range i.Metrics {
				if metric.ID == newMetric.ID {
					continue
				}
			}
			i.Metrics = append(i.Metrics, metric)
		}
	}
	r.Imports[i.ID] = i
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

func (r *ImportRepository) SaveImportMetric(_ context.Context, importID uuid.UUID, m *aggregator.ImportMetric) error {
	if r.err != nil {
		return r.err
	}

	r.m.Lock()
	defer r.m.Unlock()

	i, ok := r.Imports[importID]
	if !ok {
		return infrastructure.ErrImportNotFound
	}
	for idx, metric := range i.Metrics {
		if metric.ID == m.ID {
			i.Metrics[idx] = m
			return nil
		}
	}
	i.Metrics = append(i.Metrics, m)

	return nil
}
