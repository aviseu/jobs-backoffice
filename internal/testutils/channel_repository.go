package testutils

import (
	"context"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/google/uuid"
)

type ChannelRepository struct {
	channels map[uuid.UUID]*channel.Channel
	err      error
}

func NewChannelRepository() *ChannelRepository {
	return &ChannelRepository{channels: make(map[uuid.UUID]*channel.Channel)}
}

func (r *ChannelRepository) FailWith(err error) {
	r.err = err
}

func (r *ChannelRepository) Save(_ context.Context, ch *channel.Channel) error {
	if r.err != nil {
		return r.err
	}

	r.channels[ch.ID()] = ch
	return nil
}
