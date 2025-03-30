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

func TestImportService(t *testing.T) {
	suite.Run(t, new(ImportServiceSuite))
}

type ImportServiceSuite struct {
	suite.Suite
}

func (suite *ImportServiceSuite) Test_Success() {
	// Prepare
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	ctx := context.Background()
	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: uuid.New(),
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Now(),
	}
	ir.Imports[i.ID] = i

	// Fetch
	err := is.SetStatus(ctx, i, aggregator.ImportStatusFetching)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusFetching, ir.Imports[i.ID].Status)

	// Process Import
	err = is.SetStatus(ctx, i, aggregator.ImportStatusProcessing)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusProcessing, ir.Imports[i.ID].Status)

	// Add JobResults
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNew}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultMissing}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultMissing}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(is.SaveJobResult(ctx, i.ID, &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.Len(ir.ImportJobs(), 15)

	// Publishing
	err = is.SetStatus(ctx, i, aggregator.ImportStatusPublishing)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusPublishing, ir.Imports[i.ID].Status)

	// Completed
	err = is.MarkAsCompleted(ctx, i)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusCompleted, ir.Imports[i.ID].Status)
	suite.True(ir.Imports[i.ID].EndedAt.Valid)
	suite.True(ir.Imports[i.ID].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal(1, ir.Imports[i.ID].NewJobs())
	suite.Equal(2, ir.Imports[i.ID].MissingJobs())
	suite.Equal(3, ir.Imports[i.ID].UpdatedJobs())
	suite.Equal(4, ir.Imports[i.ID].NoChangeJobs())
	suite.Equal(5, ir.Imports[i.ID].FailedJobs())
}

func (suite *ImportServiceSuite) Test_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	ctx := context.Background()
	ch := &aggregator.Channel{
		ID:          uuid.New(),
		Name:        "airbeitnow test",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: ch.ID,
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Now(),
	}
	ir.Imports[i.ID] = i

	// Fail
	err := is.MarkAsFailed(ctx, i, errors.New("boom!"))
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusFailed, ir.Imports[i.ID].Status)
	suite.True(ir.Imports[i.ID].EndedAt.Valid)
	suite.True(ir.Imports[i.ID].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal("boom!", ir.Imports[i.ID].Error.String)
}
