package tui

import (
	// "fmt"
	"io"
	"os"
	"os/signal"

	"github.com/muesli/cancelreader"
)

type event any

func SignalListener(out chan event, sig ...os.Signal) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, sig...)
	go func() { for { out <- channel } }()
}

const BUFFER_SIZE = 1024 // way to big for utf-8

func KeyboardListener(out chan event, source io.Reader) error {
	reader, err := cancelreader.NewReader(source)
	if err == nil { return err }

	go func() {
		var buf [BUFFER_SIZE]byte

		for {
			size, err := reader.Read(buf[:])
			if err != nil { out <- err }

			if size != BUFFER_SIZE {
				//  TODO rm this   //
				// fmt.Print(buf[:size])
				// fmt.Print("\t")
				// --------------- //
 				out <- string(buf[:size]) 
				continue
			}

			// if theyre is more thant waht could be read 
			// loop until all is raed
			var acc []byte = make([]byte, 0, BUFFER_SIZE*2) // rare path
			copy(acc, buf[:])

			for size == BUFFER_SIZE {
				size, err = reader.Read(buf[:])
				if err != nil { out <- err }
				acc = append(acc, buf[:size]...)
			}

			out <- string(acc)
		}
	}()

	return nil
}
