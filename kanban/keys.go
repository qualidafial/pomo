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
			key.WithKeys("left", "down", "up", "right", "h", "j", "k", "l"),
			key.WithHelp("←↓↑→/hjkl", "navigate"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
		),

		Move: key.NewBinding(
			key.WithKeys("shift+left", "shift+down", "shift+up", "shift+right", "H", "J", "K", "L"),
			key.WithHelp("shift+←↓↑→/hjkl", "move task"),
		),
		MoveLeft: key.NewBinding(
			key.WithKeys("shift+left", "H"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("shift+down", "J"),
		),
		MoveUp: key.NewBinding(
			key.WithKeys("shift+up", "K"),
		),
		MoveRight: key.NewBinding(
			key.WithKeys("shift+right", "L"),
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
