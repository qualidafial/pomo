package kanban

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	ToggleHelp key.Binding

	NewTask    key.Binding
	EditTask   key.Binding
	DeleteTask key.Binding

	Navigate key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding

	Move      key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),

		Navigate: key.NewBinding(
			key.WithKeys("up", "down", "left", "right"),
			key.WithHelp("↑↓←→", "navigate"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
		),

		Move: key.NewBinding(
			key.WithKeys("shift+up", "shift+down", "shift+left", "shift+right"),
			key.WithHelp("shift+↑↓←→", "move task"),
		),
		MoveUp: key.NewBinding(
			key.WithKeys("shift+up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("shift+down"),
		),
		MoveLeft: key.NewBinding(
			key.WithKeys("shift+left"),
		),
		MoveRight: key.NewBinding(
			key.WithKeys("shift+right"),
		),

		NewTask: key.NewBinding(
			key.WithKeys("insert"),
			key.WithHelp("ins", "new task"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys("delete"),
			key.WithHelp("del", "delete task"),
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
		},
		{
			m.Navigate,
			m.Move,
		},
		{
			m.NewTask,
			m.EditTask,
			m.DeleteTask,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.ToggleHelp,
		m.Navigate,
		m.Move,
		m.NewTask,
		m.EditTask,
		m.DeleteTask,
	}
}
