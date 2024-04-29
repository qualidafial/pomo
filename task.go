package pomo

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

func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {}

type Task struct {
	ID     int
	Status Status
	Name   string
	Notes  string
}

type SaveTask struct {
	Task
}
