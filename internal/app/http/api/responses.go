package api

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/domain/channel"
)

type ChannelResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Integration string `json:"integration"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewChannelResponse(ch *channel.Channel) *ChannelResponse {
	return &ChannelResponse{
		ID:          ch.ID().String(),
		Name:        ch.Name(),
		Integration: ch.Integration().String(),
		Status:      ch.Status().String(),
		CreatedAt:   ch.CreatedAt().Format(time.RFC3339),
		UpdatedAt:   ch.UpdatedAt().Format(time.RFC3339),
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

func NewListChannelsResponse(channels []*channel.Channel) *ListChannelsResponse {
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

func NewListIntegrationsResponse(integrations []channel.Integration) *IntegrationsResponse {
	resp := &IntegrationsResponse{
		Integrations: make([]string, 0, len(integrations)),
	}

	for _, i := range integrations {
		resp.Integrations = append(resp.Integrations, i.String())
	}

	return resp
}
