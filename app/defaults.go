package app

import (
	"github.com/qualidafial/pomo"
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
