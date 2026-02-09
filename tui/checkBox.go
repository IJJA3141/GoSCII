package tui

import "strings"

type CheckBox struct {
	Value bool

	x, y int

	checked   string
	unchecked string
}

func (check *CheckBox) Draw(dst *strings.Builder) {
	MoveTo(dst, check.x, check.y)

	if check.Value {
		dst.WriteString(check.checked)
	} else {
		dst.WriteString(check.unchecked)
	}
}

func (check *CheckBox) Cursor(dst *strings.Builder) {
	dst.Grow(MOVE_TO_APROX_SIZE + len(SHOW_CURSOR) + len(BLINKING_UNDERLINE))

	MoveTo(dst, check.x, check.y)
	dst.WriteString(SHOW_CURSOR)
	dst.WriteString(BLINKING_UNDERLINE)
}

func (check *CheckBox) HandleKey(key string) {
	switch key {
	case KEY_ENTER:
		check.Value = !check.Value
	}
}

func (check *CheckBox) Focus() {
}

func (check *CheckBox) Blur() {
}
