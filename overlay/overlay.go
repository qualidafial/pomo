package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Overlay(background, foreground string, x, y int) string {
	lines := strings.Split(background, "\n")
	for i, replacement := range strings.Split(foreground, "\n") {
		if i >= len(lines) {
			lines = append(lines, strings.Repeat(" ", x)+replacement)
			continue
		}
		
		line := lines[i]
		
		width := lipgloss.Width(line)
		if x >= width {
			lines[i] = line + strings.Repeat(" ", x-width) + replacement
			continue
		}
		
		leftSide := lipgloss.
		replacementWidth := lipgloss.Width(replacement)
		if x+replacementWidth >= width {
			lines[i] = line[:x] + replacement + line[x+len(replacement):]
		}
	}

	return strings.Join(lines, "\n")
}
