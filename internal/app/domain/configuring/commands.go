package configuring

import "github.com/google/uuid"

type CreateChannelCommand struct {
	Name        string
	Integration string
}

func NewCreateCommand(name, integration string) *CreateChannelCommand {
	return &CreateChannelCommand{
		Name:        name,
		Integration: integration,
	}
}

type UpdateChannelCommand struct {
	Name string
	ID   uuid.UUID
}

func NewUpdateCommand(id uuid.UUID, name string) *UpdateChannelCommand {
	return &UpdateChannelCommand{
		ID:   id,
		Name: name,
	}
}
