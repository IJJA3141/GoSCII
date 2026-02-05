package tui

import "fmt"

const (
	BLINKING_CARET     = "\x1b[\x31 q"
	STEADY_CARET       = "\x1b[\x32 q"
	BLINKING_UNDERLINE = "\x1b[\x33 q"
	STEADY_UNDERLINE   = "\x1b[\x34 q"
	BLINKING_IBEAM     = "\x1b[\x35 q"
	STEADY_IBEAM       = "\x1b[\x36 q"
)

func MoveTo(x, y int) string { return fmt.Sprintf("\x1b[%d;%df", x, y) }
