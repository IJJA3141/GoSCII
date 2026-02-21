package main

import (
	"errors"
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
	file := os.Stdin

	// s, _ := term.MakeRaw(file.Fd())

	err := StartTui("", file, file)
	if err != nil {
		fmt.Print("err::")
		fmt.Println(err)
	}

	// term.Restore(file.Fd(), s)
}

const bar_width = 30
const min_width = 100
const min_height = 30
const KEY_BUFFER_SIZE = 1024

var (
	ErrNotTerm = errors.New("output wasn't a terminal")
)

type Key struct {
	Buf  []byte
	Size int
	err  error
}

var (
	box   = [...]string{"─", "│", "┌", "┐", "└", "┘", "├", "┤", "┬", "┴", "┼"}
	round = [...]string{"╭", "╮", "╯", "╰"}
)

func Box(dst *strings.Builder, x, y, width, height int) {
	tui.MoveTo(dst, x, y)
	dst.WriteString(round[0] + strings.Repeat(box[0], width-2) + round[1])
	for i := 1; i < height-1; i++ {
		tui.MoveTo(dst, x+i, y)
		dst.WriteString(box[1])
		tui.MoveTo(dst, x+i, y+width-1)
		dst.WriteString(box[1])
	}
	tui.MoveTo(dst, x+height-1, y)
	dst.WriteString(round[3] + strings.Repeat(box[0], width-2) + round[2])
}

func StartTui(image_path string, input *os.File, output *os.File) error {
	/// vars
	var running bool = true
	var width, height int
	var err error

	/// init
	s, err := term.MakeRaw(output.Fd())
	if err != nil {
		// TODO Handle error
	}

	// save screen & cursor
	output.Write([]byte("\x1b[?1049h\x1b 7"))

	// restore
	defer output.WriteString("\x1b[?1049l\x1b 8" + tui.SHOW_CURSOR + tui.BLINKING_IBEAM) // term content
	defer output.WriteString("\x1b[?25h")                                                // cursor
	defer term.Restore(output.Fd(), s)                                                   // term settings

	width, height, err = term.GetSize(output.Fd())
	if err != nil {
		// TODO Handle error
	}

	/// sig listener
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	// key listener
	r, err := cancelreader.NewReader(input)
	if err != nil {
		// TODO Handle error
	}

	keyChannel := make(chan Key, 1)
	go func() {
		buffer := make([]byte, KEY_BUFFER_SIZE)

		for {
			size, err := r.Read(buffer)
			if size != KEY_BUFFER_SIZE || err != nil {
				keyChannel <- Key{Buf: buffer, Size: size, err: err}
				continue
			}

			// if read more than buffer size
			// enter a loop to read until
			// read less than buffer size

			accSize := size
			accBuffer := make([]byte, KEY_BUFFER_SIZE)
			copy(accBuffer, buffer[:size])

			for size == KEY_BUFFER_SIZE && err == nil {
				size, err = r.Read(accBuffer)

				accSize += size
				accBuffer = append(accBuffer, buffer[:size]...)
			}

			keyChannel <- Key{Buf: accBuffer, Size: accSize, err: err}
		}
	}()

	/// temp
	// x, y := width/2, height/2
	//
	// text1 := tui.NewTextField(10, x-2, y-13)
	// text2 := tui.NewIntField(10, x-2, y+13, false)
	// text3 := tui.NewFloatField(4, x+2, y-20, false)
	//
	// check1 := tui.NewCheckBox(x+2, y+6, true)
	// check2 := tui.NewRoundCheckBox(x+2, y+8, false)
	// check3 := tui.NewSquaredCheckBox(x+2, y+10, true)
	// check4 := tui.NewLockCheckBox(x+2, y+12, true)

	bar := tui.NewBar(30, height, width-31, 1)

	focusIndex := 0
	// focusables := []tui.Focusable{&text1, &text2, &text3, &check1, &check2, &check3, &check4}
	focusables := []tui.Focusable{&bar}

	var b strings.Builder

	for running {
		b.Reset()

		// fmt.Fprintf(&b, "\x1b[48;2;%d;%d;%dm", 0x1E, 0x1C, 0x1C)
		// fmt.Fprintf(&b, "\x1b[48;2;%d;%d;%dm", 0xff, 0xff, 0xff)
		fmt.Fprintf(&b, "\x1b[39;49m")
		b.WriteString("\x1b[2J")

		// fmt.Fprintf(&b, "\x1b[48;2;%d;%d;%dm", 0x2C, 0x2B, 0x2B)
		for _, focusable := range focusables {
			// b.WriteString("[")
			focusable.Draw(&b)
			// b.WriteString("]")
		}

		drawBorder(width, height, width-30, &b)
		focusables[focusIndex].Cursor(&b)

		output.WriteString(b.String())

		select {
		case signal := <-signalChannel:
			switch signal {
			case syscall.SIGINT, syscall.SIGTERM:
				running = false

			case syscall.SIGWINCH:
				width, height, err = term.GetSize(output.Fd())
				if err != nil {
				}

				///
			}

		case ev := <-keyChannel:
			if ev.err != nil {
			}

			///

			focusables[focusIndex].HandleKey(string(ev.Buf[:ev.Size]))

			// if string(ev.Buf[:ev.Size]) == "\x1B[Z" {
			// 	focusables[focusIndex].Blur()
			//
			// 	focusIndex--
			//
			// 	if focusIndex < 0 {
			// 		focusIndex = len(focusables) - 1
			// 	}
			//
			// 	focusables[focusIndex].Focus()
			// }
			//
			// if string(ev.Buf[:ev.Size]) == "\t" {
			// 	focusables[focusIndex].Blur()
			//
			// 	focusIndex++
			//
			// 	if focusIndex >= len(focusables) {
			// 		focusIndex = 0
			// 	}
			//
			// 	focusables[focusIndex].Focus()
			// }
			//
			// if string(ev.Buf[:ev.Size]) == "\x1b" {
			// 	focusables[focusIndex].Blur()
			// }

			if ev.Buf[0] == 'q' {
				running = false
			}
		}
	}

	return err
}

func drawBorder(width, height, split int, b *strings.Builder) {
	tui.MoveTo(b, 0, 0)

	b.WriteString(round[0])

	b.WriteString(strings.Repeat(box[0], split-2))
	b.WriteString(box[8])
	b.WriteString(strings.Repeat(box[0], width-split-1))

	b.WriteString(round[1])

	for i := 2; i < height; i++ {
		tui.MoveTo(b, 0, i)
		b.WriteString(box[1])

		tui.MoveTo(b, split, i)
		b.WriteString(box[1])

		tui.MoveTo(b, width, i)
		b.WriteString(box[1])
	}

	tui.MoveTo(b, 0, height)
	b.WriteString(round[3])
	b.WriteString(strings.Repeat(box[0], split-1))
	b.WriteString(box[9])
	b.WriteString(strings.Repeat(box[0], width-split-2))
	b.WriteString(round[2])
}
