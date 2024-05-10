// Package timer provides a simple timeout component.
package timer

import (
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	lastID int
	idMtx  sync.Mutex
)

type State int

const (
	StateIdle = iota
	StateActive
	StateTimedOut
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

type Model struct {
	id    int
	state State
	// valid when state is active
	end time.Time
}

// New creates a new timer with the given timeout and nextTick interval.
func New() Model {
	return Model{
		id:    nextID(),
		state: StateIdle,
	}
}

func (m Model) ID() int {
	return m.id
}

func (m Model) State() State {
	return m.state
}

// Start starts the timer with the given end time.
func (m Model) Start(end time.Time) tea.Cmd {
	return func() tea.Msg {
		return StartMsg{
			id:  m.id,
			end: end,
		}
	}
}

// Reset resets the timer to idle
func (m Model) Reset() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{
			id: m.id,
		}
	}
}

// Remaining returns the amount of time remaining if the timer is active
func (m Model) Remaining() time.Duration {
	var remaining time.Duration
	if m.state == StateActive {
		remaining = time.Until(m.end)
		if remaining < 0 {
			remaining = 0
		}
	}
	return remaining
}

func (m Model) Idle() bool {
	return m.state == StateIdle
}

func (m Model) Active() bool {
	return m.state == StateActive
}

func (m Model) TimedOut() bool {
	return m.state == StateTimedOut
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case StartMsg:
		if msg.id != m.id {
			break
		}
		m.state = StateActive
		m.end = msg.end
		cmd = m.nextTick()
	case ResetMsg:
		if msg.id != m.id {
			break
		}
		m.state = StateIdle
		m.end = time.Time{}
	case TickMsg:
		if msg.id != m.id {
			break
		}
		if m.state == StateActive && time.Now().After(m.end) {
			m.state = StateTimedOut
			m.end = time.Time{}
			cmd = m.timeout()
		}
	}

	return m, tea.Batch(cmd, m.nextTick())
}

func (m Model) nextTick() tea.Cmd {
	if !m.Active() {
		return nil
	}

	nextTick := time.Until(m.end)%(time.Second/2) + 1
	if nextTick == 0 {
		nextTick = time.Second / 2
	}
	return tea.Tick(nextTick, func(_ time.Time) tea.Msg {
		return TickMsg{id: m.id}
	})
}

func (m Model) timeout() tea.Cmd {
	return func() tea.Msg {
		return TimeoutMsg{ID: m.id}
	}
}

// View of the timer component.
func (m Model) View() string {
	remaining := m.Remaining()

	seconds := remaining / time.Second
	nanos := remaining % time.Second
	if nanos > 0 {
		// Round up to the next whole second.
		seconds++
	}

	minutes := seconds / 60
	seconds %= 60

	hours := minutes / 60
	minutes %= 60

	// Show the colon during the upper half of each second (including the exact second).
	colon := ":"
	if nanos > 0 && nanos <= time.Second/2 {
		colon = " "
	}

	if hours > 0 {
		return fmt.Sprintf("%02d%s%02d%s%02d", hours, colon, minutes, colon, seconds)
	}
	return fmt.Sprintf("%02d%s%02d", minutes, colon, seconds)
}
