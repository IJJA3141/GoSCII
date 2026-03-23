package tui

import (
	"fmt"
	"math"
	"strings"
)

type Image interface {
	Width() int
	Height() int
	Get(x, y, width, height int) []string
}

const (
	FRAME_MIN_WIDTH  = 20
	FRAME_MIN_HEIGHT = 20
)

type Frame struct {
	width, height int
	relative      Coords
	view          Coords // top left x, y

	image Image
}

func NewFrame(a, b int, x, y int) Frame { return Frame{} }

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

func (frm *Frame) Render(b *strings.Builder, origine Coords) {
	x := origine.X + frm.relative.X
	y := origine.Y + frm.relative.Y

	if frm.image.Height() >= frm.height {
		if frm.image.Width() >= frm.width {
			//  case 0
			//
			// 	 ┌───────┐
			// 	 │┌ ─ ─ ┐│
			// 	 │|     |│
			// 	 │|     |│
			// 	 │└ ─ ─ ┘│
			// 	 └───────┘
			//

			for i, line := range frm.image.Get(frm.view.X, frm.view.Y, frm.width, frm.height) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(line)
			}

			return

		} else {
			//  case 1
			//
			//   ┌───────┐
			// ┌ ─ ─ ─ ─ ─ ┐
			// | │       │ |
			// | │       │ |
			// └ ─ ─ ─ ─ ─ ┘
			//   └───────┘

			margin := float64(frm.width - frm.image.Width())
			leftMargin := strings.Repeat(" ", int(math.Floor(margin)))
			rightMargin := strings.Repeat(" ", int(math.Ceil(margin)))

			for i, line := range frm.image.Get(0, frm.view.Y, frm.image.Width(), frm.height) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(leftMargin)
				b.WriteString(line)
				b.WriteString(rightMargin)
			}

			return
		}
	} else {
		if frm.image.Width() > frm.width {
			//  case 2
			//    ┌ ─ ─ ┐
			//   ┌|─────|┐
			//   │|     |│
			//   │|     |│
			//   │|     |│
			//   │|     |│
			//   └|─────|┘
			//    └ ─ ─ ┘

			margin := float64(frm.height-frm.image.Height()) / 2.
			emptyLine := strings.Repeat(" ", frm.width)

			var i int
			for i = range int(math.Floor(margin)) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(emptyLine)
			}

			for _, line := range frm.image.Get(frm.view.X, 0, frm.width, frm.image.Height()) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(line)
				i++
			}

			for range int(math.Ceil(margin)) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(emptyLine)
				i++
			}

			return

		} else { // CASE3
			//  case 3
			// ┌ ─ ─ ─ ─ ─ ┐
			// | ┌───────┐ |
			// | │       │ |
			// | │       │ |
			// | │       │ |
			// | │       │ |
			// | └───────┘ |
			// └ ─ ─ ─ ─ ─ ┘

			vmargin := float64(frm.height-frm.image.Height()) / 2.
			hmargin := float64(frm.width-frm.image.Width()) / 2.
			emptyLine := strings.Repeat(" ", frm.width)

			var i int
			for i = range int(math.Floor(vmargin)) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(emptyLine)
			}

			for _, line := range frm.image.Get(frm.view.X, frm.view.Y, frm.image.Width(), frm.image.Width()) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(emptyLine[:int(math.Floor(hmargin))])
				b.WriteString(line)
				b.WriteString(emptyLine[:int(math.Ceil(hmargin))])
				i++
			}

			for range int(math.Ceil(vmargin)) {
				b.WriteString(MoveTo(x, y+i))
				b.WriteString(emptyLine)
				i++
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
