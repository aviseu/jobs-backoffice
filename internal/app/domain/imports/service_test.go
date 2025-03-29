package imports_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/errs"
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
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.Len(ir.Imports, 1)
	suite.NotNil(ir.Imports[i.ID()])
	suite.Equal(chID, ir.Imports[i.ID()].ChannelID)
	suite.Equal(int(imports.StatusPending), ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].StartedAt.After(time.Now().Add(-2 * time.Second)))
	suite.False(ir.Imports[i.ID()].EndedAt.Valid)

	// Fetch
	err = s.SetStatus(ctx, i, imports.StatusFetching)
	suite.NoError(err)
	suite.Equal(int(imports.StatusFetching), ir.Imports[i.ID()].Status)

	// Process Import
	err = s.SetStatus(ctx, i, imports.StatusProcessing)
	suite.NoError(err)
	suite.Equal(int(imports.StatusProcessing), ir.Imports[i.ID()].Status)

	// Add JobResults
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusNew)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusMissing)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusMissing)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusUpdated)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusNoChange)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusFailed)))
	suite.NoError(s.SaveJobResult(ctx, imports.NewResult(uuid.New(), i.ID(), base.JobStatusFailed)))
	suite.Len(ir.JobResults, 15)

	// Publishing
	err = s.SetStatus(ctx, i, imports.StatusPublishing)
	suite.NoError(err)
	suite.Equal(int(imports.StatusPublishing), ir.Imports[i.ID()].Status)

	// Completed
	err = s.MarkAsCompleted(ctx, i)
	suite.NoError(err)
	suite.Equal(int(imports.StatusCompleted), ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].EndedAt.Valid)
	suite.True(ir.Imports[i.ID()].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal(1, ir.Imports[i.ID()].NewJobs)
	suite.Equal(2, ir.Imports[i.ID()].MissingJobs)
	suite.Equal(3, ir.Imports[i.ID()].UpdatedJobs)
	suite.Equal(4, ir.Imports[i.ID()].NoChangeJobs)
	suite.Equal(5, ir.Imports[i.ID()].FailedJobs)
}

func (suite *ServiceSuite) Test_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()
	chID := uuid.New()

	// Start
	i, err := s.Start(ctx, uuid.New(), chID)
	suite.NoError(err)
	suite.False(ir.Imports[i.ID()].EndedAt.Valid)

	// Fail
	err = s.MarkAsFailed(ctx, i, errors.New("boom!"))
	suite.NoError(err)
	suite.Equal(int(imports.StatusFailed), ir.Imports[i.ID()].Status)
	suite.True(ir.Imports[i.ID()].EndedAt.Valid)
	suite.True(ir.Imports[i.ID()].EndedAt.Time.After(time.Now().Add(-2 * time.Second)))
	suite.Equal("boom!", ir.Imports[i.ID()].Error.String)
}

func (suite *ServiceSuite) Test_FindImport_Success() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()
	id := uuid.New()
	ir.Add(imports.New(id, uuid.New()).ToDTO())

	// Execute
	i, err := s.FindImport(ctx, id)
	suite.NoError(err)

	// Success
	suite.NoError(err)
	suite.Equal(id, i.ID())
}

func (suite *ServiceSuite) Test_FindImport_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	s := imports.NewService(ir)
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
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()

	// Execute
	i, err := s.FindImport(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.ErrorIs(err, imports.ErrImportNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) Test_FindImportWithForcedMetadata_WithoutMetadata_Success() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()
	id := uuid.New()
	ir.Add(imports.New(id, uuid.New()).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusNew).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusUpdated).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusUpdated).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusNoChange).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusNoChange).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusNoChange).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusMissing).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusMissing).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusMissing).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusMissing).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusFailed).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusFailed).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusFailed).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusFailed).ToDTO())
	ir.AddResult(imports.NewResult(uuid.New(), id, base.JobStatusFailed).ToDTO())

	// Execute
	i, err := s.FindImportWithForcedMetadata(ctx, id)

	// Success
	suite.NoError(err)
	suite.Equal(id, i.ID())
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(5, i.FailedJobs())
}

func (suite *ServiceSuite) Test_FindImportWithForcedMetadata_WithMetadata_Success() {
	// Prepare
	ir := testutils.NewImportRepository()
	s := imports.NewService(ir)
	ctx := context.Background()
	id := uuid.New()
	ir.Add(imports.New(id, uuid.New(), imports.WithMetadata(1, 2, 3, 4, 5)).ToDTO())

	// Execute
	i, err := s.FindImportWithForcedMetadata(ctx, id)

	// Success
	suite.NoError(err)
	suite.Equal(id, i.ID())
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(5, i.FailedJobs())
}

func (suite *ServiceSuite) Test_FindImportWithForcedMetadata_Fail() {
	// Prepare
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	s := imports.NewService(ir)
	ctx := context.Background()

	// Execute
	i, err := s.FindImportWithForcedMetadata(ctx, uuid.New())

	// Fail
	suite.Nil(i)
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}
