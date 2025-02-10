package postgres_test

import (
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestIntegrationsSuite(t *testing.T) {
	suite.Run(t, new(IntegrationsSuite))
}

type IntegrationsSuite struct {
	testutils.IntegrationSuite
}

func (suite *IntegrationsSuite) TestIfIntegrationsTestsWork() {
	suite.NoError(suite.DB.Ping())
}
