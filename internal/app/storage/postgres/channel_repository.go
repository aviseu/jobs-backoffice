package postgres

import (
	"context"
	"fmt"

	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/jmoiron/sqlx"
)

type ChannelRepository struct {
	db *sqlx.DB
}

func NewChannelRepository(db *sqlx.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Save(ctx context.Context, ch *channel.Channel) error {
	c := fromDomainChannel(ch)
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO channels (id, name, integration, status, created_at, updated_at)
				VALUES (:id, :name, :integration, :status, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET
					name = EXCLUDED.name,
					integration = EXCLUDED.integration,
					status = EXCLUDED.status,
					updated_at = EXCLUDED.updated_at`,
		c,
	)
	if err != nil {
		return fmt.Errorf("failed to save channel %s: %w", ch.ID(), err)
	}

	return nil
}

func (r *ChannelRepository) All(ctx context.Context) ([]*channel.Channel, error) {
	var cc []*Channel
	err := r.db.SelectContext(ctx, &cc, "SELECT * FROM channels ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to get all channels: %w", err)
	}

	result := make([]*channel.Channel, 0, len(cc))
	for _, c := range cc {
		result = append(result, toDomainChannel(c))
	}

	return result, nil
}
