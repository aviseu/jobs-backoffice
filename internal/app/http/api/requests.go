package api

type CreateChannelRequest struct {
	Name        string `json:"name"`
	Integration string `json:"integration"`
}

func NewCreateChannelRequest(name, integration string) *CreateChannelRequest {
	return &CreateChannelRequest{
		Name:        name,
		Integration: integration,
	}
}
