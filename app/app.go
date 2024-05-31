package app

import (
	"fmt"
	"strconv"
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
	"github.com/qualidafial/pomo/config"
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
	modePrompt
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
	config config.Config
	store  *store.Store

	width  int
	height int

	mode mode

	pomoState pomoState
	current   pomo.Pomo
	previous  []pomo.Pomo

	dirty bool
	tag   int
	err   error

	kanban kanban.Model
	editor taskedit.Model

	prompt    prompt.Model
	onConfirm tea.Msg

	timer   timer.Model
	spinner spinner.Model
	help    help.Model

	KeyMap KeyMap
}

func New(cfg config.Config, s *store.Store) Model {
	return Model{
		config: cfg,
		store:  s,

		width:     0,
		height:    0,
		mode:      modeNormal,
		pomoState: pomoIdle,

		kanban:  kanban.New(defaultTasks()),
		timer:   timer.New(),
		spinner: spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		editor:  taskedit.New(),
		prompt:  prompt.New(),
		help:    help.New(),

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
		m.SetPrompt(fmt.Sprintf("Delete task %q?", msg.Task.Name), DeleteTaskMsg{})
	case message.TasksModifiedMsg:
		m.current.Tasks = m.kanban.Tasks()
		m.dirty = true
		m.tag++
		cmd = tea.Tick(250*time.Millisecond, func(_ time.Time) tea.Msg {
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
	case DeleteTaskMsg:
		cmd = m.kanban.Remove()
	case CancelPomoMsg:
		if m.pomoState == pomoActive {
			m.pomoState = pomoIdle
			m.current.Start = time.Time{}
			m.current.End = time.Time{}
			cmd = tea.Batch(m.timer.Reset(), m.saveState())
		}
	case CompletePomoMsg:
		if m.pomoState == pomoEnded {
			cmd = m.completePomo()
		}
	case CancelBreakMsg:
		if m.pomoState == pomoBreak || m.pomoState == pomoLongBreak {
			m.pomoState = pomoIdle
			m.current.Start = time.Time{}
			m.current.End = time.Time{}
			cmd = m.saveState()
		}
	default:
		switch m.mode {
		case modeNormal:
			m, cmd = m.updateNormal(msg)
		case modeNewTask, modeEditTask:
			m, cmd = m.updateEditing(msg)
		case modePrompt:
			m, cmd = m.updatePrompt(msg)
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

func (m *Model) SetPrompt(prompt string, onConfirm tea.Msg) {
	m.mode = modePrompt
	m.prompt.Prompt = prompt
	m.onConfirm = onConfirm
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
			m.current.End = m.current.Start.Add(m.config.PomodoroDuration)
			cmd = tea.Batch(m.timer.Start(m.current.End), m.saveState())
		case key.Matches(msg, m.KeyMap.CancelPomo):
			m.SetPrompt("Cancel pomodoro?", CancelPomoMsg{})
		case key.Matches(msg, m.KeyMap.StartBreak):
			m.SetPrompt("Complete pomodoro and start break?", CompletePomoMsg{})
		case key.Matches(msg, m.KeyMap.CancelBreak):
			m.SetPrompt("Cancel break early?", CancelBreakMsg{})
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
		task.UpdatedAt = time.Now()
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

func (m Model) updatePrompt(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case prompt.ConfirmMsg:
		if msg.ID == m.prompt.ID() {
			m.mode = modeNormal
			cmd = func() tea.Msg {
				return m.onConfirm
			}
		}
	case prompt.CancelMsg:
		if msg.ID == m.prompt.ID() {
			m.mode = modeNormal
		}
	default:
		m.prompt, cmd = m.prompt.Update(msg)
	}
	return m, cmd
}

func (m *Model) ToggleHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

func (m Model) View() string {
	m.layout()

	callToAction := m.viewCallToAction()

	var sections []string
	if callToAction != "" {
		sections = append(sections, callToAction)
	}

	sections = append(sections,
		m.kanban.View(),
		m.viewFooter(),
	)
	if m.mode == modeNormal && m.help.ShowAll {
		sections = append(sections, Help.Render(m.help.View(m)))
	}

	view := lipgloss.JoinVertical(lipgloss.Top, sections...)

	var popup string
	switch m.mode {
	case modeNewTask, modeEditTask:
		popup = m.editor.View()
	case modePrompt:
		popup = m.prompt.View()
	}
	if popup != "" {
		w, h := lipgloss.Size(popup)
		x, y := (m.width-w)/2, (m.height-h)/2
		view = overlay.Overlay(view, popup, x, y)
	}

	return view
}

func (m Model) viewCallToAction() string {
	var callToAction string
	switch m.pomoState {
	case pomoIdle:
		callToAction = "No pomodoro active."
	case pomoEnded:
		callToAction = "Your pomodoro has ended. Update tasks and start your break!"
	case pomoBreakEnded:
		callToAction = "Your break is over. Time to start another pomodoro!"
	default:
		return ""
	}
	frameWidth := CallToAction.GetHorizontalBorderSize() + CallToAction.GetHorizontalMargins()
	width := max(0, m.width-frameWidth)
	return CallToAction.Width(width).Render(callToAction)
}

func (m Model) viewFooter() string {
	var state string

	switch m.pomoState {
	case pomoIdle:
		state = "idle"
	case pomoEnded:
		state = fmt.Sprintf("pomo %d ended -- report tasks", len(m.previous)+1)
	case pomoBreakEnded:
		state = "break ended -- start another pomo"
	case pomoActive:
		state = fmt.Sprintf("pomo %d in progress", len(m.previous)+1)
	case pomoBreak:
		state = "on a break"
	case pomoLongBreak:
		state = "on a long break"
	}
	state = FooterState.Render(state)

	timer := FooterTimer.Render("🍅", m.timer.View(), "🍅")

	var pomosToday strings.Builder
	if m.config.DailyGoal > 0 && len(m.previous) >= m.config.DailyGoal {
		pomosToday.WriteString("🏆 ")
	}
	pomosToday.WriteString(strconv.Itoa(len(m.previous)))
	if m.config.DailyGoal > 0 {
		pomosToday.WriteRune('/')
		pomosToday.WriteString(strconv.Itoa(m.config.DailyGoal))
	}
	if len(m.previous) == 1 && m.config.DailyGoal == 0 {
		pomosToday.WriteString(" pomo")
	} else {
		pomosToday.WriteString(" pomos")
	}
	var pomos string
	if m.config.DailyGoal > 0 && len(m.previous) >= m.config.DailyGoal {
		pomos = FooterPomosGoal.Render(pomosToday.String())
	} else {
		pomos = FooterPomos.Render(pomosToday.String())
	}

	var errMessage string
	if m.err != nil {
		errMessage = fmt.Sprintf("error: %v", m.err)
		errMessage = FooterError.Render(errMessage)
	}

	var saveState string
	if m.dirty {
		saveState = DirtyStyle.Render(m.spinner.View() + " saving")
	} else {
		saveState = UpToDateStyle.Render("✓ saved")
	}
	saveState = FooterSaveState.Render(saveState)

	helpMessage := FooterHelp.Render("? help")

	w := lipgloss.Width
	spacerWidth := max(0, m.width-w(state)-w(timer)-w(errMessage)-w(pomos)-w(saveState)-w(helpMessage))
	spacer := strings.Repeat(" ", spacerWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		state,
		timer,
		errMessage,
		spacer,
		pomos,
		saveState,
		helpMessage,
	)
}

func (m Model) viewHelp() string {
	var help string
	if m.help.ShowAll {
		help = Help.Render(m.help.View(m))
	}
	return help
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

func (m Model) FullHelp() [][]key.Binding {
	return append(m.KeyMap.FullHelp(), m.kanban.KeyMap.FullHelp()...)
}

func (m Model) ShortHelp() []key.Binding {
	return append(m.KeyMap.ShortHelp(), m.kanban.KeyMap.ShortHelp()...)
}

func (m *Model) layout() {
	m.help.Width = m.width - Help.GetHorizontalFrameSize()

	var ctaHeight int
	if cta := m.viewCallToAction(); cta != "" {
		ctaHeight = lipgloss.Height(cta)
	}

	footerHeight := 1

	var helpHeight int
	if m.help.ShowAll {
		helpHeight = lipgloss.Height(m.viewHelp())
	}

	kanbanHeight := m.height - ctaHeight - footerHeight - helpHeight

	m.editor.SetMaxSize(m.width-2, kanbanHeight-2)

	m.kanban.SetSize(m.width, kanbanHeight)
}

func (m Model) loadState() tea.Cmd {
	current, err := m.store.GetCurrent()
	if err != nil {
		return message.Err(err)
	}

	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	previous, err := m.store.List(today)
	if err != nil {
		return message.Err(err)
	}

	return message.LoadState(current, previous)
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
	duration := m.config.BreakDuration
	if len(m.previous)%4 == 0 {
		m.pomoState = pomoLongBreak
		duration = m.config.LongBreakDuration
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

type DeleteTaskMsg struct{}

type StartPomoMsg struct{}
type CancelPomoMsg struct{}
type CompletePomoMsg struct{}
type CancelBreakMsg struct{}
