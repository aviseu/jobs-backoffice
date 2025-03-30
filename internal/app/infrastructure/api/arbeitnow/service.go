package arbeitnow

import (
	"fmt"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const endpointJobBoard = "/api/job-board-api"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	c       *client
	ch      *aggregator.Channel
	baseURL string
}

func NewService(c HTTPClient, cfg Config, ch *aggregator.Channel) *Service {
	return &Service{
		c:       newClient(c),
		baseURL: cfg.URL,
		ch:      ch,
	}
}

func (s *Service) GetJobs() ([]*aggregator.Job, error) {
	jobs := make([]*jobEntry, 0)
	page := 1
	endpoint := s.baseURL + endpointJobBoard
	for {
		resp, err := s.c.JobBoard(endpoint, s.ch)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs page %d on channel %s: %w", page, s.ch.ID, err)
		}

		jobs = append(jobs, resp.Jobs...)

		if !resp.Links.Next.Valid {
			break
		}

		endpoint = resp.Links.Next.String
		page++
	}

	result := make([]*aggregator.Job, 0, len(jobs))
	for _, j := range jobs {
		result = append(result, &aggregator.Job{
			ID:          uuid.NewSHA1(s.ch.ID, []byte(j.Slug)), // UUID V5
			ChannelID:   s.ch.ID,
			Status:      aggregator.JobStatusActive,
			URL:         j.URL,
			Title:       j.Title,
			Description: j.Description,
			Location:    j.Location,
			Remote:      j.Remote,
			PostedAt:    time.Unix(j.CreatedAt, 0),
			Source:      aggregator.IntegrationArbeitnow.String(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return result, nil
}

func (s *Service) Channel() *aggregator.Channel {
	return s.ch
}
