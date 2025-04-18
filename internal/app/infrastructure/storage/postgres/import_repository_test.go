package postgres_test

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

func TestImportRepository(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}
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
	m1ID := uuid.New()
	m2ID := uuid.New()
	m3ID := uuid.New()
	j1ID := uuid.New()
	j2ID := uuid.New()
	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	m1Time := time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)
	m2Time := time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC)
	m3Time := time.Date(2020, 1, 1, 0, 0, 5, 0, time.UTC)
	i := &aggregator.Import{
		ID:        id,
		ChannelID: chID,
		Error:     null.StringFrom("error"),
		Status:    aggregator.ImportStatusProcessing,
		StartedAt: sAt,
		EndedAt:   null.TimeFrom(eAt),
		Metrics: []*aggregator.ImportMetric{
			{ID: m1ID, JobID: j1ID, MetricType: aggregator.ImportMetricTypeNew, Err: null.NewString("", false), CreatedAt: m1Time},
			{ID: m2ID, JobID: j1ID, MetricType: aggregator.ImportMetricTypePublish, Err: null.NewString("", false), CreatedAt: m2Time},
			{ID: m3ID, JobID: j2ID, MetricType: aggregator.ImportMetricTypeError, Err: null.StringFrom("failed to process job"), CreatedAt: m3Time},
		},
	}

	// Execute
	err = r.SaveImport(context.Background(), i)

	// Execute
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM imports WHERE id = $1", i.ID)
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImport aggregator.Import
	err = suite.DB.Get(&dbImport, "SELECT * FROM imports WHERE id = $1", i.ID)
	suite.NoError(err)
	suite.Equal(i.ID, dbImport.ID)
	suite.Equal(i.ChannelID, dbImport.ChannelID)
	suite.Equal(i.Status, dbImport.Status)
	suite.True(i.StartedAt.Equal(dbImport.StartedAt))
	suite.True(i.EndedAt.Time.Equal(dbImport.EndedAt.Time))
	suite.Equal(i.Error.String, dbImport.Error.String)

	var dbImportMetrics []*aggregator.ImportMetric
	err = suite.DB.Select(&dbImportMetrics, "SELECT id, job_id, metric_type, error, created_at FROM import_metrics WHERE import_id = $1", i.ID)
	suite.NoError(err)
	suite.Len(dbImportMetrics, 3)
	m1Found := false
	m2Found := false
	m3Found := false
	for _, m := range dbImportMetrics {
		if m.ID == m1ID {
			m1Found = true
			suite.Equal(aggregator.ImportMetricTypeNew, m.MetricType)
			suite.Equal(j1ID, m.JobID)
			suite.False(m.Err.Valid)
			suite.True(m1Time.Equal(m.CreatedAt))
		}
		if m.ID == m2ID {
			m2Found = true
			suite.Equal(aggregator.ImportMetricTypePublish, m.MetricType)
			suite.Equal(j1ID, m.JobID)
			suite.False(m.Err.Valid)
			suite.True(m2Time.Equal(m.CreatedAt))
		}
		if m.ID == m3ID {
			m3Found = true
			suite.Equal(aggregator.ImportMetricTypeError, m.MetricType)
			suite.Equal(j2ID, m.JobID)
			suite.True(m.Err.Valid)
			suite.Equal("failed to process job", m.Err.String)
			suite.True(m3Time.Equal(m.CreatedAt))
		}
	}
	suite.True(m1Found)
	suite.True(m2Found)
	suite.True(m3Found)
}

