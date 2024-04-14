package pomo

type Status int

const (
	Todo Status = iota
	Doing
	Done
)

type Task struct {
	ID     int
	Status Status
	Name   string
	Notes  string
}

type SaveTask struct {
	Task
}
