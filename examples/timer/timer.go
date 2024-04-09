package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qualidafial/pomo/timer"
)

type model struct {
	timeout  time.Duration
	timer    timer.Model
	keymap   keymap
	help     help.Model
	quitting bool
}

type keymap struct {
	startStop key.Binding
	stop      key.Binding
	reset     key.Binding
	quit      key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init())
}

func (m model) View() string {
	return fmt.Sprintf("Timeout: %v\n\n%s\n\n%s",
		m.timeout,
		m.timer.View(),
		m.helpView())
}

func (m model) helpView() string {
	return m.help.ShortHelpView([]key.Binding{
		m.keymap.startStop,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m.quit()
		case key.Matches(msg, m.keymap.startStop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.reset):
			m.timer = timer.New(m.timeout)
			return m, m.timer.Init()
		}
	case timer.TimeoutMsg:
		return m.quit()
	}

	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

func (m model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func main() {
	timeout := 90 * time.Second
	if len(os.Args) == 2 {
		var err error
		timeout, err = time.ParseDuration(os.Args[1])
		if err != nil {
			log.Fatalf("invalid duration: %s", os.Args[1])
		}
	}
	m := model{
		timeout: timeout,
		timer:   timer.New(timeout),
		keymap: keymap{
			startStop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start/stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
