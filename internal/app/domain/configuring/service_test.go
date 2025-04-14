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
	dsl := testutils.NewDSL()
	cmd := configuring.NewCreateChannelCommand(
		"Channel Name",
		"arbeitnow",
	)

	// Execute
	ch, err := dsl.ConfiguringService.Create(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal("Channel Name", ch.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration)
	suite.Equal(aggregator.ChannelStatusInactive, ch.Status)
	suite.True(ch.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	suite.Equal("Channel Name", dsl.FirstChannel().Name)
	suite.Equal(aggregator.IntegrationArbeitnow, dsl.FirstChannel().Integration)
}

func (suite *ServiceSuite) Test_Create_Validation_Fail() {
	// Prepare
	dsl := testutils.NewDSL()
	cmd := configuring.NewCreateChannelCommand(
		"",
		"bad_integration",
	)

	// Execute
	_, err := dsl.ConfiguringService.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrNameIsRequired)
	suite.ErrorIs(err, configuring.ErrInvalidIntegration)
	suite.ErrorContains(err, "bad_integration")
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Create_RepositoryFail_Fail() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)
	cmd := configuring.NewCreateChannelCommand(
		"Channel Name",
		"arbeitnow",
	)

	// Execute
	_, err := dsl.ConfiguringService.Create(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Success() {
	// Prepare
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)

	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelActivated(),
			testutils.WithChannelTimestamps(cat, uat),
		),
	)
	cmd := configuring.NewUpdateChannelCommand(id, "channel 2")

	// Execute
	res, err := dsl.ConfiguringService.Update(context.Background(), cmd)

	// Assert result
	suite.NoError(err)
	suite.Equal(id, res.ID)
	suite.Equal("channel 2", res.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, res.Integration)
	suite.Equal(aggregator.ChannelStatusActive, res.Status)
	suite.True(cat.Equal(res.CreatedAt))
	suite.True(res.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert state change
	suite.NoError(err)
	suite.Equal(id, res.ID)
	ch := dsl.FirstChannel()
	suite.Equal("channel 2", ch.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration)
	suite.Equal(aggregator.ChannelStatusActive, ch.Status)
	suite.True(cat.Equal(ch.CreatedAt))
	suite.True(res.UpdatedAt.Equal(ch.UpdatedAt))
}

func (suite *ServiceSuite) Test_Update_NotFound() {
	// Prepare
	dsl := testutils.NewDSL()
	cmd := configuring.NewUpdateChannelCommand(uuid.New(), "channel 2")

	// Execute
	_, err := dsl.ConfiguringService.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_ChannelRepositoryFail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)
	cmd := configuring.NewUpdateChannelCommand(id, "channel 2")

	// Execute
	_, err := dsl.ConfiguringService.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Update_Validation_Fail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
		),
	)
	cmd := configuring.NewUpdateChannelCommand(id, "")

	// Execute
	_, err := dsl.ConfiguringService.Update(context.Background(), cmd)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrNameIsRequired)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_Success() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelDeactivated(),
		),
	)

	// Execute
	err := dsl.ConfiguringService.Activate(context.Background(), id)

	// Assert
	suite.NoError(err)
	suite.Equal(aggregator.ChannelStatusActive, dsl.FirstChannel().Status)
}

func (suite *ServiceSuite) Test_Activate_NotFound() {
	// Prepare
	dsl := testutils.NewDSL()

	// Execute
	err := dsl.ConfiguringService.Activate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Activate_ChannelRepositoryFail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelDeactivated(),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	// Execute
	err := dsl.ConfiguringService.Activate(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Deactivate_Success() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
	)

	// Execute
	err := dsl.ConfiguringService.Deactivate(context.Background(), id)

	// Assert
	suite.NoError(err)
	suite.Equal(aggregator.ChannelStatusInactive, dsl.FirstChannel().Status)
}

func (suite *ServiceSuite) Test_Deactivate_NotFound() {
	// Prepare
	dsl := testutils.NewDSL()

	// Execute
	err := dsl.ConfiguringService.Deactivate(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, configuring.ErrChannelNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_Deactivate_ChannelRepositoryFail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	// Execute
	err := dsl.ConfiguringService.Deactivate(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.False(errs.IsValidationError(err))
}
