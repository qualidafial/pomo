package board

import (
	"fmt"
	"github.com/qualidafial/pomo/overlay"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/kanban"
	"github.com/qualidafial/pomo/message"
	"github.com/qualidafial/pomo/prompt"
	"github.com/qualidafial/pomo/taskedit"
)

type state int

const (
	stateKanban state = iota
	stateNewTask
	stateEditTask
	statePromptDelete
)

type Model struct {
	KeyMap KeyMap

	width  int
	height int
	state  state

	kanban       kanban.Model
	editor       taskedit.Model
	deletePrompt prompt.Model
	help         help.Model
}

func New(tasks []pomo.Task) Model {
	h := help.New()
	h.ShowAll = false

	m := Model{
		KeyMap: DefaultKeyMap(),

		width:  0,
		height: 0,
		state:  stateKanban,

		kanban:       kanban.New(tasks),
		editor:       taskedit.New(),
		deletePrompt: prompt.New(),
		help:         h,
	}

	return m
}

func (m Model) Tasks() []pomo.Task {
	return m.kanban.Tasks()
}

func (m Model) SetTasks(tasks []pomo.Task) tea.Cmd {
	return m.kanban.SetTasks(tasks)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateKanban:
		m, cmd = m.updateKanban(msg)
	case stateNewTask, stateEditTask:
		m, cmd = m.updateEditing(msg)
	case statePromptDelete:
		return m.updatePromptDelete(msg)
	}

	return m, cmd
}

func (m Model) updateKanban(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ToggleHelp):
			m.ToggleHelp()
		default:
			m.kanban, cmd = m.kanban.Update(msg)
		}
	case message.NewTaskMsg:
		cmd = m.InputNewTask(msg.Status)
	case message.EditTaskMsg:
		cmd = m.EditTask(msg.Task)
	case message.PromptDeleteTaskMsg:
		m.PromptDeleteTask(msg.Task)
	default:
		m.kanban, cmd = m.kanban.Update(msg)
	}

	return m, cmd
}

func (m Model) updateEditing(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.SaveTaskMsg:
		task := m.editor.Task()
		if m.state == stateNewTask {
			cmd = m.kanban.AppendSelect(task)
		} else {
			cmd = m.kanban.SetTask(task)
		}
		m.state = stateKanban
	case message.CancelEditMsg:
		m.state = stateKanban
	default:
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}

func (m Model) updatePromptDelete(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case prompt.PromptResultMsg:
		if msg.ID == m.deletePrompt.ID() {
			if msg.Result {
				m.kanban.Remove()
			}
			m.state = stateKanban
		}
	default:
		m.deletePrompt, cmd = m.deletePrompt.Update(msg)
	}
	return m, cmd
}

func (m *Model) ToggleHelp() {
	m.help.ShowAll = !m.help.ShowAll
	m.layout()
}

func (m *Model) InputNewTask(status pomo.Status) tea.Cmd {
	m.state = stateNewTask
	m.editor.SetTask(pomo.Task{
		Status: status,
		Name:   "",
		Notes:  "",
	})
	return m.editor.Focus()
}

func (m *Model) EditTask(task pomo.Task) tea.Cmd {
	m.state = stateEditTask
	m.editor.SetTask(task)
	return m.editor.Focus()
}

func (m *Model) PromptDeleteTask(task pomo.Task) {
	m.state = statePromptDelete
	prompt := fmt.Sprintf("Delete task %q?", task.Name)
	m.deletePrompt.SetPrompt(prompt)
}

func (m Model) View() string {
	var popup string
	switch m.state {
	case stateNewTask, stateEditTask:
		popup = m.viewEditor()
	case statePromptDelete:
		popup = m.deletePrompt.View()
	default:
		return lipgloss.JoinVertical(lipgloss.Left,
			m.viewBoard(),
			m.viewHelp(),
		)
	}

	w, h := lipgloss.Size(popup)
	x, y := (m.width-w)/2, (m.height-h)/2
	return overlay.Overlay(m.viewBoard(), popup, x, y)
}

func (m Model) viewBoard() string {
	return m.kanban.View()
}

func (m Model) viewEditor() string {
	return m.editor.View()
}

func (m Model) viewHelp() string {
	return m.help.View(m)
}

func (m Model) FullHelp() [][]key.Binding {
	var keys [][]key.Binding
	keys = m.KeyMap.FullHelp()
	switch m.state {
	case stateKanban:
		keys = append(keys, m.kanban.KeyMap.FullHelp()...)
	case stateNewTask, stateEditTask:
		keys = append(keys, m.editor.KeyMap.FullHelp()...)
	case statePromptDelete:
		keys = append(keys, m.deletePrompt.KeyMap.FullHelp()...)
	}
	return keys
}

func (m Model) ShortHelp() []key.Binding {
	keys := m.KeyMap.ShortHelp()
	switch m.state {
	case stateKanban:
		keys = append(keys, m.kanban.KeyMap.ShortHelp()...)
	case stateNewTask, stateEditTask:
		keys = append(keys, m.editor.KeyMap.ShortHelp()...)
	case statePromptDelete:
		keys = append(keys, m.deletePrompt.KeyMap.ShortHelp()...)
	}
	return keys
}

func (m Model) SelectedTask() (pomo.Task, int) {
	return m.kanban.Task()
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.layout()
}

func (m *Model) layout() {
	helpHeight := lipgloss.Height(m.viewHelp())

	height := m.height - helpHeight

	m.kanban.SetSize(m.width, height)
	m.editor.SetMaxSize(m.width-2, height-2)
	m.help.Width = m.width
}
