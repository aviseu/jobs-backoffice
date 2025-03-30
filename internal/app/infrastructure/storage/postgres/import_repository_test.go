package postgres_test

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestImportRepository(t *testing.T) {
	suite.Run(t, new(ImportRepositorySuite))
}

type ImportRepositorySuite struct {
	testutils.PostgresSuite
}

func (suite *ImportRepositorySuite) Test_SaveImport_Success() {
	// Prepare
	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i := importing.NewImport(
		id,
		chID,
		importing.ImportWithError("error"),
		importing.ImportWithStatus(aggregator.ImportStatusProcessing),
		importing.ImportWithStartAt(sAt),
		importing.ImportWithEndAt(eAt),
	)

	// Execute
	err = r.SaveImport(context.Background(), i.ToDTO())

	// Execute
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM imports WHERE id = $1", i.ID())
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImport aggregator.Import
	err = suite.DB.Get(&dbImport, "SELECT * FROM imports WHERE id = $1", i.ID())
	suite.NoError(err)
	suite.Equal(i.ID(), dbImport.ID)
	suite.Equal(i.ChannelID(), dbImport.ChannelID)
	suite.Equal(i.Status(), dbImport.Status)
	suite.True(i.StartedAt().Equal(dbImport.StartedAt))
	suite.True(i.EndedAt().Time.Equal(dbImport.EndedAt.Time))
	suite.Equal(i.Error().String, dbImport.Error.String)
}

func (suite *ImportRepositorySuite) Test_SaveImport_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	id := uuid.New()
	chID := uuid.New()
	i := importing.NewImport(
		id,
		chID,
		importing.ImportWithError("error"),
		importing.ImportWithStatus(aggregator.ImportStatusProcessing),
		importing.ImportWithStartAt(time.Now()),
		importing.ImportWithEndAt(time.Now()),
	)

	// Execute
	err := r.SaveImport(context.Background(), i.ToDTO())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ImportRepositorySuite) Test_FindImport_Success() {
	// Prepare
	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i := importing.NewImport(
		id,
		chID,
		importing.ImportWithError("error"),
		importing.ImportWithStatus(aggregator.ImportStatusProcessing),
		importing.ImportWithStartAt(sAt),
		importing.ImportWithEndAt(eAt),
	)
	err = r.SaveImport(context.Background(), i.ToDTO())
	suite.NoError(err)
	suite.NoError(r.SaveImportJob(context.Background(), i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(r.SaveImportJob(context.Background(), i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(r.SaveImportJob(context.Background(), i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNew}))

	// Execute
	i2, err := r.FindImport(context.Background(), i.ID())

	// Assert
	suite.NoError(err)
	suite.Equal(i.ID(), i2.ID)
	suite.Equal(i.ChannelID(), i2.ChannelID)
	suite.Equal(i.Status(), i2.Status)
	suite.True(i.StartedAt().Equal(i2.StartedAt))
	suite.True(i.EndedAt().Time.Equal(i2.EndedAt.Time))
	suite.Equal(i.Error(), i2.Error)
	suite.Equal(3, len(i2.Jobs))
}

func (suite *ImportRepositorySuite) Test_FindImport_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	id := uuid.New()

	// Execute
	i, err := r.FindImport(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.Nil(i)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ImportRepositorySuite) Test_FindImport_NotFound() {
	// Prepare
	r := postgres.NewImportRepository(suite.DB)
	id := uuid.New()

	// Execute
	i, err := r.FindImport(context.Background(), id)

	// Assert
	suite.ErrorIs(err, infrastructure.ErrImportNotFound)
	suite.Nil(i)
}

func (suite *ImportRepositorySuite) Test_GetImports_Success() {
	// Prepare
	r := postgres.NewImportRepository(suite.DB)

	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
	)
	suite.NoError(err)

	id1 := uuid.New()
	sAt1 := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id1,
		chID,
		aggregator.ImportStatusProcessing,
		sAt1,
	)
	suite.NoError(err)

	id2 := uuid.New()
	sAt2 := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id2,
		chID,
		aggregator.ImportStatusProcessing,
		sAt2,
	)
	suite.NoError(err)

	id3 := uuid.New()
	sAt3 := time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id3,
		chID,
		aggregator.ImportStatusProcessing,
		sAt3,
	)
	suite.NoError(err)

	// Execute
	ii, err := r.GetImports(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(ii, 3)

	suite.Equal(id3, ii[0].ID)
	suite.Equal(id2, ii[1].ID)
	suite.Equal(id1, ii[2].ID)
}

func (suite *ImportRepositorySuite) Test_GetImports_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)

	// Execute
	ii, err := r.GetImports(context.Background())

	// Assert
	suite.Error(err)
	suite.Nil(ii)
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ImportRepositorySuite) Test_SaveImportJob_Success() {
	// Prepare
	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		aggregator.IntegrationArbeitnow,
		aggregator.ChannelStatusInactive,
	)
	suite.NoError(err)

	iID := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		iID,
		chID,
		aggregator.ImportStatusProcessing,
		time.Now(),
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	ij := &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}

	// Execute
	err = r.SaveImportJob(context.Background(), iID, ij)

	// Assert
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM import_job_results WHERE job_id = $1 and import_id = $2", ij.ID, iID)
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImportJob aggregator.ImportJob
	err = suite.DB.Get(&dbImportJob, "SELECT job_id, result FROM import_job_results WHERE job_id = $1 and import_id = $2", ij.ID, iID)
	suite.NoError(err)
	suite.Equal(aggregator.ImportJobResultUpdated, dbImportJob.Result)
}

func (suite *ImportRepositorySuite) Test_SaveImportJob_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	ij := &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}

	// Execute
	err := r.SaveImportJob(context.Background(), uuid.New(), ij)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, ij.ID.String())
	suite.ErrorContains(err, "sql: database is closed")
}
