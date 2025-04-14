package api

type CreateChannelRequest struct {
	Name        string `json:"name"`
	Integration string `json:"integration"`
}

type UpdateChannelRequest struct {
	Name string `json:"name"`
}
