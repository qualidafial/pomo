package taskview

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
	width  int
	height int

	status pomo.Status

	focused field
	summary textinput.Model
	notes   textarea.Model

	help      help.Model
	nextField key.Binding

	cancel key.Binding
	save   key.Binding
	enter  key.Binding
}

func New() Model {
	title := textinput.New()
	title.Placeholder = "Summary"

	desc := textarea.New()
	desc.Placeholder = "Notes"

	return Model{
		summary: title,
		notes:   desc,

		help: help.New(),
		save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "save"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		nextField: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
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
		return m.summary.Focus()
	case notes:
		m.summary.Blur()
		return m.notes.Focus()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.save) || key.Matches(msg, m.enter):
			return m, message.Save
		case key.Matches(msg, m.cancel):
			return m, message.Cancel
		case key.Matches(msg, m.nextField):
			cmd := m.focusField(m.focused + 1)
			return m, cmd
		}
	}

	var cmd tea.Cmd

	switch m.focused {
	case summary:
		m.summary, cmd = m.summary.Update(msg)
	case notes:
		m.notes, cmd = m.notes.Update(msg)
	}

	m.updateBindings()

	return m, cmd
}

func (m *Model) updateBindings() {
	m.save.SetEnabled(m.summary.Value() != "")
	m.enter.SetEnabled(m.save.Enabled() && m.focused == summary)
}

func (m Model) Task() pomo.Task {
	return pomo.Task{
		Status:  m.status,
		Summary: m.summary.Value(),
		Notes:   m.notes.Value(),
	}
}

func (m *Model) SetTask(task pomo.Task) {
	m.status = task.Status
	m.summary.Reset()
	m.summary.SetValue(task.Summary)

	m.notes.Reset()
	m.notes.SetValue(task.Notes)

	m.updateBindings()
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		"Summary:\n",
		m.summary.View(),
		"\nNotes:\n",
		m.notes.View(),
		"",
		m.help.ShortHelpView([]key.Binding{
			m.save, m.cancel, m.nextField,
		}),
	)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.layout()
}

func (m *Model) layout() {
	m.summary.Width = m.width

	m.notes.SetWidth(m.width)
	m.notes.SetHeight(10)
	m.help.Width = m.width
}
