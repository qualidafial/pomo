package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/kanban"
	"github.com/qualidafial/pomo/message"
	"github.com/qualidafial/pomo/overlay"
	"github.com/qualidafial/pomo/prompt"
	"github.com/qualidafial/pomo/store"
	"github.com/qualidafial/pomo/taskedit"
	"github.com/qualidafial/pomo/timer"
)

type mode int

const (
	modeNormal mode = iota
	modeNewTask
	modeEditTask
	modePromptDelete
)

type state int

const (
	// No timer is running
	stateIdle state = iota
	// A pomodoro timer is running
	stateActive
	// A pomodoro timer has finished and the user is reporting what they did
	stateReport
	// A short break timer is running
	stateBreak
	// A long break timer is running
	stateLongBreak
	// A break timer has finished
	stateBreakOver
)

type Model struct {
	store *store.Store

	width  int
	height int

	mode  mode
	state state
	pomo  pomo.Pomo
	dirty bool
	tag   int
	err   error

	timer        timer.Model
	spinner      spinner.Model
	kanban       kanban.Model
	editor       taskedit.Model
	deletePrompt prompt.Model
	help         help.Model

	KeyMap KeyMap
}

func New(s *store.Store) Model {
	return Model{
		store: s,

		width:  0,
		height: 0,
		mode:   modeNormal,
		state:  stateIdle,

		timer:        timer.New(pomodoroDuration),
		spinner:      spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		kanban:       kanban.New(defaultTasks()),
		editor:       taskedit.New(),
		deletePrompt: prompt.New(),
		help:         help.New(),

		KeyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		tea.EnterAltScreen,
		tea.DisableMouse,
		m.loadCurrentPomo,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case timer.TickMsg, timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
	case timer.TimeoutMsg:
		switch m.state {
		case stateActive:
			m.state = stateBreak
			m.timer = timer.New(breakDuration)
			cmd = m.timer.Init()
		case stateBreak:
			m.state = stateBreakOver
			m.timer = timer.Model{}
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	case message.NewTaskMsg:
		cmd = m.InputNewTask(msg.Status)
	case message.EditTaskMsg:
		cmd = m.EditTask(msg.Task)
	case message.PromptDeleteTaskMsg:
		m.PromptDeleteTask(msg.Task)
	case message.TasksModifiedMsg:
		m.dirty = true
		m.tag++
		cmd = tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return debounceSaveMsg{
				tag: m.tag,
			}
		})
		cmd = tea.Batch(cmd, m.spinner.Tick)
	case debounceSaveMsg:
		if msg.tag == m.tag {
			cmd = m.saveCurrentPomo
		}
	case errMsg:
		m.err = msg.err
		log.Errorf("%v", msg.err)
		cmd = tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return clearErrMsg{}
		})
	case loadPomoMsg:
		m.pomo = msg.pomo
		m.dirty = false
		cmd = m.kanban.SetTasks(m.pomo.Tasks)
	default:
		switch m.mode {
		case modeNormal:
			m, cmd = m.updateNormal(msg)
		case modeNewTask, modeEditTask:
			m, cmd = m.updateEditing(msg)
		case modePromptDelete:
			m, cmd = m.updatePromptDelete(msg)
		}
	}

	_, selection := m.kanban.Task()

	m.KeyMap.EditTask.SetEnabled(selection)
	m.KeyMap.DeleteTask.SetEnabled(selection)

	return m, cmd
}

func (m Model) updateNormal(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ToggleHelp):
			m.ToggleHelp()
		case key.Matches(msg, m.KeyMap.NewTask):
			cmd = message.NewTask(m.kanban.Status())
		case key.Matches(msg, m.KeyMap.EditTask):
			task, ok := m.kanban.Task()
			if ok {
				cmd = message.EditTask(task)
			}
		case key.Matches(msg, m.KeyMap.DeleteTask):
			task, ok := m.kanban.Task()
			if ok {
				cmd = message.PromptDeleteTask(task)
			}
		case key.Matches(msg, m.KeyMap.Pomo):
			switch m.state {
			case stateIdle:
				m.state = stateActive
				m.timer = timer.New(pomodoroDuration)

				cmd = tea.Batch(m.timer.Init(), m.timer.Start())
			}
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		default:
			m.kanban, cmd = m.kanban.Update(msg)
		}
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
		if m.mode == modeNewTask {
			cmd = m.kanban.AppendSelect(task)
		} else {
			cmd = m.kanban.SetTask(task)
		}
		m.mode = modeNormal
	case message.CancelEditMsg:
		m.mode = modeNormal
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
				cmd = m.kanban.Remove()
			}
			m.mode = modeNormal
		}
	default:
		m.deletePrompt, cmd = m.deletePrompt.Update(msg)
	}
	return m, cmd
}

