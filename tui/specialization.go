package tui

import (
	"fmt"
	"strconv"
)

func s(s string) bool { return ASCII_PRINTABLE_LOWER_BOUND <= s && s < ASCII_PRINTABLE_UPPER_BOUND }
func f(s string) bool { return "0" <= s && s <= "9" || s == "." }
func i(s string) bool { return "0" <= s && s <= "9" }

// specific input field constructors
func NewTextField(width, x, y int) InputField[string] {
	return InputField[string]{
		Width: width,
		X:     x,
		Y:     y,

		input: "",
		Value: "",

		format:       func(value string) string { return value },
		parse:        func(input string) (string, error) { return input, nil },
		isValidInput: s,
	}
}

func NewIntField(width, x, y int, strict bool) InputField[int] {
	var l func(string) bool
	if strict {
		l = i
	} else {
		l = s
	}

	return InputField[int]{
		Width: width,
		X:     x,
		Y:     y,

		input: "",
		Value: 0,

		format:       func(value int) string { return fmt.Sprint(value) },
		parse:        strconv.Atoi,
		isValidInput: l,
	}
}

func NewFloatField(width, x, y int, strict bool) InputField[float64] {
	var l func(string) bool
	if strict {
		l = f
	} else {
		l = s
	}

	return InputField[float64]{
		Width: width,
		X:     x,
		Y:     y,

		input: "",
		Value: 0,

		format:       func(value float64) string { return fmt.Sprintf("%f", value) },
		parse:        func(input string) (float64, error) { return strconv.ParseFloat(input, 64) },
		isValidInput: l,
	}
}

// round
func NewRoundCheckBox(x, y int, defaultValue bool) CheckBox {
	return CheckBox{
		x: x,
		y: y,

		Value: defaultValue,

		checked:   "󰄯",
		unchecked: "󰄰",
	}
}

// squared
func NewSquaredCheckBox(x, y int, defaultValue bool) CheckBox {
	return CheckBox{
		x: x,
		y: y,

		Value: defaultValue,

		checked:   "󰄮",
		unchecked: "󰄱",
	}
}

// no font
func NewCheckBox(x, y int, defaultValue bool) CheckBox {
	return CheckBox{
		x: x,
		y: y,

		Value: defaultValue,

		checked:   "x",
		unchecked: " ",
	}
}

func NewLockCheckBox(x, y int, defaultValue bool) CheckBox {
	return CheckBox{
		x: x,
		y: y,

		Value: defaultValue,

		checked:   "\x1b[1m\x1b[22m",
		unchecked: "\x1b[2m\x1b[22m",
	}
}
