package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestScheduleImportsAction(t *testing.T) {
	suite.Run(t, new(ScheduleImportsActionSuite))
}

type ScheduleImportsActionSuite struct {
	suite.Suite
}

func (suite *ScheduleImportsActionSuite) Test_Success() {
	// Prepare
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	ps := testutils.NewPubSubService()
	lbuf, log := testutils.NewLogger()
	ia := importing.NewScheduleImportAction(is, ps, log)
	s := importing.NewScheduleImportsAction(chr, ia)

	id1 := uuid.New()
	chr.Add(configuring.NewChannel(id1, "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive).ToDTO())
	id2 := uuid.New()
	chr.Add(configuring.NewChannel(id2, "channel 2", base.IntegrationArbeitnow, base.ChannelStatusActive).ToDTO())
	chr.Add(configuring.NewChannel(uuid.New(), "channel 3", base.IntegrationArbeitnow, base.ChannelStatusInactive).ToDTO())

	// Execute
	err := s.Execute(context.Background())

	// Assert
	suite.NoError(err)

	// Assert imports created
	suite.Equal(2, len(ir.Imports))
	var i1, i2 *postgres.Import
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

func (suite *ScheduleImportsActionSuite) Test_ChannelRepositoryFail() {
	// Prepare
	chr := testutils.NewChannelRepository()
	chr.FailWith(errors.New("boom!"))
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	ps := testutils.NewPubSubService()
	lbuf, log := testutils.NewLogger()
	ia := importing.NewScheduleImportAction(is, ps, log)
	s := importing.NewScheduleImportsAction(chr, ia)

	// Execute
	err := s.Execute(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to fetch active channels")
	suite.ErrorContains(err, "boom!")

	// Assert logs
	suite.Empty(lbuf)
}

func (suite *ScheduleImportsActionSuite) Test_ImportServiceFail() {
	// Prepare
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	is := importing.NewImportService(ir)
	ps := testutils.NewPubSubService()
	lbuf, log := testutils.NewLogger()
	ia := importing.NewScheduleImportAction(is, ps, log)
	s := importing.NewScheduleImportsAction(chr, ia)

	id := uuid.New()
	chr.Add(configuring.NewChannel(id, "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive).ToDTO())

	// Execute
	err := s.Execute(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "failed to start import for channel "+id.String())
	suite.ErrorContains(err, "boom!")

	// Assert logs
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())
}

func (suite *ScheduleImportsActionSuite) Test_PubSubServiceFail() {
	// Prepare
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	ps := testutils.NewPubSubService()
	ps.FailWith(errors.New("boom!"))
	lbuf, log := testutils.NewLogger()
	ia := importing.NewScheduleImportAction(is, ps, log)
	s := importing.NewScheduleImportsAction(chr, ia)

	id := uuid.New()
	chr.Add(configuring.NewChannel(id, "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive).ToDTO())

	// Execute
	err := s.Execute(context.Background())

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "ailed to publish import command for import")
	suite.ErrorContains(err, "boom!")

	// Assert logs
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())
}
