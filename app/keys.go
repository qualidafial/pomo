package app

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	StartPomo key.Binding
	PausePomo key.Binding
	Quit      key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		StartPomo: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "start pomo"),
		),
		PausePomo: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "pause pomo"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.StartPomo,
		m.PausePomo,
		m.Quit,
	}
}
