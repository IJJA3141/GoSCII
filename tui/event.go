package tui

import "github.com/muesli/cancelreader"

type Key struct {
	Buf  []byte
	Size int
	Err  error
}

const KEY_BUFFER_SIZE = 100

func Notify(c chan<- Key, r cancelreader.CancelReader) {
	if c == nil {
		// TODO Handle error
		panic("key notify channel nil")
	}

	go func() {
		buffer := make([]byte, KEY_BUFFER_SIZE)

		for {
			size, err := r.Read(buffer)
			if size != KEY_BUFFER_SIZE || err != nil {
				c <- Key{Buf: buffer, Size: size, Err: err}
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

			c <- Key{Buf: accBuffer, Size: accSize, Err: err}
		}
	}()
}
