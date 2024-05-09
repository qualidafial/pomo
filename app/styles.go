package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

var (
	UpToDateStyle = lipgloss.NewStyle().Foreground(color.Green)
	DirtyStyle    = lipgloss.NewStyle().Foreground(color.White)
	ErrorStyle    = lipgloss.NewStyle().Foreground(color.BrightRed)
)
