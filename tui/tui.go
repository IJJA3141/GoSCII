package tui

import (
	"log"
	"strings"

	"github.com/IJJA3141/GoSCII/filters"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	frame frame
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		default:
			m.frame.Update(msg)
		}

	case tea.WindowSizeMsg:
		m.frame.Resize(msg.Width, msg.Height)
	}

	return m, nil
}

func (m model) View() string {
	return strings.Join(m.frame.View(), "\n")
}

func Start(tmp filters.Ascii) {
	p := tea.NewProgram(model{
		frame: Frame(0, 0, tmp),
	}, tea.WithAltScreen())

	tea.WindowSize()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
