package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
)

var (
	UpToDateStyle = lipgloss.NewStyle().Foreground(color.BrightGreen).Bold(true)
	DirtyStyle    = lipgloss.NewStyle().Foreground(color.BrightWhite).Bold(true)

	CallToAction = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(lipgloss.Color("63")).
			Foreground(lipgloss.Color("111"))

	FooterState = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("206")).
			Foreground(lipgloss.Color("228"))
	FooterTimer = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("248"))
	FooterError = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(color.Gray).
			Foreground(color.BrightRed)
	FooterPomos = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("233")).
			Foreground(lipgloss.Color("243"))
	FooterPomosGoal = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("94")).
			Foreground(lipgloss.Color("255"))
	FooterSaveState = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230"))
	FooterHelp = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("243"))

	Help = lipgloss.NewStyle().
		Padding(1, 1, 0)
)
