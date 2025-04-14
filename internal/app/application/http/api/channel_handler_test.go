package api_test

import (
	"encoding/json"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http/api"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestChannelHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ChannelHandlerSuite))
}

type ChannelHandlerSuite struct {
	suite.Suite
}

func (suite *ChannelHandlerSuite) Test_Create_Success() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"Channel Name","integration":"arbeitnow"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	ch := dsl.FirstChannel()
	suite.Equal("Channel Name", ch.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration)
	suite.Equal(aggregator.ChannelStatusInactive, ch.Status)
	suite.True(ch.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusCreated, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID.String()+`","name":"Channel Name","integration":"arbeitnow","status":"inactive","created_at":"`+ch.CreatedAt.Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt.Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_Create_Validation_Fail() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"","integration":"bad_integration"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to find integration bad_integration: invalid integration\\nname is required\"}}\n", rr.Body.String())

	// Assert state change
	suite.Empty(dsl.Channels())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_Create_ChannelRepositoryFail_Fail() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL(
		testutils.WithChannelRepositoryError(errors.New("boom!")),
	)

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"Channel Name","integration":"arbeitnow"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Empty(dsl.Channels())

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "failed to create channel: boom!")
}

func (suite *ChannelHandlerSuite) Test_GetChannels_Success() {
	// Prepare
	suite.T().Parallel()
	id1 := uuid.New()
	id2 := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id1),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
			testutils.WithChannelActivated(),
			testutils.WithChannelTimestamps(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)),
		),
		testutils.WithChannel(
			testutils.WithChannelID(id2),
			testutils.WithChannelName("channel 2"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
			testutils.WithChannelDeactivated(),
			testutils.WithChannelTimestamps(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC)),
		),
	)

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"channels":[{"id":"`+id1.String()+`","name":"channel 1","integration":"arbeitnow","status":"active","created_at":"`+dsl.Channel(id1).CreatedAt.Format(time.RFC3339)+`","updated_at":"`+dsl.Channel(id1).UpdatedAt.Format(time.RFC3339)+`"},{"id":"`+id2.String()+`","name":"channel 2","integration":"arbeitnow","status":"inactive","created_at":"`+dsl.Channel(id2).CreatedAt.Format(time.RFC3339)+`","updated_at":"`+dsl.Channel(id2).UpdatedAt.Format(time.RFC3339)+`"}]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_GetChannels_WithCors_Success() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL(
		testutils.WithHTTPConfig(http.Config{Cors: true}),
	)

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("http://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	suite.Equal("Link", rr.Header().Get("Access-Control-Expose-Headers"))
	suite.Equal("Origin", rr.Header().Get("Vary"))
}

func (suite *ChannelHandlerSuite) Test_GetChannels_WithoutCors_Success() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL(
		testutils.WithHTTPConfig(http.Config{Cors: false}),
	)

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Empty(rr.Header().Get("Access-Control-Allow-Origin"))
	suite.Empty(rr.Header().Get("Access-Control-Expose-Headers"))
	suite.Empty(rr.Header().Get("Vary"))
}

func (suite *ChannelHandlerSuite) Test_FindChannel_Success() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
			testutils.WithChannelActivated(),
			testutils.WithChannelTimestamps(cat, uat),
		),
	)

	req, err := oghttp.NewRequest("GET", "/api/channels/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+id.String()+`","name":"channel 1","integration":"arbeitnow","status":"active","created_at":"`+cat.Format(time.RFC3339)+`","updated_at":"`+uat.Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_NotFound() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("GET", "/api/channels/"+uuid.New().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_InvalidID() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("GET", "/api/channels/invalid-uuid", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_ChannelRepositoryFail() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL(
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("GET", "/api/channels/"+uuid.New().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom")
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_Success() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
			testutils.WithChannelActivated(),
			testutils.WithChannelTimestamps(cat, uat),
		),
	)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+id.String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	ch := dsl.FirstChannel()
	suite.Equal("NewChannel Name", ch.Name)
	suite.Equal(aggregator.IntegrationArbeitnow, ch.Integration)
	suite.Equal(aggregator.ChannelStatusActive, ch.Status)
	suite.True(ch.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(ch.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+id.String()+`","name":"NewChannel Name","integration":"arbeitnow","status":"active","created_at":"`+cat.Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt.Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_NotFound() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+uuid.New().String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_InvalidID() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PATCH", "/api/channels/invalid-uuid", strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_Validation_Fail() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelActivated(),
		),
	)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+id.String(), strings.NewReader(`{"name":""}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to update channel: name is required\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	c := dsl.FirstChannel()
	suite.Equal("channel 1", c.Name)

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_ChannelRepositoryFail() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	cat := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	uat := time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelName("channel 1"),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
			testutils.WithChannelActivated(),
			testutils.WithChannelTimestamps(cat, uat),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+id.String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	ch := dsl.FirstChannel()
	suite.Equal("channel 1", ch.Name)

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom")
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_Success() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelDeactivated(),
		),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	c := dsl.FirstChannel()
	suite.Equal(aggregator.ChannelStatusActive, c.Status)
	suite.True(c.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_NotFound() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_InvalidID() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PUT", "/api/channels/invalid-uuid/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_ChannelRepositoryFail() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelDeactivated(),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	c := dsl.FirstChannel()
	suite.Equal(aggregator.ChannelStatusInactive, c.Status)

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom")
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_Success() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	c := dsl.FirstChannel()
	suite.Equal(aggregator.ChannelStatusInactive, c.Status)
	suite.True(c.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_NotFound() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_InvalidID() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PUT", "/api/channels/invalid-uuid/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_ChannelRepositoryFail() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(dsl.Channels(), 1)
	ch := dsl.FirstChannel()
	suite.Equal(aggregator.ChannelStatusActive, ch.Status)

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom")
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_Success() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
			testutils.WithChannelActivated(),
		),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/schedule", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	var resp api.ImportResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	suite.NoError(err)
	suite.Equal(id.String(), resp.ChannelID)
	iID := uuid.MustParse(resp.ID)

	// Assert state change
	suite.Len(dsl.Imports(), 1)
	i := dsl.FirstImport()
	suite.Equal(iID, i.ID)
	suite.Equal(id, i.ChannelID)
	suite.Equal(aggregator.ImportStatusPending, i.Status)

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())

	// Assert pubsub
	suite.Len(dsl.PublishedImports(), 1)
	suite.Equal(i.ID, dsl.PublishedImports()[0])
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_ChannelNotFound() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/schedule", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"channel not found"}}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())

	// Assert pubsub
	suite.Len(dsl.PublishedImports(), 0)
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_ImportRepositoryFail() {
	// Prepare
	suite.T().Parallel()
	id := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithChannel(
			testutils.WithChannelID(id),
		),
		testutils.WithImportRepositoryError(errors.New("boom")),
	)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+id.String()+"/schedule", nil)
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
	suite.Len(lines, 2)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "scheduling import for channel "+id.String())
	suite.Contains(lines[1], `"level":"ERROR"`)
	suite.Contains(lines[1], "failed to schedule import")
	suite.Contains(lines[1], "boom")

	// Assert pubsub
	suite.Empty(dsl.PublishedImports())
}
