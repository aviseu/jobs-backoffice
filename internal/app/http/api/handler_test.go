package api_test

import (
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/http"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) Test_Create_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("POST", "/api/channels", strings.NewReader(`{"name":"Channel Name","integration":"arbeitnow"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	var ch *channel.Channel
	for _, c := range r.Channels {
		ch = c
	}
	suite.Equal("Channel Name", ch.Name())
	suite.Equal(channel.IntegrationArbeitnow, ch.Integration())
	suite.Equal(channel.StatusInactive, ch.Status())
	suite.True(ch.CreatedAt().After(time.Now().Add(-2 * time.Second)))
	suite.True(ch.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusCreated, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID().String()+`","name":"Channel Name","integration":"arbeitnow","status":"inactive","created_at":"`+ch.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt().Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_Create_Validation_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_Create_RepositoryFail_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_GetChannels_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch1 := channel.New(uuid.New(), "channel 1", channel.IntegrationArbeitnow, channel.StatusActive, channel.WithTimestamps(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC)))
	r.Add(ch1)
	ch2 := channel.New(uuid.New(), "channel 2", channel.IntegrationArbeitnow, channel.StatusActive, channel.WithTimestamps(time.Date(2025, 1, 1, 0, 3, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 4, 0, 0, time.UTC)))
	r.Add(ch2)

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

func (suite *HandlerSuite) Test_FindChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

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

func (suite *HandlerSuite) Test_FindChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_FindChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_FindChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_UpdateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+ch.ID().String(), strings.NewReader(`{"name":"New Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal("New Name", c.Name())
	suite.Equal(channel.IntegrationArbeitnow, c.Integration())
	suite.Equal(channel.StatusActive, c.Status())
	suite.True(c.CreatedAt().Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+ch.ID().String()+`","name":"New Name","integration":"arbeitnow","status":"active","created_at":"`+ch.CreatedAt().Format(time.RFC3339)+`","updated_at":"`+ch.UpdatedAt().Format(time.RFC3339)+`"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_UpdateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+uuid.New().String(), strings.NewReader(`{"name":"New Name"}`))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to find channel: channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_UpdateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/invalid-uuid", strings.NewReader(`{"name":"New Name"}`))
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

func (suite *HandlerSuite) Test_UpdateChannel_Validation_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

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
	suite.Equal("channel 1", c.Name())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_UpdateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

	req, err := oghttp.NewRequest("PATCH", "/api/channels/"+ch.ID().String(), strings.NewReader(`{"name":"New Name"}`))
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
	suite.Equal("channel 1", c.Name())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *HandlerSuite) Test_ActivateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusInactive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(channel.StatusActive, c.Status())
	suite.True(c.CreatedAt().Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_ActivateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/activate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to find channel: channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_ActivateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

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

func (suite *HandlerSuite) Test_ActivateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusInactive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

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
	suite.Equal(channel.StatusInactive, c.Status())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}

func (suite *HandlerSuite) Test_DeactivateChannel_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state change
	suite.Len(r.Channels, 1)
	c := r.First()
	suite.Equal(channel.StatusInactive, c.Status())
	suite.True(c.CreatedAt().Equal(time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)))
	suite.True(c.UpdatedAt().After(time.Now().Add(-2 * time.Second)))

	// Assert response
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_DeactivateChannel_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+uuid.New().String()+"/deactivate", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal("{\"error\":{\"message\":\"failed to find channel: channel not found\"}}\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) Test_DeactivateChannel_InvalidID() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	req, err := oghttp.NewRequest("PUT", "/api/channels/invalid-uuid/deactivate", nil)
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

func (suite *HandlerSuite) Test_DeactivateChannel_Error_Fail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	r.FailWith(errors.New("boom!"))
	s := channel.NewService(r)
	h := http.APIRootHandler(s, log)

	ch := channel.New(
		uuid.New(),
		"channel 1",
		channel.IntegrationArbeitnow,
		channel.StatusActive,
		channel.WithTimestamps(
			time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		),
	)
	r.Add(ch)

	req, err := oghttp.NewRequest("PUT", "/api/channels/"+ch.ID().String()+"/deactivate", nil)
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
	suite.Equal(channel.StatusActive, c.Status())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "boom!")
}
