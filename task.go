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

func ParseStatus(s string) (Status, error) {
	switch s {
	case "todo":
		return Todo, nil
	case "doing":
		return Doing, nil
	case "done":
		return Done, nil
	default:
		return 0, fmt.Errorf("unknown status: %s", s)
	}
}

type Task struct {
	Status Status
	Name   string
	Notes  string
}

func (t Task) MarshalYAML() (any, error) {
	return task{
		Status: t.Status.String(),
		Name:   t.Name,
		Notes:  t.Notes,
	}, nil
}

func (t *Task) UnmarshalYAML(unmarshal func(any) error) error {
	var data task
	if err := unmarshal(&data); err != nil {
		return err
	}

	status, err := ParseStatus(data.Status)
	if err != nil {
		return err
	}

	*t = Task{
		Status: status,
		Name:   data.Name,
		Notes:  data.Notes,
	}
	return nil
}

type task struct {
	Status string `yaml:"status"`
	Name   string `yaml:"name"`
	Notes  string `yaml:"notes,omitempty"`
}
