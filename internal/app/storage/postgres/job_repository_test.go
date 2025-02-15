package postgres_test

import (
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/app/storage/postgres"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestJobRepository(t *testing.T) {
	suite.Run(t, new(JobRepositorySuite))
}

type JobRepositorySuite struct {
	testutils.IntegrationSuite
}

func (suite *JobRepositorySuite) Test_Save_New_Success() {
	// Prepare
	id := uuid.New()
	chID := uuid.New()
	pAt := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	j := job.New(
		id,
		chID,
		job.StatusActive,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		pAt,
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	r := postgres.NewJobRepository(suite.DB)

	// Execute
	err := r.Save(context.Background(), j)

	// Assert return
	suite.NoError(err)

	// Assert state change
	var dbJob postgres.Job
	err = suite.DB.Get(&dbJob, "SELECT * FROM jobs WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbJob.ID)
	suite.Equal(chID, dbJob.ChannelID)
	suite.Equal(job.StatusActive, job.Status(dbJob.Status))
	suite.Equal(job.PublishStatusPublished, job.PublishStatus(dbJob.PublishStatus))
	suite.Equal("https://example.com/job/id", dbJob.URL)
	suite.Equal("Software Engineer", dbJob.Title)
	suite.Equal("Job Description", dbJob.Description)
	suite.Equal("Indeed", dbJob.Source)
	suite.Equal("Amsterdam", dbJob.Location)
	suite.True(dbJob.Remote)
	suite.True(dbJob.PostedAt.Equal(pAt))
	suite.True(dbJob.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(dbJob.UpdatedAt.After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobRepositorySuite) Test_Save_Existing_Success() {
	// Prepare
	id := uuid.New()
	chID := uuid.New()
	cAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	_, err := suite.DB.Exec("INSERT INTO jobs (id, channel_id, status, publish_status, url, title, description, source, location, remote, posted_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		id,
		chID,
		int(job.StatusInactive),
		int(job.PublishStatusUnpublished),
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		cAt,
		time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
	)
	suite.NoError(err)

	pAt := time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC)
	chID2 := uuid.New()
	j := job.New(
		id,
		chID2,
		job.StatusActive,
		"https://example.com/job/id/new",
		"Software Engineer new",
		"Job Description new",
		"Indeed new",
		"Amsterdam new",
		false,
		pAt,
		job.WithPublishStatus(job.PublishStatusPublished),
	)

	r := postgres.NewJobRepository(suite.DB)

	// Execute
	err = r.Save(context.Background(), j)

	// Assert return
	suite.NoError(err)

	// Assert state change
	var count int
	err = suite.DB.Get(&count, "SELECT COUNT(*) FROM jobs WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(1, count)

	var dbJob postgres.Job
	err = suite.DB.Get(&dbJob, "SELECT * FROM jobs WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, dbJob.ID)
	suite.Equal(chID2, dbJob.ChannelID)
	suite.Equal(job.StatusActive, job.Status(dbJob.Status))
	suite.Equal(job.PublishStatusPublished, job.PublishStatus(dbJob.PublishStatus))
	suite.Equal("https://example.com/job/id/new", dbJob.URL)
	suite.Equal("Software Engineer new", dbJob.Title)
	suite.Equal("Job Description new", dbJob.Description)
	suite.Equal("Indeed new", dbJob.Source)
	suite.Equal("Amsterdam new", dbJob.Location)
	suite.False(dbJob.Remote)
	suite.True(dbJob.PostedAt.Equal(pAt))
	suite.True(dbJob.CreatedAt.Equal(cAt))
	suite.True(dbJob.UpdatedAt.After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobRepositorySuite) Test_Save_Error() {
	// Prepare
	id := uuid.New()
	chID := uuid.New()
	j := job.New(
		id,
		chID,
		job.StatusActive,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	r := postgres.NewJobRepository(suite.BadDB)

	// Execute
	err := r.Save(context.Background(), j)

	// Assert return
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *JobRepositorySuite) Test_GetByChannelID_Success() {
	// Prepare
	chID1 := uuid.New()
	jID1 := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO jobs (id, channel_id, status, publish_status, url, title, description, source, location, remote, posted_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		jID1,
		chID1,
		int(job.StatusInactive),
		int(job.PublishStatusUnpublished),
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Now(),
		time.Now(),
	)
	suite.NoError(err)
	jID2 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO jobs (id, channel_id, status, publish_status, url, title, description, source, location, remote, posted_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		jID2,
		chID1,
		int(job.StatusActive),
		int(job.PublishStatusPublished),
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		time.Now(),
		time.Now(),
	)
	suite.NoError(err)
	chID2 := uuid.New()
	jID3 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO jobs (id, channel_id, status, publish_status, url, title, description, source, location, remote, posted_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		jID3,
		chID2,
		int(job.StatusActive),
		int(job.PublishStatusPublished),
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC),
		time.Now(),
		time.Now(),
	)
	suite.NoError(err)

	r := postgres.NewJobRepository(suite.DB)

	// Execute
	jobs, err := r.GetByChannelID(context.Background(), chID1)

	// Assert return
	suite.NoError(err)
	suite.Len(jobs, 2)
	suite.Equal(jID2, jobs[0].ID())
	suite.Equal(jID1, jobs[1].ID())
}

func (suite *JobRepositorySuite) Test_GetByChannelID_Error() {
	// Prepare
	r := postgres.NewJobRepository(suite.BadDB)

	// Execute
	jobs, err := r.GetByChannelID(context.Background(), uuid.New())

	// Assert return
	suite.Nil(jobs)
	suite.Error(err)
	suite.ErrorContains(err, "sql: database is closed")
}
