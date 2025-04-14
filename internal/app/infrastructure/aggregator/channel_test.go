package aggregator_test

import (
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestChannel(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ChannelSuite))
}

type ChannelSuite struct {
	suite.Suite
}

func (suite *IntegrationSuite) Test_ChannelStatus_Success() {
	suite.Equal("inactive", aggregator.ChannelStatusInactive.String())
	suite.Equal("active", aggregator.ChannelStatusActive.String())
}
