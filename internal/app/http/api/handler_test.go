package api_test

import (
	"github.com/aviseu/jobs/internal/app/http"
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
	h := http.APIRootHandler()

	req, err := oghttp.NewRequest("POST", "/api/channels", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal("CreateChannel", rr.Body.String())
}
