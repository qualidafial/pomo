package composite

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/qualidafial/pomo/ltrim"
	"github.com/qualidafial/pomo/transparent"
)

type Option func(*options)

type options struct {
	maxWidth  int
	maxHeight int
}

func WithMaxSize(width, height int) Option {
	return func(o *options) {
		o.maxWidth = width
		o.maxHeight = height
	}
}

func defaultOptions() options {
	return options{}
}

func Render[E Element](elements []E, options ...Option) string {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	var lines []string
	for _, element := range elements {
		lines = overlayElement(lines, element, opts)
	}
	return strings.Join(lines, "\n")
}

func overlayElement(lines []string, e Element, opts options) []string {
	overlay := strings.Split(e.View(), "\n")

	xOffset, yOffset := e.Position()
	if yOffset+len(overlay) <= 0 || // above the top
		(opts.maxWidth > 0 && xOffset >= opts.maxWidth) || // off the right edge
		(opts.maxHeight > 0 && yOffset >= opts.maxHeight) { // off the left edge
		// element out of bounds at the top, right, or bottom edge
		// out of bounds off the left tbd line by line below
		return lines
	}

	yStart := yOffset
	yEnd := yOffset + len(overlay)
	if yStart < 0 {
		yStart = 0
	}
	if opts.maxHeight > 0 && yEnd > opts.maxHeight {
		yEnd = opts.maxHeight
	}

	if yEnd > len(lines) {
		newLines := make([]string, yEnd)
		copy(newLines, lines)
		lines = newLines
	}

	var transparentRune rune
	if e, ok := e.(TransparentElement); ok {
		transparentRune = e.TransparentRune()
	}

	for y := yStart; y < yEnd; y++ {
		line := overlay[y-yOffset]
		if transparentRune == 0 || !strings.ContainsRune(line, transparentRune) {
			lines[y] = overlayLine(lines[y], line, xOffset, opts)
			continue
		}

		segments, segmentOffsets := transparent.Split(line, transparentRune)
		for i, segment := range segments {
			lines[y] = overlayLine(lines[y], segment, xOffset+segmentOffsets[i], opts)
		}
	}

	return lines
}

func overlayLine(line, overlay string, offset int, opts options) string {
	if offset < 0 {
		overlay = ltrim.String(overlay, uint(-offset))
		offset = 0
	}
	ow := ansi.PrintableRuneWidth(overlay)
	if ow == 0 {
		return line
	}
	if offset+ow <= 0 {
		// overlay out of bounds off the left side
		return line
	}
	if opts.maxWidth > 0 && offset+ow > opts.maxWidth {
		overlay = truncate.String(overlay, uint(opts.maxWidth))
	}

	lw := ansi.PrintableRuneWidth(line)
	if lw < offset {
		line = padding.String(line, uint(offset))
		if len(line) == 0 {
			line = strings.Repeat(" ", offset)
		}
		return line + overlay
	}

	left := line
	right := ""
	if lw > offset {
		left = truncate.String(line, uint(offset))
	}
	if offset+ow < lw {
		right = ltrim.String(line, uint(offset+ow))
	}
	return left + overlay + right
}
