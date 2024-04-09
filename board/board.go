package board

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/message"
	"github.com/qualidafial/pomo/tasklist"
	"github.com/qualidafial/pomo/taskview"
)

type state int

const (
	stateBoard state = iota
	stateNewTask
	stateEditTask
	stateConfirmDelete
)

type Model struct {
	width  int
	height int
	state  state

	focused   pomo.Status
	taskLists []tasklist.Model

	editor taskview.Model

	fullHelp bool

	help       help.Model
	toggleHelp key.Binding
	up         key.Binding
	down       key.Binding
	left       key.Binding
	right      key.Binding
	moveUp     key.Binding
	moveDown   key.Binding
	moveLeft   key.Binding
	moveRight  key.Binding
	newTask    key.Binding
	deleteTask key.Binding
	editTask   key.Binding
	ok         key.Binding
	cancel     key.Binding
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

		state: stateBoard,

		focused: pomo.Todo,
		taskLists: []tasklist.Model{
			tasklist.New("To Do", todos),
			tasklist.New("Doing", doing),
			tasklist.New("Done", done),
		},

		editor: taskview.New(),

		help: help.New(),
		toggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		moveUp: key.NewBinding(
			key.WithKeys("shift+up"),
			key.WithHelp("shift+↑", "move up"),
		),
		moveDown: key.NewBinding(
			key.WithKeys("shift+down"),
			key.WithHelp("shift+↓", "move down"),
		),
		moveLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "move left"),
		),
		moveRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "move right"),
		),
		newTask: key.NewBinding(
			key.WithKeys("insert", "+"),
			key.WithHelp("+/ins", "new task"),
		),
		deleteTask: key.NewBinding(
			key.WithKeys("delete", "-"),
			key.WithHelp("-/del", "delete task"),
		),
		editTask: key.NewBinding(
			key.WithKeys("enter", "e"),
			key.WithHelp("e/enter", "edit task"),
		),
		ok: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "ok"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}

	for status := range m.taskLists {
		if status == int(m.focused) {
			m.taskLists[status].Focus(0)
		} else {
			m.taskLists[status].Blur()
		}
	}

	m.help.ShowAll = true

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state {
	case stateBoard:
		return m.updateBoard(msg)
	case stateNewTask, stateEditTask:
		return m.updateEditing(msg)
		//case stateConfirmDelete:
		//	return m.updateConfirmDelete(msg)
	}

	return m, nil
}

func (m Model) updateBoard(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.toggleHelp):
			m.ToggleHelp()

		case key.Matches(msg, m.newTask):
			m.InputNewTask()

		case key.Matches(msg, m.editTask):
			m.EditTask()

		//case key.Matches(msg, m.deleteTask):
		//	m.PromptDeleteTask()

		case key.Matches(msg, m.up):
			m.Up()
		case key.Matches(msg, m.down):
			m.Down()
		case key.Matches(msg, m.left):
			m.Left()
		case key.Matches(msg, m.right):
			m.Right()

		case key.Matches(msg, m.moveUp):
			cmd = m.MoveUp()
		case key.Matches(msg, m.moveDown):
			cmd = m.MoveDown()
		case key.Matches(msg, m.moveLeft):
			cmd = m.MoveLeft()
		case key.Matches(msg, m.moveRight):
			cmd = m.MoveRight()
		default:
			m.taskLists[m.focused], cmd = m.taskLists[m.focused].Update(msg)
		}
	default:
		m.taskLists[m.focused], cmd = m.taskLists[m.focused].Update(msg)
	}

	taskList := m.taskLists[m.focused]
	tasks := taskList.Tasks()
	index := taskList.Index()
	selection := index >= 0 && index < len(tasks)

	m.up.SetEnabled(index > 0)
	m.down.SetEnabled(index+1 < len(tasks))
	m.left.SetEnabled(m.focused > pomo.Todo)
	m.right.SetEnabled(m.focused < pomo.Done)

	m.moveUp.SetEnabled(selection && m.up.Enabled())
	m.moveDown.SetEnabled(selection && m.down.Enabled())
	m.moveLeft.SetEnabled(selection && m.left.Enabled())
	m.moveRight.SetEnabled(selection && m.right.Enabled())

	return m, cmd
}

func (m *Model) ToggleHelp() {
	m.fullHelp = !m.fullHelp
	m.layout()
}

func (m *Model) Up() {
	m.taskLists[m.focused].Up()
}

func (m *Model) Down() {
	m.taskLists[m.focused].Down()
}

func (m *Model) Left() {
	if m.focused > pomo.Todo {
		i := m.taskLists[m.focused].Index()
		m.taskLists[m.focused].Blur()
		m.focused--
		m.taskLists[m.focused].Focus(i)
	}
}

func (m *Model) Right() {
	if m.focused < pomo.Done {
		i := m.taskLists[m.focused].Index()
		m.taskLists[m.focused].Blur()
		m.focused++
		m.taskLists[m.focused].Focus(i)
	}
}

func (m *Model) MoveUp() tea.Cmd {
	return m.taskLists[m.focused].MoveUp()
}

