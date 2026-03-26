package main

import (
	"errors"
	"flag"
	"fmt"
	_ "image/jpeg"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IJJA3141/GoSCII/filters"
	"github.com/IJJA3141/GoSCII/io"
	"github.com/IJJA3141/GoSCII/tui"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

var in string

func init() {
	flag.StringVar(&in, "in", "./example_images/test_uwu.png", "path to the input image")
}

type State struct {
	running bool

	frameSelected bool

	bar   tui.Bar
	frame tui.Frame

	in, out *os.File
	img     *filters.RGBAPlane
}

var state State

func main() {
	// args parsing + initial initialState initialization
	flag.Parse()

	var err error
	state.img, err = io.Read(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	// logging
	f, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Ignore error from opening file
	log.SetOutput(f)
	log.Print("\n\n\n")
	defer f.Close()

	// default for now
	state.in = os.Stdin
	state.out = os.Stdin
	// state.out = os.Stdout
	state.frameSelected = false

	term_, err := tui.InitializeTerm(state.in, state.out)
	defer tui.DeinitializeTerm(state.in, state.out, term_)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}

	width, height, err := term.GetSize(int(state.out.Fd())) // TODO might need to switch with state.in
	state.img, err = state.img.LanczosResize(width, height, 3)
	log.Printf("WIDTH  -> %d\n", width)
	log.Printf("HEIGHT -> %d\n\n\n", height)
	if err != nil {
		fmt.Println(err)
		return
	}

	gray := state.img.ToGrayScale()
	gray, err = gray.BayerDithering(4)
	if err != nil {
		fmt.Println(err)
		return
	}

	state.frame, state.bar = tui.CreateUI(width, height, gray.Braille(200))

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
	state.bar.SetCoord(tui.Coords{X: width - BAR_WIDTH, Y: 0})
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
	state.frame.Render(b)
	state.bar.Render(b)
	state.bar.Cursor(b)
	state.out.WriteString(b.String())
}
