package tui

import ()

type InputField[T any] struct {
	buffer string

	isValidInput func(string) bool

	visual bool
}

// return whether the key has been consumed and whether the field submits his buffer
func (fld *InputField[T]) HandleKey(key string) (bool, bool) {
	// not really sure of the consumed one...

	switch key {
	case KEY_ENTER:
		return true, true

	case KEY_ESC:
		return true, true

	case KEY_BACKSPACE:

	case KEY_DEL:

	case KEY_RIGHT:

	case KEY_LEFT:

	case "\x01": // crtl a

	default:
		if !fld.isValidInput(key) {
			return false, false // input wasn't recognized
		}

		// handle key
	}

	return true, false
}
