package tui

import (
	"errors"
	"io"
	"os"
	"syscall"

	"github.com/charmbracelet/x/term"
)

const (
	FFRM = iota
	FBAR
	FCMD
)

type state struct {
	running    bool
	eventQueue chan event
	prevTerm   *term.State
	tty        *os.File

	width, height int
	redraw        bool

	focusIndex int

	frame   Frame
	bar     Bar
	command Command
}

var (
	ErrNotTerm = errors.New("")
)

func (this *Frame) HandleKey(key string) bool
func (this *Bar) HandleKey(key string) bool
func (this *Command) HandleKey(key string) bool

func StartTui(path string, input *os.File, output io.Writer) error {
	var err error

	state := NewState(path, input, output)
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
			var handledl bool

			switch state.focusIndex {
			case FFRM:
				handledl = state.frame.HandleKey(event)

			case FBAR:
				handledl = state.bar.HandleKey(event)

			case FCMD:
				handledl = state.command.HandleKey(event)
			}

			if !handledl {
				state.HandleKey(event)
			}
		}
	}

	return errors.Join(err, state.Close(input))
}

func NewState(imgPath string, in io.Reader, out io.Writer) state {
	return state{}
}

func (this *state) Initialize() error {
	if !term.IsTerminal(this.tty.Fd()) {
		return ErrNotTerm
	}

	var err error
	this.prevTerm, err = term.MakeRaw(this.tty.Fd())
	if err != nil {
		return err
	}

	SignalListener(this.eventQueue, syscall.SIGWINCH)
	err = KeyboardListener(this.eventQueue, this.tty)
	if err != nil {
		return err
	}

	this.running = true

	return nil
}

func (this *state) Close(input *os.File) error {
	var err error

	if this.prevTerm != nil {
		err = term.Restore(input.Fd(), this.prevTerm)
	}

	return err
}

func (this *state) HandleSignal(sig syscall.Signal) error {
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		this.running = false

	case syscall.SIGWINCH:
		var err error
		this.width, this.height, err = term.GetSize(this.tty.Fd())
		if err != nil {

			// might change
			this.running = false

			return err
		}

		this.redraw = true
	}

	return nil
}

func (this *state) HandleError(err error) error {
	this.running = false
	return err
}

func (this *state) HandleKey(key string) {
	switch key {
	case "q":
		this.running = false

	case "tab":
		this.focusIndex = min(this.focusIndex + 1)
	}
}

func (this *state) Draw() {

}
