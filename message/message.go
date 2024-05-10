package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qualidafial/pomo"
)

func TasksModified(tasks []pomo.Task) tea.Cmd {
	return func() tea.Msg {
		return TasksModifiedMsg{
			Tasks: tasks,
		}
	}
}

type TasksModifiedMsg struct {
	Tasks []pomo.Task
}
