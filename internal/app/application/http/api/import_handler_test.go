package api_test

import (
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestImportHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImportHandlerSuite))
}

type ImportHandlerSuite struct {
	suite.Suite
}

func (suite *ImportHandlerSuite) Test_List_Success() {
	// Prepare
	chID := uuid.New()
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelName("Channel Name"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(id1),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusCompleted),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)),
			testutils.WithImportEndedAt(time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC)),
			testutils.WithImportError("happened this error"),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNew, 1),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeUpdated, 2),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNoChange, 3),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeMissing, 4),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeError, 5),
		),
		testutils.WithImport(
			testutils.WithImportID(id2),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)),
		),
		testutils.WithImport(
			testutils.WithImportID(id3),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)),
		),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"imports":[{"id":"`+id1.String()+`","channel_id":"`+chID.String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"completed","started_at":"2020-01-01T00:00:03Z","ended_at":"2020-01-01T00:00:04Z","error":"happened this error","new_jobs":1,"updated_jobs":2,"no_change_jobs":3,"missing_jobs":4,"failed_jobs":5,"total_jobs":10},{"id":"`+id2.String()+`","channel_id":"`+chID.String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"pending","started_at":"2020-01-01T00:00:02Z","ended_at":null,"error":null,"new_jobs":0,"updated_jobs":0,"no_change_jobs":0,"missing_jobs":0,"failed_jobs":0,"total_jobs":0},{"id":"`+id3.String()+`","channel_id":"`+chID.String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"pending","started_at":"2020-01-01T00:00:01Z","ended_at":null,"error":null,"new_jobs":0,"updated_jobs":0,"no_change_jobs":0,"missing_jobs":0,"failed_jobs":0,"total_jobs":0}]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ImportHandlerSuite) Test_List_ImportRepositoryFail() {
	// Prepare
	chID := uuid.New()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelName("Channel Name"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(id),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusCompleted),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)),
			testutils.WithImportEndedAt(time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)),
			testutils.WithImportError("happened this error"),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNew, 1),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeUpdated, 2),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNoChange, 3),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeMissing, 4),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeError, 5),
		),
		testutils.WithImportRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "failed to get imports: boom")
}

func (suite *ImportHandlerSuite) Test_List_ChannelRepositoryFail() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "failed to get channels: boom")
}

func (suite *ImportHandlerSuite) Test_List_BadResponseWriterFail() {
	// Prepare
	chID := uuid.New()
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelName("Channel Name"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(id1),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusCompleted),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)),
			testutils.WithImportEndedAt(time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC)),
			testutils.WithImportError("happened this error"),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNew, 1),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeUpdated, 2),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNoChange, 3),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeMissing, 4),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeError, 5),
		),
		testutils.WithImport(
			testutils.WithImportID(id2),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)),
		),
		testutils.WithImport(
			testutils.WithImportID(id3),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)),
		),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := testutils.NewBadResponseWriter()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)

	// Assert log
	log := dsl.LogLines()
	suite.Len(log, 2)
	suite.Contains(log[0], `"level":"ERROR"`)
	suite.Contains(log[0], `"msg":"failed to encode response: bad response writer"`)
	suite.Contains(log[1], `"level":"ERROR"`)
	suite.Contains(log[1], `"msg":"bad response writer"`)
}

func (suite *ImportHandlerSuite) Test_Find_Success() {
	// Prepare
	chID := uuid.New()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelName("Channel Name"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(id),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusCompleted),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)),
			testutils.WithImportEndedAt(time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)),
			testutils.WithImportError("happened this error"),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNew, 1),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeUpdated, 2),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNoChange, 3),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeMissing, 4),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeError, 5),
		),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+id.String()+`","channel_id":"`+chID.String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"completed","started_at":"2020-01-01T00:00:01Z","ended_at":"2020-01-01T00:00:02Z","error":"happened this error","new_jobs":1,"updated_jobs":2,"no_change_jobs":3,"missing_jobs":4,"failed_jobs":5,"total_jobs":10}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ImportHandlerSuite) Test_Find_NotFound() {
	// Prepare
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("GET", "/api/imports/"+uuid.New().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"import not found"}}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ImportHandlerSuite) Test_Find_ImportRepositoryFail() {
	// Prepare
	chID := uuid.New()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelName("Channel Name"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(id),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusCompleted),
			testutils.WithImportStartedAt(time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)),
			testutils.WithImportEndedAt(time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)),
			testutils.WithImportError("happened this error"),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNew, 1),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeUpdated, 2),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeNoChange, 3),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeMissing, 4),
			testutils.WithImportMetrics(aggregator.ImportMetricTypeError, 5),
		),
		testutils.WithImportRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("GET", "/api/imports/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], `failed to find import `+id.String()+`: boom`)
}
