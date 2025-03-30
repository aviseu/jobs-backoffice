package scheduling_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/scheduling"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
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

func (suite *ServiceSuite) Test_ScheduleImport_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	is := scheduling.NewService(ir, chr, ps, log)

	ch := &aggregator.Channel{
		ID:          uuid.New(),
		Name:        "airbeitnow test",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}

	// Execute
	i, err := is.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.NoError(err)
	suite.NotNil(i)
	suite.Equal(ch.ID, i.ChannelID)
	suite.Equal(aggregator.ImportStatusPending, i.Status)
	suite.True(i.StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(i.EndedAt.Valid)

	// Assert state change
	suite.Len(ir.Imports, 1)
	suite.Equal(ch.ID, ir.First().ChannelID)
	suite.Equal(aggregator.ImportStatusPending, ir.First().Status)
	suite.True(ir.First().StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(ir.First().EndedAt.Valid)

	// Assert pubsub message
	suite.Len(ps.ImportIDs, 1)
	suite.Equal(i.ID, ps.ImportIDs[0])

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleImport_ImportRepositoryFailed() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom"))
	ps := testutils.NewPubSubService()
	is := scheduling.NewService(ir, chr, ps, log)

	ch := &aggregator.Channel{
		ID:          uuid.New(),
		Name:        "airbeitnow test",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}

	// Execute
	i, err := is.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.ErrorContains(err, ch.ID.String())

	// Assert state change
	suite.Len(ir.Imports, 0)

	// Assert pubsub message
	suite.Len(ps.ImportIDs, 0)

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleImport_PubSubFailed() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	ps.FailWith(errors.New("boom"))
	is := scheduling.NewService(ir, chr, ps, log)

	ch := &aggregator.Channel{
		ID:          uuid.New(),
		Name:        "airbeitnow test",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}

	// Execute
	i, err := is.ScheduleImport(context.Background(), ch)

	// Assert return
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom")
	suite.ErrorContains(err, ch.ID.String())

	// Assert state change
	suite.Len(ir.Imports, 1)
	suite.Equal(ch.ID, ir.First().ChannelID)
	suite.Equal(aggregator.ImportStatusPending, ir.First().Status)
	suite.True(ir.First().StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(ir.First().EndedAt.Valid)

	// Assert pubsub message
	suite.Len(ps.ImportIDs, 0)

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+ch.ID.String())
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	is := scheduling.NewService(ir, chr, ps, log)

	id1 := uuid.New()
	chr.Add(configuring.NewChannel(id1, "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive).ToAggregator())
	id2 := uuid.New()
	chr.Add(configuring.NewChannel(id2, "channel 2", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive).ToAggregator())
	chr.Add(configuring.NewChannel(uuid.New(), "channel 3", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusInactive).ToAggregator())

	// Execute
	err := is.ScheduleActiveChannels(context.Background())

	// Assert
	suite.NoError(err)

	// Assert imports created
	suite.Equal(2, len(ir.Imports))
	var i1, i2 *aggregator.Import
	for _, i := range ir.Imports {
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
	suite.Len(ps.ImportIDs, 2)
	suite.Equal(i1.ID, ps.ImportIDs[0])
	suite.Equal(i2.ID, ps.ImportIDs[1])

	// Assert logs
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 2)
	suite.Contains(lines[0], "scheduling import for channel "+id1.String())
	suite.Contains(lines[1], "scheduling import for channel "+id2.String())
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_ChannelRepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	chr.FailWith(errors.New("boom!"))
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	is := scheduling.NewService(ir, chr, ps, log)

	// Execute
	err := is.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to fetch active channels")
	suite.ErrorContains(err, "boom!")

	// Assert logs
	suite.Empty(lbuf)
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_ImportRepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	ps := testutils.NewPubSubService()
	is := scheduling.NewService(ir, chr, ps, log)

	id := uuid.New()
	chr.Add(configuring.NewChannel(id, "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive).ToAggregator())

	// Execute
	err := is.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to save import for channel ")
	suite.Contains(err.Error(), id.String())
	suite.ErrorContains(err, "boom!")

	// Assert logs
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], "scheduling import for channel "+id.String())
}

func (suite *ServiceSuite) Test_ScheduleActiveChannels_PubSubServiceFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	ps.FailWith(errors.New("boom!"))
	is := scheduling.NewService(ir, chr, ps, log)

	id := uuid.New()
	chr.Add(configuring.NewChannel(id, "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive).ToAggregator())

	// Execute
	err := is.ScheduleActiveChannels(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to publish import ")
	suite.Contains(err.Error(), id.String())
	suite.ErrorContains(err, "boom!")

	// Assert logs
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())
}
