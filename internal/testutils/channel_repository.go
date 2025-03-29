package testutils

import (
	"cmp"
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/google/uuid"
	"slices"
)

type ChannelRepository struct {
	Channels map[uuid.UUID]*postgres.Channel
	err      error
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		Channels: make(map[uuid.UUID]*postgres.Channel),
	}
}

func (r *ChannelRepository) First() *postgres.Channel {
	for _, ch := range r.Channels {
		return ch
	}

	return nil
}

func (r *ChannelRepository) Add(ch *postgres.Channel) {
	r.Channels[ch.ID] = ch
}

func (r *ChannelRepository) FailWith(err error) {
	r.err = err
}

func (r *ChannelRepository) All(_ context.Context) ([]*postgres.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	channels := make([]*postgres.Channel, 0, len(r.Channels))
	for _, ch := range r.Channels {
		channels = append(channels, ch)
	}

	slices.SortFunc(channels, func(a, b *postgres.Channel) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return channels, nil
}

func (r *ChannelRepository) GetActive(_ context.Context) ([]*postgres.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	channels := make([]*postgres.Channel, 0)
	for _, ch := range r.Channels {
		if ch.Status == int(base.ChannelStatusActive) {
			channels = append(channels, ch)
		}
	}

	slices.SortFunc(channels, func(a, b *postgres.Channel) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return channels, nil

}

func (r *ChannelRepository) Find(_ context.Context, id uuid.UUID) (*postgres.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	ch, ok := r.Channels[id]
	if !ok {
		return nil, postgres.ErrChannelNotFound
	}

	return ch, nil
}

func (r *ChannelRepository) Save(_ context.Context, ch *postgres.Channel) error {
	if r.err != nil {
		return r.err
	}

	r.Channels[ch.ID] = ch
	return nil
}
