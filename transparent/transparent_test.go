package transparent_test

import (
	gcolor "image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/charmbracelet/colorprofile"
	"github.com/muesli/ansi/compressor"
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
		// disabled until our ANSI compressor agrees with the one in lipgloss v2
		//{
		//	name: "torture test",
		//	input: bg(" ** the ", color.Red) +
		//		fg("quick ***br", color.Blue) +
		//		fg("own fox***jum", color.Green) +
		//		bg(fg("ps over**the ", color.White), color.BrightWhite) +
		//		fg("lazy", color.Black) +
		//		" dog ***",
		//	transparent: '*',
		//	wantSegments: []string{
		//		bg(" ", color.Red),
		//		bg(" the ", color.Red) + fg("quick ", color.Blue),
		//		fg("br", color.Blue) + fg("own fox", color.Green),
		//		fg("jum", color.Green) + bg(fg("ps over", color.White), color.BrightWhite),
		//		bg(fg("the ", color.White), color.BrightWhite) + fg("lazy", color.Black) + " dog ",
		//	},
		//	wantOffsets: []int{0, 3, 17, 29, 41},
		//},
	}

	compat.Profile = colorprofile.TrueColor
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSegments, gotOffsets := transparent.Split(tt.input, tt.transparent)

			for i := range gotSegments {
				gotSegments[i] = strings.ReplaceAll(gotSegments[i], "\x1b[0m", "\x1b[m")
			}
			for i := range tt.wantSegments {
				tt.wantSegments[i] = strings.ReplaceAll(tt.wantSegments[i], "\x1b[0m", "\x1b[m")
			}

			if !assert.Equal(t, tt.wantSegments, gotSegments) && len(tt.wantSegments) == len(gotSegments) {
				for i := range gotSegments {
					assert.Equal(t, tt.wantSegments[i], gotSegments[i])
				}
			}
			assert.Equal(t, tt.wantOffsets, gotOffsets)
		})
	}
}

func fg(s string, color gcolor.Color) string {
	return compressor.String(testStyle().Foreground(color).Render(s))
}

func bg(s string, color gcolor.Color) string {
	content := testStyle().Background(color).Render(s)
	return compressor.String(content)
}

func testStyle() lipgloss.Style {
	return lipgloss.NewStyle()
	//var buf bytes.Buffer
	//return lipgloss.NewRenderer(&buf,
	//	termenv.WithUnsafe(),
	//).NewStyle()
}
