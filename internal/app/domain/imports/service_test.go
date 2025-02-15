package imports_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/errs"
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

func (suite *ServiceSuite) Test_Success() {
	// Prepare
	r := testutils.NewImportRepository()
	s := imports.NewService(r)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.Len(r.Imports, 1)
	suite.NotNil(r.Imports[i.ID()])
	suite.Equal(chID, r.Imports[i.ID()].ChannelID())
	suite.Equal(imports.StatusPending, r.Imports[i.ID()].Status())
	suite.True(r.Imports[i.ID()].StartedAt().After(time.Now().Add(-2 * time.Second)))
	suite.False(r.Imports[i.ID()].EndedAt().Valid)

	// Fetch
	err = s.SetStatus(ctx, i, imports.StatusFetching)
	suite.NoError(err)
	suite.Equal(imports.StatusFetching, r.Imports[i.ID()].Status())

	// Process Import
	err = s.SetStatus(ctx, i, imports.StatusProcessing)
	suite.NoError(err)
	suite.Equal(imports.StatusProcessing, r.Imports[i.ID()].Status())

	// Add JobResults
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusNew)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusMissing)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusMissing)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), imports.JobStatusFailed)))
	suite.Len(r.JobResults, 15)

	// Publishing
	err = s.SetStatus(ctx, i, imports.StatusPublishing)
	suite.NoError(err)
	suite.Equal(imports.StatusPublishing, r.Imports[i.ID()].Status())

	// Completed
	err = s.MarkAsCompleted(ctx, i)
	suite.NoError(err)
	suite.Equal(imports.StatusCompleted, r.Imports[i.ID()].Status())
	suite.True(r.Imports[i.ID()].EndedAt().Valid)
	suite.True(r.Imports[i.ID()].EndedAt().Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal(1, r.Imports[i.ID()].NewJobs())
	suite.Equal(2, r.Imports[i.ID()].MissingJobs())
	suite.Equal(3, r.Imports[i.ID()].UpdatedJobs())
	suite.Equal(4, r.Imports[i.ID()].NoChangeJobs())
	suite.Equal(5, r.Imports[i.ID()].FailedJobs())
}

func (suite *ServiceSuite) Test_Fail() {
	// Prepare
	r := testutils.NewImportRepository()
	s := imports.NewService(r)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.False(r.Imports[i.ID()].EndedAt().Valid)

	// Fail
	err = s.MarkAsFailed(ctx, i, errors.New("boom!"))
	suite.NoError(err)
	suite.Equal(imports.StatusFailed, r.Imports[i.ID()].Status())
	suite.True(r.Imports[i.ID()].EndedAt().Valid)
	suite.True(r.Imports[i.ID()].EndedAt().Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal("boom!", r.Imports[i.ID()].Error().String)
}

func (suite *ServiceSuite) Test_FindImport_Success() {
	// Prepare
	r := testutils.NewImportRepository()
	s := imports.NewService(r)
	ctx := context.Background()
	id := uuid.New()
	r.Add(imports.New(id, uuid.New()))

	// Execute
	i, err := s.FindImport(ctx, id)
	suite.NoError(err)

	// Success
	suite.NoError(err)
	suite.Equal(id, i.ID())
}

func (suite *ServiceSuite) Test_FindImport_Fail() {
	// Prepare
	r := testutils.NewImportRepository()
	r.FailWith(errors.New("boom!"))
	s := imports.NewService(r)
	ctx := context.Background()

	// Execute
	i, err := s.FindImport(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_FindImport_NotFound() {
	// Prepare
	r := testutils.NewImportRepository()
	s := imports.NewService(r)
	ctx := context.Background()

	// Execute
	i, err := s.FindImport(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.ErrorIs(err, imports.ErrImportNotFound)
	suite.True(errs.IsValidationError(err))
}
