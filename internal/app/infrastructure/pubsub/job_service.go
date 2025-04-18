package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-protobuf/build/gen/events/jobs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JobService struct {
	client *client
}

func NewJobService(topic *pubsub.Topic, cfg Config) *JobService {
	return &JobService{
		client: newClient(topic, cfg.Timeout),
	}
}

func (s *JobService) PublishJobInformation(ctx context.Context, job *aggregator.Job) error {
	msg := jobs.JobInformation{
		Id:          job.ID.String(),
		ChannelId:   job.ChannelID.String(),
		Title:       job.Title,
		Description: job.Description,
		Url:         job.URL,
		Source:      job.Source,
		Location:    job.Location,
		PostedAt:    timestamppb.New(job.PostedAt),
		Remote:      job.Remote,
	}

	return s.client.publish(ctx, &msg)
}

func (s *JobService) PublishJobMissing(ctx context.Context, job *aggregator.Job) error {
	msg := jobs.JobMissing{
		Id: job.ID.String(),
	}

	return s.client.publish(ctx, &msg)
}
