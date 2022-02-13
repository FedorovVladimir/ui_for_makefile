package ui

type Command struct {
	Name        string
	description string
	Args        []arg
}

func NewCommand(name string, description string, opts ...CommandOptional) *Command {
	c := &Command{Name: name, description: description}
	for _, opt := range opts {
		opt(c)
	}
	return c
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

type arg struct {
	name        string
	description string
	Value       string
}

func (a *arg) String() string {
	return a.name + "=\"" + a.Value + "\""
}

type CommandOptional func(command *Command)

func WithArgs(args []arg) CommandOptional {
	return func(command *Command) {
		targetArgs := make([]arg, len(args))
		copy(targetArgs, args)
		command.Args = targetArgs
	}
}
