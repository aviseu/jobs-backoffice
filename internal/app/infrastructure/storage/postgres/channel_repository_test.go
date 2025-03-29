package postgres_test

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestChannelRepository(t *testing.T) {
	suite.Run(t, new(ChannelRepositorySuite))
}

type ChannelRepositorySuite struct {
	testutils.PostgresSuite
}

func (suite *ChannelRepositorySuite) Test_Save_New_Success() {
	// Prepare
	id := uuid.New()
	ch := &aggregator.Channel{
		ID:          id,
		Name:        "Channel Name",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	err := r.Save(context.Background(), ch)

	// Assert result
	suite.NoError(err)

	// Assert state change
	var dbChannel aggregator.Channel
	err = suite.DB.Get(&dbChannel, "SELECT * FROM channels WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbChannel.ID)
	suite.Equal("Channel Name", dbChannel.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, dbChannel.Integration)
	suite.Equal(aggregator.ChannelStatusActive, dbChannel.Status)
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
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
		cAt,
		cAt,
	)
	suite.NoError(err)

	uAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	ch := &aggregator.Channel{
		ID:          id,
		Name:        "Channel Name new",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
		CreatedAt:   cAt,
		UpdatedAt:   uAt,
	}
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

	var dbChannel aggregator.Channel
	err = suite.DB.Get(&dbChannel, "SELECT * FROM channels WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbChannel.ID)
	suite.Equal("Channel Name new", dbChannel.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, dbChannel.Integration)
	suite.Equal(aggregator.ChannelStatusActive, dbChannel.Status)
}

func (suite *ChannelRepositorySuite) Test_Save_Error() {
	// Prepare
	id := uuid.New()
	ch := &aggregator.Channel{
		ID:          id,
		Name:        "Channel Name",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	r := postgres.NewChannelRepository(suite.BadDB)

	// Execute
	err := r.Save(context.Background(), ch)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ChannelRepositorySuite) Test_All_Success() {
	// Prepare
	id1 := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id1,
		"Channel Name 1",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusActive,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
	)
	suite.NoError(err)

	id2 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id2,
		"Channel Name 2",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
		time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC),
	)
	suite.NoError(err)

	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	chs, err := r.All(context.Background())

	// Assert result
	suite.NoError(err)
	suite.Len(chs, 2)

	suite.Equal(id1, chs[0].ID)
	suite.Equal("Channel Name 1", chs[0].Name)
	suite.Equal(aggregator.IntegrationArbeitnow, chs[0].Integration)
	suite.Equal(aggregator.ChannelStatusActive, chs[0].Status)
	suite.True(chs[0].CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(chs[0].UpdatedAt.Equal(time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)))

	suite.Equal(id2, chs[1].ID)
	suite.Equal("Channel Name 2", chs[1].Name)
	suite.Equal(aggregator.IntegrationArbeitnow, chs[1].Integration)
	suite.Equal(aggregator.ChannelStatusInactive, chs[1].Status)
	suite.True(chs[1].CreatedAt.Equal(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC)))
	suite.True(chs[1].UpdatedAt.Equal(time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC)))
}

func (suite *ChannelRepositorySuite) Test_All_Error() {
	// Prepare
	r := postgres.NewChannelRepository(suite.BadDB)

	// Execute
	chs, err := r.All(context.Background())

	// Assert
	suite.Nil(chs)
	suite.Error(err)
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ChannelRepositorySuite) Test_GetActive_Success() {
	// Prepare
	id1 := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id1,
		"Channel Name 1",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusActive,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
	)
	suite.NoError(err)

	id2 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id2,
		"Channel Name 2",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
		time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC),
	)
	suite.NoError(err)

	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	chs, err := r.GetActive(context.Background())

	// Assert result
	suite.NoError(err)
	suite.Len(chs, 1)
	suite.Equal(id1, chs[0].ID)
}

func (suite *ChannelRepositorySuite) Test_GetActive_Error() {
	// Prepare
	r := postgres.NewChannelRepository(suite.BadDB)

	// Execute
	chs, err := r.GetActive(context.Background())

	// Assert
	suite.Nil(chs)
	suite.Error(err)
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ChannelRepositorySuite) Test_Find_Success() {
	// Prepare
	id := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id,
		"Channel Name",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusActive,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
	)
	suite.NoError(err)

	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	ch, err := r.Find(context.Background(), id)

	// Assert result
	suite.NoError(err)

	// Assert state change

	suite.Equal(id, ch.ID)
	suite.Equal("Channel Name", ch.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration)
	suite.Equal(aggregator.ChannelStatusActive, ch.Status)
	suite.True(ch.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(ch.UpdatedAt.Equal(time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)))
}

func (suite *ChannelRepositorySuite) Test_Find_NotFound() {
	// Prepare
	id := uuid.New()
	r := postgres.NewChannelRepository(suite.DB)

	// Execute
	ch, err := r.Find(context.Background(), id)

	// Assert
	suite.Nil(ch)
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorIs(err, infrastructure.ErrChannelNotFound)
}

func (suite *ChannelRepositorySuite) Test_Find_Error() {
	// Prepare
	id := uuid.New()
	r := postgres.NewChannelRepository(suite.BadDB)

	// Execute
	ch, err := r.Find(context.Background(), id)

	// Assert
	suite.Nil(ch)
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}
