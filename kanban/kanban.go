package kanban

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/message"
	"github.com/qualidafial/pomo/tasklist"
)

type Model struct {
	width  int
	height int

	status    pomo.Status
	taskLists []tasklist.Model

	KeyMap KeyMap
}

func New(tasks []pomo.Task) Model {
	var todos, doing, done []pomo.Task
	for _, t := range tasks {
		switch t.Status {
		case pomo.Todo:
			todos = append(todos, t)
		case pomo.Doing:
			doing = append(doing, t)
		case pomo.Done:
			done = append(done, t)
		}
	}

	m := Model{
		width:  0,
		height: 0,

		status: pomo.Todo,
		taskLists: []tasklist.Model{
			tasklist.New("To Do", todos),
			tasklist.New("Doing", doing),
			tasklist.New("Done", done),
		},

		KeyMap: DefaultKeyMap(),
	}

	for status := range m.taskLists {
		if status == int(m.status) {
			m.taskLists[status].Focus(0)
		} else {
			m.taskLists[status].Blur()
		}
	}

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.NewTask):
			cmd = message.NewTask(m.status)
		case key.Matches(msg, m.KeyMap.EditTask):
			task, _ := m.Task()
			cmd = message.EditTask(task)
		case key.Matches(msg, m.KeyMap.DeleteTask):
			task, _ := m.Task()
			cmd = message.PromptDeleteTask(task)
		case key.Matches(msg, m.KeyMap.Up):
			m.Up()
		case key.Matches(msg, m.KeyMap.Down):
			m.Down()
		case key.Matches(msg, m.KeyMap.Left):
			m.Left()
		case key.Matches(msg, m.KeyMap.Right):
			m.Right()

		case key.Matches(msg, m.KeyMap.MoveUp):
			cmd = m.MoveUp()
		case key.Matches(msg, m.KeyMap.MoveDown):
			cmd = m.MoveDown()
		case key.Matches(msg, m.KeyMap.MoveLeft):
			cmd = m.MoveLeft()
		case key.Matches(msg, m.KeyMap.MoveRight):
			cmd = m.MoveRight()
		default:
			m.taskLists[m.status], cmd = m.taskLists[m.status].Update(msg)
		}
	default:
		m.taskLists[m.status], cmd = m.taskLists[m.status].Update(msg)
	}

	taskList := m.taskLists[m.status]
	tasks := taskList.Tasks()
	index := taskList.Index()
	selection := index >= 0 && index < len(tasks)

	m.KeyMap.EditTask.SetEnabled(selection)
	m.KeyMap.DeleteTask.SetEnabled(selection)

	m.KeyMap.Up.SetEnabled(index > 0)
	m.KeyMap.Down.SetEnabled(index+1 < len(tasks))
	m.KeyMap.Left.SetEnabled(m.status > pomo.Todo)
	m.KeyMap.Right.SetEnabled(m.status < pomo.Done)

	m.KeyMap.Move.SetEnabled(selection)
	m.KeyMap.MoveUp.SetEnabled(selection && m.KeyMap.Up.Enabled())
	m.KeyMap.MoveDown.SetEnabled(selection && m.KeyMap.Down.Enabled())
	m.KeyMap.MoveLeft.SetEnabled(selection && m.KeyMap.Left.Enabled())
	m.KeyMap.MoveRight.SetEnabled(selection && m.KeyMap.Right.Enabled())

	return m, cmd
}

func (m Model) View() string {
	var taskLists []string
	for _, taskList := range m.taskLists {
		taskLists = append(taskLists, taskList.View())
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, taskLists...)
}

func (m Model) Status() pomo.Status {
	return m.status
}

func (m *Model) SetStatus(status pomo.Status) {
	if status < pomo.Todo || status > pomo.Done {
		return
	}
	i := m.taskLists[m.status].Index()
	m.taskLists[m.status].Blur()
	m.status = status
	m.taskLists[m.status].Focus(i)
}

func (m *Model) Up() {
	m.taskLists[m.status].Up()
}

func (m *Model) Down() {
	m.taskLists[m.status].Down()
}

func (m *Model) Left() {
	if m.status > pomo.Todo {
		m.SetStatus(m.status - 1)
	}
}

func (m *Model) Right() {
	if m.status < pomo.Done {
		m.SetStatus(m.status + 1)
	}
}

func (m *Model) MoveUp() tea.Cmd {
	return tea.Sequence(m.taskLists[m.status].MoveUp(), m.kanbanModified)
}

func (m *Model) MoveDown() tea.Cmd {
	return tea.Sequence(m.taskLists[m.status].MoveDown(), m.kanbanModified)
}

func (m *Model) MoveLeft() tea.Cmd {
	var cmd tea.Cmd

	task, index := m.taskLists[m.status].Remove()
	if index >= 0 {
		m.Left()
		task.Status = m.status
		cmd = m.taskLists[m.status].InsertSelect(index, task)
	}

	// save data changes

	return tea.Sequence(cmd, m.kanbanModified)
}

func (m *Model) MoveRight() tea.Cmd {
	var cmd tea.Cmd

	task, index := m.taskLists[m.status].Remove()
	if index >= 0 {
		m.Right()
		task.Status = m.status
		cmd = m.taskLists[m.status].InsertSelect(index, task)
	}

	// save data changes

	return tea.Sequence(cmd, m.kanbanModified)
}

func (m *Model) AppendSelect(task pomo.Task) tea.Cmd {
	m.SetStatus(task.Status)
	task.Status = m.status
	cmd := m.taskLists[m.status].AppendSelect(task)
	return tea.Sequence(cmd, m.kanbanModified)
}

func (m *Model) Remove() tea.Cmd {
	m.taskLists[m.status].Remove()
	return m.kanbanModified
}

func (m Model) Tasks() []pomo.Task {
	var tasks []pomo.Task
	for _, taskList := range m.taskLists {
		tasks = append(tasks, taskList.Tasks()...)
	}
	return tasks
}

func (m Model) SetTasks(tasks []pomo.Task) tea.Cmd {
	tasksByStatus := map[pomo.Status][]pomo.Task{}
	for _, task := range tasks {
		tasksByStatus[task.Status] = append(tasksByStatus[task.Status], task)
	}
	var cmds []tea.Cmd
	for status := pomo.Todo; status <= pomo.Done; status++ {
		cmds = append(cmds, m.taskLists[status].SetTasks(tasksByStatus[status]))
	}
	return tea.Batch(cmds...)
}

// Task returns the currently selected task
func (m Model) Task() (pomo.Task, int) {
	return m.taskLists[m.status].Selection()
}

func (m *Model) SetTask(task pomo.Task) tea.Cmd {
	index := m.taskLists[m.status].Index()
	return tea.Sequence(m.taskLists[m.status].SetTask(index, task), m.kanbanModified)
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.layout()
}

func (m *Model) layout() {
	remainingWidth := m.width
	for i := range m.taskLists {
		columnWidth := remainingWidth / (len(m.taskLists) - i)
		remainingWidth -= columnWidth
		m.taskLists[i].SetSize(columnWidth, m.height)
	}
}

func (m Model) kanbanModified() tea.Msg {
	return KanbanModifiedMsg{
		Tasks: m.Tasks(),
	}
}

type KanbanModifiedMsg struct {
	Tasks []pomo.Task
}
