package testutils

import (
	"context"
	"github.com/google/uuid"
)

type PubSubImportService struct {
	ImportIDs []uuid.UUID
	err       error
}

func NewPubSubImportService() *PubSubImportService {
	return &PubSubImportService{}
}

func (p *PubSubImportService) FailWith(err error) {
	p.err = err
}

func (p *PubSubImportService) PublishImportCommand(_ context.Context, importID uuid.UUID) error {
	if p.err != nil {
		return p.err
	}

	p.ImportIDs = append(p.ImportIDs, importID)
	return nil
}
