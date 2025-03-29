package configuring_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/errs"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
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
	s := configuring.NewService(r)
	cmd := configuring.NewCreateChannelCommand(
		"Channel Name",
		"arbeitnow",
	)

	// Execute
	ch, err := s.Create(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal("Channel Name", ch.Name())
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration())
	suite.Equal(aggregator.ChannelStatusInactive, ch.Status())
	suite.True(ch.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.Len(r.Channels, 1)
}

func (suite *ServiceSuite) Test_Create_Validation_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	cmd := configuring.NewCreateChannelCommand(
		"",
		"bad_integration",
	)

	// Execute
	_, err := s.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrNameIsRequired)
	suite.ErrorIs(err, configuring.ErrInvalidIntegration)
	suite.ErrorContains(err, "bad_integration")
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Create_RepositoryFail_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	cmd := configuring.NewCreateChannelCommand(
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

func (suite *ServiceSuite) Test_Update_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	ch := configuring.NewChannel(id, "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive, configuring.WithTimestamps(cat, uat))
	r.Add(ch.ToDTO())
	cmd := configuring.NewUpdateChannelCommand(ch.ID(), "channel 2")

	// Execute
	res, err := s.Update(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal(id, res.ID())
	suite.Equal("channel 2", res.Name())
	suite.Equal(aggregator.IntegrationArbeitnow, res.Integration())
	suite.Equal(aggregator.ChannelStatusActive, res.Status())
	suite.True(cat.Equal(res.CreatedAt()))
	suite.True(res.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.NoError(err)
	suite.Equal(id, res.ID())
	suite.Equal("channel 2", r.First().Name)
	suite.Equal(aggregator.IntegrationArbeitnow, r.First().Integration)
	suite.Equal(aggregator.ChannelStatusActive, r.First().Status)
	suite.True(cat.Equal(r.First().CreatedAt))
	suite.True(res.UpdatedAt().Equal(r.First().UpdatedAt))
}

func (suite *ServiceSuite) Test_Update_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	cmd := configuring.NewUpdateChannelCommand(uuid.New(), "channel 2")

	// Execute
	_, err := s.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	r.Add(ch.ToDTO())
	cmd := configuring.NewUpdateChannelCommand(ch.ID(), "channel 2")

	// Execute
	_, err := s.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Validation_Fail() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	r.Add(ch.ToDTO())
	cmd := configuring.NewUpdateChannelCommand(ch.ID(), "")

	// Execute
	_, err := s.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrNameIsRequired)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusInactive)
	r.Add(ch.ToDTO())

	// Execute
	err := s.Activate(context.Background(), ch.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(aggregator.ChannelStatusActive, r.First().Status)
}

func (suite *ServiceSuite) Test_Activate_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)

	// Execute
	err := s.Activate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusInactive)
	r.Add(ch.ToDTO())

	// Execute
	err := s.Activate(context.Background(), ch.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Deactivate_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	r.Add(ch.ToDTO())

	// Execute
	err := s.Deactivate(context.Background(), ch.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(aggregator.ChannelStatusInactive, r.First().Status)
}

func (suite *ServiceSuite) Test_Deactivate_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)

	// Execute
	err := s.Deactivate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Deactivate_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	r.Add(ch.ToDTO())

	// Execute
	err := s.Deactivate(context.Background(), ch.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}
