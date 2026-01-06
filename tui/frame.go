package tui

import (
	"math"
	"strings"

	"github.com/IJJA3141/GoSCII/filters"
	tea "github.com/charmbracelet/bubbletea"
)

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

type frame struct {
	// array of "screen" row
	buffer        []string
	width, height int
	x, y          int

	emptyLine  string
	hasChanged bool
	src        filters.Ascii
}

func Frame(width, height int, image filters.Ascii) frame {
	return frame{
		buffer:     make([]string, height),
		width:      width,
		height:     height,
		x:          0,
		y:          0,
		emptyLine:  strings.Repeat("\x1B[0m ", width),
		hasChanged: true,
		src:        image,
	}
}

func (this *frame) Update(msg tea.KeyMsg) {
	switch msg.String() {
	case "h", tea.KeyLeft.String():
		width := this.src.Width_()
		if width > this.width { // case 0|2
			this.x = max(0, this.x-1)
		}

	case "l", tea.KeyRight.String():
		width := this.src.Width_()
		if width > this.width { // case 0|2
			this.x = min(this.x+1, width-this.width)
		}

	case "k", tea.KeyUp.String():
		height := this.src.Height_()
		if height > this.height { // case 0|1
			this.y = max(0, this.y-1)
		}

	case "j", tea.KeyDown.String():
		height := this.src.Height_()
		if height > this.height { // case 0|1
			this.y = min(this.y+1, height-this.width)
		}

		this.hasChanged = true
	}
}

func (this *frame) Resize(width, height int) {
	// x, y position
	this.x = 0
	this.y = 0

	// update size
	this.buffer = make([]string, height)
	this.width = width
	this.height = height

	this.emptyLine = "\x1B[0m" + strings.Repeat(" ", width)

	this.hasChanged = true
}

func (this *frame) SetImage(image filters.Ascii) {
	this.x = 0
	this.y = 0
	this.src = image

	this.hasChanged = true
}

func (this *frame) View() []string {
	if this.hasChanged {

		if this.src.Height_() >= this.height {
			if this.src.Width_() >= this.width {
				//  case 0
				//
				// 	 ┌───────┐
				// 	 │┌ ─ ─ ┐│
				// 	 │|     |│
				// 	 │|     |│
				// 	 │└ ─ ─ ┘│
				// 	 └───────┘
				//
				this.buffer = this.src.Get(this.x, this.y, this.width, this.height)

			} else {
				//  case 1
				//
				//   ┌───────┐
				// ┌ ─ ─ ─ ─ ─ ┐
				// | │       │ |
				// | │       │ |
				// └ ─ ─ ─ ─ ─ ┘
				//   └───────┘
				//
				image := this.src.Get(0, this.y, this.src.Width_(), this.height)

				margin := float64(this.width-this.src.Width_()) / 2.
				leftMargin := this.emptyLine[0 : len("\x1B[0m")+int(math.Floor(margin))]
				rightMargin := this.emptyLine[0 : len("\x1B[0m")+int(math.Ceil(margin))]

				for i := range this.buffer {
					this.buffer[i] = leftMargin + image[i] + rightMargin
				}
			}
		} else {
			if this.src.Width_() > this.width {
				//  case 2
				//    ┌ ─ ─ ┐
				//   ┌|─────|┐
				//   │|     |│
				//   │|     |│
				//   │|     |│
				//   │|     |│
				//   └|─────|┘
				//    └ ─ ─ ┘
				image := this.src.Get(this.x, 0, this.width, this.src.Height_())

				margin := float64(this.height-this.src.Height_()) / 2.
				topMargin := int(math.Floor(margin))
				bottomMargin := int(math.Ceil(margin))

				for i := range topMargin {
					this.buffer[i] = this.emptyLine
				}
				copy(this.buffer[topMargin:], image)
				for i := range bottomMargin {
					this.buffer[this.height-bottomMargin+i] = this.emptyLine
				}
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
				image := this.src.Get(0, 0, this.src.Width_(), this.src.Height_())

				margin := float64(this.height-this.src.Height_()) / 2.
				topMargin := int(math.Floor(margin))
				bottomMargin := int(math.Ceil(margin))

				margin = float64(this.width-this.src.Width_()) / 2.
				leftMargin := this.emptyLine[0 : len("\x1B[0m")+int(math.Floor(margin))]
				rightMargin := this.emptyLine[0 : len("\x1B[0m")+int(math.Ceil(margin))]

				for i := range topMargin {
					this.buffer[i] = this.emptyLine
				}
				for i := range image {
					this.buffer[i+topMargin] = leftMargin + image[i] + rightMargin
				}
				for i := range bottomMargin {
					this.buffer[topMargin+len(image)+i] = this.emptyLine
				}
			}
		}
	}

	this.hasChanged = false
	return this.buffer // should always stay up to date
}
