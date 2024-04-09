package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/qualidafial/pomo"
)

const (
	pomodoroDuration  = 25 * time.Minute
	breakDuration     = 5 * time.Minute
	longBreakDuration = 15 * time.Minute
)

func defaultTasks() []pomo.Task {
	return []pomo.Task{
		{
			Status:  pomo.Todo,
			Summary: "Paint the fence",
		},
		{
			Status:  pomo.Todo,
			Summary: "foo",
		},
		{
			Status:  pomo.Todo,
			Summary: "bar",
		},
		{
			Status:  pomo.Doing,
			Summary: "Wax the car",
		},
		{
			Status:  pomo.Doing,
			Summary: "baz",
		},
		{
			Status:  pomo.Done,
			Summary: "Sand the floor",
		},
		{
			Status:  pomo.Done,
			Summary: "buz",
		},
	}
}

func defaultKeymap() keymap {
	return keymap{
		startPomo: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "start pomo"),
		),
		pausePomo: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "pause pomo"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
	}
}
