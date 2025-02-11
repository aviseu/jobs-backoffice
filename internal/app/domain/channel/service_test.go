package channel_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/errs"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
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

func (suite *ServiceSuite) Test_Create_Validation_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	cmd := channel.NewCreateCommand(
		"",
		"bad_integration",
	)

	// Execute
	_, err := s.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrNameIsRequired)
	suite.ErrorIs(err, channel.ErrInvalidIntegration)
	suite.ErrorContains(err, "bad_integration")
	suite.True(errs.IsValidationError(err))
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
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_All_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	ch1 := channel.New(uuid.New(), "channel 1", channel.IntegrationArbeitnow, channel.StatusActive)
	r.Add(ch1)
	ch2 := channel.New(uuid.New(), "channel 2", channel.IntegrationArbeitnow, channel.StatusActive)
	r.Add(ch2)

	// Execute
	chs, err := s.All(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(chs, 2)
	suite.Equal(ch1, chs[0])
	suite.Equal(ch2, chs[1])
}
