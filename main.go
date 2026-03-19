package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IJJA3141/GoSCII/tui"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

type State struct {
	running bool

	frameSelected bool

	bar   tui.Bar
	frame tui.Frame

	in, out *os.File
}

var state State

func main() {
	// args parsing + initial initialState initialization

	f, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Ignore error from opening file
	log.SetOutput(f)
	log.Print("\n\n\n")
	defer f.Close()

	// default for now
	state.in = os.Stdin
	state.out = os.Stdin
	// state.out = os.Stdout
	state.frameSelected = false

	term, err := tui.InitializeTerm(state.in, state.out)
	defer tui.DeinitializeTerm(state.in, state.out, term)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}

	state.frame, state.bar = tui.CreateUI()

	err = start()
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}

	// write to output logic
}

func start() (err error) {
	/// listeners
	// signal listener
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	// key listener
	var reader cancelreader.CancelReader
	reader, err = cancelreader.NewReader(state.in)
	if err != nil {
		return
	}

	keyChannel := make(chan tui.Key, 1)
	tui.Notify(keyChannel, reader)

	var b strings.Builder

	/// main loop
	state.running = true
	for state.running {
		/// update frame

		/// draw ui
		render(&b)

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
				err = key.Err
			} else {
				handleKey(string(key.Buf[:key.Size]))
			}
		}
	}

	return nil
}

const BAR_WIDTH = 30

func handleResize() {
	width, height, err := term.GetSize(int(state.out.Fd())) // TODO might need to switch with state.in
	errors.Join(err, state.frame.Resize(width-BAR_WIDTH, height), state.bar.Resize(BAR_WIDTH, height))
	state.bar.SetCoord(tui.Coord{X: width - BAR_WIDTH, Y: 0})
	if err != nil {
		// TODO handle erorr
	}

}

func handleKey(key string) {
	var handled bool

	// if state.frameSelected {
	// 	handled = state.frame.HandleKey(key)
	//
	// } else {
	handled = state.bar.HandleKey(key)
	// }

	if !handled {
		log.Printf("[main] %s\t%x\n", key, key)

		switch key {
		case "q", "\x03", tui.KEY_TAB:
			state.running = false
		}
	}
}

func render(b *strings.Builder) {
	b.Reset()
	b.WriteString(tui.CLEAR_SCREEN)
	state.bar.Render(b)
	state.bar.Cursor(b)
	state.out.WriteString(b.String())
}