func (suite *ImportRepositorySuite) Test_SaveImport_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	id := uuid.New()
	chID := uuid.New()
	i := &aggregator.Import{
		ID:        id,
		ChannelID: chID,
		Error:     null.StringFrom("error"),
		Status:    aggregator.ImportStatusProcessing,
	}

	// Execute
	err := r.SaveImport(context.Background(), i)

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
	i := &aggregator.Import{
		ID:        id,
		ChannelID: chID,
		Error:     null.StringFrom("error"),
		Status:    aggregator.ImportStatusProcessing,
		StartedAt: sAt,
		EndedAt:   null.TimeFrom(eAt),
	}
	err = r.SaveImport(context.Background(), i)
	suite.NoError(err)
	m1ID := uuid.New()
	m2ID := uuid.New()
	m3ID := uuid.New()
	j1ID := uuid.New()
	j2ID := uuid.New()
	j3ID := uuid.New()
	m1Time := time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)
	m2Time := time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC)
	m3Time := time.Date(2020, 1, 1, 0, 0, 5, 0, time.UTC)
	suite.NoError(r.SaveImportMetric(context.Background(), i.ID, &aggregator.ImportMetric{ID: m1ID, JobID: j1ID, MetricType: aggregator.ImportMetricTypeUpdated, Err: null.NewString("", false), CreatedAt: m1Time}))
	suite.NoError(r.SaveImportMetric(context.Background(), i.ID, &aggregator.ImportMetric{ID: m2ID, JobID: j2ID, MetricType: aggregator.ImportMetricTypeUpdated, Err: null.NewString("", false), CreatedAt: m2Time}))
	suite.NoError(r.SaveImportMetric(context.Background(), i.ID, &aggregator.ImportMetric{ID: m3ID, JobID: j3ID, MetricType: aggregator.ImportMetricTypeNew, Err: null.StringFrom("failed to process job"), CreatedAt: m3Time}))

	// Execute
	i2, err := r.FindImport(context.Background(), i.ID)

	// Assert
	suite.NoError(err)
	suite.Equal(i.ID, i2.ID)
	suite.Equal(i.ChannelID, i2.ChannelID)
	suite.Equal(i.Status, i2.Status)
	suite.True(i.StartedAt.Equal(i2.StartedAt))
	suite.True(i.EndedAt.Time.Equal(i2.EndedAt.Time))
	suite.Equal(i.Error, i2.Error)
	suite.Equal(3, len(i2.Metrics))
	var m1Found, m2Found, m3Found bool
	for _, m := range i2.Metrics {
		if m.ID == m1ID {
			m1Found = true
			suite.Equal(aggregator.ImportMetricTypeUpdated, m.MetricType)
			suite.False(m.Err.Valid)
			suite.True(m1Time.Equal(m.CreatedAt))
		}
		if m.ID == m2ID {
			m2Found = true
			suite.Equal(aggregator.ImportMetricTypeUpdated, m.MetricType)
			suite.False(m.Err.Valid)
			suite.True(m2Time.Equal(m.CreatedAt))
		}
		if m.ID == m3ID {
			m3Found = true
			suite.Equal(aggregator.ImportMetricTypeNew, m.MetricType)
			suite.True(m.Err.Valid)
			suite.Equal("failed to process job", m.Err.String)
			suite.True(m3Time.Equal(m.CreatedAt))
		}
	}
	suite.True(m1Found)
	suite.True(m2Found)
	suite.True(m3Found)
}

func (suite *ImportRepositorySuite) Test_FindImport_BadConnectionFail() {
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

func (suite *ImportRepositorySuite) Test_SaveImportMetric_New_Success() {
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
	mID := uuid.New()
	mAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	m := &aggregator.ImportMetric{ID: mID, JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeUpdated, Err: null.StringFrom("failed to process job"), CreatedAt: mAt}

	// Execute
	err = r.SaveImportMetric(context.Background(), iID, m)

	// Assert
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM import_metrics")
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImportJob aggregator.ImportMetric
	err = suite.DB.Get(&dbImportJob, "SELECT id, job_id, metric_type, error, created_at FROM import_metrics WHERE id = $1", mID)
	suite.NoError(err)
	suite.Equal(m.JobID, dbImportJob.JobID)
	suite.Equal(aggregator.ImportMetricTypeUpdated, dbImportJob.MetricType)
	suite.True(dbImportJob.Err.Valid)
	suite.Equal(m.Err.String, dbImportJob.Err.String)
	suite.True(dbImportJob.CreatedAt.Equal(mAt))
}

func (suite *ImportRepositorySuite) Test_SaveImportMetric_Existing_Success() {
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

	mID := uuid.New()
	jID := uuid.New()
	mAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	_, err = suite.DB.Exec("INSERT INTO import_metrics (id, import_id, job_id, metric_type, error, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		mID,
		iID,
		jID,
		aggregator.ImportMetricTypeError,
		null.NewString("", false),
		mAt,
	)
	suite.NoError(err)

	r := postgres.NewImportRepository(suite.DB)
	m := &aggregator.ImportMetric{ID: mID, JobID: jID, MetricType: aggregator.ImportMetricTypeUpdated, Err: null.StringFrom("failed to process job"), CreatedAt: time.Now()}

	// Execute
	err = r.SaveImportMetric(context.Background(), iID, m)

	// Assert
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM import_metrics")
	suite.NoError(err)
	suite.Equal(1, count)

	var dbImportJob aggregator.ImportMetric
	err = suite.DB.Get(&dbImportJob, "SELECT id, job_id, metric_type, error, created_at FROM import_metrics WHERE id = $1", mID)
	suite.NoError(err)
	suite.Equal(jID, dbImportJob.JobID)
	suite.Equal(aggregator.ImportMetricTypeUpdated, dbImportJob.MetricType)
	suite.True(dbImportJob.Err.Valid)
	suite.Equal(m.Err.String, dbImportJob.Err.String)
	suite.True(dbImportJob.CreatedAt.Equal(mAt))
}

func (suite *ImportRepositorySuite) Test_SaveImportMetric_Fail() {
	// Prepare
	r := postgres.NewImportRepository(suite.BadDB)
	m := &aggregator.ImportMetric{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeUpdated}

	// Execute
	err := r.SaveImportMetric(context.Background(), uuid.New(), m)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, m.ID.String())
	suite.ErrorContains(err, "sql: database is closed")
}
