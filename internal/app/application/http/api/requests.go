package api

type createChannelRequest struct {
	Name        string `json:"name"`
	Integration string `json:"integration"`
}

type updateChannelRequest struct {
	Name string `json:"name"`
}
