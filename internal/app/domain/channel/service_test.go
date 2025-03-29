package channel_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
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
	suite.Equal(base.IntegrationArbeitnow, ch.Integration())
	suite.Equal(base.ChannelStatusInactive, ch.Status())
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

	ch1 := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch1.DTO())
	ch2 := channel.New(uuid.New(), "channel 2", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch2.DTO())

	// Execute
	chs, err := s.All(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(chs, 2)
	suite.Equal(ch1, chs[0])
	suite.Equal(ch2, chs[1])
}

func (suite *ServiceSuite) Test_All_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)

	// Execute
	_, err := s.All(context.Background())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_GetActive_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	ch1 := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch1.DTO())
	ch2 := channel.New(uuid.New(), "channel 2", base.IntegrationArbeitnow, base.ChannelStatusInactive)
	r.Add(ch2.DTO())

	// Execute
	chs, err := s.GetActive(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(chs, 1)
	suite.Equal(ch1, chs[0])
}

func (suite *ServiceSuite) Test_GetActive_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)

	// Execute
	_, err := s.GetActive(context.Background())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Find_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())

	// Execute
	ch2, err := s.Find(context.Background(), ch.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(ch, ch2)
}

func (suite *ServiceSuite) Test_Find_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	// Execute
	_, err := s.Find(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Find_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())

	// Execute
	_, err := s.Find(context.Background(), ch.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	ch := channel.New(id, "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive, channel.WithTimestamps(cat, uat))
	r.Add(ch.DTO())
	cmd := channel.NewUpdateCommand(ch.ID(), "channel 2")

	// Execute
	res, err := s.Update(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal(id, res.ID())
	suite.Equal("channel 2", res.Name())
	suite.Equal(base.IntegrationArbeitnow, res.Integration())
	suite.Equal(base.ChannelStatusActive, res.Status())
	suite.True(cat.Equal(res.CreatedAt()))
	suite.True(res.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.NoError(err)
	suite.Equal(id, res.ID())
	suite.Equal("channel 2", r.First().Name)
	suite.Equal(int(base.IntegrationArbeitnow), r.First().Integration)
	suite.Equal(int(base.ChannelStatusActive), r.First().Status)
	suite.True(cat.Equal(r.First().CreatedAt))
	suite.True(res.UpdatedAt().Equal(r.First().UpdatedAt))
}

func (suite *ServiceSuite) Test_Update_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	cmd := channel.NewUpdateCommand(uuid.New(), "channel 2")

	// Execute
	_, err := s.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())
	cmd := channel.NewUpdateCommand(ch.ID(), "channel 2")

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
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())
	cmd := channel.NewUpdateCommand(ch.ID(), "")

	// Execute
	_, err := s.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrNameIsRequired)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusInactive)
	r.Add(ch.DTO())

	// Execute
	err := s.Activate(context.Background(), ch.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(int(base.ChannelStatusActive), r.First().Status)
}

func (suite *ServiceSuite) Test_Activate_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	// Execute
	err := s.Activate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusInactive)
	r.Add(ch.DTO())

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
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())

	// Execute
	err := s.Deactivate(context.Background(), ch.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(int(base.ChannelStatusInactive), r.First().Status)
}

func (suite *ServiceSuite) Test_Deactivate_NotFound() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	// Execute
	err := s.Deactivate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, channel.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Deactivate_Error() {
	// Prepare
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	ch := channel.New(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive)
	r.Add(ch.DTO())

	// Execute
	err := s.Deactivate(context.Background(), ch.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Integrations_Success() {
	// Prepare
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)

	// Execute
	integrations := s.Integrations()

	// Assert
	suite.Len(integrations, 1)
	suite.Contains(integrations, base.IntegrationArbeitnow)
}
