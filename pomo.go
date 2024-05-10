package pomo

import (
	"time"
)

type Pomo struct {
	Start time.Time `yaml:"start,omitempty"`
	End   time.Time `yaml:"end,omitempty"`
	Tasks []Task    `yaml:"tasks"`
}

func (p Pomo) MarshalYAML() (any, error) {
	var start, end string
	if !p.Start.IsZero() {
		start = p.Start.Format(time.RFC3339Nano)
	}
	if !p.End.IsZero() {
		end = p.End.Format(time.RFC3339Nano)
	}

	return pomoYaml{
		Start: start,
		End:   end,
		Tasks: p.Tasks,
	}, nil
}

func (p *Pomo) UnmarshalYAML(unmarshal func(any) error) error {
	var data pomoYaml
	var err error
	if err = unmarshal(&data); err != nil {
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

	*p = Pomo{
		Start: start,
		End:   end,
		Tasks: data.Tasks,
	}
	return nil
}

type pomoYaml struct {
	Start string `yaml:"start,omitempty"`
	End   string `yaml:"end,omitempty"`
	Tasks []Task `yaml:"tasks,omitempty"`
}
