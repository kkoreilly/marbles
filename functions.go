package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Knetic/govaluate"
	"gonum.org/v1/gonum/diff/fd"
)

// Functions are a map of named expression functions
type Functions map[string]govaluate.ExpressionFunction

// DefaultFunctions are the default functions that can be used in expressions
var DefaultFunctions = Functions{
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
	"log": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("log", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		y := math.Log(args[0].(float64)) / math.Log(args[1].(float64)) // log(v, b) = ln(v) / ln(b)
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
	"csch": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("csch", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Sinh(args[0].(float64))
		return y, nil
	},
	"sech": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("sech", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Cosh(args[0].(float64))
		return y, nil
	},
	"coth": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("coth", args, "float64")
		if err != nil {
			return 0, err
		}
		y := 1 / math.Tanh(args[0].(float64))
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
	"max": func(args ...interface{}) (interface{}, error) {
		num := math.Inf(-1)
		for _, d := range args {
			switch d := d.(type) {
			case float64:
				if d > num {
					num = d
				}
			default:
				return 0, errors.New("function max requires all number values")
			}
		}
		return num, nil
	},
	"min": func(args ...interface{}) (interface{}, error) {
		num := math.Inf(1)
		for _, d := range args {
			switch d := d.(type) {
			case float64:
				if d < num {
					num = d
				}
			default:
				return 0, errors.New("function min requires all number values")
			}
		}
		return num, nil
	},
	"avg": func(args ...interface{}) (interface{}, error) {
		var sum float64
		for _, d := range args {
			switch d := d.(type) {
			case float64:
				sum += d
			default:
				return 0, errors.New("function avg requires all number values")
			}
		}
		return sum / float64(len(args)), nil
	},
	"nmarbles": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("nmarbles", args)
		if err != nil {
			return 0, err
		}
		y := float64(TheGraph.Params.NMarbles)
		return y, nil
	},
	"inf": func(args ...interface{}) (interface{}, error) {
		err := CheckArgs("inf", args)
		if err != nil {
			return 0, err
		}
		return math.Inf(1), nil
	},
}

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

// SetFunctionsTo sets the functions of the graph to another set of functions
func (gr *Graph) SetFunctionsTo(functions Functions) {
	gr.Functions = make(Functions)
	for k, d := range functions {
		gr.Functions[k] = d
	}
}

// AddLineFunctions adds all of the line functions
func (gr *Graph) AddLineFunctions() {
	for k, ln := range gr.Lines {
		ln.SetFunctionName(k)
	}
}

// SetFunctionName sets the function name for a line and adds the function to the functions
func (ln *Line) SetFunctionName(k int) {
	if k >= len(functionNames) {
		// ln.FuncName = "unassigned"
		return
	}
	functionName := functionNames[k]
	// ln.FuncName = functionName + "(x)="
	TheGraph.Functions[functionName] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName, args, "float64")
		if err != nil {
			return 0, err
		}
		val := float64(ln.Expr.Eval(args[0].(float64), TheGraph.State.Time, ln.TimesHit))
		return val, nil
	}
	TheGraph.Functions[functionName+"'"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"'", args, "float64")
		if err != nil {
			return 0, err
		}
		val := fd.Derivative(func(x float64) float64 {
			return ln.Expr.Eval(x, TheGraph.State.Time, ln.TimesHit)
		}, args[0].(float64), &fd.Settings{
			Formula: fd.Central,
		})
		return val, nil
	}
	TheGraph.Functions[functionName+`"`] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+`"`, args, "float64")
		if err != nil {
			return 0, err
		}
		val := fd.Derivative(func(x float64) float64 {
			return ln.Expr.Eval(x, TheGraph.State.Time, ln.TimesHit)
		}, args[0].(float64), &fd.Settings{
			Formula: fd.Central2nd,
		})
		return val, nil
	}
	capitalName := strings.ToUpper(functionName)
	TheGraph.Functions[capitalName] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(capitalName, args, "float64")
		if err != nil {
			return 0, err
		}
		val := ln.Expr.Integrate(0, args[0].(float64), ln.TimesHit)
		return val, nil
	}
	TheGraph.Functions[functionName+"int"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"int", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		min := args[0].(float64)
		max := args[1].(float64)
		val := ln.Expr.Integrate(min, max, ln.TimesHit)
		return val, nil
	}
	TheGraph.Functions[functionName+"h"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"h", args, "float64")
		if err != nil {
			return 0, err
		}
		return float64(ln.TimesHit) * args[0].(float64), nil
	}
	TheGraph.Functions[functionName+"sum"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"sum", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		total := 0.0
		for i := args[0].(float64); i <= args[1].(float64); i++ {
			total += (ln.Expr.Eval(i, TheGraph.State.Time, ln.TimesHit))
		}
		return total, nil
	}
	TheGraph.Functions[functionName+"psum"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"psum", args, "float64", "float64")
		if err != nil {
			return 0, err
		}
		total := 1.0
		for i := args[0].(float64); i <= args[1].(float64); i++ {
			total *= (ln.Expr.Eval(i, TheGraph.State.Time, ln.TimesHit))
		}
		return total, nil
	}
}
