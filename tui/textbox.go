package tui

import (
	"strings"
)

type TextField struct {
	str   string
	width int
}

func (this *TextField) Draw(x, y int) string {
	if this.width > len(this.str) {
		return MoveTo(x, y) + this.str + strings.Repeat(" ", this.width-len(this.str))
	}

	return MoveTo(x, y) + this.str[len(this.str)-this.width:]
}

func (this *TextField) Cursor(x, y int) string {
	return MoveTo(x, y+min(len(this.str), this.width)) + BLINKING_IBEAM
}

func (this *TextField) HandleKey(key string) {
	if key == "\x7f" { // DEL
		if len(this.str) > 0 {
			this.str = this.str[:len(this.str)-1]
		}
	} else {
		this.str += key
	}
}
