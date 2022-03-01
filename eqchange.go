package main

import (
	"strconv"
	"strings"
)

// EquationChange type has the string that needs to be replaced and what to replace it with
type EquationChange struct {
	Old string
	New string
}

// UnreadableChangeSlice is all of the strings that should change before compiling, but the user shouldn't see
var UnreadableChangeSlice = []EquationChange{
	{"''", "dd"},
	{`"`, "dd"},
	{"'", "d"},
	{"^", "**"},
	{"√", "sqrt"},
}

// EquationChangeSlice is all of the strings that should be changed
var EquationChangeSlice = []EquationChange{
	{"**", "^"},
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
	{" == ", "=="},
	{" && ", "&&"},
	{" || ", "||"},
	{" > ", ">"},
	{" < ", "<"},
	{">", " > "},
	{"<", " < "},
	{"==", " == "},
	{"&&", " && "},
	{"||", " || "},
	{"sqrt*(", "sqrt("},
	{"cot*(", "cot("},
	{"t*an(", "tan("},
	{"a*tan(", "atan("},
	{"fact*(", "fact("},
	{"sqrt", "√"},
}

// InitEquationChangeSlice adds things that involve numbers to the EquationChangeSlice
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

// LoopEquationChangeSlice loops over the Equation Change slice and makes the replacements
func (ex *Expr) LoopEquationChangeSlice() {
	for _, d := range EquationChangeSlice {
		ex.Expr = strings.ReplaceAll(ex.Expr, d.Old, d.New)
	}
}

// LoopUnreadableChangeSlice loops over the unreadable Change slice and makes the replacements
func LoopUnreadableChangeSlice(expr string) string {
	for _, d := range UnreadableChangeSlice {
		expr = strings.ReplaceAll(expr, d.Old, d.New)
	}
	return expr
}

// func (ln *Line) CheckForDerivatives() {
// 	re := regexp.MustCompile(`\[(.*?)\]`)
// 	strs := strings.SplitAfter(ln.Expr.Expr, "]")
// 	var results []string
// 	for _, d := range strs {
// 		submatchall := re.FindAllString(d, -1)
// 		result := d
// 		for _, element := range submatchall {
// 			element = strings.ReplaceAll(element, "[", "")
// 			element = strings.ReplaceAll(element, "]", "")
// 			// ln.Expr.Expr = strings.ReplaceAll(ln.Expr.Expr, element, fmt.Sprintf("(%v, %v)", element, strings.ReplaceAll(element, "x", "x+0.001")))
// 			// fmt.Println(ln.Expr.Expr)
// 			// ln.Expr.Expr = strings.ReplaceAll(ln.Expr.Expr, "[", "")
// 			// ln.Expr.Expr = strings.ReplaceAll(ln.Expr.Expr, "]", "")
// 			// fmt.Println(ln.Expr.Expr)

// 			result = re.ReplaceAllString(d, fmt.Sprintf("(%v, %v)", element, strings.ReplaceAll(element, "x", "(x+0.001)")))
// 		}
// 		results = append(results, result)
// 	}
// 	ln.Expr.Expr = strings.Join(results, "")

// }
