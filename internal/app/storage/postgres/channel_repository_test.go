package postgres_test

import (
	"context"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/storage/postgres"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestChannelRepository(t *testing.T) {
	suite.Run(t, new(ChannelRepositorySuite))
}

type ChannelRepositorySuite struct {
	testutils.IntegrationSuite
}

func (suite *ChannelRepositorySuite) Test_Save_New_Success() {
	// Prepare
	id := uuid.New()
	ch := channel.New(
		id,
		"Channel Name",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
	)
	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	err := r.Save(context.Background(), ch)

	// Assert result
	suite.NoError(err)

	// Assert state change
	var dbChannel postgres.Channel
	err = suite.DB.Get(&dbChannel, "SELECT * FROM channels WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbChannel.ID)
	suite.Equal("Channel Name", dbChannel.Name)
	suite.Equal(channel.IntegrationArbeitnow, dbChannel.Integration)
	suite.Equal(channel.StatusActive, dbChannel.Status)
	suite.True(dbChannel.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(dbChannel.UpdatedAt.After(time.Now().Add(-2 * time.Second)))
}

func (suite *ChannelRepositorySuite) Test_Save_Existing_Success() {
	// Prepare
	id := uuid.New()
	cAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id,
		"Channel Name",
		channel.IntegrationArbeitnow,
		channel.StatusInactive,
		cAt,
		cAt,
	)
	suite.NoError(err)

	uAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	ch := channel.New(
		id,
		"Channel Name new",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(cAt, uAt),
	)
	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	err = r.Save(context.Background(), ch)

	// Assert result
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM channels WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(1, count)

	var dbChannel postgres.Channel
	err = suite.DB.Get(&dbChannel, "SELECT * FROM channels WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbChannel.ID)
	suite.Equal("Channel Name new", dbChannel.Name)
	suite.Equal(channel.IntegrationArbeitnow, dbChannel.Integration)
	suite.Equal(channel.StatusActive, dbChannel.Status)
}

func (suite *ChannelRepositorySuite) Test_Save_Error() {
	// Prepare
	id := uuid.New()
	ch := channel.New(
		id,
		"Channel Name",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
	)
	r := postgres.NewChannelRepository(suite.BadDB)

	// Execute
	err := r.Save(context.Background(), ch)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}
