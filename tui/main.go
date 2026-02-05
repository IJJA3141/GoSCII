package tui

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charmbracelet/x/term"
	"github.com/muesli/cancelreader"
)

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

func Box(x, y, width, height int) string {
	out := MoveTo(x, y) + box[2] + strings.Repeat(box[0], width-2) + box[3]
	for i := 1; i < height-1; i++ {
		out += MoveTo(x+i, y) + box[1] + MoveTo(x+i, y+width-1) + box[1]
	}
	return out + MoveTo(x+height-1, y) + box[4] + strings.Repeat(box[0], width-2) + box[5]
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

	test1 := NumBox{
		str:   "",
		width: 10,
	}

	output.Write([]byte("\x1b[?1049h"))

	for running {
		output.WriteString("\x1b[2J")
		output.WriteString(test1.Draw(height/2, width/2))
		output.WriteString(Box(height/2-1, width/2-1, 12, 3))
		output.WriteString(test1.Cursor(height/2, width/2))

		// out.WriteString(Box(1, width-29, 30, height-3))
		// out.WriteString(Box(1, 1, width-30, height))
		// out.WriteString(Box(height-2, width-29, 30, 3))
		// out.WriteString(MoveTo(height-1, width-28))
		// out.WriteString(BLINKING_IBEAM)
		// out.WriteString(MoveTo(height-1, width-28))
		// out.WriteString(str.String())

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
			test1.HandleKey(string(ev.Buf[:ev.Size]))
			if ev.Buf[0] == 'q' {
				running = false
			}
		}
	}

	output.Write([]byte("\x1b[?1049l"))
	term.Restore(output.Fd(), s)
	return err
}
