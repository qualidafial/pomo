package taskedit

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	InputField lipgloss.Style
}

func DefaultStyles() Styles {
	var s Styles
	s.InputField = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("36")).
		BorderStyle(lipgloss.NormalBorder()).
		Width(80).
		Foreground(lipgloss.Color("15"))
	return s
}
