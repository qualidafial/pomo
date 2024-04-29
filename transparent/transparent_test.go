package transparent_test

import (
	"bytes"
	"github.com/muesli/ansi/compressor"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/qualidafial/pomo/color"
	"github.com/qualidafial/pomo/transparent"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		transparent  rune
		wantSegments []string
		wantOffsets  []int
	}{
		{
			name:         "empty string",
			input:        "",
			transparent:  '*',
			wantSegments: nil,
			wantOffsets:  nil,
		},
		{
			name:        "no transparent runes",
			input:       "the quick brown fox jumps over the lazy dog",
			transparent: '*',
			wantSegments: []string{
				"the quick brown fox jumps over the lazy dog",
			},
			wantOffsets: []int{0},
		},
		{
			name:        "no transparent runes or ansi",
			input:       "the quick brown fox jumps over the lazy dog",
			transparent: '*',
			wantSegments: []string{
				"the quick brown fox jumps over the lazy dog",
			},
			wantOffsets: []int{0},
		},
		{
			name: "torture test",
			input: bg(" ** the ", color.Red) +
				fg("quick ***br", color.Blue) +
				fg("own fox***jum", color.Green) +
				bg(fg("ps over**the ", color.White), color.BrightWhite) +
				fg("lazy", color.Black) +
				" dog ***",
			transparent: '*',
			wantSegments: []string{
				bg(" ", color.Red),
				bg(" the ", color.Red) + fg("quick ", color.Blue),
				fg("br", color.Blue) + fg("own fox", color.Green),
				fg("jum", color.Green) + bg(fg("ps over", color.White), color.BrightWhite),
				bg(fg("the ", color.White), color.BrightWhite) + fg("lazy", color.Black) + " dog ",
			},
			wantOffsets: []int{0, 3, 17, 29, 41},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSegments, gotOffsets := transparent.Split(tt.input, tt.transparent)
			assert.Equal(t, tt.wantSegments, gotSegments)
			assert.Equal(t, tt.wantOffsets, gotOffsets)
		})
	}
}

func fg(s string, color lipgloss.TerminalColor) string {
	return compressor.String(testStyle().Foreground(color).Render(s))
}

func bg(s string, color lipgloss.TerminalColor) string {
	return compressor.String(testStyle().Background(color).Render(s))
}

func testStyle() lipgloss.Style {
	var buf bytes.Buffer
	return lipgloss.NewRenderer(&buf,
		termenv.WithProfile(termenv.TrueColor),
		termenv.WithUnsafe(),
	).NewStyle()
}
