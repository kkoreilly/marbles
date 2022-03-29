// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gist"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
	"gonum.org/v1/gonum/diff/fd"
)

// Graph contains the lines and parameters of a graph
type Graph struct {
	Params Params `view:"-" desc:"the parameters for updating the marbles"`
	Lines  Lines  `view:"-" desc:"the lines of the graph -- can have any number"`
}

// Line represents one line with an equation etc
type Line struct {
	FunctionName string     `desc:"Name of the function, use to refer to it in other equations. Can't be changed" json:"-"`
	Expr         Expr       `width:"50" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	GraphIf      Expr       `width:"50" desc:"Graph this line if this condition is true. Ex: x>3"`
	Bounce       Expr       `min:"0" max:"2" step:".05" desc:"how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all"`
	LineColors   LineColors `desc:"Line color and colorswitch" view:"no-inline"`
	TimesHit     int        `view:"-" json:"-"`
	Changes      bool       `view:"-" json:"-"`
}

// Params is the parameters of the graph
type Params struct {
	NMarbles         int        `min:"1" max:"10000" step:"10" desc:"number of marbles"`
	Width            float64    `min:"0" max:"10" step:"1" desc:"length of spawning zone for marbles, set to 0 for all spawn in a column"`
	MarbleStartPos   mat32.Vec2 `desc:"Marble starting position"`
	NSteps           int        `step:"10" desc:"number of steps to take when running, set to negative 1 to run until stopped"`
	StartSpeed       float64    `min:"0" max:"2" step:".05" desc:"Coordinates per unit of time"`
	UpdtRate         Param      `desc:"how fast to move along velocity vector -- lower = smoother, more slow-mo"`
	Gravity          Param      `desc:"how fast it accelerates down"`
	TimeStep         Param      `desc:"how fast time increases"`
	Time             float64    `view:"-" json:"-" inactive:"+" desc:"time in msecs since starting"`
	MinSize          mat32.Vec2
	MaxSize          mat32.Vec2
	TrackingSettings GraphTrackingSettings `view:"no-inline"`
}

// Param is the type of certain parameters that can change over time and x
type Param struct {
	Expr    Expr    `width:"50" label:""`
	Changes bool    `view:"-"`
	BaseVal float64 `view:"-"`
}

// GraphTrackingSettings contains tracking line settings and a bool for whether the graph should override user settings.
type GraphTrackingSettings struct {
	Override         bool             `label:"Override user tracking settings. If false, TrackingSettings has no effect"`
	TrackingSettings TrackingSettings `view:"inline"`
}

// LineColors contains the color and colorswitch for a line
type LineColors struct {
	Color       gist.Color `desc:"color to draw the line in" view:"no-inline"`
	ColorSwitch gist.Color `desc:"Switch the color of the marble that hits this line" view:"no-inline"`
}

// Lines is a collection of lines
type Lines []*Line

// colors is all of the colors that are used for marbles and default lines
var colors = []string{"black", "red", "blue", "green", "purple", "brown", "orange"}

var basicFunctionList = []string{}

// functionNames has all of the supported function names, in order
var functionNames = []string{"f", "g", "b", "c", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "u", "v", "w", "y", "z"}

// currentFile is the last saved or opened file, used for the save button
var currentFile string

// stop is used to tell RunMarbles to stop
var stop = false

// last evaluated x value
var currentX float64

// TheGraph is current graph
var TheGraph Graph

// KiTGraph is there to have the toolbar
var KiTGraph = kit.Types.AddType(&Graph{}, GraphProps)

// X and Y axis are the x and y axis
var xAxis, yAxis *svg.Line

// GraphProps define the ToolBar for overall app
var GraphProps = ki.Props{
	"ToolBar": ki.PropSlice{
		// {Name: "OpenJSON", Value: ki.Props{
		// 	"label": "Open",
		// 	"desc":  "Opens line equations and params from a .json file.",
		// 	"icon":  "file-open",
		// 	"Args": ki.PropSlice{
		// 		{Name: "File Name", Value: ki.Props{
		// 			"ext":     ".json",
		// 			"default": "savedGraphs/",
		// 		}},
		// 	},
		// }},
		// {Name: "SaveLast", Value: ki.Props{
		// 	"label": "Save",
		// 	"desc":  "Save line equations and params to the last opened / saved file.",
		// 	"icon":  "file-save",
		// }},
		// {Name: "SaveJSON", Value: ki.Props{
		// 	"label": "Save As...",
		// 	"desc":  "Saves line equations and params to a .json file.",
		// 	"icon":  "file-save",
		// 	"Args": ki.PropSlice{
		// 		{Name: "File Name", Value: ki.Props{
		// 			"ext":     ".json",
		// 			"default": "savedGraphs/",
		// 		}},
		// 	},
		// }},
		// {Name: "OpenAutoSave", Value: ki.Props{
		// 	"label": "Open Autosaved",
		// 	"desc":  "Opens the most recently graphed set of equations and parameters.",
		// 	"icon":  "file-open",
		// }},
		// {Name: "sep-ctrl", Value: ki.BlankProp{}},
		{Name: "Graph", Value: ki.Props{
			"desc": "updates graph for current equations",
			"icon": "file-image",
		}},
		{Name: "Run", Value: ki.Props{
			"desc":            "runs the marbles for NSteps",
			"icon":            "run",
			"no-update-after": true,
		}},
		{Name: "Stop", Value: ki.Props{
			"desc":            "runs the marbles for NSteps",
			"icon":            "stop",
			"no-update-after": true,
		}},
		// {Name: "Jump", Value: ki.Props{
		// 	"label": "Jump",
		// 	"desc":  "Jump n frames forward",
		// 	"icon":  "run",
		// 	"Args": ki.PropSlice{
		// 		{Name: "n", Value: ki.Props{}},
		// 	},
		// }},
		{Name: "Step", Value: ki.Props{
			"desc":            "steps the marbles for one step",
			"icon":            "step-fwd",
			"no-update-after": true,
		}},
		// {Name: "sep-ctrl", Value: ki.BlankProp{}},
		// {Name: "Upload", Value: ki.Props{
		// 	"label": "Upload Graph",
		// 	"desc":  "Allows other people to download your graph. Enter a name for your graph",
		// 	"icon":  "file-upload",
		// 	"Args": ki.PropSlice{
		// 		{Name: "Name", Value: ki.Props{}},
		// 	},
		// }},
		// {Name: "Download", Value: ki.Props{
		// 	"label": "Download Graph",
		// 	"desc":  "Download a graph from the global database",
		// 	"icon":  "file-download",
		// }},
		{Name: "sep-ctrl", Value: ki.BlankProp{}},
		{Name: "SelectNextMarble", Value: ki.Props{
			"label":           "Next Marble",
			"desc":            "selects the next marble",
			"icon":            "forward",
			"no-update-after": true,
			"shortcut":        gi.KeyFunFocusNext,
		}},
		{Name: "StopSelecting", Value: ki.Props{
			"label":           "Unselect",
			"desc":            "stops selecting the marble",
			"icon":            "stop",
			"no-update-after": true,
		}},
		{Name: "TrackSelectedMarble", Value: ki.Props{
			"label":           "Track",
			"desc":            "toggles track for the currently selected marble",
			"icon":            "edit",
			"no-update-after": true,
			"shortcut":        gi.KeyFunTranspose,
		}},
		{Name: "sep-ctrl", Value: ki.BlankProp{}},
		{Name: "AddLine", Value: ki.Props{
			"label": "Add New Line",
			"desc":  "Adds a new line",
			"icon":  "plus",
		}},
	},
}

// Defaults sets the default parameters and lines for the graph, specified in settings
func (gr *Graph) Defaults() {
	gr.Params.Defaults()
	gr.Lines.Defaults()
}

// Graph updates graph for current equations, and resets marbles too
func (gr *Graph) Graph() {
	if runningMarbles {
		return
	}
	InitCoords()
	gr.CompileExprs()
	ResetMarbles()
	if problemWithCompile {
		problemWithCompile = false
		return
	}
	gr.Params.Time = 0
	problemWithEval = false
	rand.Seed(time.Now().UnixNano())
	randNum = rand.Float64()
	errorText.SetText("")
	gr.Lines.Graph(false)
	gr.AutoSave()
}

// Run runs the marbles for NSteps
func (gr *Graph) Run() {
	go RunMarbles()
}

// Stop stops the marbles
func (gr *Graph) Stop() {
	stop = true
	runningMarbles = false
}

// // Jump jumps n frames forward
// func (gr *Graph) Jump(n int) {
// 	if runningMarbles {
// 		return
// 	}
// 	Jump(n)
// }

// Step does one step update of marbles
func (gr *Graph) Step() {
	if runningMarbles {
		return
	}
	UpdateMarbles()
	TheGraph.Params.Time += TheGraph.Params.TimeStep.Eval(0)
}

// SelectNextMarble calls select next marble
func (gr *Graph) SelectNextMarble() {
	SelectNextMarble()
}

// StopSelecting stops selecting current marble
func (gr *Graph) StopSelecting() {
	var updt bool
	if !runningMarbles {
		updt = svgGraph.UpdateStart()
	}
	if selectedMarble != -1 {
		svgMarbles.Child(selectedMarble).SetProp("stroke", "none")
		selectedMarble = -1
	}
	if !runningMarbles {
		svgGraph.UpdateEnd(updt)
		svgGraph.SetNeedsFullRender()
	}
}

// TrackSelectedMarble toggles track for the currently selected marble
func (gr *Graph) TrackSelectedMarble() {
	if selectedMarble == -1 {
		return
	}
	Marbles[selectedMarble].ToggleTrack(selectedMarble)
}

// AddLine adds a new blank line
func (gr *Graph) AddLine() {
	newLine := &Line{"", Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, LineColors{gist.NilColor, gist.NilColor}, 0, false}
	// newLine.Defaults(rand.Intn(10))
	gr.Lines = append(gr.Lines, newLine)
}

// Reset resets the graph to its starting position (one default line and default params)
func (gr *Graph) Reset() {
	currentFile = ""
	UpdateCurrentFileText()
	gr.Lines = nil
	gr.Lines.Defaults()
	gr.Params.Defaults()
}

// CompileExprs gets the lines of the graph ready for graphing
func (gr *Graph) CompileExprs() {
	for k, ln := range gr.Lines {
		ln.Changes = false
		if ln.Expr.Expr == "" {
			ln.Expr.Expr = TheSettings.LineDefaults.Expr
		}
		if ln.LineColors.Color == gist.NilColor {
			if TheSettings.LineDefaults.LineColors.Color == gist.White {
				color, _ := gist.ColorFromName(colors[k%len(colors)])
				ln.LineColors.Color = color
			} else {
				ln.LineColors.Color = TheSettings.LineDefaults.LineColors.Color
			}
		}
		if ln.LineColors.ColorSwitch == gist.NilColor {
			ln.LineColors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
		}
		// if ln.LineColors.Color == gist.NilColor {
		// 	white, _ := gist.ColorFromName("white")
		// 	if TheSettings.LineDefaults.LineColors.Color == white {
		// 		black, _ := gist.ColorFromName("black")
		// 		ln.LineColors.Color = black
		// 	} else {
		// 		ln.LineColors.Color = TheSettings.LineDefaults.LineColors.Color
		// 	}
		// }
		if ln.Bounce.Expr == "" {
			ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
		}
		if ln.GraphIf.Expr == "" {
			ln.GraphIf.Expr = TheSettings.LineDefaults.GraphIf
		}
		ln.SetFunctionName(k)
		if CheckIfChanges(ln.Expr.Expr) || CheckIfChanges(ln.GraphIf.Expr) || CheckIfChanges(ln.Bounce.Expr) {
			ln.Changes = true
		}
		ln.TimesHit = 0

		// ln.CheckForDerivatives()
		ln.Compile()
	}
	gr.Params.UpdtRate.Compile()
	gr.Params.Gravity.Compile()
	gr.Params.TimeStep.Compile()
}

//SetFunctionName sets the function name for a line and adds the function to the functions
func (ln *Line) SetFunctionName(k int) {
	if k >= len(functionNames) {
		ln.FunctionName = "unassigned"
		return
	}
	functionName := functionNames[k]
	ln.FunctionName = functionName + "(x) = "
	DefaultFunctions[functionName] = func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), functionName)
		if !ok {
			return 0, err
		}
		val := float64(ln.Expr.Eval(args[0].(float64), TheGraph.Params.Time, ln.TimesHit))
		return val, nil
	}
	DefaultFunctions[functionName+"d"] = func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), functionName+"d")
		if !ok {
			return 0, err
		}
		val := fd.Derivative(func(x float64) float64 {
			return ln.Expr.Eval(x, TheGraph.Params.Time, ln.TimesHit)
		}, args[0].(float64), &fd.Settings{
			Formula: fd.Central,
		})
		return val, nil
	}
	DefaultFunctions[functionName+"dd"] = func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), functionName+"dd")
		if !ok {
			return 0, err
		}
		val := fd.Derivative(func(x float64) float64 {
			return ln.Expr.Eval(x, TheGraph.Params.Time, ln.TimesHit)
		}, args[0].(float64), &fd.Settings{
			Formula: fd.Central2nd,
		})
		return val, nil
	}
	capitalName := strings.ToUpper(functionName)
	DefaultFunctions[capitalName] = func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(1, len(args), capitalName)
		if !ok {
			return 0, err
		}
		val := ln.Expr.Integrate(0, args[0].(float64), ln.TimesHit)
		return val, nil
	}
	DefaultFunctions[functionName+"i"] = func(args ...interface{}) (interface{}, error) {
		ok, err := CheckArgs(2, len(args), functionName+"i")
		if !ok {
			return 0, err
		}
		min := args[0].(float64)
		max := args[1].(float64)
		val := ln.Expr.Integrate(min, max, ln.TimesHit)
		return val, nil
	}
}

// CheckIfChanges checks if an equation changes over time
func CheckIfChanges(expr string) bool {
	for _, d := range basicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if strings.Contains(expr, "a") || strings.Contains(expr, "h") || strings.Contains(expr, "t") {
		return true
	}
	strs := before(expr, "'")
	for _, d := range strs {
		for k, fn := range functionNames {
			if d == fn && k < len(TheGraph.Lines) {
				return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
			}
		}
	}
	strs = before(expr, `"`)
	for _, d := range strs {
		for k, fn := range functionNames {
			if d == fn && k < len(TheGraph.Lines) {
				return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
			}
		}
	}
	strs = before(expr, "i")
	for _, d := range strs {
		for k, fn := range functionNames {
			if d == fn && k < len(TheGraph.Lines) {
				return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
			}
		}
	}
	strs = before(expr, "(")
	for _, d := range strs {
		for k, fn := range functionNames {
			if (d == fn || d == strings.ToUpper(fn)) && k < len(TheGraph.Lines) {
				return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
			}
		}
	}
	return false
}
func before(str, substr string) []string {
	result := []string{}
	for {
		pos := strings.Index(str, substr)
		if pos == -1 {
			return result
		}
		result = append(result, str[0:pos])
		str = strings.Replace(str, substr, "", 1)
	}

}

