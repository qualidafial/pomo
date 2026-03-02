package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/truncate"
	"github.com/qualidafial/pomo/skip"
)

func New(background, foreground tea.Model) Model {
	return Model{
		background: background,
		foreground: foreground,
	}
}

type Model struct {
	background tea.Model
	foreground tea.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ResultMsg:
		return m.background, Result(msg.Result)
	default:
		m.foreground, cmd = m.foreground.Update(msg)
	}

	return m, cmd
}

func (m Model) View() tea.View {
	view := m.background.View()

	bg := view.Content
	fg := m.foreground.View().Content

	bw, bh := lipgloss.Size(bg)
	fw, fh := lipgloss.Size(fg)

	// foreground completely hides background
	if fw >= bw && fh >= bh {
		view.SetContent(fg)
		return view
	}

	bgLines := strings.Split(bg, "\n")

	// foreground is wider than background
	if fw >= bw {
		top := (bh - fh) / 2
		bottom := top + fh
		bgUpper := strings.Join(bgLines[0:top], "\n")
		bgLower := strings.Join(bgLines[bottom:fh], "\n")

		view.SetContent(
			lipgloss.JoinVertical(lipgloss.Left,
				bgUpper,
				fg,
				bgLower),
		)
		return view
	}

	// foreground is taller than background
	if fh >= bh {
		left := (bw - fw) / 2
		right := left + fw

		var bgLeft []string
		var bgRight []string
		for _, line := range bgLines {
			bgLeft = append(bgLeft, truncate.String(line, uint(left)))
			bgRight = append(bgRight, skip.String(line, uint(right)))
		}

		view.SetContent(
			lipgloss.JoinHorizontal(lipgloss.Center,
				strings.Join(bgLeft, "\n"),
				fg,
				strings.Join(bgRight, "\n")),
		)
		return view
	}

	// foreground is shorter and narrower than background
	top := (bh - fh) / 2
	bottom := top + fh
	left := (bw - fw) / 2
	right := left + fw

	bgTop := strings.Join(bgLines[:top], "\n")
	bgMiddle := bgLines[top:bottom]
	bgBottom := strings.Join(bgLines[bottom:], "\n")

	var bgLeft []string
	var bgRight []string
	for _, line := range bgMiddle {
		bgLeft = append(bgLeft, truncate.String(line, uint(left)))
		bgRight = append(bgRight, skip.String(line, uint(right)))
	}

	view.SetContent(
		lipgloss.JoinVertical(lipgloss.Left,
			bgTop,
			lipgloss.JoinHorizontal(lipgloss.Left,
				strings.Join(bgLeft, "\n"),
				fg,
				strings.Join(bgRight, "\n")),
			bgBottom,
		),
	)
	return view
}

func Result(result any) tea.Cmd {
	return func() tea.Msg {
		return ResultMsg{
			Result: result,
		}
	}
}

type ResultMsg struct {
	Result any
}
