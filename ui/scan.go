package ui

import "github.com/charmbracelet/bubbles/list"

func ScanCommands() map[string][]*Command {
	c := map[string][]*Command{}
	c["Общие"] = []*Command{NewCommand("command_1", "command_1 description")}
	return c
}

func ConvertListCommandToListModel(title string, commands []*Command) list.Model {
	var items []list.Item
	for _, c := range commands {
		items = append(items, c)
	}
	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = title
	return m
}
