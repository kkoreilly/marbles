package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Knetic/govaluate"
)

// Expr is an expression
type Expr struct {
	Expr   string                         `label:"y=" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	Val    *govaluate.EvaluableExpression `view:"-" json:"-"`
	Params map[string]interface{}         `view:"-" json:"-"`
}

// factorial variables
const lim = 100

var facts [lim]float64

// Functions that can be used in expressions
var functions = map[string]govaluate.ExpressionFunction{
	"cos": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "cos")
		if !ok {
			return float64(0), err
		}
		y := math.Cos(args[0].(float64))
		return y, nil
	},
	"sin": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "sin")
		if !ok {
			return float64(0), err
		}
		y := math.Sin(args[0].(float64))
		return y, nil
	},
	"tan": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "tan")
		if !ok {
			return float64(0), err
		}
		y := math.Tan(args[0].(float64))
		return y, nil
	},
	"pow": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(2, len(args), "pow")
		if !ok {
			return float64(0), err
		}
		y := math.Pow(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"abs": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "abs")
		if !ok {
			return float64(0), err
		}
		y := math.Abs(args[0].(float64))
		return y, nil
	},
	"fact": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "fact")
		if !ok {
			return float64(0), err
		}
		y := FactorialMemoization(int(args[0].(float64)))
		return y, nil
	},
	"ceil": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "ceil")
		if !ok {
			return float64(0), err
		}
		y := math.Ceil(args[0].(float64))
		return y, nil
	},
	"floor": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "floor")
		if !ok {
			return float64(0), err
		}
		y := math.Floor(args[0].(float64))
		return y, nil
	},
	"mod": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(2, len(args), "mod")
		if !ok {
			return float64(0), err
		}
		y := math.Mod(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"rand": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "rand")
		if !ok {
			return float64(0), err
		}
		y := float64(rand.Float32()) * args[0].(float64)
		return y, nil
	},
	"sqrt": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "sqrt")
		if !ok {
			return float64(0), err
		}
		y := math.Sqrt(args[0].(float64))
		return y, nil
	},
	"ln": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "ln")
		if !ok {
			return float64(0), err
		}
		y := math.Log(args[0].(float64))
		return y, nil
	},
	"csc": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "csc")
		if !ok {
			return float64(0), err
		}
		y := 1 / math.Sin(args[0].(float64))
		return y, nil
	},
	"sec": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "sec")
		if !ok {
			return float64(0), err
		}
		y := 1 / math.Cos(args[0].(float64))
		return y, nil
	},
	"cot": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "cot")
		if !ok {
			return float64(0), err
		}
		y := 1 / math.Tan(args[0].(float64))
		return y, nil
	},
	"asin": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "asin")
		if !ok {
			return float64(0), err
		}
		y := math.Asin(args[0].(float64))
		return y, nil
	},
	"acos": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "acos")
		if !ok {
			return float64(0), err
		}
		y := math.Acos(args[0].(float64))
		return y, nil
	},
	"atan": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "atan")
		if !ok {
			return float64(0), err
		}
		y := math.Atan(args[0].(float64))
		return y, nil
	},
	"ifb": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(5, len(args), "ifb")
		if !ok {
			return float64(0), err
		}
		if (args[0].(float64) > args[1].(float64)) && (args[0].(float64) < args[2].(float64)) {
			return args[3].(float64), nil
		}
		return args[4].(float64), nil
	},
	"ife": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(4, len(args), "ife")
		if !ok {
			return float64(0), err
		}
		if args[0].(float64) == args[1].(float64) {
			return args[2].(float64), nil
		}
		return args[3].(float64), nil
	},
	"d": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "d")
		if !ok {
			return float64(0), err
		}
		inc := 0.001
		ln := Gr.Lines[int(args[0].(float64))]
		val1 := float64(ln.Expr.Eval(currentX, Gr.Params.Time, ln.TimesHit))
		val2 := float64(ln.Expr.Eval(currentX+float32(inc), Gr.Params.Time, ln.TimesHit))
		return Deriv(val1, val2, inc), nil
	},
	"f": func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), "d")
		if !ok {
			return float64(0), err
		}
		ln := Gr.Lines[int(args[0].(float64))]
		val := float64(ln.Expr.Eval(currentX, Gr.Params.Time, ln.TimesHit))
		return val, nil
	},
}

// Deriv takes the derivative given value 1 and two, and the difference in x between them
func Deriv(val1, val2, inc float64) float64 {
	return (val2 - val1) / inc
}

// CheckArgs checks if a function is passed the right number of arguments.
func CheckArgs(needed, have int, name string) (bool, error) {
	if needed != have {
		return false, fmt.Errorf("function %v needs %v arguments, not %v arguments", name, needed, have)
	}
	return true, nil
}

// Compile gets an expression ready for evaluation.
func (ex *Expr) Compile() error {
	var err error
	ex.Val, err = govaluate.NewEvaluableExpressionWithFunctions(ex.Expr, functions)
	if HandleError(err) {
		ex.Val = nil
		return err
	}
	if ex.Params == nil {
		ex.Params = make(map[string]interface{}, 2)
	}
	ex.Params["pi"] = math.Pi
	ex.Params["e"] = math.E
	return err
}

// Eval gives the y value of the function for given x, t and h value
func (ex *Expr) Eval(x, t float32, h int) float32 {
	if ex.Expr == "" {
		return 0
	}
	currentX = x
	ex.Params["x"] = float64(x)
	ex.Params["t"] = float64(t)
	ex.Params["a"] = float64(10 * math.Sin(float64(t)))
	ex.Params["h"] = float64(h)
	yi, err := ex.Val.Evaluate(ex.Params)
	if HandleError(err) {
		problemWithEval = true
		return 0
	}
	y := float32(yi.(float64))
	return y
}

// FactorialMemoization is used to take the factorial for the fact() function
func FactorialMemoization(n int) (res float64) {
	if n < 0 {
		return 1
	}
	if facts[n] != 0 {
		res = facts[n]
		return res
	}
	if n > 0 {
		res = float64(n) * FactorialMemoization(n-1)
		return res
	}
	return 1
}
