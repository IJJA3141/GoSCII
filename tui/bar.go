package tui

import "strings"

const (
	MARGIN = 2
	M      = 5
)

type Bar struct {
	width, height int
	x, y          int

	fWidth, fHeight *InputField[int]
	color           *CheckBox

	focus           int
	child           []Focusable
	hasFocusedChild bool
}

func NewBar(width, height int, x, y int) Bar {
	fWidth := NewIntField(width-2*MARGIN, MARGIN+x, MARGIN+y, false)
	fHeight := NewIntField(width-2*MARGIN, MARGIN+x, MARGIN+y+M, false)
	fColor := NewCheckBox(MARGIN+x, MARGIN+2*M, true)

	bar := Bar{
		width:  width,
		height: height,
		x:      x,
		y:      y,

		fWidth:  &fWidth,
		fHeight: &fHeight,
		color:   &fColor,

		focus:           0,
		hasFocusedChild: false,
	}

	bar.child = []Focusable{bar.fWidth, bar.fHeight, bar.color}

	return bar
}

func (bar *Bar) Draw(b *strings.Builder) {
	bar.fWidth.Draw(b)
	bar.fHeight.Draw(b)
	bar.color.Draw(b)

	MoveTo(b, bar.x, bar.y)
}

func (bar *Bar) Cursor(b *strings.Builder) {
	if bar.hasFocusedChild {
		bar.child[bar.focus].Cursor(b)
	} else {
		b.WriteString(HIDE_CURSOR)
	}
}

func (bar *Bar) HandleKey(key string) {
	switch key {
	case KEY_SHIFT_TAB:
		bar.hasFocusedChild = true
		bar.child[bar.focus].Blur()

		bar.focus--

		if bar.focus < 0 {
			bar.focus = len(bar.child) - 1
		}

		bar.child[bar.focus].Focus()

	case KEY_TAB:
		bar.hasFocusedChild = true
		bar.child[bar.focus].Blur()

		bar.focus++

		if bar.focus >= len(bar.child) {
			bar.focus = 0
		}

		bar.child[bar.focus].Focus()

	case KEY_ESC:
		bar.child[bar.focus].Blur()
		bar.hasFocusedChild = false

	default:
		bar.child[bar.focus].HandleKey(key)
	}
}

func (bar *Bar) Focus() {
	bar.hasFocusedChild = true
}
func (bar *Bar) Blur() {
	bar.hasFocusedChild = false
}
