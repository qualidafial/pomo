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

type timerState int

const (
	statePaused = iota
	stateRunning
	stateTimedOut
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

type StartStopMsg struct {
	ID      int
	running bool
}

type TickMsg struct {
	ID    int
	colon bool
}

type TimeoutMsg struct {
	ID int
}

type Model struct {
	id    int
	state timerState
	// valid when state is active
	end time.Time
	// valid when state is paused
	remaining time.Duration
}

// New creates a new timer with the given timeout and tick interval.
func New(timeout time.Duration) Model {
	return Model{
		id:        nextID(),
		state:     statePaused,
		remaining: timeout,
	}
}

func (m Model) ID() int {
	return m.id
}

func (m Model) Paused() bool {
	return m.state == statePaused
}

func (m Model) Running() bool {
	return m.state == stateRunning
}

func (m Model) TimedOut() bool {
	return m.state == stateTimedOut
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.state == stateTimedOut {
		return m, nil
	}

	switch msg := msg.(type) {
	case StartStopMsg:
		if msg.ID != 0 && msg.ID != m.id {
			break
		}

		switch {
		case msg.running && m.state == statePaused:
			m.state = stateRunning
			m.end = time.Now().Add(m.remaining)
			m.remaining = 0
		case !msg.running && m.state == stateRunning:
			m.state = statePaused
			m.remaining = time.Until(m.end)
			m.end = time.Time{}
		}
	case TickMsg:
		if msg.ID != 0 && msg.ID != m.id {
			break
		}

		if m.state == stateRunning && time.Now().After(m.end) {
			m.end = time.Time{}
			m.state = stateTimedOut
			return m, tea.Batch(m.timedOut())
		}
	}

	return m, tea.Batch(m.tick())
}

// View of the timer component.
func (m Model) View() string {
	switch m.state {
	case stateRunning:
		return timerView(time.Until(m.end))
	case statePaused:
		return timerView(m.remaining)
	}
	return timerView(0)
}

func timerView(remaining time.Duration) string {
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

// Start resumes the timer. Has no effect if the timer has timed out.
func (m *Model) Start() tea.Cmd {
	return m.startStop(true)
}

// Stop pauses the timer. Has no effect if the timer has timed out.
func (m *Model) Stop() tea.Cmd {
	return m.startStop(false)
}

// Toggle stops the timer if it's running and starts it if it's stopped.
func (m *Model) Toggle() tea.Cmd {
	return m.startStop(!m.Running())
}

func (m Model) tick() tea.Cmd {
	if !m.Running() {
		return nil
	}

	nextTick := time.Until(m.end)%(time.Second/2) + 1
	if nextTick == 0 {
		nextTick = time.Second / 2
	}
	return tea.Tick(nextTick, func(_ time.Time) tea.Msg {
		return TickMsg{ID: m.id}
	})
}

func (m Model) timedOut() tea.Cmd {
	if m.state != stateTimedOut {
		return nil
	}
	return func() tea.Msg {
		return TimeoutMsg{ID: m.id}
	}
}

func (m Model) startStop(v bool) tea.Cmd {
	if m.state == stateRunning && !v ||
		m.state == statePaused && v {
		return func() tea.Msg {
			return StartStopMsg{ID: m.id, running: v}
		}
	}
	return nil
}
