package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

var (
	UpToDateStyle = lipgloss.NewStyle().Foreground(color.Green)
	ErrorStyle    = lipgloss.NewStyle().Foreground(color.BrightRed)
)
