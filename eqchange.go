package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
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
	{"sqrt", "√"},
	{"pi", "π"},
}

// PrepareExpr prepares an expression by looping both equation change slices
func (ex *Expr) PrepareExpr(functionsArg map[string]govaluate.ExpressionFunction) (string, map[string]govaluate.ExpressionFunction) {
	functions := make(map[string]govaluate.ExpressionFunction)
	for name, function := range functionsArg {
		functions[name] = function
	}
	ex.LoopEquationChangeSlice()
	params := []string{"π", "e", "x", "a", "t", "h"}
	ex.Expr = strings.ReplaceAll(ex.Expr, "true", "(0==0)") // prevent true and false from being interpreted as functions
	ex.Expr = strings.ReplaceAll(ex.Expr, "false", "(1==0)")
	for fname := range functions { // if there is a function name and no parentheses after, put parentheses around the next character
		for _, pname := range params {
			ex.Expr = strings.ReplaceAll(ex.Expr, fname+pname, fname+"("+pname+")")
		}
		for n := 0; n < 10; n++ {
			ns := strconv.Itoa(n)
			ex.Expr = strings.ReplaceAll(ex.Expr, fname+ns, fname+"("+ns+")")
		}
	}
	expr := LoopUnreadableChangeSlice(ex.Expr)
	ex.Expr = strings.ReplaceAll(ex.Expr, "(0==0)", "true")
	ex.Expr = strings.ReplaceAll(ex.Expr, "(1==0)", "false")
	for fname := range functions { // do again so things like sqrt and f'(x) work
		for _, pname := range params {
			expr = strings.ReplaceAll(expr, fname+pname, fname+"("+pname+")")
		}
		for n := 0; n < 10; n++ {
			ns := strconv.Itoa(n)
			expr = strings.ReplaceAll(expr, fname+ns, fname+"("+ns+")")
		}
	}
	i := 0
	functionsToDelete := []string{}
	functionsToAdd := make(map[string]govaluate.ExpressionFunction)
	for name, function := range functions { // to prevent issues with the equation, all functions are turned into zfunctionindexz. z is just a letter that isn't used in anything else.
		newName := fmt.Sprintf("z%vz", i)
		expr = strings.ReplaceAll(expr, name+"(", newName+"(")
		functionsToAdd[newName] = function
		functionsToDelete = append(functionsToDelete, name)
		i++
	}
	for name, function := range functionsToAdd {
		functions[name] = function
	}
	for _, name := range functionsToDelete {
		delete(functions, name)
	}
	for n := 0; n < 10; n++ { // if the expression contains a number and then a parameter or a function right after, then change it to multiply the number and the parameter/function
		ns := strconv.Itoa(n)
		for _, pname := range params {
			expr = strings.ReplaceAll(expr, ns+pname, ns+"*"+pname)
			expr = strings.ReplaceAll(expr, ns+"("+pname, ns+"*("+pname)
		}
		for fname := range functions {
			expr = strings.ReplaceAll(expr, ns+fname, ns+"*"+fname)
			expr = strings.ReplaceAll(expr, ns+"("+fname, ns+"*("+fname)
		}
	}
	for _, pname := range params { // if the expression contains a parameter before another parameter or a function, make it multiply
		for _, pname1 := range params {
			for strings.Contains(expr, pname+pname1) {
				expr = strings.ReplaceAll(expr, pname+pname1, pname+"*"+pname1)
				expr = strings.ReplaceAll(expr, pname+"("+pname1, pname+"*("+pname1)
			}
		}
		for fname := range functions {
			expr = strings.ReplaceAll(expr, pname+fname, pname+"*"+fname)
			expr = strings.ReplaceAll(expr, pname+"("+fname, pname+"*("+fname)
		}
	}

	return expr, functions
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
