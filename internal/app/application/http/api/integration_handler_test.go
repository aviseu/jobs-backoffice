package api_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationHandler(t *testing.T) {
	suite.Run(t, new(IntegrationHandlerSuite))
}

type IntegrationHandlerSuite struct {
	suite.Suite
}

func (suite *IntegrationHandlerSuite) Test_ListIntegrations_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	r := testutils.NewChannelRepository()
	s := configuring.NewService(r)
	h := http.APIRootHandler(s, nil, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("GET", "/api/integrations", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"integrations":["arbeitnow"]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}
