package tui

import (
	"log"
	"strings"

	"github.com/IJJA3141/GoSCII/filters"
	tea "github.com/charmbracelet/bubbletea"
)

type mode = int

const (
	NORMAL = iota
	VISUAL
	INSERT
	COMMAND
)

const (
	BG        = "255;255;255"
	NORMAL_FG = "0;255;0"
	NORMAL_BG = "100;100;100"

	VISUAL_FG = "0;155;0"
	VISUAL_BG = "100;255;0"

	INSERT_FG = "0;255;100"
	INSERT_BG = "0;255;255"

	COMMAND_FG = "0;255;200"
	COMMAND_BG = "200;255;0"
)

type model struct {
	frame   frame
	editor  editor
	command command

	width, height int
	menuWidth     int
	mode          mode

	stack []any
}

func (m model) Init() tea.Cmd {
	return nil
}

func (this *model) Run() tea.Cmd {
	switch this.command.cmd {
	case "q":
		return tea.Quit

	case "invert":
		switch plane := this.stack[len(this.stack)-1].(type) {
		case filters.GrayScalePlane:
			pln := plane.Inverse()
			this.stack = append(this.stack, pln)
			this.frame.SetImage(pln.Braille(255 / 2))
		}
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.Resize(m.width-m.menuWidth, m.height)
		m.editor.Height(m.height - 2)
		m.command.Width(m.menuWidth)
		m.editor.Width(m.menuWidth)

	case tea.KeyMsg:
		switch m.mode {
		case NORMAL:
			switch msg.String() {
			case "v":
				m.mode = VISUAL

			case ":":
				m.command.Init()
				m.mode = COMMAND

			case "i":
				m.mode = INSERT

			default:
				m.editor.Update(msg)
			}

		case INSERT:
			switch msg.String() {
			case tea.KeyEsc.String():
				m.mode = NORMAL

			default:
				m.editor.Update(msg)
			}

		case VISUAL:
			switch msg.String() {
			case ":":
				m.command.Init()
				m.mode = COMMAND

			case tea.KeyEsc.String():
				m.mode = NORMAL

			default:
				m.frame.Update(msg)
			}

		case COMMAND:
			switch msg.String() {
			case tea.KeyEsc.String():
				m.command.Kill()
				m.mode = NORMAL

			case tea.KeyEnter.String():
				m.mode = NORMAL
				return m, m.Run()

			default:
				m.command.Update(msg)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width < 4 || m.height < 4 {
		return ""
	}

	var str strings.Builder
	frame := m.frame.View()
	editor := m.editor.View()

	for i := range m.height - 2 {
		str.WriteString(frame[i] + "\x1b[0m|" + editor[i] + "\n")
	}

	str.WriteString(frame[m.height-2] + "\x1b[0m|")
	switch m.mode {
	case NORMAL:
		str.WriteString("\x1b[38;2;" + NORMAL_FG + "m\x1b[48;2;" + NORMAL_BG + "m NORMAL \x1b[48;2;" + BG + "m" + strings.Repeat(" ", m.menuWidth-9))

	case VISUAL:
		str.WriteString("\x1b[38;2;" + VISUAL_FG + "m\x1b[48;2;" + VISUAL_BG + "m VISUAL \x1b[48;2;" + BG + "m" + strings.Repeat(" ", m.menuWidth-9))

	case INSERT:
		str.WriteString("\x1b[38;2;" + INSERT_FG + "m\x1b[48;2;" + INSERT_BG + "m INSERT \x1b[48;2;" + BG + "m" + strings.Repeat(" ", m.menuWidth-9))

	case COMMAND:
		str.WriteString("\x1b[38;2;" + COMMAND_FG + "m\x1b[48;2;" + COMMAND_BG + "m COMMAND \x1b[48;2;" + BG + "m" + strings.Repeat(" ", m.menuWidth-10))
	}
	str.WriteString("\x1b[0m\n")

	str.WriteString(frame[m.height-1] + "\x1b[0m|" + m.command.View())
	return str.String()
}

func Start(image filters.Ascii) {
	p := tea.NewProgram(model{
		frame:   Frame(0, 0, image),
		editor:  Editor(),
		command: Command(),

		width: 0, height: 0,
		mode:      NORMAL,
		menuWidth: 55,
	}, tea.WithAltScreen())

	tea.WindowSize()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
