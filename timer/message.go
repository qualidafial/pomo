package timer

import (
	"time"
)

type StartMsg struct {
	id  int
	end time.Time
}

type ResetMsg struct {
	id int
}

type TickMsg struct {
	id int
}

type TimeoutMsg struct {
	ID int
}