// InitBasicFunctionList adds all of the basic functions to a list
func InitBasicFunctionList() {
	for k := range DefaultFunctions {
		basicFunctionList = append(basicFunctionList, k)
	}
	basicFunctionList = append(basicFunctionList, "true", "false")
}

// Compile compiles all of the expressions in a line
func (ln *Line) Compile() {
	ln.Expr.Compile()
	ln.Bounce.Compile()
	ln.GraphIf.Compile()
}

// Defaults sets the line to the defaults specified in settings
func (ln *Line) Defaults(lidx int) {
	ln.Expr.Expr = TheSettings.LineDefaults.Expr
	if TheSettings.LineDefaults.LineColors.Color == gist.White {
		color, _ := gist.ColorFromName(colors[lidx%len(colors)])
		ln.LineColors.Color = color
	} else {
		ln.LineColors.Color = TheSettings.LineDefaults.LineColors.Color
	}
	ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
	ln.GraphIf.Expr = TheSettings.LineDefaults.GraphIf
	ln.LineColors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
}

// Defaults makes the lines and then defaults them
func (ls *Lines) Defaults() {
	*ls = make(Lines, 1, 10)
	ln := Line{}
	(*ls)[0] = &ln
	ln.Defaults(0)

}

// Graph graphs the lines
func (ls *Lines) Graph(fromMarbles bool) {
	var updt bool
	if !fromMarbles {
		updt = svgGraph.UpdateStart()
	}
	svgGraph.ViewBox.Min = TheGraph.Params.MinSize
	svgGraph.ViewBox.Size = TheGraph.Params.MaxSize.Sub(TheGraph.Params.MinSize)
	gmin = TheGraph.Params.MinSize
	gmax = TheGraph.Params.MaxSize
	gsz = TheGraph.Params.MaxSize.Sub(TheGraph.Params.MinSize)
	nln := len(*ls)
	if svgLines.NumChildren() != nln {
		svgLines.SetNChildren(nln, svg.KiT_Path, "line")
	}
	for i, ln := range *ls {
		if !ln.Changes && fromMarbles { // If the line doesn't change over time then we don't need to keep graphing it while running marbles
			continue
		}
		ln.Graph(i, fromMarbles)
	}
	if !fromMarbles {
		svgGraph.UpdateEnd(updt)
	}
}

