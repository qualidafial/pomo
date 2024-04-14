package pomo

import (
	"time"
)

type State int

const (
	PomoActive State = iota
	PomoPaused
	PomoDone
	PomoAbandoned
)

type Pomo struct {
	ID        int
	State     State
	Start     time.Time
	End       *time.Time
	Remaining *time.Duration
	Tasks     []Task
}
