package tui

import (
	"errors"
	"fmt"
	"strings"
)

type CheckBox struct {
	checked bool

	icon [2]string // icon[0] = false, icon[1] = true

	coords Coords
	label  string

	submit func(bool) bool

	buffer  strings.Builder
	focused bool
}

func (cbx *CheckBox) Focus() {
	cbx.focused = true
}

func (cbx *CheckBox) Blur() {
	cbx.checked = cbx.submit(cbx.checked)
	cbx.focused = false
}

func (cbx *CheckBox) Resize(width, height int) (err error) {
	if width < len(cbx.label)+max(len(cbx.icon[0]), len(cbx.icon[1]))+2 {
		errors.Join(err, fmt.Errorf("width=%d is too small for InputField", width))
	}

	if height < 1 {
		errors.Join(err, fmt.Errorf("width=%d is too big for InputField", width))
	}

	return
}

func (cbx *CheckBox) Cursor(b *strings.Builder, _ Coords) { b.WriteString(HIDE_CURSOR) }

func (cbx *CheckBox) Render(b *strings.Builder, coord Coords) {
	b.WriteString(MoveTo(coord.X+cbx.coords.X, coord.Y+cbx.coords.Y))
	b.WriteString(cbx.label)

	if cbx.focused {
		b.WriteRune('[')
	} else {
		b.WriteRune(' ')
	}

	if cbx.checked {
		b.WriteString(cbx.icon[1])
	} else {
		b.WriteString(cbx.icon[0])
	}

	if cbx.focused {
		b.WriteRune(']')
	}
}

func (cbx *CheckBox) HandleKey(key string) bool {
	if !cbx.focused {
		cbx.focused = true
		return true
	}

	if key == KEY_ENTER || key == " " {
		cbx.checked = !cbx.checked
		return true
	}

	return false
}
