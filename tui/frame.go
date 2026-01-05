package tui

import (
	"math"
	"strings"

	"github.com/IJJA3141/GoSCII/filters"
	tea "github.com/charmbracelet/bubbletea"
)

const speed = 2

type Frame struct {
	img           filters.Stampable
	width, height int
	x, y          int
	buffer        []string
	xCrop, yCrop  bool
}

func NewFrame(width, height int) Frame {
	return Frame{
		img:    nil,
		width:  width,
		height: height,
		x:      0,
		y:      0,
		buffer: make([]string, width*height),
		xCrop:  false,
		yCrop:  false,
	}
}

func (frame *Frame) SetImage(img filters.Stampable) {
	frame.img = img
	frame.x = 0
	frame.y = 0
}

func (frame *Frame) Resize(width, height int) {
	frame.width = width
	frame.height = height
	frame.x = 0
	frame.y = 0
	frame.xCrop = width > frame.width
	frame.yCrop = height > frame.height
}

func (frame *Frame) Update(msg tea.Msg) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyLeft.String(), "h":
			frame.x = max(0, frame.x-1*speed)
			return

		case tea.KeyRight.String(), "l":
			width, _, _ := frame.img.Stamp()

			if frame.width >= width {
				frame.x = 0
				return
			}

			frame.x = min(width-frame.width, frame.x+1*speed)
			return

		case tea.KeyUp.String(), "k":
			frame.y = max(0, frame.y-1)
			return

		case tea.KeyDown.String(), "j":
			_, height, _ := frame.img.Stamp()

			if frame.height >= height {
				frame.y = 0
				return
			}

			frame.y = min(height-frame.height, frame.y+1*speed)
			return
		}
	}
}

func (frame *Frame) View() []string {
	width, height, stamp := frame.img.Stamp()

	xCrop := width > frame.width
	yCrop := height > frame.height

	out := make([]string, frame.height)
	emptyLine := strings.Repeat(" ", frame.width)

	if yCrop {
		start := (height - frame.height + frame.y) / 2

		for j := range out {
			if xCrop {
				x := (width - frame.width + frame.x) / 2
				out[j] = strings.Join(stamp[j+start][x:x+frame.width], "")
			} else {
				diff := float64(frame.width-width) / 2.
				leftPadding := emptyLine[0:int(math.Floor(diff))]
				rightPadding := emptyLine[0:int(math.Ceil(diff))]

				out[j] = leftPadding
				out[j] += strings.Join(stamp[j+start], "")
				out[j] += rightPadding
			}
		}

	} else {
		diff := float64(frame.height-height) / 2.
		top := int(math.Floor(diff))

		for j := range top {
			out[j] = emptyLine
		}

		for j := range height {
			if xCrop {
				x := (width - frame.width) / 2
				out[j] = strings.Join(stamp[j][x:x+frame.width], "")
			} else {
				diff := float64(frame.width-width) / 2.
				leftPadding := emptyLine[0:int(math.Floor(diff))]
				rightPadding := emptyLine[0:int(math.Ceil(diff))]

				out[j] = leftPadding
				out[j] += strings.Join(stamp[j], "")
				out[j] += rightPadding
			}
		}

		bottom := int(math.Ceil(diff))
		for j := range bottom {
			out[top+height+j] = emptyLine
		}
	}

	return out
}
