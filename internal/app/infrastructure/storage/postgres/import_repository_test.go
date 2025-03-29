package postgres_test

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
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
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i := imports.New(
		id,
		chID,
		imports.WithError("error"),
		imports.WithStatus(base.ImportStatusProcessing),
		imports.WithStartAt(sAt),
		imports.WithEndAt(eAt),
		imports.WithMetadata(1, 2, 3, 4, 5),
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

	var dbImport postgres.Import
	err = suite.DB.Get(&dbImport, "SELECT * FROM imports WHERE id = $1", i.ID())
	suite.NoError(err)
	suite.Equal(i.ID(), dbImport.ID)
	suite.Equal(i.ChannelID(), dbImport.ChannelID)
	suite.Equal(i.Status(), dbImport.Status)
	suite.True(i.StartedAt().Equal(dbImport.StartedAt))
	suite.True(i.EndedAt().Time.Equal(dbImport.EndedAt.Time))
	suite.Equal(i.Error().String, dbImport.Error.String)
	suite.Equal(i.NewJobs(), dbImport.NewJobs)
	suite.Equal(i.UpdatedJobs(), dbImport.UpdatedJobs)
	suite.Equal(i.NoChangeJobs(), dbImport.NoChangeJobs)
	suite.Equal(i.MissingJobs(), dbImport.MissingJobs)
	suite.Equal(i.FailedJobs(), dbImport.FailedJobs)
}

func (suite *ImportRepositorySuite) Test_SaveImport_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	id := uuid.New()
	chID := uuid.New()
	i := imports.New(
		id,
		chID,
		imports.WithError("error"),
		imports.WithStatus(base.ImportStatusProcessing),
		imports.WithStartAt(time.Now()),
		imports.WithEndAt(time.Now()),
		imports.WithMetadata(1, 2, 3, 4, 5),
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
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i := imports.New(
		id,
		chID,
		imports.WithError("error"),
		imports.WithStatus(base.ImportStatusProcessing),
		imports.WithStartAt(sAt),
		imports.WithEndAt(eAt),
		imports.WithMetadata(1, 2, 3, 4, 5),
	)
	err = r.SaveImport(context.Background(), i.ToDTO())
	suite.NoError(err)

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
	suite.Equal(i.NewJobs(), i2.NewJobs)
	suite.Equal(i.UpdatedJobs(), i2.UpdatedJobs)
	suite.Equal(i.NoChangeJobs(), i2.NoChangeJobs)
	suite.Equal(i.MissingJobs(), i2.MissingJobs)
	suite.Equal(i.FailedJobs(), i2.FailedJobs)
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
	suite.ErrorIs(err, postgres.ErrImportNotFound)
	suite.Nil(i)
}

func (suite *ImportRepositorySuite) Test_GetImports_Success() {
	// Prepare
	r := postgres.NewImportRepository(suite.DB)

	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
	)
	suite.NoError(err)

	id1 := uuid.New()
	sAt1 := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id1,
		chID,
		base.ImportStatusProcessing,
		sAt1,
	)
	suite.NoError(err)

	id2 := uuid.New()
	sAt2 := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id2,
		chID,
		base.ImportStatusProcessing,
		sAt2,
	)
	suite.NoError(err)

	id3 := uuid.New()
	sAt3 := time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		id3,
		chID,
		base.ImportStatusProcessing,
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
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
	)
	suite.NoError(err)

	iID := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		iID,
		chID,
		base.ImportStatusProcessing,
		time.Now(),
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	jr := imports.NewResult(uuid.New(), iID, base.JobStatusUpdated)

	// Execute
	err = r.SaveImportJob(context.Background(), jr.ToDTO())

	// Assert
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM import_job_results WHERE job_id = $1 and import_id = $2", jr.JobID(), jr.ImportID())
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImportJobResult postgres.ImportJobResult
	err = suite.DB.Get(&dbImportJobResult, "SELECT * FROM import_job_results WHERE job_id = $1 and import_id = $2", jr.JobID(), jr.ImportID())
	suite.NoError(err)
	suite.Equal(base.JobStatusUpdated, dbImportJobResult.Result)
}

func (suite *ImportRepositorySuite) Test_SaveImportJob_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	jr := imports.NewResult(uuid.New(), uuid.New(), base.JobStatusUpdated)

	// Execute
	err := r.SaveImportJob(context.Background(), jr.ToDTO())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, jr.JobID().String())
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *ImportRepositorySuite) Test_GetJobsByImportID_Success() {
	// Prepare
	chID := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO channels (id, name, integration, status) VALUES ($1, $2, $3, $4)",
		chID,
		"Channel Name",
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
	)
	suite.NoError(err)

	iID1 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		iID1,
		chID,
		base.ImportStatusProcessing,
		time.Now(),
	)
	suite.NoError(err)

	iID2 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO imports (id, channel_id, status, started_at) VALUES ($1, $2, $3, $4)",
		iID2,
		chID,
		base.ImportStatusProcessing,
		time.Now(),
	)
	suite.NoError(err)

	_, err = suite.DB.Exec("INSERT INTO import_job_results (import_id, job_id, result) VALUES ($1, $2, $3)",
		iID1,
		uuid.New(),
		base.JobStatusUpdated,
	)
	suite.NoError(err)
	_, err = suite.DB.Exec("INSERT INTO import_job_results (import_id, job_id, result) VALUES ($1, $2, $3)",
		iID1,
		uuid.New(),
		base.JobStatusUpdated,
	)
	suite.NoError(err)
	_, err = suite.DB.Exec("INSERT INTO import_job_results (import_id, job_id, result) VALUES ($1, $2, $3)",
		iID1,
		uuid.New(),
		base.JobStatusUpdated,
	)
	suite.NoError(err)
	_, err = suite.DB.Exec("INSERT INTO import_job_results (import_id, job_id, result) VALUES ($1, $2, $3)",
		iID2,
		uuid.New(),
		base.JobStatusUpdated,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)

	// Execute
	jobs, err := r.GetJobsByImportID(context.Background(), iID1)

	// Assert
	suite.NoError(err)
	suite.Len(jobs, 3)
}

func (suite *ImportRepositorySuite) Test_GetJobsByImportID_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)

	// Execute
	jobs, err := r.GetJobsByImportID(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.Nil(jobs)
	suite.ErrorContains(err, "sql: database is closed")
}
