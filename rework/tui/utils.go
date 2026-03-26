package tui

import "fmt"

type Coords struct{ X, Y int } // way more than needed

func (c Coords) MoveTo() string          { return fmt.Sprintf("\x1b[%d;%df", c.X, c.Y) }
func (c Coords) Add(other Coords) Coords { return Coords{c.X + other.X, c.Y + other.Y} }

func increase(i, bound int) int {
	if i < bound {
		return i + 1
	} else {
		return 0
	}
}
func decrease(i, bound int) int {
	if i < 1 {
		return bound
	} else {
		return i - 1
	}
}
