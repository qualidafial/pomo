package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/board"
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
	width  int
	height int

	state  state
	keymap keymap
	timer  timer.Model
	board  board.Model
	help   help.Model
}

func New() Model {
	return Model{
		keymap: defaultKeymap(),
		help:   help.New(),
		timer:  timer.New(pomodoroDuration),
		board:  board.New(defaultTasks()),
	}
}

type keymap struct {
	startPomo key.Binding
	pausePomo key.Binding
	quit      key.Binding
}

func (m keymap) Bindings() []key.Binding {
	return []key.Binding{
		m.startPomo,
		m.pausePomo,
		m.quit,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		tea.EnterAltScreen,
		tea.DisableMouse,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.startPomo):
			switch m.state {
			case stateIdle:
				m.state = stateActive
				m.timer = timer.New(pomodoroDuration)

				cmd = tea.Batch(m.timer.Init(), m.timer.Start())
			}
		case key.Matches(msg, m.keymap.pausePomo):
			cmd = m.timer.Toggle()
		case key.Matches(msg, m.keymap.quit):
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

	m.keymap.startPomo.SetEnabled(m.state == stateIdle || m.state == stateBreakOver)
	m.keymap.pausePomo.SetEnabled(m.state == stateActive)

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
	b.WriteString(m.help.ShortHelpView(m.keymap.Bindings()))

	return b.String()
}

func (m *Model) layout() {
	timerHeight := 1

	m.board.SetSize(m.width, m.height-timerHeight)
}
