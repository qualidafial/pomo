package prompt

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

type Styles struct {
	Frame  lipgloss.Style
	Prompt lipgloss.Style
	Help   lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Frame: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(color.Cyan),
		Prompt: lipgloss.NewStyle(),
		Help:   lipgloss.NewStyle(),
	}
}
