package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Knetic/govaluate"
)

// Expression
type Expr struct {
	Expr   string                         `label:"y=" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	Val    *govaluate.EvaluableExpression `view:"-" json:"-"`
	Params map[string]interface{}         `view:"-" json:"-"`
}

const LIM = 100

var facts [LIM]float64

var functions = map[string]govaluate.ExpressionFunction{
	"cos": func(args ...interface{}) (interface{}, error) {
		y := math.Cos(args[0].(float64))
		return y, nil
	},
	"sin": func(args ...interface{}) (interface{}, error) {
		y := math.Sin(args[0].(float64))
		return y, nil
	},
	"tan": func(args ...interface{}) (interface{}, error) {
		y := math.Tan(args[0].(float64))
		return y, nil
	},
	"pow": func(args ...interface{}) (interface{}, error) {
		y := math.Pow(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"abs": func(args ...interface{}) (interface{}, error) {
		y := math.Abs(args[0].(float64))
		return y, nil
	},
	"fact": func(args ...interface{}) (interface{}, error) {
		y := FactorialMemoization(int(args[0].(float64)))
		return y, nil
	},
	"ceil": func(args ...interface{}) (interface{}, error) {
		y := math.Ceil(args[0].(float64))
		return y, nil
	},
	"floor": func(args ...interface{}) (interface{}, error) {
		y := math.Floor(args[0].(float64))
		return y, nil
	},
	"mod": func(args ...interface{}) (interface{}, error) {
		y := math.Mod(args[0].(float64), args[1].(float64))
		return y, nil
	},
	"rand": func(args ...interface{}) (interface{}, error) {
		y := float64(rand.Float32()) * args[0].(float64)
		return y, nil
	},
	"sqrt": func(args ...interface{}) (interface{}, error) {
		y := math.Sqrt(args[0].(float64))
		return y, nil
	},
	"ln": func(args ...interface{}) (interface{}, error) {
		y := math.Log(args[0].(float64))
		return y, nil
	},
	"csc": func(args ...interface{}) (interface{}, error) {
		y := 1 / math.Sin(args[0].(float64))
		return y, nil
	},
	"sec": func(args ...interface{}) (interface{}, error) {
		y := 1 / math.Cos(args[0].(float64))
		return y, nil
	},
	"cot": func(args ...interface{}) (interface{}, error) {
		y := 1 / math.Tan(args[0].(float64))
		return y, nil
	},
	"asin": func(args ...interface{}) (interface{}, error) {
		y := math.Asin(args[0].(float64))
		return y, nil
	},
	"acos": func(args ...interface{}) (interface{}, error) {
		y := math.Acos(args[0].(float64))
		return y, nil
	},
	"atan": func(args ...interface{}) (interface{}, error) {
		y := math.Atan(args[0].(float64))
		return y, nil
	},
	"ifb": func(args ...interface{}) (interface{}, error) {
		if (args[0].(float64) > args[1].(float64)) && (args[0].(float64) < args[2].(float64)) {
			return args[3].(float64), nil
		}
		return args[4].(float64), nil
	},
	"ife": func(args ...interface{}) (interface{}, error) {
		if args[0].(float64) == args[1].(float64) {
			return args[2].(float64), nil
		}
		return args[3].(float64), nil
	},
}

func (ex *Expr) Compile() error {
	var err error
	ex.Val, err = govaluate.NewEvaluableExpressionWithFunctions(ex.Expr, functions)
	if err != nil {
		ex.Val = nil
		fmt.Printf("Error: %v \n", err)
	}
	if ex.Params == nil {
		ex.Params = make(map[string]interface{}, 2)
	}
	return err
}

// Eval gives the y value of the function for given x, t and h value
func (ex *Expr) Eval(x, t float32, h int) float32 {
	if ex.Expr == "" {
		return 0
	}
	ex.Params["x"] = float64(x)
	ex.Params["t"] = float64(t)
	ex.Params["a"] = float64(10 * math.Sin(float64(t)))
	ex.Params["h"] = float64(h)
	yi, _ := ex.Val.Evaluate(ex.Params)
	y := float32(yi.(float64))
	return y
}

// Used to take the factorial for the fact() function
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
