package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
	"github.com/qualidafial/pomo/modal"
	"log"
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
	return tea.EnterAltScreen
}

func (m background) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
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

func (m background) View() string {
	return lipgloss.NewStyle().
		Background(color.Blue).
		Width(m.width).
		Height(m.height).
		Render(fmt.Sprintf("background %dx%d\n\nresult: %v", m.width, m.height, m.result))
}

type foreground struct {
	width, height int
}

func (m foreground) Init() tea.Cmd {
	return nil
}

func (m foreground) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, modal.Result("run away")
		case "enter":
			return m, modal.Result("let's do this")
		}
	}
	return m, nil
}

func (m foreground) View() string {
	return lipgloss.NewStyle().
		Background(color.Gray).
		Foreground(color.BrightCyan).
		Width(m.width).
		Height(m.height).
		Render(fmt.Sprintf("foreground %dx%d", m.width, m.height))
}
