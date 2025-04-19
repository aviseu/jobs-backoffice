package aggregator_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestImport(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImportSuite))
}

type ImportSuite struct {
	suite.Suite
}

func (suite *ImportSuite) Test_ImportStatus_Success() {
	suite.Equal("pending", aggregator.ImportStatusPending.String())
	suite.Equal("fetching", aggregator.ImportStatusFetching.String())
	suite.Equal("processing", aggregator.ImportStatusProcessing.String())
	suite.Equal("publishing", aggregator.ImportStatusPublishing.String())
	suite.Equal("completed", aggregator.ImportStatusCompleted.String())
	suite.Equal("failed", aggregator.ImportStatusFailed.String())
}

func (suite *ImportSuite) Test_ImportMetric_Success() {
	suite.Equal("new", aggregator.ImportMetricTypeNew.String())
	suite.Equal("updated", aggregator.ImportMetricTypeUpdated.String())
	suite.Equal("no_change", aggregator.ImportMetricTypeNoChange.String())
	suite.Equal("missing", aggregator.ImportMetricTypeMissing.String())
	suite.Equal("error", aggregator.ImportMetricTypeError.String())
	suite.Equal("publish", aggregator.ImportMetricTypePublish.String())
	suite.Equal("late_publish", aggregator.ImportMetricTypeLatePublish.String())
	suite.Equal("missing_publish", aggregator.ImportMetricTypeMissingPublish.String())
}

func (suite *ImportSuite) Test_Import_NoMetadata_Success() {
	// Prepare
	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: uuid.New(),
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
		Metrics: []*aggregator.ImportMetric{
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeNew},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeUpdated},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeUpdated},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeNoChange},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeNoChange},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeNoChange},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissing},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissing},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissing},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissing},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeError},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeError},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeError},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeError},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeError},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypePublish},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeLatePublish},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeLatePublish},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissingPublish},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissingPublish},
			{ID: uuid.New(), JobID: uuid.New(), MetricType: aggregator.ImportMetricTypeMissingPublish},
		},
	}

	// Assert
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(6, i.TotalJobs())
	suite.Equal(5, i.Errors())
	suite.Equal(1, i.Published())
	suite.Equal(2, i.LatePublished())
	suite.Equal(3, i.MissingPublished())
}

func (suite *ImportSuite) Test_Import_Metadata_Success() {
	// Prepare
	i := &aggregator.Import{
		ID:        uuid.New(),
		ChannelID: uuid.New(),
		Status:    aggregator.ImportStatusPending,
		StartedAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
		Metadata: &aggregator.ImportMetadata{
			New:              1,
			Updated:          2,
			NoChange:         3,
			Missing:          4,
			Errors:           5,
			Published:        1,
			LatePublished:    2,
			MissingPublished: 3,
		},
	}

	// Assert
	suite.Equal(1, i.NewJobs())
	suite.Equal(2, i.UpdatedJobs())
	suite.Equal(3, i.NoChangeJobs())
	suite.Equal(4, i.MissingJobs())
	suite.Equal(6, i.TotalJobs())
	suite.Equal(5, i.Errors())
	suite.Equal(1, i.Published())
	suite.Equal(2, i.LatePublished())
	suite.Equal(3, i.MissingPublished())
}
