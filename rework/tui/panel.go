package tui

import (
	"errors"
	"fmt"
)

const PANEL_MIN_WIDTH = 25

type Panel struct {
	focus int
	width int

	relative Coords

	elements []Element
}

//=================
//	constructor

func NewPanel(relative Coords) Panel {
	width := InputField[int]{}
	height := InputField[int]{}
	ratio := CheckBox{}

	width.submit = func(i int) {
		if ratio.checked {
			r := float64(height.value) / float64(width.value)
			r *= float64(i)
			height.SetValue(int(r))
		}

		width.SetValue(i)
	}

	height.submit = func(i int) {
		if ratio.checked {
			r := float64(width.value) / float64(height.value)
			r *= float64(i)
			width.SetValue(int(r))
		}

		width.SetValue(i)
	}

	return Panel{
		focus:    0,
		relative: relative,
		elements: []Element{&width, &height, &ratio},
	}
}

//============
//  Layout

// not shure about that...
func (pnl *Panel) Resize(width, height int, relative Coords) (err error) {
	if height < len(pnl.elements) {
		return fmt.Errorf("Failed to resize Panel, height %d is too small.", height)
	}

	if width != PANEL_MIN_WIDTH {
		return fmt.Errorf("Tryed to change width of panel with constant width.")
	}

	for _, element := range pnl.elements {
		err = errors.Join(element.Resize(PANEL_MIN_WIDTH, 1))
	}

	if err == nil {
		pnl.relative = relative
	}

	return err
}

func (pnl *Panel) Render(sb *ScreenBuffer) {
	// TODO opti to only render modified elements
	for _, element := range pnl.elements {
		element.Render(sb, pnl.relative)
	}
}

func (pnl *Panel) Cursor(sb *ScreenBuffer) {
	pnl.elements[pnl.focus].Cursor(sb, pnl.relative)
}

//=================
//  Interaction

func (br *Panel) HandleKey(key string) bool {
	handled := br.elements[br.focus].HandleKey(key)

	if !handled {
		handled = true

		switch key {
		case KEY_TAB:
			br.elements[br.focus].Blur()
			br.focus = increase(br.focus, len(br.elements)-1)
			br.elements[br.focus].Focus()

		case KEY_SHIFT_TAB:
			br.elements[br.focus].Blur()
			br.focus = decrease(br.focus, len(br.elements)-1)
			br.elements[br.focus].Focus()

		default:
			handled = false
		}
	}

	// if handled {
	// 	log.Printf("[bar] %s\t%x\n", key, key)
	// }

	return handled
}
