package tui

import "strings"

type Coord struct{ X, Y int }

type Element interface {
	HandleKey(string) bool
	Resize(int, int) error
	Render(*strings.Builder, Coord)
	Cursor(*strings.Builder, Coord)
	Focus()
	Blur()
}
