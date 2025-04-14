package aggregator_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationSuite))
}

type IntegrationSuite struct {
	suite.Suite
}

func (suite *IntegrationSuite) Test_ParseIntegration_Success() {
	// Execute
	i, ok := aggregator.ParseIntegration("arbeitnow")

	// Assert
	suite.True(ok)
	suite.Equal(aggregator.IntegrationArbeitnow, i)
}

func (suite *IntegrationSuite) Test_ParseIntegration_Error() {
	// Execute
	_, ok := aggregator.ParseIntegration("invalid")

	// Assert
	suite.False(ok)
}

func (suite *IntegrationSuite) Test_ListIntegrations_Success() {
	// Execute
	list := aggregator.ListIntegrations()

	// Assert
	suite.Len(list, 1)
	suite.Equal(aggregator.IntegrationArbeitnow, list[0])
}

func (suite *IntegrationSuite) Test_Integration_Success() {
	suite.Equal("arbeitnow", aggregator.IntegrationArbeitnow.String())
}
