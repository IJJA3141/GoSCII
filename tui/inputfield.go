package tui

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type InputField[T any] struct {
	width, minWidth, maxWidth int

	coords Coord
	label  string

	buffer strings.Builder
	valid  string

	parse  func(string) (T, error)
	accept func(string) bool
	submit func(T) string

	offset  int
	cursor  int
	visual  bool
	focused bool
}

func (fld *InputField[T]) Focus() {
	fld.v()
	fld.focused = true
}

func (fld *InputField[T]) Blur() {
	input := fld.buffer.String()
	newValue, err := fld.parse(input)

	fld.buffer.Reset()

	if err == nil {
		fld.valid = fld.submit(newValue)
	}

	// clear buffer
	fld.buffer.WriteString(fld.valid)

	// reset view
	fld.cursor = 0
	fld.offset = 0
	fld.visual = false
	fld.focused = false
}

func (fld *InputField[T]) Resize(width, height int) (err error) {
	if width < fld.minWidth {
		errors.Join(err, fmt.Errorf("width=%d is too small for InputField", width))
	} else if fld.maxWidth < width {
		errors.Join(err, fmt.Errorf("width=%d is too big for InputField", width))
	}

	// if height < fld.minHeight {
	// 	errors.Join(err, fmt.Errorf("height=%d is too small for InputField", height))
	// } else if fld.maxHeight < height {
	// 	errors.Join(err, fmt.Errorf("height=%d is too big for InputField", height))
	// }

	return
}

func (fld *InputField[T]) Cursor(b *strings.Builder, parent Coord) {
	if fld.visual || !fld.focused {
		b.WriteString(HIDE_CURSOR)
		return
	}

	x := parent.X + fld.coords.X + len(fld.label) + 1
	b.WriteString(MoveTo(x+fld.cursor, parent.Y+fld.coords.Y))
	b.WriteString(BLINKING_IBEAM)
	b.WriteString(SHOW_CURSOR)
}

func (fld *InputField[T]) HandleKey(key string) bool {
	if !fld.focused {
		fld.focused = true
		return true
	}

	handled := true

	switch key {
	case KEY_ENTER:
		fld.Blur()

	case KEY_BACKSPACE:
		if fld.visual {
			fld.clear()
		} else {
			fld.backspace()
		}

	case KEY_DEL:
		if fld.visual {
			fld.clear()
		} else {
			fld.del()
		}

	case KEY_RIGHT:
		fld.cursorRight()

	case KEY_LEFT:
		fld.cursorLeft()

	case "\x01": // crtl a
		fld.Focus()

	default:
		handled = fld.accept(key)
		if handled {
			fld.append(key)
		}
	}

	if handled {
		log.Printf("[field] %s\t%x\n", key, key)
	}

	return handled
}

func (fld *InputField[T]) Render(b *strings.Builder, coord Coord) {
	viewEnd := min(fld.offset+fld.width, fld.buffer.Len())

	b.WriteString(MoveTo(coord.X+fld.coords.X, coord.Y+fld.coords.Y))
	b.WriteString(fld.label)

	if fld.focused {
		b.WriteRune('[')
	} else {
		b.WriteRune(' ')
	}

	if fld.visual {
		b.WriteString(HILIGHT_START)
	}

	b.WriteString(fld.buffer.String()[fld.offset:viewEnd])
	b.WriteString(strings.Repeat(" ", fld.width-viewEnd+fld.offset))

	if fld.visual {
		b.WriteString(HILIGHT_END)
	}

	if fld.focused {
		b.WriteRune(']')
	} else {
		b.WriteRune(' ')
	}
}

//
// private functions
//

func (field *InputField[T]) cursorLeft() {
	if field.visual {
		field.visual = false
		field.offset = 0
		field.cursor = 0

	} else if field.cursor > 0 {
		field.cursor--

	} else if field.offset > 0 {
		field.offset--
	}
}

func (field *InputField[T]) cursorRight() {
	if field.visual {
		field.visual = false

		if field.buffer.Len() > field.width {
			field.offset = field.buffer.Len() - field.width
			field.cursor = field.width

		} else {
			field.offset = 0
			field.cursor = field.buffer.Len()
		}

	} else if field.buffer.Len() < field.width+field.offset {
		field.cursor = min(field.cursor+1, field.buffer.Len()-field.offset)

	} else if field.cursor >= field.width {
		field.offset = min(field.offset+1, field.buffer.Len()-field.width)

	} else {
		field.cursor++
	}
}

func (field *InputField[T]) clear() {
	field.buffer.Reset()
	field.cursor = 0
	field.offset = 0
	field.visual = false
}

func (field *InputField[T]) backspace() {
	// remove
	index := field.cursor + field.offset

	if index == 0 {
		return
	}

	input := field.buffer.String()
	field.buffer.Reset()
	field.buffer.WriteString(input[:index-1] + input[index:])

	// position
	if field.offset > 0 {
		field.offset--

	} else if field.cursor > 0 {
		field.cursor--
	}
}

func (field *InputField[T]) del() {
	// remove
	index := field.cursor + field.offset
	input := field.buffer.String()
	field.buffer.Reset()
	field.buffer.WriteString(input[:index] + input[index+1:])

	// position
	if field.offset > 0 { // => len field.input > field.width
		field.offset--

		if field.cursor < field.width {
			field.cursor++
		}
	}
}

func (field *InputField[T]) append(input string) {
	if field.visual {
		field.clear()
	}

	i := field.offset + field.cursor
	input2 := field.buffer.String()
	field.buffer.Reset()
	field.buffer.WriteString(input2[:i] + input + input2[i:])

	if field.cursor < field.width {
		field.cursor++
	} else {
		field.offset = min(field.offset+1, field.buffer.Len())
	}
}

func (fld *InputField[T]) v() {
	fld.visual = fld.buffer.Len() != 0
}
