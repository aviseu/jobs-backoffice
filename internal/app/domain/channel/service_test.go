package channel_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) Test_Create_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	cmd := channel.NewCreateCommand(
		"Channel Name",
		"arbeitnow",
	)

	// Execute
	ch, err := s.Create(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal("Channel Name", ch.Name())
	suite.Equal(channel.IntegrationArbeitnow, ch.Integration())
	suite.Equal(channel.StatusInactive, ch.Status())
	suite.True(ch.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.Len(r.Channels, 1)
}

func (suite *ServiceSuite) Test_Create_InvalidIntegration_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	cmd := channel.NewCreateCommand(
		"Channel Name",
		"bad_integration",
	)

	// Execute
	_, err := s.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrInvalidIntegration)
	suite.ErrorContains(err, "bad_integration")
}

func (suite *ServiceSuite) Test_Create_NameIsRequired_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	cmd := channel.NewCreateCommand(
		"",
		"arbeitnow",
	)

	// Execute
	_, err := s.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrNameIsRequired)
}

func (suite *ServiceSuite) Test_Create_RepositoryFail_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	cmd := channel.NewCreateCommand(
		"Channel Name",
		"arbeitnow",
	)

	// Execute
	_, err := s.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
}
