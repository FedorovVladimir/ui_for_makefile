package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
)

const (
	stateMenu = iota
	stateRunCommand
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title   string
	desc    string
	command command
}

func (i item) Title() string { return i.title }

func (i item) Description() string { return i.desc }

func (i item) FilterValue() string { return i.title }

type model struct {
	width  int
	height int
	state  int

	page  int
	lists []list.Model

	currentCommand command
}

func (m *model) Init() tea.Cmd {
	m.state = stateMenu
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == stateMenu {
				m.currentCommand = m.lists[m.page].SelectedItem().(item).command
				m.runCommand()
				return m, tea.Quit
			}
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
	if m.state == stateMenu {
		s := docStyle.Render(m.lists[m.page].View())
		s += m.ViewMenu()
		return s
	}
	return "111"
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

func (m *model) runCommand() {
	//fmt.Println("command in")
	//defer func() {
	//	fmt.Println("command out")
	//}()
	//cmd := exec.Command("make", m.currentCommand.name)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//p = tea.NewProgram(m, tea.WithOutput(nil), tea.WithInput(nil))
	//err := cmd.Run()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("tab enter for out")
}

type command struct {
	name    string
	hasArgs bool
	text    string
}

func main() {
	m := &model{lists: []list.Model{}}

	itemsTab1 := []list.Item{
		item{title: "Макс", desc: "Вперед к победе", command: command{name: "command_1"}},
		item{title: "item 2", desc: "description 2", command: command{name: "command_1"}},
	}
	m.lists = append(m.lists, list.New(itemsTab1, list.NewDefaultDelegate(), 0, 0))
	m.lists[0].Title = "Tab 1"

	itemsTab2 := []list.Item{
		item{title: "item 1", desc: "description 1", command: command{name: "command_1"}},
		item{title: "item 2", desc: "description 2", command: command{name: "command_1"}},
	}
	m.lists = append(m.lists, list.New(itemsTab2, list.NewDefaultDelegate(), 0, 0))
	m.lists[1].Title = "Tab 2"

	for {
		p := tea.NewProgram(m, tea.WithAltScreen())
		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		cmd := exec.Command("make", "command_1")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		p = tea.NewProgram(m, tea.WithOutput(nil), tea.WithInput(nil))
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}
