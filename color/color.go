package color

import (
	"math"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	// Black is the ANSI 16 color black.
	Black lipgloss.Color = "0"

	// Red is the ANSI 16 color red.
	Red lipgloss.Color = "1"

	// Green is the ANSI 16 color green.
	Green lipgloss.Color = "2"

	// Yellow is the ANSI 16 color yellow.
	Yellow lipgloss.Color = "3"

	// Blue is the ANSI 16 color blue.
	Blue lipgloss.Color = "4"

	// Magenta is the ANSI 16 color magenta.
	Magenta lipgloss.Color = "5"

	// Cyan is the ANSI 16 color cyan.
	Cyan lipgloss.Color = "6"

	// White is the ANSI 16 color white.
	White lipgloss.Color = "7"

	// Gray is the ANSI 16 color "bright black", or gray.
	Gray lipgloss.Color = "8"

	// BrightRed is the ANSI 16 color bright red.
	BrightRed lipgloss.Color = "9"

	// BrightGreen is the ANSI 16 color bright green.
	BrightGreen lipgloss.Color = "10"

	// BrightYellow is the ANSI 16 color bright yellow.
	BrightYellow lipgloss.Color = "11"

	// BrightBlue is the ANSI 16 color bright blue.
	BrightBlue lipgloss.Color = "12"

	// BrightMagenta is the ANSI 16 color bright magenta.
	BrightMagenta lipgloss.Color = "13"

	// BrightCyan is the ANSI 16 color bright cyan.
	BrightCyan lipgloss.Color = "14"

	// BrightWhite is the ANSI 16 color bright white.
	BrightWhite lipgloss.Color = "15"
)

// ANSI256ColorCube returns an ANSI 256 color from the 6x6x6 color cube range.
// Color intensities are rounded to the nearest ~16.6% (1/6th) increment.
func ANSI256ColorCube(red, green, blue float64) lipgloss.Color {
	r := rangeToInt(red, 5)
	g := rangeToInt(green, 5)
	b := rangeToInt(blue, 5)
	number := 16 + 36*r + 6*g + b
	return lipgloss.Color(strconv.Itoa(number))
}

// ANSI256Grayscale returns an ANSI 256 color from the grayscale range.
// Intensity is rounded to the nearest ~4.16% (1/24th) increment.
func ANSI256Grayscale(intensity float64) lipgloss.Color {
	number := 232 + rangeToInt(intensity, 23)
	return lipgloss.Color(strconv.Itoa(number))
}

// rangeToInt converts a float in the range [0.0-1.0] to an integer in the range
// [0..max].
func rangeToInt(f float64, max int) int {
	// Clamp to [0..1].
	if f <= 0 || math.IsNaN(f) {
		return 0
	}
	if f >= 1 {
		return max
	}

	// If we just multiply and round, then the numbers [1..max-1] will be
	// twice as likely to be selected as 0 and max.

	// Convert the range [0..1] range to [-0.5..max+0.5] and round to nearest integer.
	// This removes the rounding bias so that 0 and max have an equal chance of
	// being selected.
	f = f*float64(max+1) - 0.5
	i := int(math.Round(f))
	if i < 0 {
		return 0
	}
	if i > max {
		return max
	}
	return i
}
