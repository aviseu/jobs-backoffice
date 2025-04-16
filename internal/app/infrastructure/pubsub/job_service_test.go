package pubsub_test

import (
	gpubsub "cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/pubsub"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/aviseu/jobs-protobuf/build/gen/events/jobs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"sync"
	"testing"
	"time"
)

func TestJobService(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(JobServiceSuite))
}

type JobServiceSuite struct {
	testutils.PubSubSuite
}

func (suite *JobServiceSuite) Test_PublishJob_Success() {
	// Prepare
	ctx := context.Background()
	s := pubsub.NewJobService(suite.JobTopic, pubsub.Config{Timeout: 1 * time.Second})

	var resp jobs.JobInformation
	subCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := suite.JobSubscription.Receive(subCtx, func(ctx context.Context, msg *gpubsub.Message) {
			if err := proto.Unmarshal(msg.Data, &resp); err != nil {
				suite.Fail(fmt.Errorf("failed to unmarshal message: %w", err).Error())
			}
			msg.Ack()
			cancel()
		})
		suite.NoError(err)
		wg.Done()
	}()

	job := &aggregator.Job{
		ID:          uuid.New(),
		ChannelID:   uuid.New(),
		Title:       "Test Job",
		Description: "Test Description",
		URL:         "https://example.com",
		Source:      "Test Source",
		Location:    "Test Location",
		PostedAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Remote:      true,
	}

	// Execute
	err := s.PublishJobInformation(ctx, job)

	// Assert result
	suite.NoError(err)

	// Assert message received by subscription attached to the topic
	wg.Wait()
	suite.Equal(job.ID.String(), resp.Id)
	suite.Equal(job.ChannelID.String(), resp.ChannelId)
	suite.Equal("Test Job", resp.Title)
	suite.Equal("Test Description", resp.Description)
	suite.Equal("https://example.com", resp.Url)
	suite.Equal("Test Source", resp.Source)
	suite.Equal("Test Location", resp.Location)
	suite.True(resp.PostedAt.AsTime().Equal(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)))
	suite.Equal(job.Remote, resp.Remote)
}

func (suite *JobServiceSuite) Test_PublishImport_ConnectionFailed() {
	// Prepare
	ctx := context.Background()
	s := pubsub.NewJobService(suite.BadJobTopic, pubsub.Config{Timeout: 1 * time.Second})

	job := &aggregator.Job{
		ID:          uuid.New(),
		ChannelID:   uuid.New(),
		Title:       "Test Job",
		Description: "Test Description",
		URL:         "https://example.com",
		Source:      "Test Source",
		Location:    "Test Location",
		PostedAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Remote:      true,
	}

	// Execute
	err := s.PublishJobInformation(ctx, job)

	// Assert result
	suite.Error(err)
	suite.ErrorContains(err, "failed to publish pubsub message")
}
