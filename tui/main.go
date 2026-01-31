package tui

import (
	"errors"
	"os"
	"syscall"

	"github.com/charmbracelet/x/term"
)

const bar_width = 30
const min_width = 100
const min_height = 30

type state struct {
	running  bool
	drawable bool

	eventQueue chan event

	prevTerm *term.State

	width, height int

	input  *os.File
	output *os.File

	frame Frame
	bar   Bar
}

var (
	ErrNotTerm = errors.New("output wasn't a terminal")
)

func StartTui(image_path string, input *os.File, output *os.File) error {
	var err error

	state := state{input: input, output: output, eventQueue: make(chan event)}

	err = state.Initialize()
	if err != nil {
		return err
	}

	for state.running {
		state.Draw()

		event := <-state.eventQueue
		switch event := event.(type) {
		case syscall.Signal:
			err = state.HandleSignal(event)

		case error:
			err = state.HandleError(event)

		case string:
			state.HandleKey(event)
		}
	}

	return errors.Join(err, state.Close())
}

func (this *state) Initialize() error {
	if !term.IsTerminal(this.output.Fd()) {
		return ErrNotTerm
	}

	var err error
	this.prevTerm, err = term.MakeRaw(this.input.Fd())
	if err != nil {
		return err
	}

	this.width, this.height, err = term.GetSize(this.input.Fd())
	if err != nil {
		return err
	}

	SignalListener(this.eventQueue, syscall.SIGWINCH)
	err = KeyboardListener(this.eventQueue, this.input)
	if err != nil {
		return err
	}

	this.output.Write([]byte("\x1b[2J\x1b[H"))

	this.running = true
	return nil
}

func (this *state) Close() error {
	var err error

	if this.prevTerm != nil {
		err = term.Restore(this.output.Fd(), this.prevTerm)
	}

	return err
}

func (frame *Frame) resize(width, height int) {}
func (bar *Bar) resize(width, height int)     {}

func (this *state) HandleSignal(sig syscall.Signal) error {
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		this.output.Write([]byte("!"))
		this.running = false

	case syscall.SIGWINCH:
		var err error
		this.width, this.height, err = term.GetSize(this.input.Fd())
		if err != nil {

			// might change
			this.running = false

			return err
		}

		this.drawable = this.width >= min_height && this.height >= min_height
		if this.drawable {
			this.frame.resize(this.width-bar_width, this.height)
			this.bar.resize(bar_width, this.height)
		}
	}

	return nil
}

func (this *state) HandleError(err error) error {
	this.running = false
	return err
}

func (this *state) HandleKey(key string) {
	print("new key")

	if key == "q" {
		this.running = false
	}

	this.bar.HandleKey(key) // for now
}

func (this *state) Draw() {
	// clear canvas

	// draw child

	// set cursor

	this.output.Write([]byte("Hello world!\n" + this.bar.Draw(5, 5)))
}
