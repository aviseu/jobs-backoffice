package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
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
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	is := importing.NewService(ir, ps, log)

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

func (suite *ServiceSuite) Test_ScheduleImport_RepositoryFailed() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom"))
	ps := testutils.NewPubSubService()
	is := importing.NewService(ir, ps, log)

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
	ir := testutils.NewImportRepository()
	ps := testutils.NewPubSubService()
	ps.FailWith(errors.New("boom"))
	is := importing.NewService(ir, ps, log)

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
