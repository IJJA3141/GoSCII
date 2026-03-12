package tui

import (
	"fmt"
	"strings"
)

type Focus int

const (
	F_WIDTH Focus = iota
	F_HEIGHT
	F_RATIO

	_F_SIZE_
)

type Bar struct {
	index Focus

	size  sizeModule
	state state_t
}

type state_t struct {
	width, height int
}

func NewBar(width, height int, x, y int) Bar { return Bar{} }

func (br *Bar) Draw(b *strings.Builder)
func (br *Bar) Cursor(b *strings.Builder)
func (br *Bar) Resize(width, height int)

func (br *Bar) HandleKey(key string) bool {
	switch br.index {
	case F_HEIGHT:
		handled, newHeight := br.size.height.HandleKey(key)

		if br.size.ratio.Value {
			ratio := float64(br.state.width) / float64(br.state.height)
			br.state.width = int(ratio * float64(newHeight))
		}

		br.state.height = newHeight

		if handled {
			return true
		}

	case F_WIDTH:
		handled, newWidth := br.size.width.HandleKey(key)

		if br.size.ratio.Value {
			ratio := float64(br.state.width) / float64(br.state.height)
			br.state.height = int(float64(newWidth) / ratio)
		}

		br.state.width = newWidth
		if handled {
			return true
		}

	case F_RATIO:
		if br.size.ratio.HandleKey(key) {
			return true
		}

	default:
		panic(fmt.Sprintf("unexpected tui.Focus: %#v", br.index))
	}

	switch key {
	case KEY_TAB:
		br.index++
		if br.index == _F_SIZE_ {
			br.index = 0
		}

	case KEY_SHIFT_TAB:
		br.index--
		if br.index < 0 {
			br.index = _F_SIZE_ - 1
		}

	default:
		return false
	}

	return true
}

type sizeModule struct {
	width, height InputField[int]
	ratio         CheckBox
}
