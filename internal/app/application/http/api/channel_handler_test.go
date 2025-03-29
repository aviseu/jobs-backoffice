package api_test

import (
	"encoding/json"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http/api"
	"github.com/aviseu/jobs-backoffice/internal/app/domain"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/storage/postgres"
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
	suite.Run(t, new(ChannelHandlerSuite))
}

type ChannelHandlerSuite struct {
	suite.Suite
}

func (suite *ChannelHandlerSuite) Test_Create_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"Channel Name","integration":"arbeitnow"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	var ch *postgres.Channel
	for _, c := range r.Channels {
		ch = c
	}
	suite.Equal("Channel Name", ch.Name)
	suite.Equal(base.IntegrationArbeitnow, ch.Integration)
	suite.Equal(base.ChannelStatusInactive, ch.Status)
	suite.True(ch.CreatedAt.After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusCreated, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID.String()+`","name":"Channel Name","integration":"arbeitnow","status":"inactive","created_at":"`+ch.CreatedAt.Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt.Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_Create_Validation_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"","integration":"bad_integration"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to find integration bad_integration: invalid integration\\nname is required\"}}\n", rr.Body.String())

	// Assert state change
	suite.Empty(r.Channels)

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_Create_RepositoryFail_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"Channel Name","integration":"arbeitnow"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Empty(r.Channels)

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "failed to create channel: boom!")
}

func (suite *ChannelHandlerSuite) Test_GetChannels_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch1 := configuring.NewChannel(uuid.New(), "channel 1", base.IntegrationArbeitnow, base.ChannelStatusActive, configuring.WithTimestamps(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)))
	r.Add(ch1.ToDTO())
	ch2 := configuring.NewChannel(uuid.New(), "channel 2", base.IntegrationArbeitnow, base.ChannelStatusActive, configuring.WithTimestamps(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC)))
	r.Add(ch2.ToDTO())

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"channels":[{"id":"`+ch1.ID().String()+`","name":"channel 1","integration":"arbeitnow","status":"active","created_at":"`+ch1.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+ch1.UpdatedAt().Format(time.RFC3339)+`"},{"id":"`+ch2.ID().String()+`","name":"channel 2","integration":"arbeitnow","status":"active","created_at":"`+ch2.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+ch2.UpdatedAt().Format(time.RFC3339)+`"}]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_GetChannels_WithCors_Success() {
	// Prepare
	_, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{Cors: true}, log)

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("http://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	suite.Equal("Link", rr.Header().Get("Access-Control-Expose-Headers"))
	suite.Equal("Origin", rr.Header().Get("Vary"))
}

func (suite *ChannelHandlerSuite) Test_GetChannels_WithoutCors_Success() {
	// Prepare
	_, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{Cors: false}, log)

	req, err := oghttp.NewRequest("GET", "/api/channels", nil)
	suite.NoError(err)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Empty(rr.Header().Get("Access-Control-Allow-Origin"))
	suite.Empty(rr.Header().Get("Access-Control-Expose-Headers"))
	suite.Empty(rr.Header().Get("Vary"))
}

func (suite *ChannelHandlerSuite) Test_FindChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("GET", "/api/channels/"+ch.ID().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID().String()+`","name":"channel 1","integration":"arbeitnow","status":"active","created_at":"`+ch.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt().Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("GET", "/api/channels/"+uuid.New().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("GET", "/api/channels/invalid-uuid", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_FindChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("GET", "/api/channels/"+uuid.New().String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+ch.ID().String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal("NewChannel Name", c.Name)
	suite.Equal(base.IntegrationArbeitnow, c.Integration)
	suite.Equal(base.ChannelStatusActive, c.Status)
	suite.True(c.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID().String()+`","name":"NewChannel Name","integration":"arbeitnow","status":"active","created_at":"`+ch.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+c.UpdatedAt.Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+uuid.New().String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/invalid-uuid", strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_Validation_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+ch.ID().String(), strings.NewReader(`{"name":""}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to update channel: name is required\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal("channel 1", c.Name)

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_UpdateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+ch.ID().String(), strings.NewReader(`{"name":"NewChannel Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal("channel 1", c.Name)

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(base.ChannelStatusActive, c.Status)
	suite.True(c.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/invalid-uuid/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_ActivateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusInactive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(base.ChannelStatusInactive, c.Status)

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(base.ChannelStatusInactive, c.Status)
	suite.True(c.CreatedAt.Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt.After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/invalid-uuid/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to parse post uuid invalid-uuid: invalid UUID length: 12\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ChannelHandlerSuite) Test_DeactivateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, r, nil, nil, http.Config{}, log)

	ch := configuring.NewChannel(
		uuid.New(),
		"channel 1",
		base.IntegrationArbeitnow,
		base.ChannelStatusActive,
		configuring.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"Internal Server Error\"}}\n", rr.Body.String())

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(base.ChannelStatusActive, c.Status)

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	ps := testutils.NewPubSubService()
	ia := domain.NewScheduleImportAction(is, ps, log)
	h := http.APIRootHandler(chs, chr, is, ia, http.Config{}, log)

	chID := uuid.New()
	ch := configuring.NewChannel(chID, "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive)
	chr.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/schedule", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	var resp api.ImportResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	suite.NoError(err)
	suite.Equal(ch.ID().String(), resp.ChannelID)
	iID := uuid.MustParse(resp.ID)

	// Assert state change
	suite.Len(ir.Imports, 1)
	i, ok := ir.Imports[iID]
	suite.True(ok)
	suite.Equal(iID, i.ID)
	suite.Equal(chID, i.ChannelID)
	suite.Equal(base.ImportStatusPending, i.Status)

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "scheduling import for channel "+chID.String())

	// Assert pubsub
	suite.Len(ps.ImportIDs, 1)
	suite.Equal(i.ID, ps.ImportIDs[0])
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_ChannelNotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	is := importing.NewImportService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	ps := testutils.NewPubSubService()
	ia := domain.NewScheduleImportAction(is, ps, log)
	h := http.APIRootHandler(chs, chr, is, ia, http.Config{}, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/schedule", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"channel not found: channel not found"}}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf)

	// Assert pubsub
	suite.Len(ps.ImportIDs, 0)
}

func (suite *ChannelHandlerSuite) Test_ScheduleImport_ImportRepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	is := importing.NewImportService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	ps := testutils.NewPubSubService()
	ia := domain.NewScheduleImportAction(is, ps, log)
	h := http.APIRootHandler(chs, chr, is, ia, http.Config{}, log)

	chID := uuid.New()
	ch := configuring.NewChannel(chID, "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive)
	chr.Add(ch.ToDTO())

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/schedule", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 2)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "scheduling import for channel "+chID.String())
	suite.Contains(lines[1], `"level":"ERROR"`)
	suite.Contains(lines[1], "failed to schedule import")
	suite.Contains(lines[1], "boom!")

	// Assert pubsub
	suite.Len(ps.ImportIDs, 0)
}
