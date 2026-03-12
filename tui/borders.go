package tui

import "strings"

var (
	SQUARED_CORNER = [...]rune{'┌', '┐', '└', '┘'}
	ROUND_CORNER   = [...]rune{'╭', '╮', '╰', '╯'}
	BORDERS        = [...]rune{'─', '│', '├', '┤', '┬', '┴', '┼'}
)

type Borders struct {
	width, height, barWidth int
	barLeft                 bool
}

func NewBorders(width, height, barWidth int, barLeft bool) Borders {
	return Borders{
		width:    width,
		height:   height,
		barWidth: barWidth,
		barLeft:  barLeft,
	}
}
func (brs *Borders) Draw(b *strings.Builder) {
	if brs.barLeft {
		MoveTo(b, 0, 0)
		b.WriteRune(SQUARED_CORNER[0])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.barWidth))
		b.WriteRune(BORDERS[4])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.width-brs.barWidth-3))
		b.WriteRune(SQUARED_CORNER[1])

		for i := 1; i < brs.height-2; i++ {
			MoveTo(b, 0, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, brs.barWidth+2, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, brs.width, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, 0, i)
		}

		MoveTo(b, brs.height, 0)
		b.WriteRune(SQUARED_CORNER[2])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.barWidth))
		b.WriteRune(BORDERS[5])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.width-brs.barWidth-3))
		b.WriteRune(SQUARED_CORNER[3])

	} else {
		MoveTo(b, 0, 0)
		b.WriteRune(SQUARED_CORNER[0])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.width-brs.barWidth-3))
		b.WriteRune(BORDERS[4])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.barWidth))
		b.WriteRune(SQUARED_CORNER[1])

		for i := 1; i < brs.height-2; i++ {
			MoveTo(b, 0, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, brs.barWidth+2, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, brs.width, i)
			b.WriteRune(BORDERS[1])

			MoveTo(b, 0, i)
		}

		MoveTo(b, brs.height, 0)
		b.WriteRune(SQUARED_CORNER[2])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.width-brs.barWidth-3))
		b.WriteRune(BORDERS[5])
		b.WriteString(strings.Repeat(string(BORDERS[0]), brs.barWidth))
		b.WriteRune(SQUARED_CORNER[3])
	}
}

func (brs *Borders) Resize(width, height int) {
	brs.height = height
	brs.width = width
}
