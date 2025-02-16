package testutils

import (
	"context"
	"github.com/google/uuid"
)

type PubSubService struct {
	ImportIDs []uuid.UUID
	err       error
}

func NewPubSubService() *PubSubService {
	return &PubSubService{}
}

func (p *PubSubService) FailWith(err error) {
	p.err = err
}

func (p *PubSubService) PublishImportCommand(_ context.Context, importID uuid.UUID) error {
	if p.err != nil {
		return p.err
	}

	p.ImportIDs = append(p.ImportIDs, importID)
	return nil
}
