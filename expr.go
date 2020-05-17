package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"math"
	"math/rand"
)

type Expr struct {
	Expr   string                         `label:"y=" desc:"equation: use 'x' for the x value, and must use * for multiplication, and start with 0 for decimal numbers (0.01 instead of .01)"`
	Val    *govaluate.EvaluableExpression `view:"-" json:"-"`
	Params map[string]interface{}         `view:"-" json:"-"`
}

func (ex *Expr) Compile() error {
	var err error
	fmt.Printf("expr: %v \n", ex)
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

func (ex *Expr) Eval(x, t float32) float32 {
	// fmt.Printf("Ex: %v \n", ex.Expr)
	if ex.Expr == "" {
		return 0
	}
	ex.Params["x"] = float64(x)
	ex.Params["t"] = float64(t)
	yi, _ := ex.Val.Evaluate(ex.Params)
	y := float32(yi.(float64))
	return y
}

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
}

const LIM = 100

var facts [LIM]float64

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
