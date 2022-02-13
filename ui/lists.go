package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ModelLists struct {
	width  int
	height int
	state  int

	page  int
	Lists []list.Model

	CurrentCommand *Command
}

func (m *ModelLists) Init() tea.Cmd {
	return nil
}

func (m *ModelLists) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			item := m.Lists[m.page].SelectedItem()
			m.CurrentCommand = item.(*Command)
			return m, tea.Quit
		case "right":
			m.page++
			if m.page == len(m.Lists) {
				m.page = 0
			}
			top, right, bottom, left := docStyle.GetMargin()
			m.Lists[m.page].SetSize(m.width-left-right, m.height-top-bottom-1)
		case "left":
			m.page--
			if m.page == -1 {
				m.page = len(m.Lists) - 1
			}
			top, right, bottom, left := docStyle.GetMargin()
			m.Lists[m.page].SetSize(m.width-left-right, m.height-top-bottom-1)
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.Lists[m.page].SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.Lists[m.page], cmd = m.Lists[m.page].Update(msg)
	return m, cmd
}

func (m *ModelLists) View() string {
	s := docStyle.Render(m.Lists[m.page].View())
	s += m.ViewMenu()
	return s
}

func (m *ModelLists) ViewMenu() string {
	style := lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	s := "\n"
	for i, l := range m.Lists {
		if m.page == i {
			s += style.Render(" " + l.Title + " ")
			continue
		}
		s += " " + l.Title + " "
	}
	return s
}
