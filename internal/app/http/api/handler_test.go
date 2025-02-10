package api_test

import (
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/http"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
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

	req, err := oghttp.NewRequest("POST", "/api/channels", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal("CreateChannel", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}
