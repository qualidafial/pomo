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

	Prompt string

	help help.Model
}

func New() Model {
	return Model{
		Styles: DefaultStyles(),
		KeyMap: DefaultKeyMap(),

		id:       nextID(),
		maxWidth: 80,
		Prompt:   "",

		help: help.New(),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case initMsg:
		if m.id == msg.id {
			m.Prompt = msg.prompt
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Yes):
			cmd = func() tea.Msg {
				return ConfirmMsg{
					ID: m.id,
				}
			}
		case key.Matches(msg, m.KeyMap.No):
			cmd = func() tea.Msg {
				return CancelMsg{
					ID: m.id,
				}
			}
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
		wordwrap.String(m.Prompt, width),
	)
}

func (m Model) viewHelp() string {
	return m.Styles.Help.Render(m.help.View(m.KeyMap))
}

func (m Model) ID() int {
	return m.id
}

type initMsg struct {
	id      int
	prompt  string
	confirm tea.Msg
	cancel  tea.Msg
}

type ConfirmMsg struct {
	ID int
}

type CancelMsg struct {
	ID int
}
