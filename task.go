package pomo

import (
	"fmt"
)

type Status int

const (
	Todo Status = iota
	Doing
	Done
)

func (s Status) String() string {
	switch s {
	case Todo:
		return "todo"
	case Doing:
		return "doing"
	case Done:
		return "done"
	default:
		return "unknown"
	}
}

func (s Status) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	switch str {
	case "todo":
		*s = Todo
	case "doing":
		*s = Doing
	case "done":
		*s = Done
	default:
		return fmt.Errorf("unknown task status '%s'", s)
	}

	return nil
}

type Task struct {
	Status Status `yaml:"status"`
	Name   string `yaml:"name"`
	Notes  string `yaml:"notes,omitempty"`
}
