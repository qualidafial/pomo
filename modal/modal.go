package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/qualidafial/pomo/skip"
	"strings"
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

func (m Model) View() string {
	background := m.background.View()
	foreground := m.foreground.View()

	bw, bh := lipgloss.Size(background)
	fw, fh := lipgloss.Size(foreground)

	// foreground completely hides background
	if fw >= bw && fh >= bh {
		return foreground
	}

	bgLines := strings.Split(background, "\n")

	// foreground is wider than background
	if fw >= bw {
		top := (bh - fh) / 2
		bottom := top + fh
		bgUpper := strings.Join(bgLines[0:top], "\n")
		bgLower := strings.Join(bgLines[bottom:fh], "\n")
		return lipgloss.JoinVertical(lipgloss.Left,
			bgUpper,
			foreground,
			bgLower)
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

		lipgloss.JoinHorizontal(lipgloss.Center,
			strings.Join(bgLeft, "\n"),
			foreground,
			strings.Join(bgRight, "\n"))
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

	return lipgloss.JoinVertical(lipgloss.Left,
		bgTop,
		lipgloss.JoinHorizontal(lipgloss.Left,
			strings.Join(bgLeft, "\n"),
			foreground,
			strings.Join(bgRight, "\n")),
		bgBottom,
	)
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
