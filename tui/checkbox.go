package tui

import (
	"fmt"
	"strings"
)

type CheckBox struct {
	label  string
	icons  [2]string // icon[0] = false, icon[1] = true
	coords Coords

	submit func(bool) bool

	checked bool
	focused bool
}

// Layout & Rendering

func (cbx *CheckBox) Resize(width, height int) error {
	if width < len(cbx.label)+max(len(cbx.icons[0]), len(cbx.icons[1]))+2 {
		return fmt.Errorf("Failed to resize CheckBox, width %d is too small.", width)
	}

	if height < 1 {
		return fmt.Errorf("Failed to resize CheckBox, height %d is too small.", width)
	}

	return nil
}

func (cbx *CheckBox) Render(b *strings.Builder, coord Coords) {
	b.WriteString(MoveTo(coord.X+cbx.coords.X, coord.Y+cbx.coords.Y))

	if cbx.focused {
		b.WriteRune('[')
	} else {
		b.WriteRune(' ')
	}

	b.WriteString(cbx.label)

	if cbx.checked {
		b.WriteString(cbx.icons[1])
	} else {
		b.WriteString(cbx.icons[0])
	}

	if cbx.focused {
		b.WriteRune(']')
	} else {
		b.WriteRune(' ')
	}
}

func (cbx *CheckBox) Cursor(b *strings.Builder, _ Coords) { b.WriteString(HIDE_CURSOR) }

// Interaction
func (cbx *CheckBox) HandleKey(key string) bool {
	if key != "\n" && key != KEY_ESC && !cbx.focused {
		cbx.focused = true
		return true
	}

	if key == KEY_ENTER || key == " " {
		cbx.checked = !cbx.checked
		return true
	}

	return false
}

// State Management
func (cbx *CheckBox) Focus() { cbx.focused = true }
func (cbx *CheckBox) Blur() {
	cbx.checked = cbx.submit(cbx.checked)
	cbx.focused = false
}
