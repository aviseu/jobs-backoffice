package api_test

import (
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/http"
	"github.com/aviseu/jobs/internal/testutils"
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
