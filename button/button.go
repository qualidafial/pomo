package button

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	blurred = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 2).
		Margin(1, 2).
		Bold(true).
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E"))

	focused = blurred.Copy().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		Underline(true)
)

type Model struct {
	title        string
	minWidth     int
	focused      bool
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
}

func New(title string) Model {
	return Model{
		title:        title,
		focusedStyle: focused.Copy(),
		blurredStyle: blurred.Copy(),
	}
}

func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) SetTitle(title string) {
	m.title = title
	m.layout()
}

func (m *Model) SetMinWidth(minWidth int) {
	m.minWidth = minWidth
	m.layout()
}

func (m *Model) layout() {
	contentWidth := m.minWidth - focused.GetHorizontalMargins()
	if contentWidth < len(m.title) {
		contentWidth = len(m.title)
	}
	m.focusedStyle.Width(contentWidth)
	m.blurredStyle.Width(contentWidth)
}

func (m Model) View() string {
	style := blurred
	if m.focused {
		style = focused
	}

	return style.Render(m.title)
}
