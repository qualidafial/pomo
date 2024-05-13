package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/kanban"
	"github.com/qualidafial/pomo/message"
	"github.com/qualidafial/pomo/overlay"
	"github.com/qualidafial/pomo/prompt"
	"github.com/qualidafial/pomo/store"
	"github.com/qualidafial/pomo/taskedit"
	"github.com/qualidafial/pomo/timer"
)

type mode int

const (
	modeNormal mode = iota
	modeNewTask
	modeEditTask
	modePromptDelete
)

type pomoState int

const (
	// No timer is running
	pomoIdle pomoState = iota
	// A pomodoro timer is running
	pomoActive
	// A pomodoro timer has finished and the user is reporting what they did
	pomoEnded
	// A short break timer is running
	pomoBreak
	// A long break timer is running
	pomoLongBreak
	// A break timer has finished
	pomoBreakEnded
)

type Model struct {
	store *store.Store

	width  int
	height int

	mode mode

	pomoState pomoState
	current   pomo.Pomo
	previous  []pomo.Pomo

	dirty bool
	tag   int
	err   error

	timer        timer.Model
	spinner      spinner.Model
	kanban       kanban.Model
	editor       taskedit.Model
	deletePrompt prompt.Model
	help         help.Model

	KeyMap KeyMap
}

func New(s *store.Store) Model {
	return Model{
		store: s,

		width:     0,
		height:    0,
		mode:      modeNormal,
		pomoState: pomoIdle,

		timer:        timer.New(),
		spinner:      spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		kanban:       kanban.New(defaultTasks()),
		editor:       taskedit.New(),
		deletePrompt: prompt.New(),
		help:         help.New(),

		KeyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.DisableMouse,
		m.loadState(),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case timer.StartMsg, timer.ResetMsg, timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
	case timer.TimeoutMsg:
		err := beeep.Beep(0, 0)
		if err != nil {
			log.Error("sending beep on timer expiration", "err", err)
		}
		switch m.pomoState {
		case pomoActive:
			m.pomoState = pomoEnded
			cmd = m.timer.Reset()
			err := beeep.Notify("pomo", "Pomodoro completed! Update your task statuses and start your break!", "")
			if err != nil {
				log.Error("sending notification at end of pomodoro", "err", err)
			}
		case pomoBreak, pomoLongBreak:
			m.pomoState = pomoBreakEnded
			beeep.Notify("pomo", "Break's over! Time to start another pomodoro!", "")
			log.Error("sending notification at end of break", "err", err)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	case message.NewTaskMsg:
		cmd = m.InputNewTask(msg.Status)
	case message.EditTaskMsg:
		cmd = m.EditTask(msg.Task)
	case message.PromptDeleteTaskMsg:
		m.PromptDeleteTask(msg.Task)
	case message.TasksModifiedMsg:
		m.current.Tasks = m.kanban.Tasks()
		m.dirty = true
		m.tag++
		cmd = tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return debounceSaveMsg{
				tag: m.tag,
			}
		})
		cmd = tea.Batch(cmd, m.spinner.Tick)
	case debounceSaveMsg:
		if msg.tag == m.tag {
			cmd = m.saveState()
		}
	case message.ErrMsg:
		m.err = msg.Err
		log.Errorf("%v", msg.Err)
		cmd = tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return clearErrMsg{}
		})

	case message.LoadStateMsg:
		m.current = msg.Current
		m.previous = msg.Previous

		now := time.Now()
		year, month, day := now.Date()
		today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

		cmd = m.kanban.SetTasks(m.current.Tasks)

		// infer current pomodoro state from start/end dates:
		// state        start          end
		// ==================================
		// idle         zero           n/a
		// break        future         n/a
		// active       past           future
		// ended        past           past
		// idle         before today   zero (a break that ended yesterday)
		// break ended  earlier today  zero (a break whose resume time has passed)
		switch {
		case m.current.Start.IsZero():
			// no start date: idle
			m.pomoState = pomoIdle
		case m.current.Start.Compare(now) >= 0:
			// start > now: on a break
			m.pomoState = pomoBreak

			// infer break vs long break by the number of completed pomos today
			if len(m.previous) > 0 && len(m.previous)%4 == 0 {
				m.pomoState = pomoLongBreak
			}

			cmd = tea.Batch(cmd, m.timer.Start(m.current.Start))
		// start < now guaranteed from here on
		case m.current.End.After(now):
			// start < now < end: in an active pomodoro
			m.pomoState = pomoActive
			cmd = tea.Batch(cmd, m.timer.Start(m.current.End))
		case !m.current.End.IsZero():
			// start < end < now: pomodoro has ended
			m.pomoState = pomoEnded
		// end is guaranteed empty from here on
		case m.current.Start.After(today):
			// today < start < now: break whose resume time was earlier today
			m.pomoState = pomoBreakEnded
		default:
			// today < start: break whose resume time was before today
			m.pomoState = pomoIdle
			m.current.Start = time.Time{}
			m.current.End = time.Time{}
		}

		m.dirty = false
	default:
		switch m.mode {
		case modeNormal:
			m, cmd = m.updateNormal(msg)
		case modeNewTask, modeEditTask:
			m, cmd = m.updateEditing(msg)
		case modePromptDelete:
			m, cmd = m.updatePromptDelete(msg)
		}
	}

	_, selection := m.kanban.Task()

	m.KeyMap.StartPomo.SetEnabled(m.pomoState == pomoIdle || m.pomoState == pomoBreakEnded)
	m.KeyMap.CancelPomo.SetEnabled(m.pomoState == pomoActive)
	m.KeyMap.StartBreak.SetEnabled(m.pomoState == pomoEnded)
	m.KeyMap.CancelBreak.SetEnabled(m.pomoState == pomoBreak || m.pomoState == pomoLongBreak)

	m.KeyMap.EditTask.SetEnabled(selection)
	m.KeyMap.DeleteTask.SetEnabled(selection)

	return m, cmd
}

