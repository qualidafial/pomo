package app

import (
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
	"github.com/qualidafial/pomo/store"
	"github.com/qualidafial/pomo/timer"
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

	state state
	pomo  pomo.Pomo
	dirty bool
	tag   int
	err   error

	timer   timer.Model
	spinner spinner.Model
	board   kanban.Model
	help    help.Model

	KeyMap
}

func New(s *store.Store) Model {
	return Model{
		store: s,

		width:  0,
		height: 0,

		state: stateIdle,

		timer:   timer.New(pomodoroDuration),
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		board:   kanban.New(defaultTasks()),
		help:    help.New(),

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
	case errMsg:
		m.err = msg.err
		log.Errorf("%v", msg.err)
		cmd = tea.Tick(5*time.Second, func(_ time.Time) tea.Msg {
			return clearErrMsg{}
		})
	case loadPomoMsg:
		m.pomo = msg.pomo
		m.dirty = false
		cmd = m.board.SetTasks(m.pomo.Tasks)
	case kanban.KanbanModifiedMsg:
		m.dirty = true
		m.tag++
		cmd = tea.Tick(500*time.Millisecond, func(_ time.Time) tea.Msg {
			return debounceSaveMsg{
				tag: m.tag,
			}
		})
	case debounceSaveMsg:
		if msg.tag == m.tag {
			cmd = m.saveCurrentPomo
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.StartPomo):
			switch m.state {
			case stateIdle:
				m.state = stateActive
				m.timer = timer.New(pomodoroDuration)

				cmd = tea.Batch(m.timer.Init(), m.timer.Start())
			}
		case key.Matches(msg, m.KeyMap.PausePomo):
			cmd = m.timer.Toggle()
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		default:
			m.board, cmd = m.board.Update(msg)
		}
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
	default:
		m.board, cmd = m.board.Update(msg)
	}

	m.KeyMap.StartPomo.SetEnabled(m.state == stateIdle || m.state == stateBreakOver)
	m.KeyMap.PausePomo.SetEnabled(m.state == stateActive)

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.viewHeader(),
		m.board.View())
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
	b.WriteString(m.help.ShortHelpView(m.KeyMap.ShortHelp()))
	b.WriteRune(' ')

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("✕ " + m.err.Error()))
	} else if m.dirty {
		b.WriteString(m.spinner.View() + "saving..")
	} else {
		b.WriteString(UpToDateStyle.Render("✓ up to date"))
	}

	view := b.String()

	return view
}

func (m *Model) layout() {
	timerHeight := 1

	m.board.SetSize(m.width, m.height-timerHeight)
}

func (m Model) saveCurrentPomo() tea.Msg {
	p := m.pomo
	p.Tasks = m.board.Tasks()
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
