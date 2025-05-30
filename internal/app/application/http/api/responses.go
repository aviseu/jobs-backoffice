package api

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"gopkg.in/guregu/null.v3"
)

type ChannelResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Integration string `json:"integration"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewChannelResponse(ch *aggregator.Channel) *ChannelResponse {
	return &ChannelResponse{
		ID:          ch.ID.String(),
		Name:        ch.Name,
		Integration: ch.Integration.String(),
		Status:      ch.Status.String(),
		CreatedAt:   ch.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ch.UpdatedAt.Format(time.RFC3339),
	}
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		},
	}
}

type ListChannelsResponse struct {
	Channels []*ChannelResponse `json:"channels"`
}

func NewListChannelsResponse(channels []*aggregator.Channel) *ListChannelsResponse {
	resp := &ListChannelsResponse{
		Channels: make([]*ChannelResponse, 0, len(channels)),
	}

	for _, ch := range channels {
		resp.Channels = append(resp.Channels, NewChannelResponse(ch))
	}

	return resp
}

type IntegrationsResponse struct {
	Integrations []string `json:"integrations"`
}

func NewListIntegrationsResponse(integrations []aggregator.Integration) *IntegrationsResponse {
	resp := &IntegrationsResponse{
		Integrations: make([]string, 0, len(integrations)),
	}

	for _, i := range integrations {
		resp.Integrations = append(resp.Integrations, i.String())
	}

	return resp
}

type ImportResponse struct {
	ID               string      `json:"id"`
	ChannelID        string      `json:"channel_id"`
	ChannelName      string      `json:"channel_name"`
	Integration      string      `json:"integration"`
	Status           string      `json:"status"`
	StartedAt        string      `json:"started_at"`
	EndedAt          null.String `json:"ended_at"`
	Error            null.String `json:"error"`
	NewJobs          int         `json:"new_jobs"`
	UpdatedJobs      int         `json:"updated_jobs"`
	NoChangeJobs     int         `json:"no_change_jobs"`
	MissingJobs      int         `json:"missing_jobs"`
	TotalJobs        int         `json:"total_jobs"`
	Errors           int         `json:"errors"`
	Published        int         `json:"published"`
	LatePublished    int         `json:"late_published"`
	MissingPublished int         `json:"missing_published"`
}

func NewImportResponse(i *aggregator.Import, ch *aggregator.Channel) *ImportResponse {
	ended := null.NewString("", false)
	if i.EndedAt.Valid {
		ended = null.StringFrom(i.EndedAt.Time.Format(time.RFC3339))
	}

	return &ImportResponse{
		ID:               i.ID.String(),
		ChannelID:        i.ChannelID.String(),
		ChannelName:      ch.Name,
		Integration:      ch.Integration.String(),
		Status:           i.Status.String(),
		StartedAt:        i.StartedAt.Format(time.RFC3339),
		EndedAt:          ended,
		Error:            i.Error,
		NewJobs:          i.NewJobs(),
		UpdatedJobs:      i.UpdatedJobs(),
		NoChangeJobs:     i.NoChangeJobs(),
		MissingJobs:      i.MissingJobs(),
		TotalJobs:        i.TotalJobs(),
		Errors:           i.Errors(),
		Published:        i.Published(),
		LatePublished:    i.LatePublished(),
		MissingPublished: i.MissingPublished(),
	}
}

type ImportsResponse struct {
	Imports []*ImportResponse `json:"imports"`
}

func NewImportsResponse(imports []*aggregator.Import, channels []*aggregator.Channel) *ImportsResponse {
	resp := &ImportsResponse{
		Imports: make([]*ImportResponse, 0, len(imports)),
	}

	for _, i := range imports {
		for _, ch := range channels {
			if ch.ID == i.ChannelID {
				resp.Imports = append(resp.Imports, NewImportResponse(i, ch))
				break
			}
		}
	}

	return resp
}
