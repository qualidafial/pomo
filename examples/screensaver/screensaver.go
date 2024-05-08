package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qualidafial/pomo/color"
	"github.com/qualidafial/pomo/overlay"
)

const (
	fps = 20

	minColor = 9
	maxColor = 15
)

var (
	//go:embed dvd.txt
	dvd string

	background = color.ANSI256Grayscale(0.075)
)

func main() {
	count := 1
	if len(os.Args) > 1 {
		var err error
		count, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("usage: %s [count]\n", os.Args[0])
		}
	}
	var floaters []floater
	for i := 0; i < count; i++ {
		floaters = append(floaters, floater{
			content:     dvd,
			foreground:  minColor + rand.IntN(8),
			x:           rand.IntN(100),
			y:           rand.IntN(30),
			dx:          (-1 + 2*rand.IntN(2)) * (1 + rand.IntN(4)),
			dy:          (-1 + 2*rand.IntN(2)) * (1 + rand.IntN(2)),
			transparent: ' ',
		})
	}
	if _, err := tea.NewProgram(model{
		floaters: floaters,
	}).Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	width    int
	height   int
	floaters []floater
	paused   bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.tick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmd = tea.Quit
		case " ":
			m.paused = !m.paused
			if !m.paused {
				cmd = m.tick()
			}
		}
	case time.Time:
		if m.paused {
			break
		}

		for i := range m.floaters {
			m.floaters[i] = m.floaters[i].Tick(m.width, m.height)
		}

		return m, m.tick()
	}

	return m, cmd
}

func (m model) tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) View() string {
	var elements []overlay.Element
	elements = append(elements, overlay.DefaultElement{
		X: 0,
		Y: 0,
		Content: lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Background(background).
			Render(),
	})
	for _, f := range m.floaters {
		elements = append(elements, f)
	}
	return overlay.Composite(elements,
		overlay.WithMaxSize(m.width, m.height))
}

type floater struct {
	content       string
	width, height int
	x, y, dx, dy  int
	foreground    int
	transparent   rune
}

func (f floater) Tick(maxWidth, maxHeight int) floater {
	if f.width == 0 && f.height == 0 {
		f.width, f.height = lipgloss.Size(f.content)
	}
	f.x = clamp(f.x+f.dx, 0, maxWidth-f.width)
	f.y = clamp(f.y+f.dy, 0, maxHeight-f.height)

	bounce := false
	if f.x == 0 {
		f.dx = -f.dx
		bounce = true
	}
	if f.x == maxWidth-f.width {
		f.dx = -f.dx
		bounce = true
	}
	if f.y == 0 {
		f.dy = -f.dy
		bounce = true
	}
	if f.y == maxHeight-f.height {
		f.dy = -f.dy
		bounce = true
	}

	if bounce {
		f.foreground++
		if f.foreground < minColor || f.foreground > maxColor {
			f.foreground = minColor
		}
	}

	return f
}

func (f floater) Position() (x, y int) {
	return f.x, f.y
}

func (f floater) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(strconv.Itoa(f.foreground))).
		Background(background).
		Render(f.content)
}

func (f floater) Transparent() rune {
	return f.transparent
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
