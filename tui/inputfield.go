package tui

import (
	"fmt"
	"log"
	"strings"
)

type InputField[T any] struct {
	width    int // only the width of the input
	minWidth int // should be bigger than `len(label) + 3`

	coords Coords
	label  string

	buffer []rune
	Value  T

	format func(T) string
	parse  func(string) (T, error)
	accept func(string) bool
	submit func(T) T

	offset  int
	cursor  int
	visual  bool
	focused bool
}

// Layout & Rendering

func (fld *InputField[T]) Resize(width, height int) error {
	if height < 1 {
		return fmt.Errorf("Failed to resize InputField, height %d is too small.", height)
	}

	if width < fld.minWidth {
		return fmt.Errorf("Failed to resize InputField, width %d is too small.", width)
	}

	fld.width = width - len(fld.label) - 2 // 2 for the selectors [ ]
	fld.offset = 0

	return nil
}

func (fld *InputField[T]) Render(b *strings.Builder, parent Coords) {
	// Move top left
	b.WriteString(MoveTo(parent.X+fld.coords.X, parent.Y+fld.coords.Y))

	// [ ?
	if fld.focused {
		b.WriteRune('[')
	} else {
		b.WriteRune(' ')
	}

	// field name
	b.WriteString(fld.label)

	// start highlight
	if fld.visual {
		b.WriteString(HIGHLIGHT_START)
	}

	// Safe rune slicing
	length := len(fld.buffer)
	viewEnd := min(fld.offset+fld.width, length)

	if fld.offset < length {
		b.WriteString(string(fld.buffer[fld.offset:viewEnd]))
	}

	// Pad remaining space
	visibleLen := viewEnd - fld.offset
	if visibleLen < fld.width {
		b.WriteString(strings.Repeat(" ", fld.width-visibleLen))
	}

	// end highlight
	if fld.visual {
		b.WriteString(HIGHLIGHT_END)
	}

	// ] ?
	if fld.focused {
		b.WriteRune(']')
	} else {
		b.WriteRune(' ')
	}
}

func (fld *InputField[T]) Cursor(b *strings.Builder, parent Coords) {
	if fld.visual || !fld.focused {
		b.WriteString(HIDE_CURSOR)
		return
	}

	x := parent.X + fld.coords.X + len(fld.label) + 1 // 1 for [
	b.WriteString(MoveTo(x+fld.cursor, parent.Y+fld.coords.Y))
	b.WriteString(BLINKING_IBEAM)
	b.WriteString(SHOW_CURSOR)
}

// Interaction
func (fld *InputField[T]) HandleKey(key string) bool {
	accept := fld.accept(key)
	handled := true

	if accept && !fld.focused {
		fld.Focus()
		return true
	}

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
			fld.delete()
		}

	case KEY_RIGHT:
		fld.cursorRight()

	case KEY_LEFT:
		fld.cursorLeft()

	case "\x01": // ctrl+a
		fld.Focus()

	default:
		if accept {
			fld.insert(key)
		} else {
			handled = false
		}
	}

	if handled {
		log.Printf("[field] %s\t%x\n", fld.label, key)
	}

	return handled
}

// State Management
func (fld *InputField[T]) Focus() {
	fld.visual = len(fld.buffer) != 0
	fld.focused = true
}

func (fld *InputField[T]) Blur() {
	// get buffer input
	input := string(fld.buffer)
	value, err := fld.parse(input)

	// if parsed successfully try to submit and update with new value
	// else return to previous valid value
	if err == nil {
		fld.Value = fld.submit(value)
	}

	// display new value
	fld.buffer = []rune(fld.format(fld.Value))

	// remove focus
	fld.focused = false
	fld.visual = false
	fld.cursor = 0
}

// private methods

func (fld *InputField[T]) cursorLeft() {
	if fld.visual {
		fld.visual = false
		fld.offset = 0
		fld.cursor = 0

	} else if fld.cursor > 0 {
		fld.cursor--

	} else if fld.offset > 0 {
		fld.offset--
	}
}

func (fld *InputField[T]) cursorRight() {
	length := len(fld.buffer)

	if fld.visual {
		fld.visual = false

		if length > fld.width {
			fld.offset = length - fld.width
			fld.cursor = fld.width

		} else {
			fld.offset = 0
			fld.cursor = length
		}

	} else if length < fld.width+fld.offset {
		fld.cursor = min(fld.cursor+1, length-fld.offset)

	} else if fld.cursor >= fld.width {
		fld.offset = min(fld.offset+1, length-fld.width)

	} else {
		fld.cursor++
	}
}

func (fld *InputField[T]) clear() {
	fld.buffer = nil // ?
	fld.cursor = 0
	fld.offset = 0
	fld.visual = false
}

func (fld *InputField[T]) backspace() {
	index := fld.cursor + fld.offset
	if index == 0 || index > len(fld.buffer) {
		return
	}

	// remove rune before index
	fld.buffer = append(fld.buffer[:index-1], fld.buffer[index:]...)

	// update position
	if fld.offset > 0 {
		fld.offset--

	} else if fld.cursor > 0 {
		fld.cursor--
	}
}

func (fld *InputField[T]) delete() {
	index := fld.cursor + fld.offset
	if index >= len(fld.buffer) {
		return
	}

	// remove rune after index
	fld.buffer = append(fld.buffer[:index], fld.buffer[index+1:]...)

	// update position
	if fld.offset > 0 && fld.cursor < fld.width { // ?
		fld.offset--
		fld.cursor++
	}
}

func (fld *InputField[T]) insert(s string) {
	if fld.visual {
		fld.clear()
	}

	// insert runes
	runes := []rune(s)
	length := len(runes)
	index := fld.offset + fld.cursor

	buffer := make([]rune, 0, len(fld.buffer)+length)
	buffer = append(buffer, fld.buffer[:index]...)
	buffer = append(buffer, runes...)
	buffer = append(buffer, fld.buffer[index:]...)

	fld.buffer = buffer

	// update position
	fld.cursor += length
	if fld.cursor > fld.width {
		fld.offset += fld.cursor - fld.width
		fld.cursor = fld.width
	}
}
