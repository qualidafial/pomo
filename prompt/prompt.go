package prompt

import (
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	lastID  int
	idMutex sync.Mutex
)

func nextID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	lastID++
	return lastID
}

type Model struct {
	Style  lipgloss.Style
	KeyMap KeyMap

	id     int
	width  int
	prompt string

	help help.Model
}

func New() Model {
	return Model{
		Style:  DefaultStyle(),
		KeyMap: DefaultKeyMap(),

		id:     nextID(),
		width:  0,
		prompt: "",

		help: help.New(),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Yes):
			cmd = m.result(true)
		case key.Matches(msg, m.KeyMap.No):
			cmd = m.result(false)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	style := m.Style.Copy()
	maxWidth := m.width - style.GetHorizontalBorderSize() - style.GetHorizontalPadding()
	style.MaxWidth(maxWidth)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.Style.Render(m.prompt),
		m.help.View(m.KeyMap),
	)
}

func (m Model) ID() int {
	return m.id
}

func (m *Model) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m Model) result(v bool) tea.Cmd {
	return func() tea.Msg {
		return PromptResultMsg{
			ID:     m.id,
			Result: v,
		}
	}
}

type PromptResultMsg struct {
	ID     int
	Result bool
}
