package tui

import "strings"

type Frame struct{}

func NewFrame(a, b int, x, y int) Frame { return Frame{} }
func (f *Frame) Draw(b *strings.Builder)
func (f *Frame) Cursor(b *strings.Builder)
func (f *Frame) Resize(width, height int)
func (f *Frame) HandleKey(key string) bool
