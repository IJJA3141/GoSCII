package tui

import (
	"fmt"
	"log"
)

type InputField[T any] struct {
	width, minWidth int // refers only to the width of the input filed (doesn't include `len(label)` nor '[]'

	relative Coords
	label    string

	buffer []rune
	value  T // ?

	accept func(string) bool // might change string to key type
	parse  func(string) (T, error)
	submit func(T)
	format func(T) string

	offset  int
	cursor  int
	visual  bool
	focused bool
}

//============
//  Layout

func (fld *InputField[T]) Resize(width, height int) error {
	if height < 1 {
		return fmt.Errorf("Failed to resize InputField, height %d is too small.", height)
	}

	if width < len(fld.label)+2+fld.minWidth {
		return fmt.Errorf("Failed to resize InputField, width %d is too small.", width)
	}

	fld.width = width - len(fld.label) - 2 // 2 for the selectors [ ]
	fld.offset = 0

	return nil
}

func (fld *InputField[T]) Render(sb *ScreenBuffer, origine Coords) {
	// Move top left
	sb.WriteString(origine.Add(fld.relative).MoveTo())

	// [ ?
	if fld.focused {
		sb.WriteRune('[')
	} else {
		sb.WriteRune(' ')
	}

	// field name
	sb.WriteString(fld.label)

	// start highlight
	if fld.visual {
		sb.WriteString(HIGHLIGHT_START)
	}

	// Safe rune slicing
	length := len(fld.buffer)
	viewEnd := min(fld.offset+fld.width, length)

	if fld.offset < length {
		sb.WriteRunes(fld.buffer[fld.offset:viewEnd])
	}

	// Pad remaining space
	visibleLen := viewEnd - fld.offset
	if visibleLen < fld.width {
		sb.WriteRuneN(' ', fld.width-visibleLen)
	}

	// end highlight
	if fld.visual {
		sb.WriteString(HIGHLIGHT_END)
	}

	// ] ?
	if fld.focused {
		sb.WriteRune(']')
	} else {
		sb.WriteRune(' ')
	}
}

func (fld *InputField[T]) Cursor(b *ScreenBuffer, origine Coords) {
	if fld.visual || !fld.focused {
		b.WriteString(HIDE_CURSOR)
		return
	}

	offset := Coords{fld.cursor + len(fld.label) + 1, 0} // 1 for [
	cursor := origine.Add(fld.relative).Add(offset)

	b.WriteString(cursor.MoveTo())
	b.WriteString(BLINKING_IBEAM)
	b.WriteString(SHOW_CURSOR)
}

//================
//  Interaction

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

//======================
//  State Management

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
		fld.SetValue(value)
		fld.submit(value)
	}

	// remove focus
	fld.focused = false
	fld.visual = false
	fld.cursor = 0
}

func (fld * InputField[T]) SetValue(value T) {
	fld.value = value
	fld.buffer = []rune(fld.format(value))
}

//=====================
//  Local functions

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
