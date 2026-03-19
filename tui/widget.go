package tui

import "strings"

type Coords struct{ X, Y int }

type Widget interface {
	// Layout & Rendering
	Resize(width, height int) error
	Render(b *strings.Builder, parent Coords)
	Cursor(b *strings.Builder, parent Coords)

	// Interaction
	HandleKey(key string) bool

	// State Management
	Focus()
	Blur()
}
