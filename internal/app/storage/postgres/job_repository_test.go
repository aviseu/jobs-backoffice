package postgres_test

import (
	"context"
	"github.com/aviseu/jobs/internal/app/domain"
	"github.com/aviseu/jobs/internal/app/storage/postgres"
	"github.com/aviseu/jobs/internal/testutils"
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
	pAt := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	j := domain.NewJob(
		id,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		pAt,
	)
	r := postgres.NewJobRepository(suite.DB)

	// Execute
	err := r.Save(context.Background(), j)

	// Assert return
	suite.NoError(err)

	// Assert state change
	var job postgres.Job
	err = suite.DB.Get(&job, "SELECT * FROM jobs WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, job.ID)
	suite.Equal("https://example.com/job/id", job.URL)
	suite.Equal("Software Engineer", job.Title)
	suite.Equal("Job Description", job.Description)
	suite.Equal("Indeed", job.Source)
	suite.Equal("Amsterdam", job.Location)
	suite.True(job.Remote)
	suite.True(job.PostedAt.Equal(pAt))
	suite.True(job.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(job.UpdatedAt.After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobRepositorySuite) Test_Save_Existing_Success() {
	// Prepare
	id := uuid.New()
	cAt := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	_, err := suite.DB.Exec("INSERT INTO jobs (id, url, title, description, source, location, remote, posted_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		id,
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
	j := domain.NewJob(
		id,
		"https://example.com/job/id/new",
		"Software Engineer new",
		"Job Description new",
		"Indeed new",
		"Amsterdam new",
		false,
		pAt,
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

	var job postgres.Job
	err = suite.DB.Get(&job, "SELECT * FROM jobs WHERE id = $1", id)
	suite.NoError(err)
	suite.Equal(id, job.ID)
	suite.Equal("https://example.com/job/id/new", job.URL)
	suite.Equal("Software Engineer new", job.Title)
	suite.Equal("Job Description new", job.Description)
	suite.Equal("Indeed new", job.Source)
	suite.Equal("Amsterdam new", job.Location)
	suite.False(job.Remote)
	suite.True(job.PostedAt.Equal(pAt))
	suite.True(job.CreatedAt.Equal(cAt))
	suite.True(job.UpdatedAt.After(time.Now().Add(-2 * time.Second)))
}

func (suite *JobRepositorySuite) Test_Save_Error() {
	// Prepare
	id := uuid.New()
	j := domain.NewJob(
		id,
		"https://example.com/job/id",
		"Software Engineer",
		"Job Description",
		"Indeed",
		"Amsterdam",
		true,
		time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
	)
	r := postgres.NewJobRepository(suite.BadDB)

	// Execute
	err := r.Save(context.Background(), j)

	// Assert return
	suite.Error(err)
	suite.ErrorContains(err, id.String())
	suite.ErrorContains(err, "sql: database is closed")
}
