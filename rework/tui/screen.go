package tui

import (
	"io"
	"unicode/utf8"
)

type ScreenBuffer struct {
	buf []byte
}

func (sb *ScreenBuffer) Reset() {
	clear(sb.buf) // ?
	sb.buf = sb.buf[:0]
}

func (sb *ScreenBuffer) WriteString(s string) {
	sb.buf = append(sb.buf, s...)
}

func (sb *ScreenBuffer) WriteRune(r rune) {
	sb.buf = utf8.AppendRune(sb.buf, r)
}

func (sb *ScreenBuffer) WriteRunes(rs []rune) {
	for _, r := range rs {
		sb.buf = utf8.AppendRune(sb.buf, r)
	}
}

func (sb *ScreenBuffer) WriteRuneN(r rune, n int) {
	for range n {
		sb.WriteRune(r)
	}
}

func (sb ScreenBuffer) Flush(o io.Writer) (int, error) {
	return o.Write(sb.buf)
}
