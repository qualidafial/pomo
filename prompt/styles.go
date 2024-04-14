package prompt

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

func DefaultStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color.BrightCyan).
		Padding(0, 1).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color.Gray)
}
