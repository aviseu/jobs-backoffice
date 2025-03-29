package base_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestChannel(t *testing.T) {
	suite.Run(t, new(ChannelSuite))
}

type ChannelSuite struct {
	suite.Suite
}

func (suite *ChannelSuite) Test_ParseIntegration_Success() {
	// Execute
	i, ok := base.ParseIntegration("arbeitnow")

	// Assert
	suite.True(ok)
	suite.Equal(base.IntegrationArbeitnow, i)
}

func (suite *ChannelSuite) Test_ParseIntegration_Error() {
	// Execute
	_, ok := base.ParseIntegration("invalid")

	// Assert
	suite.False(ok)
}
