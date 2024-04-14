package prompt

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Yes key.Binding
	No  key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Yes: key.NewBinding(
			key.WithKeys("y", "enter"),
			key.WithHelp("y/enter", "yes"),
		),
		No: key.NewBinding(
			key.WithKeys("n", "esc"),
			key.WithHelp("n/esc", "no"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.Yes,
			m.No,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Yes,
		m.No,
	}
}
