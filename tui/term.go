package tui

import (
	"os"

	"golang.org/x/term"
)

func InitializeTerm(in, out *os.File) (*term.State, error) {
	/// term settings
	// set term in raw mod and save old state
	state, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return nil, err
	}

	// save screen & cursor
	out.Write([]byte("\x1b[?1049h\x1b 7"))

	return state, nil
}

func DeinitializeTerm(in, out *os.File, state *term.State) {
	// restore
	out.WriteString("\x1b[?1049l\x1b 8")                        // term content
	out.WriteString("\x1b[?25h" + SHOW_CURSOR + BLINKING_IBEAM) // cursor
	term.Restore(int(in.Fd()), state)
}

func CreateUI(width, height int, img Image) (Frame, Bar) {
	return NewFrame(width-BAR_MIN_WIDTH, height, Coords{0, 0}, img), NewBar(Coords{X: width - BAR_MIN_WIDTH, Y: 0})
}
