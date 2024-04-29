package pomo

import (
	"fmt"
	"time"
)

type State int

const (
	StateIdle State = iota
	StateActive
	StatePaused
	StateDone
	StateAbandoned
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateActive:
		return "active"
	case StatePaused:
		return "paused"
	case StateDone:
		return "done"
	case StateAbandoned:
		return "abandoned"
	default:
		return "unknown"
	}
}

func (s State) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *State) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	switch str {
	case "idle":
		*s = StateIdle
	case "active":
		*s = StateActive
	case "paused":
		*s = StatePaused
	case "done":
		*s = StateDone
	case "abandoned":
		*s = StateAbandoned
	default:
		return fmt.Errorf("unknown pomodoro state '%s'", s)
	}
	return nil
}

type Pomodoro struct {
	State    State       `yaml:"state"`
	Activity []TimeRange `yaml:"activity"`
	Tasks    []Task      `yaml:"tasks"`
}

type TimeRange struct {
	Start time.Time  `yaml:"start"`
	End   *time.Time `yaml:"end,omitempty"`
}
