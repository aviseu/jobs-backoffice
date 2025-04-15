package arbeitnow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
)

const ChannelHeader = "X-Channel-Id"

type client struct {
	c HTTPClient
}

func newClient(c HTTPClient) *client {
	return &client{
		c: c,
	}
}

func (c *client) JobBoard(endpoint string, ch *aggregator.Channel) (*jobBoardResponse, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for jobBoard: %w", err)
	}
	req.Header.Set(ChannelHeader, ch.ID.String())

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobBoard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get job board: %w", c.handleFailedResponse(resp))
	}

	var jobsResponse jobBoardResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w: %s", err, resp.Body)
	}

	return &jobsResponse, nil
}

func (*client) handleFailedResponse(resp *http.Response) error {
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if len(content) == 0 {
		return fmt.Errorf("failed to request with http code %d and no body", resp.StatusCode)
	}

	return fmt.Errorf("failed to request with http code %d and body: %s", resp.StatusCode, content)
}