func (m Model) updateNormal(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ToggleHelp):
			m.ToggleHelp()
		case key.Matches(msg, m.KeyMap.NewTask):
			cmd = message.NewTask(m.kanban.Status())
		case key.Matches(msg, m.KeyMap.EditTask):
			task, ok := m.kanban.Task()
			if ok {
				cmd = message.EditTask(task)
			}
		case key.Matches(msg, m.KeyMap.DeleteTask):
			task, ok := m.kanban.Task()
			if ok {
				cmd = message.PromptDeleteTask(task)
			}
		case key.Matches(msg, m.KeyMap.StartPomo):
			m.pomoState = pomoActive
			m.current.Start = time.Now()
			m.current.End = m.current.Start.Add(pomodoroDuration)
			cmd = tea.Batch(m.timer.Start(m.current.End), m.saveState())
		case key.Matches(msg, m.KeyMap.CancelPomo):
			m.pomoState = pomoIdle
			m.current.Start = time.Time{}
			m.current.End = time.Time{}
			cmd = tea.Batch(m.timer.Reset(), m.saveState())
		case key.Matches(msg, m.KeyMap.StartBreak):
			cmd = m.completePomo()
		case key.Matches(msg, m.KeyMap.CancelBreak):
			m.pomoState = pomoIdle
			m.current.Start = time.Time{}
			m.current.End = time.Time{}
			cmd = m.saveState()
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		default:
			m.kanban, cmd = m.kanban.Update(msg)
		}
	default:
		m.kanban, cmd = m.kanban.Update(msg)
	}
	return m, cmd
}

func (m Model) updateEditing(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.SaveTaskMsg:
		task := m.editor.Task()
		if m.mode == modeNewTask {
			cmd = m.kanban.AppendSelect(task)
		} else {
			cmd = m.kanban.SetTask(task)
		}
		m.mode = modeNormal
	case message.CancelEditMsg:
		m.mode = modeNormal
	default:
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}

func (m Model) updatePromptDelete(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case prompt.PromptResultMsg:
		if msg.ID == m.deletePrompt.ID() {
			if msg.Result {
				cmd = m.kanban.Remove()
			}
			m.mode = modeNormal
		}
	default:
		m.deletePrompt, cmd = m.deletePrompt.Update(msg)
	}
	return m, cmd
}

func (m *Model) ToggleHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

