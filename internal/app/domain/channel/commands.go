package channel

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
