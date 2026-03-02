package message

import (
	tea "charm.land/bubbletea/v2"
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
