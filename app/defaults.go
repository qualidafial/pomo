package app

import (
	"time"

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
			Status: pomo.Todo,
			Name:   "Paint the fence",
		},
		{
			Status: pomo.Todo,
			Name:   "foo",
		},
		{
			Status: pomo.Todo,
			Name:   "bar",
		},
		{
			Status: pomo.Doing,
			Name:   "Wax the car",
		},
		{
			Status: pomo.Doing,
			Name:   "baz",
		},
		{
			Status: pomo.Done,
			Name:   "Sand the floor",
		},
		{
			Status: pomo.Done,
			Name:   "buz",
		},
	}
}
