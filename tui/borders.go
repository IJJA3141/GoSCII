package tui

import "strings"

type Borders struct {
	width, height, split int
}

func NewBorders(width, height, split int) Borders {
	return Borders{
		width:  width,
		height: height,
		split:  split,
	}
}
func (brs *Borders) Draw(b *strings.Builder)
func (brs *Borders) Resize(width, height int)
