package ui

type Command struct {
	Name        string
	description string
}

func NewCommand(name string, description string) *Command {
	return &Command{Name: name, description: description}
}

func (c Command) Title() string {
	return c.Name
}

func (c Command) Description() string {
	return c.description
}

func (c Command) FilterValue() string {
	return c.Name + " " + c.description
}
