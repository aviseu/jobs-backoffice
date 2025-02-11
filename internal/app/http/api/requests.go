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

type UpdateChannelRequest struct {
	Name string `json:"name"`
}

func NewUpdateChannelRequest(name string) *UpdateChannelRequest {
	return &UpdateChannelRequest{
		Name: name,
	}
}
