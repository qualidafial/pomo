package board

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	ToggleHelp key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.ToggleHelp,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.ToggleHelp,
	}
}
