package main

import (
	"fmt"
	"os"
	// "os/signal"
	// "strings"
	// "syscall"

	"github.com/IJJA3141/GoSCII/tui"
	// "github.com/charmbracelet/x/term"
	// "github.com/muesli/cancelreader" // We may want to move away from this: using epoll instead of poll is suboptimal.
)

func main() {
	file := os.Stdin

	// s, _ := term.MakeRaw(file.Fd())

	err := tui.StartTui("", file, file)
	if err != nil {
		fmt.Print("err::")
		fmt.Println(err)
	}

	// term.Restore(file.Fd(), s)
}

// type event any
//
// type resizeEvent struct {
// 	width, height int
// }
//
// type kbdEvent struct {
// 	key string
// }
//
// func start1(out chan event) {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, syscall.SIGINT, syscall.SIGWINCH)
//
// 	go func() {
// 		for {
// 			<-c
// 			width, height, err := term.GetSize(os.Stdin.Fd())
// 			if err != nil {
// 			}
//
// 			out <- resizeEvent{width, height}
// 		}
// 	}()
// }
//
// func start2(out chan event) {
// 	file := os.Stdin
// 	r, err := cancelreader.NewReader(file)
// 	if err != nil {
// 	}
//
// 	go func() {
// 		for {
// 			var buf [1024]byte
// 			n, err := r.Read(buf[:])
// 			if err != nil {
// 			}
//
// 			out <- kbdEvent{string(buf[:n])}
// 			// t := binary.LittleEndian.Uint16(buf[:])
// 			// fmt.Print(t)
// 			// fmt.Print("\t")
// 			// if t == 768 ||
// 			// 	buf[0] == 'q' {
// 			// 	break
// 			// }
// 		}
// 	}()
// }
//
// func main() {
// 	file := os.Stdin
//
// 	s, _ := term.MakeRaw(file.Fd())
//
// 	width, height, _ := term.GetSize(os.Stdin.Fd())
//
// 	print("\x1b[2J\x1b[H")
//
// 	str := strings.Repeat("A", width) + "\r\n"
// 	for range height {
// 		print(str)
// 	}
//
// 	// go func() {
// 	// 	time.Sleep(50000)
// 	// 	str := strings.Repeat(" ", 15)
// 	// 	for i := range 5 {
// 	// 		print("\x1b["+fmt.Sprint(i+5)+";10H", str)
// 	// 	}
// 	// }()
// 	// go func() {
// 	// 	str := strings.Repeat("-", 15)
// 	// 	for i := range 10 {
// 	// 		print("\x1b["+fmt.Sprint(i+5)+";15H", str)
// 	// 	}
// 	// }()
//
// 	event_channel := make(chan event, 1)
//
// 	// term events
// 	start1(event_channel)
//
// 	// user events
// 	start2(event_channel)
//
// 	var u = 5
// 	var o = 15
//
// exit:
// 	for {
// 		// wait event
// 		// update with event
// 		// new render
//
// 		switch event := (<-event_channel).(type) {
// 		case resizeEvent:
// 			fmt.Printf("%dx%d", event.width, event.height)
//
// 		case kbdEvent:
// 			if event.key == "q" {
// 				break exit
// 			} else if event.key == "n" {
// 				str := strings.Repeat(" ", 15)
// 				for i := range 10 {
// 					print("\x1b["+fmt.Sprint(i+u)+";"+fmt.Sprint(o)+"H", str)
// 				}
// 				u += 20
// 				if u > height {
// 					u = 0
// 					o += 20
// 				}
// 			} else {
// 				print([]byte(event.key))
// 				print("\t")
// 			}
// 		}
// 	}
//
// 	term.Restore(file.Fd(), s)
// }
