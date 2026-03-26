package tui

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const ( // rng for now
	BAR_MIN_WIDTH  = 30
	BAR_MIN_HEIGHT = 20
)

type state struct {
	ratio         bool
	width, height int
}

type Bar struct {
	focus  int
	coords Coords

	inputs []Widget
	state  state
}

func (br *Bar) SetCoord(coord Coords) { br.coords = coord }

func s(s string) bool { return ASCII_PRINTABLE_LOWER_BOUND <= s && s < ASCII_PRINTABLE_UPPER_BOUND }
func IsVisible(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// func f(s string) bool { return "0" <= s && s <= "9" || s == "." }
// func i(s string) bool { return "0" <= s && s <= "9" }

func NewBar(coord Coords) Bar {
	var state state

	widthInput := InputField[int]{
		coords: Coords{X: 0, Y: 2},
		label:  "W ",

		width: 10, minWidth: 10,
		format: func(i int) string { return fmt.Sprint(i) },
		parse:  strconv.Atoi,
		accept: s,

		submit: func(width int) int {
			if state.ratio {
				state.height = int(float64(state.height) / float64(state.width) * float64(width))
				state.width = width
				return width

			} else {
				state.width = width
				return width
			}
		},
	}

	heightInput := InputField[int]{
		width: 10, minWidth: 10,
		coords: Coords{X: 0, Y: 4},
		label:  "H ",

		format: func(i int) string { return fmt.Sprint(i) },
		parse:  strconv.Atoi,
		accept: s,

		submit: func(i int) int { return max(-8000, min(i, 0)) },
	}

	ratioCBX := CheckBox{
		checked: false,
		icons:   [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coords{X: 0, Y: 6},
		label:  "Ratio ",

		submit: func(b bool) bool { return b },
	}

	return Bar{
		focus:  0,
		coords: coord,
		inputs: []Widget{&widthInput, &heightInput, &ratioCBX},
	}
}

func (br *Bar) Render(b *strings.Builder) {
	for i, input := range br.inputs {
		if i != br.focus {
			input.Render(b, br.coords)
		}
	}

	br.inputs[br.focus].Render(b, br.coords)
}

func (br *Bar) Cursor(b *strings.Builder) {
	br.inputs[br.focus].Cursor(b, br.coords)
}

func (br *Bar) Resize(width, height int) (err error) {
	if BAR_MIN_WIDTH < width || BAR_MIN_HEIGHT < height {
		errors.Join(err, fmt.Errorf("w = %d and h = %d are to small for bar", width, height))
	}

	for _, child := range br.inputs {
		errors.Join(err, child.Resize(width, 1))
	}

	return
}

func (br *Bar) HandleKey(key string) bool {
	handled := br.inputs[br.focus].HandleKey(key)

	if !handled {
		handled = true

		switch key {
		case KEY_TAB:
			br.inputs[br.focus].Blur()
			br.focus = increase(br.focus, len(br.inputs)-1)
			br.inputs[br.focus].Focus()

		case KEY_SHIFT_TAB:
			br.inputs[br.focus].Blur()
			br.focus = decrease(br.focus, len(br.inputs)-1)
			br.inputs[br.focus].Focus()

		default:
			handled = false
		}
	}

	if handled {
		log.Printf("[bar] %s\t%x\n", key, key)
	}

	return handled
}
