package scheduling_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) Test_ScheduleImport_Success() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelActivated(),
		),
	)
	ch := dsl.FirstChannel()

	// Execute
	i, err := dsl.SchedulingService.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.NoError(err)
	suite.NotNil(i)
	suite.Equal(ch.ID, i.ChannelID)
	suite.Equal(aggregator.ImportStatusPending, i.Status)
	suite.True(i.StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(i.EndedAt.Valid)

	// Assert state change
	suite.Len(dsl.Imports(), 1)
	dbImport := dsl.FirstImport()
	suite.Equal(ch.ID, dbImport.ChannelID)
	suite.Equal(aggregator.ImportStatusPending, dbImport.Status)
	suite.True(dbImport.StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(dbImport.EndedAt.Valid)

	// Assert pubsub message
	suite.Len(dsl.PublishedImports(), 1)
	suite.Equal(i.ID, dsl.PublishedImports()[0])

	// Assert log
	logs := dsl.LogLines()
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleImport_ImportRepositoryFailed() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelActivated(),
		),
		testutils.WithImportRepositoryError(errors.New("boom")),
	)
	ch := dsl.FirstChannel()

	// Execute
	i, err := dsl.SchedulingService.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.ErrorContains(err, ch.ID.String())

	// Assert state change
	suite.Len(dsl.Imports(), 0)

	// Assert pubsub message
	suite.Len(dsl.PublishedImports(), 0)

	// Assert log
	logs := dsl.LogLines()
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleImport_PubSubFailed() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelActivated(),
		),
		testutils.WithPubSubServiceError(errors.New("boom")),
	)
	ch := dsl.FirstChannel()

	// Execute
	i, err := dsl.SchedulingService.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.ErrorContains(err, ch.ID.String())

	// Assert state change
	suite.Len(dsl.Imports(), 1)
	dbImport := dsl.FirstImport()
	suite.Equal(ch.ID, dbImport.ChannelID)
	suite.Equal(aggregator.ImportStatusPending, dbImport.Status)
	suite.True(dbImport.StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(dbImport.EndedAt.Valid)

	// Assert pubsub message
	suite.Len(dsl.PublishedImports(), 0)

	// Assert log
	logs := dsl.LogLines()
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_Success() {
	// Prepare
	id1 := uuid.New()
	id2 := uuid.New()

	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id1),
			testutils.WithChannelActivated(),
		),
		testutils.WithChannel(
			testutils.WithChannelID(id2),
			testutils.WithChannelActivated(),
		),
		testutils.WithChannel(
			testutils.WithChannelDeactivated(),
		),
	)

	// Execute
	err := dsl.SchedulingService.ScheduleActiveChannels(context.Background())

	// Assert
	suite.NoError(err)

	// Assert imports created
	suite.Equal(2, len(dsl.Imports()))
	var i1, i2 *aggregator.Import
	for _, i := range dsl.Imports() {
		if i.ChannelID == id1 {
			i1 = i
		}
		if i.ChannelID == id2 {
			i2 = i
		}
	}
	suite.NotNil(i1)
	suite.NotNil(i2)

	// Assert correct values published
	publishedImportIDs := dsl.PublishedImports()
	suite.Len(publishedImportIDs, 2)
	i1Exists := false
	for _, id := range publishedImportIDs {
		if id == i1.ID {
			i1Exists = true
		}
	}
	suite.True(i1Exists)
	i2Exists := false
	for _, id := range publishedImportIDs {
		if id == i2.ID {
			i2Exists = true
		}
	}
	suite.True(i2Exists)

	// Assert logs
	lines := dsl.LogLines()
	suite.Len(lines, 2)
	id1Logged := false
	id2Logged := false
	for _, line := range lines {
		if strings.Contains(line, "scheduling import for channel "+id1.String()) {
			id1Logged = true
		}
		if strings.Contains(line, "scheduling import for channel "+id2.String()) {
			id2Logged = true
		}
	}
	suite.True(id1Logged)
	suite.True(id2Logged)
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_ChannelRepositoryFail() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelActivated(),
		),
		testutils.WithChannel(
			testutils.WithChannelActivated(),
		),
		testutils.WithChannel(
			testutils.WithChannelDeactivated(),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	// Execute
	err := dsl.SchedulingService.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to fetch active channels")
	suite.ErrorContains(err, "boom")

	// Assert logs
	suite.Len(dsl.LogLines(), 0)
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_ImportRepositoryFail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
		testutils.WithImportRepositoryError(errors.New("boom")),
	)

	// Execute
	err := dsl.SchedulingService.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to save import for channel ")
	suite.Contains(err.Error(), id.String())
	suite.ErrorContains(err, "boom")

	// Assert logs
	logs := dsl.LogLines()
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+id.String())
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_PubSubServiceFail() {
	// Prepare
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
		testutils.WithPubSubServiceError(errors.New("boom")),
	)

	// Execute
	err := dsl.SchedulingService.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to publish import ")
	suite.Contains(err.Error(), id.String())
	suite.ErrorContains(err, "boom")

	// Assert logs
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())
}
