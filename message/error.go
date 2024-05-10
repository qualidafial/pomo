package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Err(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{
			Err: err,
		}
	}
}

type ErrMsg struct {
	Err error
}
