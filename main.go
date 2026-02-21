package goscii

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IJJA3141/GoSCII/tui"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/cancelreader"
)

func main() {
	// args parsing + initial initialState initialization
	state = State{
		running: false,
	}

	// default for now
	file := os.Stdin
	state.input = file
	state.output = file

	err := start()
	if err != nil {
		fmt.Println(state.err) // TODO ERR
	}

	// write to output logic
}

/// TUI

const BAR_WIDTH = 30
const MIN_WIDTH = 100
const MIN_HEIGHT = 35

type UI int
type Mode int

const (
	DISPLAY_NORMAL UI = iota
	DISPLAY_ERR
)

const (
	INPUT_NORMAL Mode = iota
	INPUT_VISUAL
	INPUT_NONE
)

type State struct {
	running       bool
	width, height int // runes
	err           error
	prevTerm      *term.State
	input, output *os.File

	borders tui.Borders
	frame   tui.Frame
	bar     tui.Bar

	display     UI
	interaction Mode

	tuiBuffer strings.Builder
}

var state State

func start() error {
	/// term settings
	// set term in raw mod and save old state
	state.prevTerm, state.err = term.MakeRaw(state.output.Fd())
	if state.err != nil {
		// TODO Handle error
		return state.err
	}

	// save screen & cursor
	state.output.Write([]byte("\x1b[?1049h\x1b 7"))

	// restore
	defer state.output.WriteString("\x1b[?1049l\x1b 8")                                // term content
	defer state.output.WriteString("\x1b[?25h" + tui.SHOW_CURSOR + tui.BLINKING_IBEAM) // cursor
	defer term.Restore(state.output.Fd(), state.prevTerm)

	// get term initial width and height
	state.width, state.height, state.err = term.GetSize(state.output.Fd())
	if state.err != nil {
		// TODO Handle error
		return state.err
	}

	/// listeners
	// signal listener
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	// key listener
	var reader cancelreader.CancelReader
	reader, state.err = cancelreader.NewReader(state.input)
	if state.err != nil {
		// TODO Handle error
		return state.err
	}

	keyChannel := make(chan tui.Key, 1)
	tui.Notify(keyChannel, reader)

	/// UI elements
	var barleft bool
	barleft = true
	if barleft {
		state.borders = tui.NewBorders(state.width, state.height, state.width-BAR_WIDTH)
		state.frame = tui.NewFrame(state.width-BAR_WIDTH-2, state.height-2, 1+BAR_WIDTH+1, 1)
		state.bar = tui.NewBar(BAR_WIDTH, state.height-2, 1, 1)

	} else {
		state.borders = tui.NewBorders(state.width, state.height, state.width-BAR_WIDTH)
		state.frame = tui.NewFrame(state.width-BAR_WIDTH-2, state.height-2, 1, 1)
		state.bar = tui.NewBar(BAR_WIDTH, state.height-2, state.width-BAR_WIDTH-1, 1)
	}

	/// main loop
	for state.running {
		/// update frame

		/// draw ui
		drawUI()

		/// await event
		select {
		case signal := <-signalChannel:
			switch signal {
			case syscall.SIGINT, syscall.SIGTERM:
				state.running = false

			case syscall.SIGWINCH:
				handleResize()
			}

		case key := <-keyChannel:
			if key.Err != nil {
				// TODO Handle error
				state.err = key.Err
			} else {
				handleKey(string(key.Buf[:key.Size]))
			}
		}
	}

	return nil
}

func drawUI() {
	// reset writer
	state.tuiBuffer.Reset()
	state.tuiBuffer.WriteString("\x1b[39;49m\x1b[2J")

	// elements
	switch state.display {
	case DISPLAY_NORMAL:
		state.borders.Draw(&state.tuiBuffer)
		state.frame.Draw(&state.tuiBuffer)
		state.bar.Draw(&state.tuiBuffer)

	case DISPLAY_ERR:
		// TODO handle err
	}

	// cursor
	switch state.interaction {
	case INPUT_NORMAL:
		state.bar.Cursor(&state.tuiBuffer)

	case INPUT_VISUAL:
		state.frame.Cursor(&state.tuiBuffer)

	default:
		state.tuiBuffer.WriteString(tui.HIDE_CURSOR)
	}

	// write
	state.output.WriteString(state.tuiBuffer.String())
}

func handleResize() {
	state.width, state.height, state.err = term.GetSize(state.output.Fd())
	if state.err != nil {
		return
	}

	if state.width < MIN_WIDTH || state.height < MIN_HEIGHT {
		// TODO Handle error
		state.err = fmt.Errorf("to small")
	}

	state.borders.Resize(state.width, state.height)
	state.frame.Resize(state.width, state.height)
	state.bar.Resize(state.width, state.height)
}

const CRTL = rune('a') - 0x01

func handleKey(key string) {
	switch state.interaction {
	case INPUT_NORMAL:
		if state.bar.HandleKey(key) {
			return
		}

		switch key {
		case "q", tui.KEY_ESC:
			state.running = false

		case string(CRTL + 'h'):
			state.interaction = INPUT_VISUAL
		}

	case INPUT_VISUAL:
		if state.frame.HandleKey(key) {
			return
		}

		switch key {
		case "q", tui.KEY_ESC:
			state.running = false

		case string(CRTL + 'l'):
			state.interaction = INPUT_NORMAL
		}

	default:
		switch key {
		case "q", tui.KEY_ESC:
			state.running = false
		}
	}
}
