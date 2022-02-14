package main

import (
	"fmt"
	"github.com/FedorovVladimir/ui_for_makefile/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
)

func main() {
	commands, err := ui.ScanCommands()
	if err != nil {
		panic(err)
	}

	mLists := &ui.ModelLists{Lists: []list.Model{}}
	for k, v := range commands {
		mLists.Lists = append(mLists.Lists, ui.ConvertListCommandToListModel(k, v))
	}

	mInput := ui.InitialModel()

	for {
		p := tea.NewProgram(mLists, tea.WithAltScreen())
		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		if mLists.CurrentCommand == nil {
			break
		}

		for i, arg := range mLists.CurrentCommand.Args {
			if arg.Value == "" {
				p := tea.NewProgram(mInput)
				if err := p.Start(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
				mLists.CurrentCommand.Args[i].Value = mInput.TextInput.Value()
				mInput.TextInput.SetValue("")
			}
		}

		var cmd *exec.Cmd
		if len(mLists.CurrentCommand.Args) == 0 {
			cmd = exec.Command("make", mLists.CurrentCommand.Name)
		}
		if len(mLists.CurrentCommand.Args) == 1 {
			arg := mLists.CurrentCommand.Args[0].String()
			cmd = exec.Command("make", mLists.CurrentCommand.Name, arg)
			mLists.CurrentCommand.Args[0].Value = ""
			mLists.CurrentCommand = nil
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", err)
			os.Exit(2)
		}
	}
}
