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
	var values []string
	var args []arg

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
		if strings.HasPrefix(line, "## list: ") {
			comment := strings.TrimPrefix(line, "## list: ")
			paths := strings.Split(comment, ",")
			for _, path := range paths {
				values = append(values, strings.TrimSpace(path))
			}
			continue
		}
		if strings.HasPrefix(line, "## ") && strings.Contains(line, ":") {
			comment := strings.TrimPrefix(line, "## ")
			paths := strings.Split(comment, ":")
			argName := strings.TrimSpace(paths[0])
			argDescription := strings.TrimSpace(paths[1])
			args = append(args, arg{
				name:        argName,
				description: argDescription,
				value:       "",
			})
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
			if len(args) > 0 {
				if len(values) > 0 {
					for _, value := range values {
						args[0].value = value
						c[tabName] = append(c[tabName], NewCommand(commandName, description+" "+value, WithArgs(args)))
					}
				} else {
					c[tabName] = append(c[tabName], NewCommand(commandName, description, WithArgs(args)))
				}
			} else {
				c[tabName] = append(c[tabName], NewCommand(commandName, description))
			}
			description = ""
			values = []string{}
			args = args[:0]
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
