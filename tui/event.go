package tui

import (
	// "fmt"
	"io"
	"os"
	"os/signal"

	"github.com/muesli/cancelreader"
)

type event any

func SignalListener(out chan event, running *bool, sig ...os.Signal) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, sig...)
	go func() {
		for *running {
			out <- <-channel
		}
	}()
}

func KeyboardListener(out chan event, running *bool, source io.Reader) error {
	r, err := cancelreader.NewReader(source)
	if err != nil {
		return err
	}

	go func() {
		// var buf [1024]byte
		var buf [4]byte

		// the end  of input longer than 1024 bytes will be cut off
		for *running {
			size, err := r.Read(buf[:])
			if err != nil {
				out <- err
			} else {
				out <- struct {
					Buf  []byte
					Size int
				}{buf[:], size}
			}
		}
	}()

	return nil
}
