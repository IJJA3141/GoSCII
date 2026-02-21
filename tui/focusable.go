package tui

import "strings"

type Focusable interface {
	Draw(b *strings.Builder)
	Cursor(b *strings.Builder)
	HandleKey(key string)
	Focus()
	Blur()
}
