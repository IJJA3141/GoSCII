package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type editor struct {
	width, height int
}

func Editor() editor {
	return editor{}
}

func (this *editor) View() []string {
	out := make([]string, this.height)
	empty := strings.Repeat(" ", this.width)

	for i := range out {
		out[i] = empty
	}

	return out
}

func (this *editor) Update(msg tea.KeyMsg) {
}

func (this *editor) Height(height int) {
	this.height = height
}

func (this *editor) Width(width int) {
	this.width = width
}

func (this *editor) AddLine() {
}

type command struct {
	width int
	focus bool
	cmd   string
}

func Command() command {
	return command{}
}

func (this *command) Width(width int) {
	this.width = width
}

func (this *command) View() string {
	if this.focus {
		return ":" + this.cmd + strings.Repeat(" ", this.width-1-len(this.cmd))
	}

	return strings.Repeat(" ", this.width)
}

func (this *command) Update(msg tea.KeyMsg) {
	switch msg.String() {
	case tea.KeyBackspace.String():
		this.cmd = this.cmd[:max(0, len(this.cmd)-1)]
	
	default:
		this.cmd += msg.String()
	}
}

func (this *command) Init() {
	this.cmd = ""
	this.focus = true
}

func (this *command) Kill() {
	this.focus = false
}
