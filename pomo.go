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

func ParseState(s string) (State, error) {
	switch s {
	case "idle":
		return StateIdle, nil
	case "active":
		return StateActive, nil
	case "done":
		return StateDone, nil
	case "abandoned":
		return StateAbandoned, nil
	default:
		return 0, fmt.Errorf("invalid state: %s", s)
	}
}

type Pomo struct {
	State    State
	Start    time.Time     `yaml:"start,omitempty"`
	End      time.Time     `yaml:"end,omitempty"`
	Duration time.Duration `yaml:"duration"`
	Tasks    []Task        `yaml:"tasks"`
}

func (p Pomo) MarshalYAML() (any, error) {
	var start, end, duration string
	if !p.Start.IsZero() {
		start = p.Start.Format(time.RFC3339Nano)
	}
	if !p.End.IsZero() {
		end = p.End.Format(time.RFC3339Nano)
	}
	if p.Duration != 0 {
		duration = p.Duration.String()
	}

	return pomoYaml{
		State:    p.State.String(),
		Start:    start,
		End:      end,
		Duration: duration,
		Tasks:    p.Tasks,
	}, nil
}

func (p *Pomo) UnmarshalYAML(unmarshal func(any) error) error {
	var data pomoYaml
	if err := unmarshal(&data); err != nil {
		return err
	}

	state, err := ParseState(data.State)
	if err != nil {
		return err
	}

	var start time.Time
	if data.Start != "" {
		start, err = time.Parse(time.RFC3339Nano, data.Start)
		if err != nil {
			return nil
		}
	}

	var end time.Time
	if data.End != "" {
		end, err = time.Parse(time.RFC3339Nano, data.End)
		if err != nil {
			return nil
		}
	}

	var duration time.Duration
	if data.Duration != "" {
		duration, err = time.ParseDuration(data.Duration)
		if err != nil {
			return nil
		}
	}
	*p = Pomo{
		State:    state,
		Start:    start,
		End:      end,
		Duration: duration,
		Tasks:    data.Tasks,
	}
	return nil
}

type pomoYaml struct {
	State    string `yaml:"state"`
	Start    string `yaml:"start,omitempty"`
	End      string `yaml:"end,omitempty"`
	Duration string `yaml:"duration,omitempty"`
	Tasks    []Task `yaml:"tasks,omitempty"`
}
