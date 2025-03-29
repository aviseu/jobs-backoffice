package configuring

import "github.com/google/uuid"

type CreateCommand struct {
	Name        string
	Integration string
}

func NewCreateCommand(name, integration string) *CreateCommand {
	return &CreateCommand{
		Name:        name,
		Integration: integration,
	}
}

type UpdateCommand struct {
	Name string
	ID   uuid.UUID
}

func NewUpdateCommand(id uuid.UUID, name string) *UpdateCommand {
	return &UpdateCommand{
		ID:   id,
		Name: name,
	}
}
