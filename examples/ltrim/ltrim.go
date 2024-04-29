package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/qualidafial/pomo/color"
	"github.com/qualidafial/pomo/skip"
)

func main() {
	s := "toma🍅🍅🍅toes"
	for i := 0; i <= lipgloss.Width(s); i++ {
		padding.NewWriter(uint(i), func(w io.Writer) {
			w.Write([]byte(" "))
		})
		left := padding.String(truncate.String(background(s, color.Cyan), uint(i)), uint(i))
		right := skip.String(background(s, color.Red), uint(i))
		fmt.Printf("%s%s\n", left, right)
	}
}

func foreground(s string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().Foreground(color).Render(s)
}

func background(s string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().Background(color).Render(s)
}
