package testutils

import (
	"cmp"
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"slices"
)

type ChannelRepository struct {
	Channels map[uuid.UUID]*aggregator.Channel
	err      error
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{
		Channels: make(map[uuid.UUID]*aggregator.Channel),
	}
}

func (r *ChannelRepository) First() *aggregator.Channel {
	for _, ch := range r.Channels {
		return ch
	}

	return nil
}

func (r *ChannelRepository) Add(ch *aggregator.Channel) {
	r.Channels[ch.ID] = ch
}

func (r *ChannelRepository) FailWith(err error) {
	r.err = err
}

func (r *ChannelRepository) All(_ context.Context) ([]*aggregator.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	channels := make([]*aggregator.Channel, 0, len(r.Channels))
	for _, ch := range r.Channels {
		channels = append(channels, ch)
	}

	slices.SortFunc(channels, func(a, b *aggregator.Channel) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return channels, nil
}

func (r *ChannelRepository) GetActive(_ context.Context) ([]*aggregator.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	channels := make([]*aggregator.Channel, 0)
	for _, ch := range r.Channels {
		if ch.Status == aggregator.ChannelStatusActive {
			channels = append(channels, ch)
		}
	}

	slices.SortFunc(channels, func(a, b *aggregator.Channel) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return channels, nil

}

func (r *ChannelRepository) Find(_ context.Context, id uuid.UUID) (*aggregator.Channel, error) {
	if r.err != nil {
		return nil, r.err
	}

	ch, ok := r.Channels[id]
	if !ok {
		return nil, infrastructure.ErrChannelNotFound
	}

	return ch, nil
}

func (r *ChannelRepository) Save(_ context.Context, ch *aggregator.Channel) error {
	if r.err != nil {
		return r.err
	}

	r.Channels[ch.ID] = ch
	return nil
}
