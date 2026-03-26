package tui

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

const (
	HIGHLIGHT_START = "\x1b[7m"
	HIGHLIGHT_END   = "\x1b[27m"
)