// Graph graphs a single line
func (ln *Line) Graph(lidx int, fromMarbles bool) {
	path := svgLines.Child(lidx).(*svg.Path)
	path.SetProp("fill", "none")
	path.SetProp("stroke", ln.LineColors.Color)
	ps := ""
	start := true
	skipped := false
	for x := gmin.X; x < gmax.X; x += ginc.X {
		if problemWithEval {
			return
		}
		fx := float64(x)
		y := ln.Expr.Eval(fx, TheGraph.Params.Time, ln.TimesHit)
		GraphIf := ln.GraphIf.EvalBool(fx, y, TheGraph.Params.Time, ln.TimesHit)
		if GraphIf && gmin.Y < float32(y) && gmax.Y > float32(y) {
			if start || skipped {
				ps += fmt.Sprintf("M %v %v ", x, y)
				start, skipped = false, false
			} else {
				ps += fmt.Sprintf("L %v %v ", x, y)
			}
		} else {
			skipped = true
		}
	}
	path.SetData(ps)
}

// InitCoords makes the x and y axis
func InitCoords() {
	updt := svgGraph.UpdateStart()
	svgCoords.DeleteChildren(true)

	xAxis = svg.AddNewLine(svgCoords, "xAxis", TheGraph.Params.MinSize.X, 0, TheGraph.Params.MaxSize.X, 0)
	xAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	yAxis = svg.AddNewLine(svgCoords, "yAxis", 0, TheGraph.Params.MinSize.Y, 0, TheGraph.Params.MaxSize.Y)
	yAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	svgGraph.UpdateEnd(updt)
}

