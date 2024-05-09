package kanban

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
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
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.Navigate,
			m.Move,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Navigate,
		m.Move,
	}
}
