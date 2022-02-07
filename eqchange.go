package main

import (
	"strconv"
	"strings"
)

// Equation Change type has the string that needs to be replaced and what to replace it with
type EquationChange struct {
	Old string
	New string
}

// Equation Change Slice is all of the strings that should be changed
var EquationChangeSlice = []EquationChange{
	{"^", "**"},
	{"x(", "x*("},
	{"t(", "t*("},
	{"a(", "a*("},
	{"ax", "a*x"},
	{"xa", "x*a"},
	{"tx", "t*x"},
	{"xt", "x*t"},
	{"at", "a*t"},
	{"ta", "t*a"},
	{"aa", "a*a"},
	{"tt", "t*t"},
	{"xx", "x*x"},
	{"h(", "h*("},
	{"hx", "h*x"},
	{"xh", "x*h"},
	{"ah", "a*h"},
	{"ha", "h*a"},
	{"th", "t*h"},
	{"ht", "h*t"},
	{"+.", "+0."},
	{"-.", "-0."},
	{"*.", "*0."},
	{"/.", "/0."},
	{")(", ")*("},
	{"sqrt*(", "sqrt("},
	{"cot*(", "cot("},
	{"t*an(", "tan("},
	{"a*tan(", "atan("},
	{"fact*(", "fact("},
}

// Init Equation Change Slice adds things that involve numbers to the EquationChangeSlice
func InitEquationChangeSlice() {
	for i := 0; i < 10; i++ {
		is := strconv.Itoa(i)
		EquationChangeSlice = append(EquationChangeSlice,
			EquationChange{is + "(", is + "*("},
			EquationChange{is + "x", is + "*x"},
			EquationChange{is + "t", is + "*t"},
			EquationChange{is + "a", is + "*a"},
			EquationChange{is + "h", is + "*h"},
		)
	}
}

// Replace Equation Change Slice loops over the Equation Change slice and makes the replacements
func (ln *Line) LoopEquationChangeSlice() {
	for _, d := range EquationChangeSlice {
		ln.Expr.Expr = strings.Replace(ln.Expr.Expr, d.Old, d.New, -1)
	}
}
