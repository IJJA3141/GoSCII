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
const buf_size = 1024

var (
	ErrNotTerm = errors.New("output wasn't a terminal")
)

type a struct {
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
	}

	width, height, err = term.GetSize(output.Fd())
	if err != nil {
	}

	/// sig listener
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	// key listener
	r, err := cancelreader.NewReader(input)
	if err != nil {
	}

	// var str strings.Builder

	// the end  of input longer than 1024 bytes will be cut off
	keyChannel := make(chan a, 1)
	go func() {
		buf := make([]byte, buf_size)

		for {
			size, err := r.Read(buf)
			if size != buf_size {
				keyChannel <- a{Buf: buf, Size: size, err: err}
				continue
			}

			acc := make([]byte, buf_size)
			ss := size
			for ss == buf_size {
				ss, err = r.Read(acc)
				if err != nil {
					break
				}

				size += ss
				buf = append(buf, acc[:ss]...)
			}

			keyChannel <- a{Buf: buf, Size: size, err: err}
		}
	}()

	x, y := height/2, width/2

	text1 := tui.NewTextField(10, x-2, y-13)
	text2 := tui.NewIntField(10, x-2, y+13, false)
	text3 := tui.NewFloatField(4, x+2, y-20, false)

	check1 := tui.NewCheckBox(x+2, y+6, true)
	check2 := tui.NewRoundCheckBox(x+2, y+8, false)
	check3 := tui.NewSquaredCheckBox(x+2, y+10, true)
	check4 := tui.NewLockCheckBox(x+2, y+12, true)

	index := 0
	fArr := []tui.Focusable{&text1, &text2, &text3, &check1, &check2, &check3, &check4}

	output.Write([]byte("\x1b[?1049h"))
	var b strings.Builder

	for running {
		b.Reset()

		b.WriteString("\x1b[2J")

		text1.Draw(&b)
		text2.Draw(&b)
		text3.Draw(&b)

		check1.Draw(&b)
		check2.Draw(&b)
		check3.Draw(&b)
		check4.Draw(&b)

		// Box(&b, text1.X-1, text1.Y-1, text1.Width+2, 3)
		// Box(&b, text2.X-1, text2.Y-1, text2.Width+2, 3)
		// Box(&b, text3.X-1, text3.Y-1, text3.Width+2, 3)

		fArr[index].Cursor(&b)

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

			fArr[index].HandleKey(string(ev.Buf[:ev.Size]))

			if string(ev.Buf[:ev.Size]) == "\x1B[Z" {
				fArr[index].Blur()

				index--

				if index < 0 {
					index = len(fArr) - 1
				}

				fArr[index].Focus()
			}

			if string(ev.Buf[:ev.Size]) == "\t" {
				fArr[index].Blur()

				index++

				if index >= len(fArr) {
					index = 0
				}

				fArr[index].Focus()
			}

			if string(ev.Buf[:ev.Size]) == "\x1b" {
				fArr[index].Blur()
			}

			if ev.Buf[0] == 'q' {
				running = false
			}
		}
	}

	output.Write([]byte("\x1b[?1049l"))
	term.Restore(output.Fd(), s)
	return err
}