func (m *Model) MoveDown() tea.Cmd {
	return m.taskLists[m.focused].MoveDown()
}

func (m *Model) MoveLeft() tea.Cmd {
	var cmd tea.Cmd

	task, index := m.taskLists[m.focused].Remove()
	if index >= 0 {
		m.Left()
		task.Status = m.focused
		cmd = m.taskLists[m.focused].InsertSelect(index, task)
	}

	// save data changes

	return cmd
}

func (m *Model) MoveRight() tea.Cmd {
	var cmd tea.Cmd

	task, index := m.taskLists[m.focused].Remove()
	if index >= 0 {
		m.Right()
		task.Status = m.focused
		cmd = m.taskLists[m.focused].InsertSelect(index, task)
	}

	// save data changes

	return cmd
}

func (m *Model) InputNewTask() tea.Cmd {
	m.state = stateNewTask
	m.editor.SetTask(pomo.Task{
		Status:  m.focused,
		Summary: "",
		Notes:   "",
	})
	return m.editor.Focus()
}

func (m *Model) EditTask() {
	if task, index := m.SelectedTask(); index >= 0 {
		m.state = stateEditTask
		m.editor.SetTask(task)
		m.editor.Focus()
	}
}

//func (m *Model) PromptDeleteTask() {
//	if t, index := m.SelectedTask(); index >= 0 {
//		m.editor.Prompt = "Delete task " + t.Title + "? (y/n)"
//		m.state = stateConfirmDelete
//	}
//}

func (m Model) updateEditing(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.SaveMsg:
		task := m.editor.Task()
		if m.state == stateNewTask {
			cmd = m.taskLists[m.focused].AppendSelect(task)
		} else {
			index := m.taskLists[m.focused].Index()
			cmd = m.taskLists[m.focused].SetTask(index, task)
		}
		m.state = stateBoard
	case message.CancelMsg:
		m.state = stateBoard
	default:
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}

//func (m Model) updateConfirmDelete(msg tea.Msg) (Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch {
//		case key.Matches(msg, m.ok) || strings.EqualFold(msg.String(), "y"):
//			_, _ = m.taskLists[m.focused].Remove()
//			m.state = stateBoard
//			m.editor.Prompt = ""
//			m.editor.SetValue("")
//		case key.Matches(msg, m.cancel) || strings.EqualFold(msg.String(), "n"):
//			m.state = stateBoard
//			m.editor.Prompt = ""
//			m.editor.SetValue("")
//		}
//	}
//
//	return m, nil
//}

func (m Model) View() string {
	switch m.state {
	case stateNewTask, stateEditTask:
		return m.viewEditor()
	//case stateConfirmDelete:
	//return m.viewConfirmDelete()
	default:
		return lipgloss.JoinVertical(lipgloss.Top,
			m.viewBoard(),
			m.viewHelp(),
		)

	}
}

func (m Model) viewBoard() string {
	var taskLists []string
	for _, taskList := range m.taskLists {
		taskLists = append(taskLists, taskList.View())
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, taskLists...)
}

func (m Model) viewEditor() string {
	return m.editor.View()
}

func (m Model) viewHelp() string {
	if m.fullHelp {
		return m.help.FullHelpView(m.fullKeyBindings())
	}
	return m.help.ShortHelpView(m.flatKeyBindings())
}

func (m Model) fullKeyBindings() [][]key.Binding {
	var keys [][]key.Binding
	keys = append(keys, []key.Binding{
		m.toggleHelp,
	})
	switch m.state {
	case stateBoard:
		keys = append(keys, []key.Binding{
			m.up,
			m.down,
			m.left,
			m.right,
		})
		keys = append(keys, []key.Binding{
			m.moveUp,
			m.moveDown,
			m.moveLeft,
			m.moveRight,
		})
		keys = append(keys, []key.Binding{
			m.newTask,
			m.deleteTask,
			m.editTask,
		})
	case stateNewTask, stateEditTask, stateConfirmDelete:
		keys = append(keys, []key.Binding{
			m.ok,
			m.cancel,
		})
	}
	return keys
}

func (m Model) flatKeyBindings() []key.Binding {
	var keys []key.Binding
	for _, k := range m.fullKeyBindings() {
		keys = append(keys, k...)
	}
	return keys
}

func (m Model) SelectedTask() (pomo.Task, int) {
	return m.taskLists[m.focused].Selection()
}

func (m *Model) SetSize(w, h int) {
	tea.Printf("board: %dx%d", w, h)
	m.width = w
	m.height = h
	m.layout()
}

func (m *Model) layout() {
	height := m.height

	helpHeight := 1
	if m.fullHelp {
		helpHeight = 4
	}
	height -= helpHeight

	remainingWidth := m.width
	for i := range m.taskLists {
		width := remainingWidth / (len(m.taskLists) - i)
		remainingWidth -= width
		m.taskLists[i].SetSize(width, height)
	}
	m.editor.SetSize(m.width, height)
	m.help.Width = m.width
}
