package tui

import (
	"fmt"
	"io"
	"strings"
)

const (
	BAR_WIDTH = 10
	BAR_MARGIN = "   "
)

type Bar struct {
	out io.Writer

	width, height int
	lock bool
	fIndex int
}

func (this *Bar) Update(key string)

func (this *Bar) Draw(x, y int) {
	str := "\x1b[" + fmt.Sprint(y) + ";" + fmt.Sprint(x) + "H"
	str += "Size:"
	str += "\x1b[" + fmt.Sprint(y+1) + ";" + fmt.Sprint(x) + "H"
	str += fmt.Sprintf("", 


	this.out.Write([]byte(str))
}
