package transparent

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

// Split scans a string for the transparent rune, and returns the non-
// transparent segments with their horizontal cell offsets.
func Split(s string, transparent rune) ([]string, []int) {
	if s != "" && !strings.ContainsRune(s, transparent) {
		return []string{s}, []int{0}
	}

	w := NewWriter(transparent)
	_, _ = w.Write([]byte(s))
	w.Flush()

	return w.Segments()
}

type Writer struct {
	transparent rune

	ansiWriter *ansi.Writer
	buf        bytes.Buffer
	ansi       bool

	offset    int
	inSegment bool
	segments  []string
	offsets   []int
}

func NewWriter(transparent rune) *Writer {
	w := &Writer{
		transparent: transparent,
	}
	w.ansiWriter = &ansi.Writer{
		Forward: &w.buf,
	}
	return w
}

func (w *Writer) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == ansi.Marker {
			// ANSI escape sequence
			w.ansi = true
		} else if w.ansi {
			if ansi.IsTerminator(c) {
				w.ansi = false
			}
		} else {
			if c != w.transparent {
				if !w.inSegment {
					w.offsets = append(w.offsets, w.offset)
					w.buf.Reset()
					w.ansiWriter.RestoreAnsi()
					w.inSegment = true
				}
			} else if w.inSegment {
				w.Flush()
			}
			w.offset += runewidth.RuneWidth(c)
		}

		_, err := w.ansiWriter.Write([]byte(string(c)))
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}

func (w *Writer) Flush() {
	if w.inSegment {
		w.ansiWriter.ResetAnsi()
		w.segments = append(w.segments, w.buf.String())
		w.buf.Reset()
		w.inSegment = false
	}
}

func (w *Writer) Segments() ([]string, []int) {
	return w.segments, w.offsets
}
