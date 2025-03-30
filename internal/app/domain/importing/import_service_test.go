package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/errs"
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
	s := importing.NewImportService(ir)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.Len(ir.Imports, 1)
	suite.NotNil(ir.Imports[i.ID()])
	suite.Equal(chID, ir.Imports[i.ID()].ChannelID)
	suite.Equal(aggregator.ImportStatusPending, ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(ir.Imports[i.ID()].EndedAt.Valid)

	// Fetch
	err = s.SetStatus(ctx, i, aggregator.ImportStatusFetching)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusFetching, ir.Imports[i.ID()].Status)

	// Process Import
	err = s.SetStatus(ctx, i, aggregator.ImportStatusProcessing)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusProcessing, ir.Imports[i.ID()].Status)

	// Add JobResults
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNew}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultMissing}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultMissing}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultUpdated}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultNoChange}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.NoError(s.SaveJobResult(ctx, i.ID(), &aggregator.ImportJob{ID: uuid.New(), Result: aggregator.ImportJobResultFailed}))
	suite.Len(ir.ImportJobs(), 15)

	// Publishing
	err = s.SetStatus(ctx, i, aggregator.ImportStatusPublishing)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusPublishing, ir.Imports[i.ID()].Status)

	// Completed
	err = s.MarkAsCompleted(ctx, i)
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusCompleted, ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].EndedAt.Valid)
	suite.True(ir.Imports[i.ID()].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal(1, ir.Imports[i.ID()].NewJobs())
	suite.Equal(2, ir.Imports[i.ID()].MissingJobs())
	suite.Equal(3, ir.Imports[i.ID()].UpdatedJobs())
	suite.Equal(4, ir.Imports[i.ID()].NoChangeJobs())
	suite.Equal(5, ir.Imports[i.ID()].FailedJobs())
}

func (suite *ImportServiceSuite) Test_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := importing.NewImportService(ir)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.False(ir.Imports[i.ID()].EndedAt.Valid)

	// Fail
	err = s.MarkAsFailed(ctx, i, errors.New("boom!"))
	suite.NoError(err)
	suite.Equal(aggregator.ImportStatusFailed, ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].EndedAt.Valid)
	suite.True(ir.Imports[i.ID()].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal("boom!", ir.Imports[i.ID()].Error.String)
}

func (suite *ImportServiceSuite) Test_FindImport_Success() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := importing.NewImportService(ir)
	ctx := context.Background()
	id := uuid.New()
	ir.AddImport(importing.NewImport(id, uuid.New()).ToDTO())

	// Execute
	i, err := s.FindImport(ctx, id)
	suite.NoError(err)

	// Success
	suite.NoError(err)
	suite.Equal(id, i.ID())
}

func (suite *ImportServiceSuite) Test_FindImport_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	s := importing.NewImportService(ir)
	ctx := context.Background()

	// Execute
	i, err := s.FindImport(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ImportServiceSuite) Test_FindImport_NotFound() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := importing.NewImportService(ir)
	ctx := context.Background()

	// Execute
	i, err := s.FindImport(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.ErrorIs(err, importing.ErrImportNotFound)
	suite.True(errs.IsValidationError(err))
}
