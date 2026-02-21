package tui

import "strings"

type Bar struct{}

func NewBar(width, height int, x, y int) Bar { return Bar{} }
func (br *Bar) Draw(b *strings.Builder)
func (br *Bar) Cursor(b *strings.Builder)
func (br *Bar) Resize(width, height int)
func (br *Bar) HandleKey(key string) bool
