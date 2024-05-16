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
			key.WithHelp("←/h", "navigate left"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "navigate down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "navigate up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "navigate right"),
		),

		Move: key.NewBinding(
			key.WithKeys("shift+left", "shift+down", "shift+up", "shift+right", "H", "J", "K", "L"),
			key.WithHelp("shift+←↓↑→/hjkl", "move task"),
		),
		MoveLeft: key.NewBinding(
			key.WithKeys("shift+left", "H"),
			key.WithHelp("shift+←/h", "move left"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("shift+down", "J"),
			key.WithHelp("shift+↓/j", "move down"),
		),
		MoveUp: key.NewBinding(
			key.WithKeys("shift+up", "K"),
			key.WithHelp("shift+↑/k", "move up"),
		),
		MoveRight: key.NewBinding(
			key.WithKeys("shift+right", "L"),
			key.WithHelp("shift+→/l", "move right"),
		),
	}
}

func (m KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.Left,
			m.Down,
			m.Up,
			m.Right,
		},
		{
			m.MoveLeft,
			m.MoveDown,
			m.MoveUp,
			m.MoveRight,
		},
	}
}

func (m KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Navigate,
		m.Move,
	}
}
