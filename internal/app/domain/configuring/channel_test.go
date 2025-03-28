package configuring_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	suite.Run(t, new(ChannelSuite))
}

type ChannelSuite struct {
	suite.Suite
}

func (suite *ChannelSuite) Test_Success() {
	// Prepare
	id := uuid.New()

	// Execute
	ch := configuring.New(
		id,
		"Channel Name",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
	)

	// Assert
	suite.Equal(id, ch.ID())
	suite.Equal("Channel Name", ch.Name())
	suite.Equal("arbeitnow", ch.Integration().String())
	suite.Equal("active", ch.Status().String())
	suite.True(ch.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt().After(time.Now().Add(-2 * time.Second)))
}

func (suite *ChannelSuite) Test_WithTimestamps_Success() {
	// Prepare
	id := uuid.New()
	cAt := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)

	// Execute
	ch := configuring.New(
		id,
		"Channel Name",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(cAt, uAt),
	)

	// Assert
	suite.Equal(id, ch.ID())
	suite.Equal("Channel Name", ch.Name())
	suite.Equal("arbeitnow", ch.Integration().String())
	suite.Equal("active", ch.Status().String())
	suite.True(ch.CreatedAt().Equal(cAt))
	suite.True(ch.UpdatedAt().Equal(uAt))
}
