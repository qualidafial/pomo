package composite

// Element is an element of a composite view.
type Element interface {
	// Position returns the x and y offset of the element
	Position() (x int, y int)
	// View returns the rendered content of the element
	View() string
}

// TransparentElement is an element that can have transparent cells.
type TransparentElement interface {
	Element

	// TransparentRune returns the rune that will be treated as transparent
	// when rendering the element in a composite.
	TransparentRune() rune
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
