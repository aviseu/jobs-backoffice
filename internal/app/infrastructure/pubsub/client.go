package pubsub

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"
)

type client struct {
	topic   *pubsub.Topic
	timeout time.Duration
}

func newClient(topic *pubsub.Topic, timeout time.Duration) *client {
	return &client{
		topic:   topic,
		timeout: timeout,
	}
}

func (c *client) publish(ctx context.Context, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	result := c.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish pubsub message: %w", err)
	}

	return nil
}
