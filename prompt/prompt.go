package prompt

import (
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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
	Styles Styles
	KeyMap KeyMap

	id       int
	maxWidth int
	prompt   string

	help help.Model
}

func New() Model {
	return Model{
		Styles: DefaultStyles(),
		KeyMap: DefaultKeyMap(),

		id:       nextID(),
		maxWidth: 80,
		prompt:   "",

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
	return m.Styles.Frame.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			m.viewPrompt(),
			"",
			m.viewHelp(),
		),
	)
}

func (m Model) viewPrompt() string {
	width := m.maxWidth - m.Styles.Frame.GetHorizontalFrameSize()
	return m.Styles.Prompt.Render(
		wordwrap.String(m.prompt, width),
	)
}

func (m Model) viewHelp() string {
	return m.Styles.Help.Render(m.help.View(m.KeyMap))
}

func (m *Model) SetPrompt(prompt string) {
	m.prompt = prompt
}

func (m *Model) SetMaxWidth(w int) {
	if w > 80 {
		w = 80
	}
	m.maxWidth = w
}

func (m Model) ID() int {
	return m.id
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
