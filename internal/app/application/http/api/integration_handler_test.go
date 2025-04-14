package api_test

import (
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationHandlerSuite))
}

type IntegrationHandlerSuite struct {
	suite.Suite
}

func (suite *IntegrationHandlerSuite) Test_ListIntegrations_Success() {
	// Prepare
	suite.T().Parallel()
	dsl := testutils.NewDSL()

	req, err := oghttp.NewRequest("GET", "/api/integrations", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.APIServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"integrations":["arbeitnow"]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(dsl.LogLines())
}
