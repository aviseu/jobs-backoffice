package configuring

import "github.com/google/uuid"

type CreateChannelCommand struct {
	Name        string
	Integration string
}

func NewCreateChannelCommand(name, integration string) *CreateChannelCommand {
	return &CreateChannelCommand{
		Name:        name,
		Integration: integration,
	}
}

type UpdateChannelCommand struct {
	Name string
	ID   uuid.UUID
}

func NewUpdateChannelCommand(id uuid.UUID, name string) *UpdateChannelCommand {
	return &UpdateChannelCommand{
		ID:   id,
		Name: name,
	}
}