// Defaults sets the graph parameters to the default settings
func (pr *Params) Defaults() {
	pr.NMarbles = TheSettings.GraphDefaults.NMarbles
	pr.MarbleStartPos = TheSettings.GraphDefaults.MarbleStartPos
	pr.NSteps = TheSettings.GraphDefaults.NSteps
	pr.StartSpeed = TheSettings.GraphDefaults.StartSpeed
	pr.UpdtRate = TheSettings.GraphDefaults.UpdtRate
	pr.Gravity = TheSettings.GraphDefaults.Gravity
	pr.TimeStep = TheSettings.GraphDefaults.TimeStep
	pr.MinSize = TheSettings.GraphDefaults.MinSize
	pr.MaxSize = TheSettings.GraphDefaults.MaxSize
	pr.Width = TheSettings.GraphDefaults.Width
	pr.TrackingSettings = TheSettings.GraphDefaults.TrackingSettings
}

// BasicDefaults sets the default defaults for the graph parameters
func (pr *Params) BasicDefaults() {
	pr.NMarbles = 10
	pr.MarbleStartPos = mat32.Vec2{X: 0, Y: 10}
	pr.NSteps = -1
	pr.StartSpeed = 0
	pr.UpdtRate.Expr.Expr = ".02"
	pr.Gravity.Expr.Expr = "0.1"
	pr.TimeStep.Expr.Expr = "0.01"
	pr.MinSize = mat32.Vec2{X: -10, Y: -10}
	pr.MaxSize = mat32.Vec2{X: 10, Y: 10}
	pr.Width = 0
	pr.TrackingSettings.Override = false
	pr.TrackingSettings.TrackingSettings.Defaults()
}

// Eval evaluates a parameter
func (pr *Param) Eval(x float64) float64 {
	if !pr.Changes {
		return pr.BaseVal
	}
	return pr.Expr.Eval(x, TheGraph.Params.Time, 0)
}

// Compile compiles evalexpr and sets changes
func (pr *Param) Compile() {
	pr.Expr.Compile()
	if CheckIfChanges(pr.Expr.Expr) || strings.Contains(pr.Expr.Expr, "x") {
		pr.Changes = true
	} else {
		pr.BaseVal = pr.Expr.Eval(0, 0, 0)
	}
}
