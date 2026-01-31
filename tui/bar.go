package tui

import (
	"fmt"
	"io"
	"strconv"
)

const (
	BAR_WIDTH  = 10
	BAR_MARGIN = "   "
)

type Bar struct {
	out io.Writer

	width, height   int
	swidth, sheight string
	lock            bool
	fIndex          int
}

func (this *Bar) HandleKey(key string) {
	_, err := strconv.Atoi(key)
	if err == nil {
		this.swidth += key
	} else if key == string([]byte{0x0, 0xB}) { // not shure of that one // tab
		
	}
}

func (this *Bar) Draw(x, y int) string {
	str := "\x1b[" + fmt.Sprint(y) + ";" + fmt.Sprint(x) + "H"
	str += "Size:"
	str += "\x1b[" + fmt.Sprint(y+1) + ";" + fmt.Sprint(x) + "H"
	str += this.swidth
	return str
}
