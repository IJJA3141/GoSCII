package tui

import (
	"fmt"
	"strings"
)

const MOVE_TO_APROX_SIZE = len("\x1b[xxx;xxxf")

func MoveTo(dst *strings.Builder, x, y int) { fmt.Fprintf(dst, "\x1b[%d;%df", x, y) }

// keys
const (
	KEY_ENTER     = "\r"
	KEY_BACKSPACE = "\x7f" // not really
	KEY_DEL       = "\x1b[3~"
	KEY_UP        = "\x1b[A"
	KEY_DOWN      = "\x1b[B"
	KEY_RIGHT     = "\x1b[C"
	KEY_LEFT      = "\x1b[D"
	KEY_ESC       = "\x1b["
)

// cursor control
const (
	BLINKING_CARET     = "\x1b[\x31 q"
	STEADY_CARET       = "\x1b[\x32 q"
	BLINKING_UNDERLINE = "\x1b[\x33 q"
	STEADY_UNDERLINE   = "\x1b[\x34 q"
	BLINKING_IBEAM     = "\x1b[\x35 q"
	STEADY_IBEAM       = "\x1b[\x36 q"

	HIDE_CURSOR = "\x1b[?25l"
	SHOW_CURSOR = "\x1b[?25h"
)

// ascii bounds
const (
	ASCII_PRINTABLE_LOWER_BOUND = "\x20"
	ASCII_PRINTABLE_UPPER_BOUND = "\x7f"
)

// ansi cmd
const (
	HILIGHT_START = "\x1b[7m"
	HILIGHT_END   = "\x1b[27m"
)
