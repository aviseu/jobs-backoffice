package arbeitnow

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/google/uuid"
)

const endpointJobBoard = "/api/job-board-api"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	c       *client
	ch      *channel.Channel
	baseURL string
}

func NewService(c HTTPClient, cfg Config, ch *channel.Channel) *Service {
	return &Service{
		c:       newClient(c),
		baseURL: cfg.URL,
		ch:      ch,
	}
}

func (s *Service) GetJobs() ([]*job.Job, error) {
	jobs := make([]*jobEntry, 0)
	page := 1
	endpoint := s.baseURL + endpointJobBoard
	for {
		resp, err := s.c.JobBoard(endpoint, s.ch)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs page %d on channel %s: %w", page, s.ch.ID(), err)
		}

		jobs = append(jobs, resp.Jobs...)

		if !resp.Links.Next.Valid {
			break
		}

		endpoint = resp.Links.Next.String
		page++
	}

	result := make([]*job.Job, 0, len(jobs))
	for _, j := range jobs {
		result = append(result, job.New(
			uuid.NewSHA1(s.ch.ID(), []byte(j.Slug)), // UUID V5
			s.ch.ID(),
			job.StatusActive,
			j.URL,
			j.Title,
			j.Description,
			s.ch.Integration().String(),
			j.Location,
			j.Remote,
			time.Unix(j.CreatedAt, 0),
		))
	}

	return result, nil
}

func (s *Service) Channel() *channel.Channel {
	return s.ch
}
