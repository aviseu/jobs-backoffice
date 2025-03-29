package base_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

type IntegrationSuite struct {
	suite.Suite
}

func (suite *IntegrationSuite) Test_ParseIntegration_Success() {
	// Execute
	i, ok := base.ParseIntegration("arbeitnow")

	// Assert
	suite.True(ok)
	suite.Equal(base.IntegrationArbeitnow, i)
}

func (suite *IntegrationSuite) Test_ParseIntegration_Error() {
	// Execute
	_, ok := base.ParseIntegration("invalid")

	// Assert
	suite.False(ok)
}

func (suite *IntegrationSuite) Test_ListIntegrations_Success() {
	// Execute
	list := base.ListIntegrations()

	// Assert
	suite.Len(list, 1)
	suite.Equal(base.IntegrationArbeitnow, list[0])
}
