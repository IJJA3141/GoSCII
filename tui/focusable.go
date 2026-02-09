package tui

import "strings"

type Focusable interface {
	Cursor(b *strings.Builder)
	HandleKey(string)
	Focus()
	Blur()
}
