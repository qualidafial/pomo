package pomo

import (
	"fmt"
	"time"
)

type State int

const (
	StateIdle State = iota
	StateActive
	StateDone
	StateAbandoned
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateActive:
		return "active"
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
	State    State         `yaml:"state"`
	Start    DateTime      `yaml:"start,omitempty"`
	End      DateTime      `yaml:"end,omitempty"`
	Duration time.Duration `yaml:"duration"`
	Tasks    []Task        `yaml:"tasks"`
}

type DateTime time.Time

func (dt DateTime) String() string {
	return time.Time(dt).UTC().Format(time.RFC3339)
}

func (dt DateTime) MarshalYAML() (any, error) {
	return dt.String(), nil
}

func (dt *DateTime) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d Duration) MarshalYAML() (any, error) {
	return d.String(), nil
}

func (d *Duration) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	t, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = Duration(t)
	return nil
}
