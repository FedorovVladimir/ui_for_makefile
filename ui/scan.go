package ui

import (
	"bufio"
	"github.com/charmbracelet/bubbles/list"
	"os"
	"strings"
)

func ScanCommands() (map[string][]*Command, error) {
	c := map[string][]*Command{}
	tabName := "default"
	description := ""

	file, err := os.Open("Makefile")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## group: ") {
			tabName = strings.TrimPrefix(line, "## group: ")
			continue
		}
		if strings.HasPrefix(line, "## ") {
			if description == "" {
				description = strings.TrimPrefix(line, "## ")
				continue
			}
		}
		if strings.HasSuffix(line, ":") {
			commandName := strings.TrimSuffix(line, ":")
			if description == "" {
				description = commandName
			}
			c[tabName] = append(c[tabName], NewCommand(commandName, description))
			description = ""
			continue
		}
	}

	return c, nil
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
