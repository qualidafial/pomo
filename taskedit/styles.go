package taskedit

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

type Styles struct {
	Frame lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Frame: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(color.Cyan),
	}
}
