package tasklist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo"
)

type Item struct {
	pomo.Task
}

func item(t pomo.Task) Item {
	return Item{
		Task: t,
	}
}

func (i Item) Title() string {
	return i.Task.Summary
}

func (i Item) Description() string {
	return i.Task.Notes
}

func (i Item) FilterValue() string {
	return i.Task.Summary + i.Task.Notes
}

type Model struct {
	width, height int
	focused       bool

	focusedBorder lipgloss.Style
	defaultBorder lipgloss.Style

	list list.Model
}

func New(title string, tasks []pomo.Task) Model {
	delegate := list.NewDefaultDelegate()

	list := list.New(nil, delegate, 0, 0)
	list.Title = title
	list.SetShowHelp(false)

	m := Model{
		width:   0,
		height:  0,
		focused: false,

		list: list,

		defaultBorder: lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()),
		focusedBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")),
	}

	m.SetTasks(tasks)
	m.Blur()

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *Model) Focus(index int) {
	m.focused = true
	m.layout()
	delegate := list.NewDefaultDelegate()
	m.list.SetDelegate(delegate)
	m.Select(index)
}

func (m *Model) Blur() {
	m.focused = false
	m.layout()
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.NormalTitle
	delegate.Styles.SelectedDesc = delegate.Styles.NormalDesc
	m.list.SetDelegate(delegate)
}

func (m Model) Tasks() []pomo.Task {
	items := m.list.Items()
	tasks := make([]pomo.Task, len(items))
	for i, item := range items {
		tasks[i] = item.(Item).Task
	}
	return tasks
}

func (m *Model) SetTasks(tasks []pomo.Task) tea.Cmd {
	items := make([]list.Item, len(tasks))
	for index, task := range tasks {
		items[index] = item(task)
	}

	return m.list.SetItems(items)
}

func (m Model) Task(index int) (pomo.Task, bool) {
	items := m.list.Items()
	if index < 0 || index >= len(items) {
		return pomo.Task{}, false
	}
	return items[index].(Item).Task, true
}

func (m *Model) SetTask(index int, task pomo.Task) tea.Cmd {
	return m.list.SetItem(index, item(task))
}

func (m Model) Selection() (pomo.Task, int) {
	index := m.list.Index()
	item, ok := m.Task(index)
	if !ok {
		return pomo.Task{}, -1
	}
	return item, index
}

func (m Model) Index() int {
	return m.list.Index()
}

func (m Model) Count() int {
	return len(m.list.Items())
}

func (m *Model) Select(index int) {
	count := m.Count()
	if index >= count {
		index = count - 1
	}
	if index < 0 {
		index = 0
	}
	m.list.Select(index)
}

func (m *Model) Remove() (pomo.Task, int) {
	task, index := m.Selection()
	if index >= 0 {
		m.list.RemoveItem(m.list.Index())
		m.Select(index)
	}
	return task, index
}

func (m *Model) Insert(index int, task pomo.Task) tea.Cmd {
	return m.list.InsertItem(index, item(task))
}

func (m *Model) InsertSelect(index int, task pomo.Task) tea.Cmd {
	cmd := m.Insert(index, task)
	m.Select(index)
	return cmd
}

func (m *Model) Append(task pomo.Task) tea.Cmd {
	items := m.list.Items()
	items = append(items, item(task))
	cmd := m.list.SetItems(items)
	return cmd
}

func (m *Model) AppendSelect(task pomo.Task) tea.Cmd {
	return m.InsertSelect(m.Count(), task)
}

func (m *Model) Up() {
	m.list.Select(m.list.Index() - 1)
}

func (m *Model) Down() {
	m.list.Select(m.list.Index() + 1)
}

func (m *Model) MoveUp() tea.Cmd {
	i := m.list.Index()
	items := m.list.Items()
	if len(items) > 1 && i > 0 {
		items[i], items[i-1] = items[i-1], items[i]
		cmd := m.list.SetItems(items)
		m.list.Select(i - 1)
		return cmd
	}
	return nil
}

func (m *Model) MoveDown() tea.Cmd {
	i := m.list.Index()
	items := m.list.Items()
	if len(items) > 1 && i+1 < len(items) {
		items[i], items[i+1] = items[i+1], items[i]
		cmd := m.list.SetItems(items)
		m.list.Select(i + 1)
		return cmd
	}
	return nil
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.layout()
}

func (m *Model) layout() {
	width := m.width - m.defaultBorder.GetHorizontalBorderSize()
	height := m.height - m.defaultBorder.GetVerticalBorderSize()

	m.defaultBorder = m.defaultBorder.
		Width(width).
		Height(height)
	m.focusedBorder = m.focusedBorder.
		Width(width).
		Height(height)
	m.list.SetWidth(width)
	m.list.SetHeight(height)
}

func (m Model) View() string {
	return m.borderStyle().Render(m.list.View())
}

func (m Model) borderStyle() lipgloss.Style {
	if m.focused {
		return m.focusedBorder
	}
	return m.defaultBorder
}