func (m *Model) ToggleHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

func (m Model) View() string {
	m.layout()

	var popup string
	switch m.mode {
	case modeNewTask, modeEditTask:
		popup = m.editor.View()
	case modePromptDelete:
		popup = m.deletePrompt.View()
	default:
		return lipgloss.JoinVertical(lipgloss.Top,
			m.viewHeader(),
			m.kanban.View(),
			m.help.View(m),
		)
	}

	background := lipgloss.JoinVertical(lipgloss.Top,
		m.viewHeader(),
		m.kanban.View(),
		// hide main help when popup is visible
	)

	w, h := lipgloss.Size(popup)
	x, y := (m.width-w)/2, (m.height-h)/2
	return overlay.Overlay(background, popup, x, y)
}

func (m Model) viewHeader() string {
	var b strings.Builder

	b.WriteString("🍅 ")
	switch m.state {
	case stateIdle:
		b.WriteString("idle")
	case stateReport:
		b.WriteString("done. report tasks")
	default:
		b.WriteString(m.timer.View())
		if m.state == stateBreak {
			b.WriteString(" (break)")
		} else if m.state == stateLongBreak {
			b.WriteString(" (long break)")
		}
	}

	b.WriteString(" 🍅 ")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("✕ " + m.err.Error()))
	} else if m.dirty {
		b.WriteString(DirtyStyle.Render(m.spinner.View() + " unsaved changes"))
	} else {
		b.WriteString(UpToDateStyle.Render("✓ saved"))
	}

	view := b.String()

	return view
}

func (m *Model) InputNewTask(status pomo.Status) tea.Cmd {
	m.mode = modeNewTask
	m.editor.SetTask(pomo.Task{
		Status: status,
		Name:   "",
		Notes:  "",
	})
	return m.editor.Focus()
}

func (m *Model) EditTask(task pomo.Task) tea.Cmd {
	m.mode = modeEditTask
	m.editor.SetTask(task)
	return m.editor.Focus()
}

func (m *Model) PromptDeleteTask(task pomo.Task) {
	m.mode = modePromptDelete
	prompt := fmt.Sprintf("Delete task %q?", task.Name)
	m.deletePrompt.SetPrompt(prompt)
}

func (m Model) FullHelp() [][]key.Binding {
	return append(m.KeyMap.FullHelp(), m.kanban.KeyMap.FullHelp()...)
}

func (m Model) ShortHelp() []key.Binding {
	return append(m.KeyMap.ShortHelp(), m.kanban.KeyMap.ShortHelp()...)
}

func (m *Model) layout() {
	m.help.Width = m.width

	headerHeight := 1
	helpHeight := lipgloss.Height(m.help.View(m))
	kanbanHeight := m.height - headerHeight - helpHeight

	m.editor.SetMaxSize(m.width-2, kanbanHeight-2)

	m.kanban.SetSize(m.width, kanbanHeight)
}

func (m Model) saveCurrentPomo() tea.Msg {
	p := m.pomo
	p.Tasks = m.kanban.Tasks()
	err := m.store.SaveCurrent(p)
	if err != nil {
		return errMsg{err}
	}
	return loadPomoMsg{p}
}

func (m Model) loadCurrentPomo() tea.Msg {
	p, err := m.store.GetCurrent()
	if err != nil {
		return errMsg{err}
	}
	return loadPomoMsg{p}
}

type debounceSaveMsg struct {
	tag int
}

type loadPomoMsg struct {
	pomo pomo.Pomo
}

type errMsg struct {
	err error
}

type clearErrMsg struct{}
