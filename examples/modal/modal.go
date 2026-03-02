package main

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/qualidafial/pomo/color"
	"github.com/qualidafial/pomo/modal"
)

func main() {
	if _, err := tea.NewProgram(background{}).Run(); err != nil {
		log.Fatal(err)
	}
}

type background struct {
	width, height int

	result any

	modal tea.Model
}

func (m background) Init() tea.Cmd {
	return nil
}

func (m background) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.modal = modal.New(m, foreground{
				width:  m.width / 3,
				height: m.height / 4,
			})
			return m.modal, tea.ClearScreen
		case "ctrl+c":
			cmd = tea.Quit
		}
	case modal.ResultMsg:
		m.result = msg.Result
	}

	return m, cmd
}

func (m background) View() tea.View {
	var v tea.View
	v.SetContent(lipgloss.NewStyle().
		Background(color.Blue).
		Width(m.width).
		Height(m.height).
		Render(fmt.Sprintf("background %dx%d\n\nresult: %v", m.width, m.height, m.result)))
	v.AltScreen = true
	return v
}

type foreground struct {
	width, height int
}

func (m foreground) Init() tea.Cmd {
	return nil
}

func (m foreground) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, modal.Result("run away")
		case "enter":
			return m, modal.Result("let's do this")
		}
	}
	return m, nil
}

func (m foreground) View() tea.View {
	var v tea.View
	v.SetContent(lipgloss.NewStyle().
		Background(color.Gray).
		Foreground(color.BrightCyan).
		Width(m.width).
		Height(m.height).
		Render(fmt.Sprintf("foreground %dx%d", m.width, m.height)),
	)
	return v
}
