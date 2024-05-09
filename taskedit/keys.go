package taskedit

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	NextField key.Binding
	PrevField key.Binding

	Save   key.Binding
	Enter  key.Binding
	Cancel key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		NextField: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		PrevField: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev"),
		),

		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "save"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.Save,
			m.Cancel,
		},
		{
			m.NextField,
			m.PrevField,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Save,
		m.Cancel,
		m.NextField,
		m.PrevField,
	}
}
