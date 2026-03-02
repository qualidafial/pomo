package taskedit

import (
	"charm.land/lipgloss/v2"
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
