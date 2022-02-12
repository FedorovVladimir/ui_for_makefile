package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"learn_console_ui/ui"
	"os"
	"os/exec"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	width  int
	height int
	state  int

	page  int
	lists []list.Model

	currentCommand *ui.Command
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			item := m.lists[m.page].SelectedItem()
			m.currentCommand = item.(*ui.Command)
			return m, tea.Quit
		case "right":
			m.page++
			if m.page == len(m.lists) {
				m.page = 0
			}
			top, right, bottom, left := docStyle.GetMargin()
			m.lists[m.page].SetSize(m.width-left-right, m.height-top-bottom-1)
		case "left":
			m.page--
			if m.page == -1 {
				m.page = len(m.lists) - 1
			}
			top, right, bottom, left := docStyle.GetMargin()
			m.lists[m.page].SetSize(m.width-left-right, m.height-top-bottom-1)
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.lists[m.page].SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.lists[m.page], cmd = m.lists[m.page].Update(msg)
	return m, cmd
}

func (m *model) View() string {
	s := docStyle.Render(m.lists[m.page].View())
	s += m.ViewMenu()
	return s
}

func (m *model) ViewMenu() string {
	style := lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	s := "\n"
	for i, l := range m.lists {
		if m.page == i {
			s += style.Render(" " + l.Title + " ")
			continue
		}
		s += " " + l.Title + " "
	}
	return s
}

func main() {
	commands, err := ui.ScanCommands()
	if err != nil {
		panic(err)
	}
	m := &model{lists: []list.Model{}}
	for k, v := range commands {
		m.lists = append(m.lists, ui.ConvertListCommandToListModel(k, v))
	}

	for {
		p := tea.NewProgram(m, tea.WithAltScreen())
		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		var cmd *exec.Cmd
		if len(m.currentCommand.Args) == 0 {
			cmd = exec.Command("make", m.currentCommand.Name)
		}
		if len(m.currentCommand.Args) == 1 {
			arg := m.currentCommand.Args[0].String()
			cmd = exec.Command("make", m.currentCommand.Name, arg)
		}
		fmt.Println(cmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", err)
			os.Exit(2)
		}
	}
}
