package app

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/board"
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
	err   error

	timer timer.Model
	board board.Model
	help  help.Model

	KeyMap
}

func New(s *store.Store) Model {
	return Model{
		store: s,

		width:  0,
		height: 0,

		state: stateIdle,

		help:  help.New(),
		timer: timer.New(pomodoroDuration),
		board: board.New(defaultTasks()),

		KeyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		tea.EnterAltScreen,
		tea.DisableMouse,
		m.loadCurrentPomo,
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
	case pomoMsg:
		m.pomo = msg.pomo
		cmd = m.board.SetTasks(m.pomo.Tasks)
	case board.BoardModifiedMsg:
		cmd = m.saveCurrentPomo
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
	default:
		m.board, cmd = m.board.Update(msg)
	}

	m.KeyMap.StartPomo.SetEnabled(m.state == stateIdle || m.state == stateBreakOver)
	m.KeyMap.PausePomo.SetEnabled(m.state == stateActive)

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.pomodoroView(),
		m.board.View())
}

func (m Model) pomodoroView() string {
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

	if m.err != nil {
		b.WriteString(" ERROR: ")
		b.WriteString(m.err.Error())
	}

	return b.String()
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
	return pomoMsg{p}
}

func (m Model) loadCurrentPomo() tea.Msg {
	p, err := m.store.GetCurrent()
	if err != nil {
		return errMsg{err}
	}
	return pomoMsg{p}
}

type pomoChangedMsg struct {
}

type pomoMsg struct {
	pomo pomo.Pomo
}

type errMsg struct {
	err error
}

type clearErrMsg struct{}
