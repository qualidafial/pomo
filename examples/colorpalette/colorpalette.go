package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	style := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(5)

	fmt.Println("0-16: ANSI 16 colors (4-bit)")
	fmt.Println("color number = red + 2*green + 4*blue + 8*bright")
	printColors(0, 16, 8, true, style)
	printColors(0, 16, 8, false, style)
	fmt.Println()

	// 16..231: 6x6x6 color cube
	// blue: add 0-5
	// green: add 0-5 * 6
	// red: add 0-5 * 36
	fmt.Println("16-255: ANSI 256 colors (8-bit)")
	fmt.Println("16 to 231: 6x6x6 color cube")
	fmt.Println("color number = 16 + blue + 6*green + 36*red")
	printColors(16, 232, 36, true, style)
	printColors(16, 232, 36, false, style)
	fmt.Println()

	fmt.Println("232 to 255: grayscale")
	fmt.Println("color number = 232 + 23*brightness")
	printColors(232, 256, 24, true, style)
	printColors(232, 256, 24, false, style)
}

func printColors(from, to int, perLine int, foreground bool, style lipgloss.Style) {
	for i := from; i < to; i++ {
		color := lipgloss.Color(strconv.Itoa(i))
		altColor := lipgloss.Color("0")

		if foreground {
			style = style.Foreground(color).Background(altColor)
		} else {
			style = style.Background(color).Foreground(altColor)
		}
		fmt.Print(style.Render(string(color)))

		if i > from && (i-from+1)%perLine == 0 {
			fmt.Println()
		}
	}
}
