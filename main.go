package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
)

const (
	stateMenu = iota
	stateLoading
	stateInputText
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

	spinner   spinner.Model
	textInput textinput.Model

	currentCommand command
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateInputText {
			switch msg.Type {
			case tea.KeyEnter:
				m.state = stateLoading
				go m.runCommand()
				return m, m.spinner.Tick
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.currentCommand = m.lists[m.page].SelectedItem().(item).command
			if m.currentCommand.hasArgs {
				m.state = stateInputText
				return m, textinput.Blink
			} else {
				m.state = stateLoading
				go m.runCommand()
				return m, m.spinner.Tick
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

	if m.state == stateLoading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	if m.state == stateInputText {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.lists[m.page], cmd = m.lists[m.page].Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.state == stateLoading {
		return fmt.Sprintf("\n\n   %s Loading...\n\n", m.spinner.View())
	}
	if m.state == stateInputText {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			m.currentCommand.text,
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}
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

func (m *model) runCommand() {
	var c *exec.Cmd
	if m.currentCommand.hasArgs {
		c = exec.Command("make", m.currentCommand.name, fmt.Sprintf(`arg="%s"`, m.textInput.Value()))
		m.textInput.SetValue("")
	} else {
		c = exec.Command("make", m.currentCommand.name)
	}
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	m.state = stateMenu
}

type command struct {
	name    string
	hasArgs bool
	text    string
}

func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "you text"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	m := &model{spinner: s, textInput: ti, lists: []list.Model{}}

	itemsTab1 := []list.Item{
		item{title: "item 1", desc: "description 1", command: command{name: "command_with_arg", hasArgs: true, text: "Input text for file:"}},
		item{title: "item 2", desc: "description 2", command: command{name: "command_2"}},
		item{title: "item 3", desc: "description 3", command: command{name: "command_3"}},
	}
	m.lists = append(m.lists, list.New(itemsTab1, list.NewDefaultDelegate(), 0, 0))
	m.lists[0].Title = "Tab 1"

	itemsTab2 := []list.Item{
		item{title: "item 1", desc: "description 1", command: command{name: "command_1"}},
		item{title: "item 2", desc: "description 2", command: command{name: "command_2"}},
		item{title: "item 3", desc: "description 3", command: command{name: "command_3"}},
	}
	m.lists = append(m.lists, list.New(itemsTab2, list.NewDefaultDelegate(), 0, 0))
	m.lists[0].Title = "Tab 2"

	itemsTab3 := []list.Item{
		item{title: "item 1", desc: "description 1", command: command{name: "command_1"}},
		item{title: "item 2", desc: "description 2", command: command{name: "command_2"}},
		item{title: "item 3", desc: "description 3", command: command{name: "command_3"}},
	}
	m.lists = append(m.lists, list.New(itemsTab3, list.NewDefaultDelegate(), 0, 0))
	m.lists[0].Title = "Tab 3"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