func (m Model) View() string {
	m.layout()

	var popup string
	switch m.mode {
	case modeNewTask, modeEditTask:
		popup = m.editor.View()
	case modePromptDelete:
		popup = m.deletePrompt.View()
	default:
		return lipgloss.JoinVertical(lipgloss.Top,
			m.viewHeader(),
			m.kanban.View(),
			m.help.View(m),
		)
	}

	background := lipgloss.JoinVertical(lipgloss.Top,
		m.viewHeader(),
		m.kanban.View(),
		// hide main help when popup is visible
	)

	w, h := lipgloss.Size(popup)
	x, y := (m.width-w)/2, (m.height-h)/2
	return overlay.Overlay(background, popup, x, y)
}

func (m Model) viewHeader() string {
	var b strings.Builder

	b.WriteString("🍅 ")
	switch m.pomoState {
	case pomoIdle:
		b.WriteString("idle")
	case pomoEnded:
		b.WriteString("done. report tasks")
	case pomoBreakEnded:
		b.WriteString("break's over!")
	default:
		b.WriteString(m.timer.View())
		if m.pomoState == pomoBreak {
			b.WriteString(" (break)")
		} else if m.pomoState == pomoLongBreak {
			b.WriteString(" (long break)")
		}
	}

	b.WriteString(" 🍅 ")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("✕ " + m.err.Error()))
	} else if m.dirty {
		b.WriteString(DirtyStyle.Render(m.spinner.View() + " unsaved changes"))
	} else {
		b.WriteString(UpToDateStyle.Render("✓ saved"))
	}

	view := b.String()

	return view
}

func (m *Model) InputNewTask(status pomo.Status) tea.Cmd {
	m.mode = modeNewTask
	m.editor.SetTask(pomo.Task{
		Status: status,
		Name:   "",
		Notes:  "",
	})
	return m.editor.Focus()
}

func (m *Model) EditTask(task pomo.Task) tea.Cmd {
	m.mode = modeEditTask
	m.editor.SetTask(task)
	return m.editor.Focus()
}

func (m *Model) PromptDeleteTask(task pomo.Task) {
	m.mode = modePromptDelete
	prompt := fmt.Sprintf("Delete task %q?", task.Name)
	m.deletePrompt.SetPrompt(prompt)
}

func (m Model) FullHelp() [][]key.Binding {
	return append(m.KeyMap.FullHelp(), m.kanban.KeyMap.FullHelp()...)
}

func (m Model) ShortHelp() []key.Binding {
	return append(m.KeyMap.ShortHelp(), m.kanban.KeyMap.ShortHelp()...)
}

func (m *Model) layout() {
	m.help.Width = m.width

	headerHeight := 1
	helpHeight := lipgloss.Height(m.help.View(m))
	kanbanHeight := m.height - headerHeight - helpHeight

	m.editor.SetMaxSize(m.width-2, kanbanHeight-2)

	m.kanban.SetSize(m.width, kanbanHeight)
}

func (m Model) loadState() tea.Cmd {
	current, err := m.store.GetCurrent()
	if err != nil {
		return message.Err(err)
	}

	return message.LoadState(current, nil)
}

func (m *Model) saveState() tea.Cmd {
	err := m.store.SaveCurrent(m.current)
	if err != nil {
		return message.Err(err)
	}
	return message.LoadState(m.current, m.previous)
}

func (m *Model) completePomo() tea.Cmd {
	var incomplete, workedOn []pomo.Task
	for _, task := range m.kanban.Tasks() {
		if task.Status < pomo.Done {
			incomplete = append(incomplete, task)
		}
		if task.Status > pomo.Todo {
			workedOn = append(workedOn, task)
		}
	}

	completed := pomo.Pomo{
		Start: m.current.Start,
		End:   m.current.End,
		Tasks: workedOn,
	}
	err := m.store.SavePomo(completed)
	if err != nil {
		return message.Err(fmt.Errorf("saving pomodoro: %w", err))
	}

	m.previous = append(m.previous, completed)

	m.pomoState = pomoBreak
	duration := breakDuration
	if len(m.previous)%4 == 0 {
		m.pomoState = pomoLongBreak
		duration = longBreakDuration
	}
	breakEnd := time.Now().Add(duration)
	m.current.Start = breakEnd
	m.current.End = time.Time{}
	m.current.Tasks = incomplete

	err = m.store.SaveCurrent(m.current)
	if err != nil {
		return message.Err(fmt.Errorf("updating current pomodoro: %w", err))
	}

	return tea.Batch(m.timer.Start(breakEnd), m.kanban.SetTasks(incomplete))
}

type debounceSaveMsg struct {
	tag int
}

type clearErrMsg struct{}
