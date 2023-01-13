package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qualidafial/pomo/set"
)

type mode int

const (
	modeIdle mode = iota
	modeActive
	modeReport
	modeBreak
	modeLongBreak

	pomodoroDuration  = 25 * time.Minute
	breakDuration     = 5 * time.Minute
	longBreakDuration = 15 * time.Minute
)

type todo struct {
	text     string
	complete bool
}

type model struct {
	mode mode

	start    time.Time
	duration time.Duration
	running  bool

	todos     []todo
	cursor    int
	selection set.Set[int]
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oh shit, we encountered an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		todos: []todo{
			{text: "Wax the car"},
			{text: "Paint the fence"},
			{text: "Sand the floor"},
			{text: "Mow the lawn"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.selection = m.selection.Clone()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "ctrl+up": // move up
			if m.cursor > 0 {
				m.todos[m.cursor], m.todos[m.cursor-1] = m.todos[m.cursor-1], m.todos[m.cursor]
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}

		case "ctrl+down": // move down
			if m.cursor < len(m.todos)-1 {
				m.todos[m.cursor], m.todos[m.cursor+1] = m.todos[m.cursor+1], m.todos[m.cursor]
				m.cursor++
			}

		// case "a": // add task

		case "c": // toggle complete
			if len(m.todos) > m.cursor {
				m.todos[m.cursor].complete = !m.todos[m.cursor].complete
			}

		case " ": // toggle selected task
			if len(m.todos) > m.cursor {
				m.selection.Toggle(m.cursor)
			}

		case "enter":
			switch m.mode {
			case modeIdle:
				now := time.Now()

				m.mode = modeActive
				m.start = &now
				m.duration = pomodoroDuration
				m.paused = false
			}
			if len(m.todos) > m.cursor {

			}

		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Pomodoro timer: ")
	switch m.mode {
	case modeIdle:
		b.WriteString("idle")
	case modeReport:
		b.WriteString("done. report tasks")
	default:
		remaining := m.duration
		if m.start != nil {
			remaining -= time.Now().Sub(*m.start)
		}
		minutes := int(remaining / time.Minute)
		seconds := int((remaining % time.Minute) / time.Second)
		fmt.Fprintf(&b, "%2d:%2d", minutes, seconds)
		if m.mode == modeBreak {
			b.WriteString(" (break)")
		} else if m.mode == modeLongBreak {
			b.WriteString(" (long break)")
		}
	}
	b.WriteString("\n\n")

	for i, todo := range m.todos {
		if i == m.cursor {
			b.WriteString("=> ")
		} else {
			b.WriteString("   ")
		}

		if m.selection.Contains(i) {
			b.WriteString("* ")
		} else {
			b.WriteString("  ")
		}

		if todo.complete {
			b.WriteString("[x] ")
		} else {
			b.WriteString("[ ] ")
		}

		b.WriteString(todo.text)
		b.WriteString("\n")
	}

	return b.String()
}
