package taskedit

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/message"
)

type field int

const (
	summary field = iota
	notes
)

type Model struct {
	KeyMap KeyMap
	Styles Styles

	maxWidth  int
	maxHeight int

	status pomo.Status

	focused field
	name    textinput.Model
	notes   textarea.Model

	help help.Model
}

func New() Model {
	styles := DefaultStyles()

	title := textinput.New()
	title.Placeholder = "name here"

	notes := textarea.New()
	notes.ShowLineNumbers = false
	notes.Placeholder = "notes here"

	return Model{
		Styles: styles,
		KeyMap: DefaultKeyMap(),

		name:  title,
		notes: notes,

		help: help.New(),
	}
}

func (m *Model) Focus() tea.Cmd {
	return m.focusField(summary)
}

func (m *Model) focusField(f field) tea.Cmd {
	if f > notes {
		f = summary
	}
	if f < summary {
		f = notes
	}
	m.focused = f

	switch f {
	case summary:
		m.notes.Blur()
		return m.name.Focus()
	case notes:
		m.name.Blur()
		return m.notes.Focus()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Save) || key.Matches(msg, m.KeyMap.Enter):
			return m, message.SaveTask(m.Task())
		case key.Matches(msg, m.KeyMap.Cancel):
			return m, message.CancelEdit
		case key.Matches(msg, m.KeyMap.NextField):
			cmd := m.focusField(m.focused + 1)
			return m, cmd
		}
	}

	var cmd tea.Cmd

	switch m.focused {
	case summary:
		m.name, cmd = m.name.Update(msg)
	case notes:
		m.notes, cmd = m.notes.Update(msg)
	}

	m.enableKeys()

	return m, cmd
}

func (m *Model) enableKeys() {
	m.KeyMap.Save.SetEnabled(m.name.Value() != "")
	m.KeyMap.Enter.SetEnabled(m.KeyMap.Save.Enabled() && m.focused == summary)
}

func (m Model) Task() pomo.Task {
	return pomo.Task{
		Status: m.status,
		Name:   m.name.Value(),
		Notes:  m.notes.Value(),
	}
}

func (m *Model) SetTask(task pomo.Task) {
	m.status = task.Status
	m.name.Reset()
	m.name.SetValue(task.Name)

	m.notes.Reset()
	m.notes.SetValue(task.Notes)

	m.enableKeys()
}

func (m Model) View() string {
	return m.Styles.Frame.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			m.viewName(),
			"",
			m.viewNotes(),
			"",
			m.viewHelp(),
		),
	)
}

func (m Model) viewName() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Name:",
		m.name.View(),
	)
}

func (m Model) viewNotes() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Notes:",
		m.notes.View(),
	)
}

func (m Model) viewHelp() string {
	return m.help.View(m.KeyMap)
}

func (m *Model) SetMaxSize(maxWidth, maxHeight int) {
	if maxWidth > 80 {
		maxWidth = 80
	}
	m.maxWidth = maxWidth
	m.maxHeight = maxHeight
	m.layout()
}

func (m *Model) layout() {
	w := m.maxWidth - m.Styles.Frame.GetHorizontalFrameSize()
	h := m.maxHeight - m.Styles.Frame.GetVerticalFrameSize()

	m.name.Width = w - 3

	notesHeight := h - 4
	if notesHeight > 4 {
		notesHeight = 4
	}
	m.notes.SetWidth(w)
	m.notes.SetHeight(notesHeight)

	m.help.Width = w
}
