package overlay

import (
	"github.com/qualidafial/pomo/transparent"
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/qualidafial/pomo/skip"
)

func Overlay(bg, fg string, x, y int, opts ...Option) string {
	options := applyOptions(opts)
	bgLines := splitLines(bg)
	fgLines := splitLines(fg)
	bgLines = overlayLines(bgLines, fgLines, x, y, 0, options)
	return joinLines(bgLines)
}

func OverlayTransparent(bg, fg string, x, y int, transparent rune, opts ...Option) string {
	options := applyOptions(opts)
	bgLines := splitLines(bg)
	fgLines := splitLines(fg)
	bgLines = overlayLines(bgLines, fgLines, x, y, transparent, options)
	return joinLines(bgLines)
}

func overlayLines(bg []string, fg []string, x, y int, transparentRune rune, opts overlayOptions) []string {
	if y+len(fg) <= 0 || // above the top
		(opts.maxWidth > 0 && x >= opts.maxWidth) || // off the right edge
		(opts.maxHeight > 0 && y >= opts.maxHeight) { // off the left edge
		// element out of bounds at the top, right, or bottom edge
		// out of bounds off the left tbd line by line below
		return bg
	}

	rowStart := y
	rowEnd := y + len(fg)
	if rowStart < 0 {
		rowStart = 0
	}
	if opts.maxHeight > 0 && rowEnd > opts.maxHeight {
		rowEnd = opts.maxHeight
	}

	if rowEnd > len(bg) {
		newLines := make([]string, rowEnd)
		copy(newLines, bg)
		bg = newLines
	}

	for row := rowStart; row < rowEnd; row++ {
		line := fg[row-y]
		if transparentRune == 0 || !strings.ContainsRune(line, transparentRune) {
			bg[row] = overlayLine(bg[row], line, x, opts)
			continue
		}

		segments, segmentOffsets := transparent.Split(line, transparentRune)
		for i, segment := range segments {
			bg[row] = overlayLine(bg[row], segment, x+segmentOffsets[i], opts)
		}
	}

	return bg
}

func overlayLine(bg, fg string, offset int, opts overlayOptions) string {
	if offset < 0 {
		fg = skip.String(fg, uint(-offset))
		offset = 0
	}
	ow := ansi.PrintableRuneWidth(fg)
	if ow == 0 {
		return bg
	}
	if offset+ow <= 0 {
		// overlayLines out of bounds off the left side
		return bg
	}
	if opts.maxWidth > 0 && offset+ow > opts.maxWidth {
		fg = truncate.String(fg, uint(opts.maxWidth))
	}

	lw := ansi.PrintableRuneWidth(bg)
	if lw < offset {
		bg = padding.String(bg, uint(offset))
		if len(bg) == 0 {
			bg = strings.Repeat(" ", offset)
		}
		return bg + fg
	}

	left := bg
	right := ""
	if lw > offset {
		left = truncate.String(bg, uint(offset))
	}
	if offset+ow < lw {
		right = skip.String(bg, uint(offset+ow))
	}
	return left + fg + right
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func joinLines(s []string) string {
	return strings.Join(s, "\n")
}
