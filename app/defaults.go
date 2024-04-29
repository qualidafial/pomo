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
			Notes:  "Up, down, up down",
		},
		{
			Status: pomo.Doing,
			Name:   "Wax the car",
			Notes:  "Wax on, wax off",
		},
		{
			Status: pomo.Done,
			Name:   "Sand the floor",
			Notes:  "Use little circles",
		},
	}
}
