package tui

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const ( // rng for now
	MIN_WIDTH  = 20
	MIN_HEIGHT = 20
)

type Bar struct {
	focus int
	coord Coord

	inputs []Element
}

func (br *Bar) SetCoord(coord Coord) { br.coord = coord }

func s(s string) bool { return ASCII_PRINTABLE_LOWER_BOUND <= s && s < ASCII_PRINTABLE_UPPER_BOUND }

// func f(s string) bool { return "0" <= s && s <= "9" || s == "." }
// func i(s string) bool { return "0" <= s && s <= "9" }

func NewBar(coord Coord) Bar {

	widthInput := InputField[int]{
		width: 10, minWidth: 0, maxWidth: 30,
		coords: Coord{X: 0, Y: 0},
		label:  "W ",

		parse:  strconv.Atoi,
		accept: s,
		submit: func(i int) string {
			return fmt.Sprint(i)
		},
	}

	heightInput := InputField[int]{
		width: 10, minWidth: 0, maxWidth: 30,
		coords: Coord{X: 15, Y: 0},
		label:  "H ",

		parse:  strconv.Atoi,
		accept: s,
		submit: func(i int) string {
			return fmt.Sprint(i)
		},
	}

	ratioCBX := CheckBox{
		checked: false,
		icon:    [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coord{X: 29, Y: 0},
		label:  "",

		submit: func(b bool) bool { return b },
	}

	a := CheckBox{
		checked: false,
		icon:    [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coord{X: 0, Y: 2},
		label:  "",

		submit: func(b bool) bool { return b },
	}

	b := CheckBox{
		checked: false,
		icon:    [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coord{X: 2, Y: 2},
		label:  "",

		submit: func(b bool) bool { return b },
	}

	c := CheckBox{
		checked: false,
		icon:    [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coord{X: 4, Y: 2},
		label:  "",

		submit: func(b bool) bool { return b },
	}

	d := CheckBox{
		checked: false,
		icon:    [2]string{"\x1b[2m\x1b[22m", "\x1b[1m\x1b[22m"},

		coords: Coord{X: 6, Y: 2},
		label:  "",

		submit: func(b bool) bool { return b },
	}

	return Bar{
		focus:  0,
		coord:  coord,
		inputs: []Element{&widthInput, &heightInput, &ratioCBX, &a, &b, &c, &d},
	}
}

func (br *Bar) Render(b *strings.Builder) {
	for i, input := range br.inputs {
		if i != br.focus {
			input.Render(b, br.coord)
		}
	}

	br.inputs[br.focus].Render(b, br.coord)
}

func (br *Bar) Cursor(b *strings.Builder) {
	br.inputs[br.focus].Cursor(b, br.coord)
}

func (br *Bar) Resize(width, height int) (err error) {
	if MIN_WIDTH < width || MIN_HEIGHT < height {
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
