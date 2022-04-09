package main

import (
	"fmt"
	"math"

	"github.com/Knetic/govaluate"
	"gonum.org/v1/gonum/integrate"
)

// Expr is an expression
type Expr struct {
	Expr   string                         `label:"" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	Val    *govaluate.EvaluableExpression `view:"-" json:"-"`
	Params map[string]interface{}         `view:"-" json:"-"`
}

// factorial variables
const lim = 100

var randNum float64

var facts [lim]float64

// DefaultFunctions that can be used in expressions
var DefaultFunctions = map[string]govaluate.ExpressionFunction{
	"cos": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("cos", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Cos(args[0].(float64))
		return y, nil
	},
	"sin": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("sin", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Sin(args[0].(float64))
		return y, nil
	},
	"tan": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("tan", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Tan(args[0].(float64))
		return y, nil
	},
	"pow": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("pow", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		y := math.Pow(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"abs": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("abs", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Abs(args[0].(float64))
		return y, nil
	},
	"fact": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("fact", args, "float64")
		if err != nil {
			return 0, err
		}
		y := FactorialMemoization(int(args[0].(float64)))
		return y, nil
	},
	"ceil": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("ceil", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Ceil(args[0].(float64))
		return y, nil
	},
	"floor": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("floor", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Floor(args[0].(float64))
		return y, nil
	},
	"mod": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("mod", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		y := math.Mod(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"rand": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("rand", args, "float64")
		if err != nil {
			return 0, err
		}
		y := randNum * args[0].(float64)
		return y, nil
	},
	"sqrt": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("sqrt", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Sqrt(args[0].(float64))
		return y, nil
	},
	"ln": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("ln", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Log(args[0].(float64))
		return y, nil
	},
	"csc": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("csc", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Sin(args[0].(float64))
		return y, nil
	},
	"sec": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("sec", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Cos(args[0].(float64))
		return y, nil
	},
	"cot": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("cot", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Tan(args[0].(float64))
		return y, nil
	},
	"if": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("if", args, "bool", "float64", "float64")
		if err != nil {
			return 0, err
		}
		if args[0].(bool) {
			return args[1].(float64), nil
		}
		return args[2].(float64), nil
	},
	"arcsin": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arcsin", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Asin(args[0].(float64))
		return y, nil
	},
	"arccos": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccos", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Acos(args[0].(float64))
		return y, nil
	},
	"arctan": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arctan", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Atan(args[0].(float64))
		return y, nil
	},
	"sinh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("sinh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Sinh(args[0].(float64))
		return y, nil
	},
	"cosh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("cosh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Cosh(args[0].(float64))
		return y, nil
	},
	"tanh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("tanh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Tanh(args[0].(float64))
		return y, nil
	},
	"arcsinh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arcsinh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Asinh(args[0].(float64))
		return y, nil
	},
	"arccosh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccosh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Acosh(args[0].(float64))
		return y, nil
	},
	"arctanh": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arctanh", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Atanh(args[0].(float64))
		return y, nil
	},
	"arcsec": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arcsec", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Acos(1 / args[0].(float64))
		return y, nil
	},
	"arccsc": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccsc", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Asin(1 / args[0].(float64))
		return y, nil
	},
	"arccot": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccot", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Atan(1 / args[0].(float64))
		if args[0].(float64) < 0 {
			y += math.Pi
		}

		return y, nil
	},
	"arcsech": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arcsech", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Acosh(1 / args[0].(float64))
		return y, nil
	},
	"arccsch": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccsch", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Asinh(1 / args[0].(float64))
		return y, nil
	},
	"arccoth": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("arccoth", args, "float64")
		if err != nil {
			return 0, err
		}
		y := math.Atanh(1 / args[0].(float64))
		return y, nil
	},
}

// Integrate returns the integral of an expression
func (ex *Expr) Integrate(min, max float64, h int) float64 {
	var vals []float64
	sign := float64(1)
	diff := max - min
	if diff == 0 {
		return 0
	}
	if diff < 0 {
		diff = -diff
		sign = -1
		min, max = max, min
	}
	accuracy := 16
	dx := diff / float64(accuracy)
	for x := min; x <= max; x += dx {
		vals = append(vals, ex.Eval(x, TheGraph.State.Time, h))
	}
	if len(vals) != accuracy+1 {
		vals = append(vals, ex.Eval(max, TheGraph.State.Time, h))
	}
	val := integrate.Romberg(vals, dx)
	return sign * val
}

// Deriv takes the derivative given value 1 and two, and the difference in x between them
// func Deriv(val1, val2, inc float64) float64 {
// 	return (val2 - val1) / inc
// }

// CheckArgs checks if a function is passed the right number of arguments, and the right type of arguments.
func CheckArgs(name string, have []interface{}, want ...string) error {
	if len(have) != len(want) {
		return fmt.Errorf("function %v needs %v arguments, not %v arguments", name, len(want), len(have))
	}
	for i, d := range want {
		if d != fmt.Sprintf("%T", have[i]) {
			return fmt.Errorf("function %v needs %v. %v does not work", name, want, have)
		}
	}
	return nil
}

// Compile gets an expression ready for evaluation.
func (ex *Expr) Compile() error {
	expr, functions := ex.PrepareExpr(DefaultFunctions)
	var err error
	ex.Val, err = govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if HandleError(err) {
		ex.Val = nil
		return err
	}
	if ex.Params == nil {
		ex.Params = make(map[string]interface{}, 2)
	}
	ex.Params["π"] = math.Pi
	ex.Params["e"] = math.E
	return err
}

// Eval gives the y value of the function for given x, t and h value
func (ex *Expr) Eval(x, t float64, h int) float64 {
	if ex.Expr == "" {
		return 0
	}
	currentX = x
	ex.Params["x"] = x
	ex.Params["t"] = t
	ex.Params["a"] = 10 * math.Sin(t)
	ex.Params["h"] = h
	yi, err := ex.Val.Evaluate(ex.Params)
	if HandleError(err) {
		return 0
	}
	switch yi.(type) {
	case float64:
		return yi.(float64)
	default:
		TheGraph.Stop()
		HandleError(fmt.Errorf("expression %v is invalid, it is a %T value, should be a float64 value", ex.Expr, yi))
		return 0
	}
}

// EvalWithY calls eval but with a y value set
func (ex *Expr) EvalWithY(x, t float64, h int, y float64) float64 {
	ex.Params["y"] = y
	return ex.Eval(x, t, h)
}

// EvalBool checks if a statement is true based on the x, y, t and h values
func (ex *Expr) EvalBool(x, y, t float64, h int) bool {
	if ex.Expr == "" {
		return true
	}
	currentX = x
	ex.Params["x"] = x
	ex.Params["t"] = t
	ex.Params["a"] = 10 * math.Sin(t)
	ex.Params["h"] = h
	ex.Params["y"] = y
	ri, err := ex.Val.Evaluate(ex.Params)
	if HandleError(err) {
		return true
	}
	switch ri.(type) {
	case bool:
		return ri.(bool)
	default:
		TheGraph.Stop()
		HandleError(fmt.Errorf("expression %v is invalid, it is a %T value, should be a bool value", ex.Expr, ri))
		return false
	}
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
