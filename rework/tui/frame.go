package tui

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Image interface {
	Dimensions() (int, int)
	Get(x, y, width, height int) []string
}

const (
	FRAME_MIN_WIDTH  = 20
	FRAME_MIN_HEIGHT = 20
)

type Frame struct {
	width, height int
	origine       Coords
	view          Coords // top left x, y

	image Image
}

func NewFrame(width, height int, origine Coords, img Image) Frame {
	return Frame{
		width:   width,
		height:  height,
		origine: origine,
		view:    Coords{0, 0},
		image:   img,
	}
}

// Layout & Rendering
func (frm *Frame) Resize(width, height int) error {
	if height < FRAME_MIN_HEIGHT {
		return fmt.Errorf("Failed to resize InputField, height %d is too small.", height)
	}

	if width < FRAME_MIN_WIDTH {
		return fmt.Errorf("Failed to resize InputField, width %d is too small.", width)
	}

	frm.height = height
	frm.width = width
	frm.view.X = 0
	frm.view.Y = 0

	return nil
}

// Display cases
//
// frame ─ ─ ─
// image ─────
//
// case 0        case 1        case 2        case 3
//	                             ┌ ─ ─ ┐    ┌ ─ ─ ─ ─ ─ ┐
//	┌───────┐     ┌───────┐     ┌|─────|┐   | ┌───────┐ |
//	│┌ ─ ─ ┐│   ┌ ─ ─ ─ ─ ─ ┐   │|     |│   | │       │ |
//	│|     |│   | │       │ |   │|     |│   | │       │ |
//	│|     |│   | │       │ |   │|     |│   | │       │ |
//	│└ ─ ─ ┘│   └ ─ ─ ─ ─ ─ ┘   │|     |│   | │       │ |
//	└───────┘     └───────┘     └|─────|┘   | └───────┘ |
//	                             └ ─ ─ ┘    └ ─ ─ ─ ─ ─ ┘

func (frm *Frame) Render(b *ScreenBuffer) {
	coords := frm.origine

	width, height := frm.image.Dimensions()

	if height >= frm.height {
		if width >= frm.width {
			log.Println("case 0")
			//  case 0
			//
			// 	 ┌───────┐
			// 	 │┌ ─ ─ ┐│
			// 	 │|     |│
			// 	 │|     |│
			// 	 │└ ─ ─ ┘│
			// 	 └───────┘
			//

			for _, line := range frm.image.Get(frm.view.X, frm.view.Y, frm.width, frm.height) {
				b.WriteString(coords.MoveTo())
				b.WriteString(line)
				coords.Y++
			}

			return

		} else {
			log.Println("case 1")
			//  case 1
			//
			//   ┌───────┐
			// ┌ ─ ─ ─ ─ ─ ┐
			// | │       │ |
			// | │       │ |
			// └ ─ ─ ─ ─ ─ ┘
			//   └───────┘

			margin := float64(frm.width - width)
			leftMargin := strings.Repeat(" ", int(math.Floor(margin)))
			rightMargin := strings.Repeat(" ", int(math.Ceil(margin)))

			for _, line := range frm.image.Get(0, frm.view.Y, width, frm.height) {
				b.WriteString(coords.MoveTo())
				b.WriteString(leftMargin)
				b.WriteString(line)
				b.WriteString(rightMargin)
				coords.Y++
			}

			return
		}
	} else {
		if width > frm.width {
			log.Println("case 2")
			//  case 2
			//    ┌ ─ ─ ┐
			//   ┌|─────|┐
			//   │|     |│
			//   │|     |│
			//   │|     |│
			//   │|     |│
			//   └|─────|┘
			//    └ ─ ─ ┘

			margin := float64(frm.height-height) / 2.
			emptyLine := strings.Repeat(" ", frm.width)

			for range int(math.Floor(margin)) {
				b.WriteString(coords.MoveTo())
				b.WriteString(emptyLine)
				coords.Y++
			}

			for _, line := range frm.image.Get(frm.view.X, 0, frm.width, height) {
				b.WriteString(coords.MoveTo())
				b.WriteString(line)
				coords.Y++
			}

			for range int(math.Ceil(margin)) {
				b.WriteString(coords.MoveTo())
				b.WriteString(emptyLine)
				coords.Y++
			}

			return

		} else { // CASE3
			log.Println("case 3")
			//  case 3
			// ┌ ─ ─ ─ ─ ─ ┐
			// | ┌───────┐ |
			// | │       │ |
			// | │       │ |
			// | │       │ |
			// | │       │ |
			// | └───────┘ |
			// └ ─ ─ ─ ─ ─ ┘

			vmargin := float64(frm.height-height) / 2.
			hmargin := float64(frm.width-width) / 2.
			emptyLine := strings.Repeat(" ", frm.width)

			for range int(math.Floor(vmargin)) {
				b.WriteString(coords.MoveTo())
				b.WriteString(emptyLine)
				coords.Y++
			}

			for _, line := range frm.image.Get(frm.view.X, frm.view.Y, width, height) {
				b.WriteString(coords.MoveTo())
				b.WriteString(emptyLine[:int(math.Floor(hmargin))])
				b.WriteString(line)
				b.WriteString(emptyLine[:int(math.Ceil(hmargin))])
				coords.Y++
			}

			for range int(math.Ceil(vmargin)) {
				b.WriteString(coords.MoveTo())
				b.WriteString(emptyLine)
				coords.Y++
			}

			return
		}
	}
}

func (frm *Frame) Cursor(b *strings.Builder, parent Coords) {}

// Interaction
func (frm *Frame) HandleKey(key string) bool { return false }

// State Management
func (frm *Frame) Focus() {}
func (frm *Frame) Blur()  {}
