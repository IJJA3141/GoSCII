package tui

import (
	"strings"
)

type InputField[T any] struct {
	Value T // always valid

	Width int
	X, Y  int

	input string

	format       func(T) string
	parse        func(string) (T, error)
	isValidInput func(string) bool

	cursor    int
	viewStart int
	visual    bool
}

func (field *InputField[T]) Draw(dst *strings.Builder) {
	viewEnd := min(field.viewStart+field.Width, len(field.input))

	dst.Grow(field.Width + MOVE_TO_APROX_SIZE + len(HILIGHT_START) + len(HILIGHT_END))

	MoveTo(dst, field.X, field.Y)

	if field.visual {
		dst.WriteString(HILIGHT_START)
		dst.WriteString(field.input[field.viewStart:viewEnd])
		dst.WriteString(HILIGHT_END)

	} else {
		dst.WriteString(field.input[field.viewStart:viewEnd])
	}

	dst.WriteString(strings.Repeat(" ", max(0, field.Width-len(field.input))))
}

func (field *InputField[T]) Cursor(dst *strings.Builder) {
	if field.visual {
		dst.WriteString(HIDE_CURSOR)
		return
	}

	dst.Grow(MOVE_TO_APROX_SIZE + len(BLINKING_IBEAM) + len(SHOW_CURSOR))

	MoveTo(dst, field.X+field.cursor, field.Y)
	dst.WriteString(BLINKING_IBEAM)
	dst.WriteString(SHOW_CURSOR)
}

func (field *InputField[T]) HandleKey(key string) {
	switch key {
	case KEY_ENTER:
		field.Blur()

	case KEY_BACKSPACE:
		if field.visual {
			field.clear()
		} else if field.cursor+field.viewStart > 0 {
			field.backspace()
		}

	case KEY_DEL:
		if field.visual {
			field.clear()
		} else if field.cursor+field.viewStart < len(field.input) {
			field.del()
		}

	case KEY_RIGHT:
		field.cursorRight()

	case KEY_LEFT:
		field.cursorLeft()

	case "\x01": // crtl a
		field.Focus()

	default:
		field.append(key)
	}
}

func (field *InputField[T]) Focus() {
	field.visual = len(field.input) != 0
}

func (field *InputField[T]) Blur() {
	val, err := field.parse(field.input)
	if err == nil {
		field.Value = val

	}

	field.cursor = 0
	field.viewStart = 0
	field.input = field.format(field.Value)
	field.visual = false
}

//
// private functions
//

func (field *InputField[T]) cursorLeft() {
	if field.visual {
		field.visual = false
		field.viewStart = 0
		field.cursor = 0

	} else if field.cursor > 0 {
		field.cursor--

	} else if field.viewStart > 0 {
		field.viewStart--
	}
}

func (field *InputField[T]) cursorRight() {
	if field.visual {
		field.visual = false

		if len(field.input) > field.Width {
			field.viewStart = len(field.input) - field.Width
			field.cursor = field.Width

		} else {
			field.viewStart = 0
			field.cursor = len(field.input)
		}

	} else if len(field.input) < field.Width+field.viewStart {
		field.cursor = min(field.cursor+1, len(field.input)-field.viewStart)

	} else if field.cursor >= field.Width {
		field.viewStart = min(field.viewStart+1, len(field.input)-field.Width)

	} else {
		field.cursor++
	}
}

func (field *InputField[T]) clear() {
	field.input = ""
	field.cursor = 0
	field.viewStart = 0
	field.visual = false
}

func (field *InputField[T]) backspace() {
	// remove
	index := field.cursor + field.viewStart
	field.input = field.input[:index-1] + field.input[index:]

	// position
	if field.viewStart > 0 {
		field.viewStart--

	} else if field.cursor > 0 {
		field.cursor--
	}
}

func (field *InputField[T]) del() {
	// remove
	index := field.cursor + field.viewStart
	field.input = field.input[:index] + field.input[index+1:]

	// position
	if field.viewStart > 0 { // => len field.input > field.width
		field.viewStart--

		if field.cursor < field.Width {
			field.cursor++
		}
	}
}

func (field *InputField[T]) append(input string) {
	if field.isValidInput(input) {
		if field.visual {
			field.input = input
			field.visual = false
			field.cursor = min(len(field.input), field.Width)
			field.viewStart = max(len(field.input)-field.Width, 0)

		} else {
			i := field.viewStart + field.cursor
			field.input = field.input[:i] + input + field.input[i:]

			if field.cursor < field.Width {
				field.cursor++
			} else {
				field.viewStart = min(field.viewStart+1, len(field.input))
			}
		}
	}
}
