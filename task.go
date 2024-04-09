package pomo

type Status int

const (
	Todo Status = iota
	Doing
	Done
)

type Task struct {
	Status  Status
	Summary string
	Notes   string
}

type SaveTask struct {
	Task
}
