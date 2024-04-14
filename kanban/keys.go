package kanban

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	NewTask    key.Binding
	EditTask   key.Binding
	DeleteTask key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	MoveUp    key.Binding
	MoveDown  key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
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

		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),

		MoveUp: key.NewBinding(
			key.WithKeys("shift+up"),
			key.WithHelp("shift+↑", "move up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("shift+down"),
			key.WithHelp("shift+↓", "move down"),
		),
		MoveLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "move left"),
		),
		MoveRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "move right"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.NewTask,
			m.EditTask,
			m.DeleteTask,
		},
		{
			m.Up,
			m.Down,
			m.Left,
			m.Right,
		},
		{
			m.MoveUp,
			m.MoveDown,
			m.MoveLeft,
			m.MoveRight,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.NewTask,
		m.EditTask,
		m.DeleteTask,
		m.Up,
		m.Down,
		m.Left,
		m.Right,
		m.MoveUp,
		m.MoveDown,
		m.MoveLeft,
		m.MoveRight,
	}
}
