package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChannelRepository struct {
	db *sqlx.DB
}

func NewChannelRepository(db *sqlx.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Save(ctx context.Context, ch *aggregator.Channel) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO channels (id, name, integration, status, created_at, updated_at)
				VALUES (:id, :name, :integration, :status, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET
					name = EXCLUDED.name,
					integration = EXCLUDED.integration,
					status = EXCLUDED.status,
					updated_at = EXCLUDED.updated_at`,
		ch,
	)
	if err != nil {
		return fmt.Errorf("failed to save channel %s: %w", ch.ID, err)
	}

	return nil
}

func (r *ChannelRepository) All(ctx context.Context) ([]*aggregator.Channel, error) {
	var result []*aggregator.Channel
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM channels ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to get all channels: %w", err)
	}

	return result, nil
}

func (r *ChannelRepository) GetActive(ctx context.Context) ([]*aggregator.Channel, error) {
	var result []*aggregator.Channel
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM channels WHERE status = $1 ORDER BY name", aggregator.ChannelStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get active channels: %w", err)
	}

	return result, nil
}

func (r *ChannelRepository) Find(ctx context.Context, id uuid.UUID) (*aggregator.Channel, error) {
	var c aggregator.Channel
	err := r.db.GetContext(ctx, &c, "SELECT * FROM channels WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find channel %s: %w", id, infrastructure.ErrChannelNotFound)
		}

		return nil, fmt.Errorf("failed to find channel %s: %w", id, err)
	}

	return &c, nil
}
