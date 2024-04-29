package overlay

// Element is an element of a overlay view.
type Element interface {
	// Position returns the x and y offset of the element
	Position() (x int, y int)
	// View returns the rendered content of the element
	View() string
}

// TransparentElement is an element that can have transparent cells.
type TransparentElement interface {
	Element

	// Transparent returns the rune that will be treated as transparent
	// when rendering the element in a overlay.
	Transparent() rune
}

// DefaultElement is a simple implementation of the Element interface.
type DefaultElement struct {
	X, Y    int
	Content string
}

func (e DefaultElement) Position() (x int, y int) {
	return e.X, e.Y
}

func (e DefaultElement) View() string {
	return e.Content
}

func Composite[E Element](elements []E, opts ...Option) string {
	options := applyOptions(opts)
	var lines []string
	for _, element := range elements {
		lines = overlayElement(lines, element, options)
	}
	return joinLines(lines)
}

func overlayElement(bg []string, element Element, opts overlayOptions) []string {
	fg := splitLines(element.View())

	x, y := element.Position()
	var transparentRune rune
	if element, ok := element.(TransparentElement); ok {
		transparentRune = element.Transparent()
	}
	return overlayLines(bg, fg, x, y, transparentRune, opts)
}
