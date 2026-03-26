package tui

type Element interface {
	// Layout & Rendering
	Resize(w, h int) error
	Render(sb *ScreenBuffer, origine Coords)
	Cursor(sb *ScreenBuffer, origine Coords)

	// Interaction
	HandleKey(key string) bool

	// Inner State Management
	Focus()
	Blur()
}
