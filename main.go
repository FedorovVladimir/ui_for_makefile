package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title   string
	desc    string
	command string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	page   int
	lists  []list.Model
	width  int
	height int

	loading bool
	spinner spinner.Model
	err     error
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
			item, _ := m.lists[m.page].SelectedItem().(item)
			m.loading = true
			go m.runCommand(item.command)
			return m, m.spinner.Tick
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
	default:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.lists[m.page], cmd = m.lists[m.page].Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.loading {
		return fmt.Sprintf("\n\n   %s Loading...\n\n", m.spinner.View())
	}
	s := docStyle.Render(m.lists[m.page].View())
	s += m.ViewMenu()
	return s
}

func (m *model) runCommand(command string) {
	c := exec.Command("make", command)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	m.loading = false
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
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	itemsTab1 := []list.Item{
		item{title: "item 1", desc: "description 1", command: "command_1"},
		item{title: "item 2", desc: "description 2", command: "command_2"},
		item{title: "item 3", desc: "description 3", command: "command_3"},
	}
	itemsTab2 := []list.Item{
		item{title: "item 1", desc: "description 1", command: "command_1"},
		item{title: "item 2", desc: "description 2", command: "command_2"},
		item{title: "item 3", desc: "description 3", command: "command_3"},
	}
	itemsTab3 := []list.Item{
		item{title: "item 1", desc: "description 1", command: "command_1"},
		item{title: "item 2", desc: "description 2", command: "command_2"},
		item{title: "item 3", desc: "description 3", command: "command_3"},
	}
	m := &model{
		lists: []list.Model{
			list.New(itemsTab1, list.NewDefaultDelegate(), 0, 0),
			list.New(itemsTab3, list.NewDefaultDelegate(), 0, 0),
			list.New(itemsTab2, list.NewDefaultDelegate(), 0, 0),
		},
		spinner: s,
	}
	m.lists[0].Title = "Tab 1"
	m.lists[1].Title = "Tab 2"
	m.lists[2].Title = "Tab 3"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
