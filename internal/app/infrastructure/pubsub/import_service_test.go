package pubsub_test

import (
	gpubsub "cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/aviseu/jobs-protobuf/build/gen/commands/imports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"sync"
	"testing"
	"time"
)

func TestImportService(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ImportServiceSuite))
}

type ImportServiceSuite struct {
	testutils.PubSubSuite
}

func (suite *ImportServiceSuite) Test_PublishImport_Success() {
	// Prepare
	ctx := context.Background()
	s := pubsub.NewImportService(suite.ImportTopic, pubsub.Config{Timeout: 1 * time.Second})
	importID := uuid.New()

	var resp imports.ExecuteImportChannel
	subCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := suite.ImportSubscription.Receive(subCtx, func(ctx context.Context, msg *gpubsub.Message) {
			if err := proto.Unmarshal(msg.Data, &resp); err != nil {
				suite.Fail(fmt.Errorf("failed to unmarshal message: %w", err).Error())
			}
			msg.Ack()
			cancel()
		})
		suite.NoError(err)
		wg.Done()
	}()

	// Execute
	err := s.PublishImportCommand(ctx, importID)

	// Assert result
	suite.NoError(err)

	// Assert message received by subscription attached to the topic
	wg.Wait()
	suite.Equal(importID.String(), resp.ImportId)
}

func (suite *ImportServiceSuite) Test_PublishImport_ConnectionFailed() {
	// Prepare
	ctx := context.Background()
	s := pubsub.NewImportService(suite.BadImportTopic, pubsub.Config{Timeout: 1 * time.Second})
	importID := uuid.New()

	// Execute
	err := s.PublishImportCommand(ctx, importID)

	// Assert result
	suite.Error(err)
	suite.ErrorContains(err, "failed to publish pubsub message")
}
