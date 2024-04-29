package skip

import (
	"bytes"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

func String(s string, width uint) string {
	return string(Bytes([]byte(s), width))
}

func Bytes(b []byte, width uint) []byte {
	w := NewWriter(width)
	_, _ = w.Write(b)
	return w.Bytes()
}

type Writer struct {
	width uint

	ansiWriter *ansi.Writer
	buf        bytes.Buffer
	ansi       bool
}

func NewWriter(width uint) *Writer {
	w := &Writer{
		width: width,
	}
	w.ansiWriter = &ansi.Writer{
		Forward: &w.buf,
	}
	return w
}

// Write skips content at the given printable cell width, leaving any
// ansi sequences intact.
func (w *Writer) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == ansi.Marker {
			// ANSI escape sequence
			w.ansi = true
		} else if w.ansi {
			if ansi.IsTerminator(c) {
				w.ansi = false
			}
		} else if w.width > 0 {
			cw := uint(runewidth.RuneWidth(c))
			if cw > w.width {
				// wide character at ltrim boundary took an extra cell
				if w.ansiWriter.LastSequence() != "" {
					w.ansiWriter.ResetAnsi()
				}
				for cw > w.width {
					_, _ = w.ansiWriter.Write([]byte(" "))
					w.width++
				}
			}
			w.width -= cw

			if w.width == 0 {
				w.ansiWriter.RestoreAnsi()
			}
			continue
		}

		_, err := w.ansiWriter.Write([]byte(string(c)))
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}

func (w *Writer) String() string {
	return w.buf.String()
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}
