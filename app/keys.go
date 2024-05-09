package app

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	ToggleHelp key.Binding

	Quit key.Binding

	Pomo key.Binding

	NewTask    key.Binding
	EditTask   key.Binding
	DeleteTask key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),

		Pomo: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pomo…"),
		),

		NewTask: key.NewBinding(
			key.WithKeys("+", "insert"),
			key.WithHelp("+/ins", "new task"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys("-", "delete", "backspace"),
			key.WithHelp("-/del", "delete task"),
		),
		EditTask: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "edit task"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.ToggleHelp,
			m.Quit,
			m.Pomo,
		},
		{
			m.NewTask,
			m.DeleteTask,
			m.EditTask,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.ToggleHelp,
		m.Quit,
		m.Pomo,
		m.NewTask,
		m.DeleteTask,
		m.EditTask,
	}
}
