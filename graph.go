// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gist"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
	"github.com/goki/pi/complete"
	"gonum.org/v1/gonum/diff/fd"
)

// Graph contains the lines and parameters of a graph
type Graph struct {
	Params    Params    `view:"-" desc:"the parameters for updating the marbles"`
	Lines     Lines     `view:"-" desc:"the lines of the graph -- can have any number"`
	Marbles   []*Marble `view:"-" json:"-"`
	State     State     `view:"-" json:"-"`
	Functions Functions `view:"-" json:"-"`
}

// State has the state of the graph
type State struct {
	Running bool
	Time    float64
	Error   error
	File    string
}

// Line represents one line with an equation etc
type Line struct {
	Expr     Expr       `width:"70" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	GraphIf  Expr       `width:"50" desc:"Graph this line if this condition is true. Ex: x>3"`
	Bounce   Expr       `width:"30" min:"0" max:"2" step:".05" desc:"how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all"`
	Colors   LineColors `desc:"Line color and colorswitch" view:"no-inline"`
	TimesHit int        `view:"-" json:"-"`
	Changes  bool       `view:"-" json:"-"`
}

// Params is the parameters of the graph
type Params struct {
	NMarbles         int        `min:"1" max:"10000" step:"10" desc:"number of marbles"`
	Width            float64    `min:"0" step:"1" desc:"length of spawning zone for marbles, set to 0 for all spawn in a column"`
	MarbleStartPos   mat32.Vec2 `desc:"Marble starting position"`
	NSteps           int        `step:"10" desc:"number of steps to take when running, set to negative 1 to run until stopped"`
	StartSpeed       float64    `min:"0" max:"2" step:".05" desc:"Coordinates per unit of time"`
	UpdtRate         Param      `desc:"how fast to move along velocity vector -- lower = smoother, more slow-mo"`
	TimeStep         Param      `desc:"how fast time increases"`
	YForce           Param      `label:"Y Force (Gravity)" desc:"how fast it accelerates down"`
	XForce           Param      `label:"X Force (Wind)" desc:"how fast the marbles move side to side without collisions, set to 0 for no movement"`
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

var completeWords = []string{}

// functionNames has all of the supported function names, in order
var functionNames = []string{"f", "g", "b", "c", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "u", "v", "w"}

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
		{Name: "Step", Value: ki.Props{
			"desc":            "steps the marbles for one step",
			"icon":            "step-fwd",
			"no-update-after": true,
		}},
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
			"label":    "Add New Line",
			"desc":     "Adds a new line",
			"icon":     "plus",
			"shortcut": "Control+M",
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
	if gr.State.Running {
		return
	}
	gr.State.Error = nil
	gr.SetFunctionsTo(DefaultFunctions)
	InitCoords()
	gr.CompileExprs()
	ResetMarbles()
	if gr.State.Error != nil {
		return
	}
	gr.State.Time = 0
	rand.Seed(time.Now().UnixNano())
	randNum = rand.Float64()
	gr.Lines.Graph(false)
	SetCompleteWords(TheGraph.Functions)
	if gr.State.Error == nil {
		errorText.SetText("Graphed successfully")
	}
}

// AutoGraph is used to graph the function when something is changed
func (gr *Graph) AutoGraph() {
	updt := svgGraph.UpdateStart()
	TheGraph.Graph()
	svgGraph.SetNeedsFullRender()
	svgGraph.UpdateEnd(updt)
}

// Run runs the marbles for NSteps
func (gr *Graph) Run() {
	go RunMarbles()
}

// Stop stops the marbles
func (gr *Graph) Stop() {
	gr.State.Running = false
}

// Step does one step update of marbles
func (gr *Graph) Step() {
	if gr.State.Running {
		return
	}
	UpdateMarbles()
	TheGraph.State.Time += TheGraph.Params.TimeStep.Eval(0, 0)
}

// SelectNextMarble calls select next marble
func (gr *Graph) SelectNextMarble() {
	SelectNextMarble()
}

// StopSelecting stops selecting current marble
func (gr *Graph) StopSelecting() {
	var updt bool
	if !gr.State.Running {
		updt = svgGraph.UpdateStart()
	}
	if selectedMarble != -1 {
		svgMarbles.Child(selectedMarble).SetProp("stroke", "none")
		selectedMarble = -1
	}
	if !gr.State.Running {
		svgGraph.UpdateEnd(updt)
		svgGraph.SetNeedsFullRender()
	}
}

// TrackSelectedMarble toggles track for the currently selected marble
func (gr *Graph) TrackSelectedMarble() {
	if selectedMarble == -1 {
		return
	}
	gr.Marbles[selectedMarble].ToggleTrack(selectedMarble)
}

// AddLine adds a new blank line
func (gr *Graph) AddLine() {
	k := len(gr.Lines)
	var color gist.Color
	if TheSettings.LineDefaults.LineColors.Color == gist.White {
		color, _ = gist.ColorFromName(colors[k%len(colors)])
	} else {
		color = TheSettings.LineDefaults.LineColors.Color
	}
	newLine := &Line{Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, LineColors{color, gist.White}, 0, false}
	gr.Lines = append(gr.Lines, newLine)
}

// Reset resets the graph to its starting position (one default line and default params)
func (gr *Graph) Reset() {
	gr.State.File = ""
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
		if ln.Colors.Color == gist.NilColor {
			if TheSettings.LineDefaults.LineColors.Color == gist.White {
				color, _ := gist.ColorFromName(colors[k%len(colors)])
				ln.Colors.Color = color
			} else {
				ln.Colors.Color = TheSettings.LineDefaults.LineColors.Color
			}
		}
		if ln.Colors.ColorSwitch == gist.NilColor {
			ln.Colors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
		}
		if ln.Bounce.Expr == "" {
			ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
		}
		if ln.GraphIf.Expr == "" {
			ln.GraphIf.Expr = TheSettings.LineDefaults.GraphIf
		}
		if CheckCircular(ln.Expr.Expr, k) {
			HandleError(errors.New("circular logic detected"))
			return
		}
		ln.SetFunctionName(k)
		if CheckIfChanges(ln.Expr.Expr) || CheckIfChanges(ln.GraphIf.Expr) || CheckIfChanges(ln.Bounce.Expr) {
			ln.Changes = true
		}
		ln.TimesHit = 0
		ln.Compile()
	}
	gr.Params.UpdtRate.Compile()
	gr.Params.YForce.Compile()
	gr.Params.XForce.Compile()
	gr.Params.TimeStep.Compile()
}

// CheckCircular checks if an expr references itself
func CheckCircular(expr string, k int) bool {
	if CheckIfReferences(expr, k) {
		return true
	}
	for i := range functionNames {
		if CheckIfReferences(expr, i) {
			return CheckCircular(TheGraph.Lines[i].Expr.Expr, k)
		}
	}
	return false
}

// CheckIfReferences checks if an expr references a given function
func CheckIfReferences(expr string, k int) bool {
	sort.Slice(basicFunctionList, func(i, j int) bool {
		return len(basicFunctionList[i]) > len(basicFunctionList[j])
	})
	for _, d := range basicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if k >= len(functionNames) || k >= len(TheGraph.Lines) {
		return false
	}
	funcName := functionNames[k]
	if strings.Contains(expr, funcName) || strings.Contains(expr, strings.ToUpper(funcName)) {
		return true
	}
	return false
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
		err := CheckArgs(functionName+"d", args, "float64")
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
		err := CheckArgs(functionName+"dd", args, "float64")
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
	TheGraph.Functions[functionName+"i"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"i", args, "float64", "float64")
		if err != nil {
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
	for k := range functionNames {
		if CheckIfReferences(expr, k) {
			return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
		}
	}
	return false
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
		ln.Colors.Color = color
	} else {
		ln.Colors.Color = TheSettings.LineDefaults.LineColors.Color
	}
	ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
	ln.GraphIf.Expr = TheSettings.LineDefaults.GraphIf
	ln.Colors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
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
	path.SetProp("stroke", ln.Colors.Color)
	ps := ""
	start := true
	skipped := false
	for x := gmin.X; x < gmax.X; x += ginc.X {
		if TheGraph.State.Error != nil {
			return
		}
		fx := float64(x)
		y := ln.Expr.Eval(fx, TheGraph.State.Time, ln.TimesHit)
		GraphIf := ln.GraphIf.EvalBool(fx, y, TheGraph.State.Time, ln.TimesHit)
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
	pr.YForce = TheSettings.GraphDefaults.YForce
	pr.XForce = TheSettings.GraphDefaults.XForce
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
	pr.TimeStep.Expr.Expr = "0.01"
	pr.YForce.Expr.Expr = "-0.1"
	pr.XForce.Expr.Expr = "0"
	pr.MinSize = mat32.Vec2{X: -10, Y: -10}
	pr.MaxSize = mat32.Vec2{X: 10, Y: 10}
	pr.Width = 0
	pr.TrackingSettings.Override = false
	pr.TrackingSettings.TrackingSettings.Defaults()
}

// Eval evaluates a parameter
func (pr *Param) Eval(x, y float64) float64 {
	if !pr.Changes {
		return pr.BaseVal
	}
	return pr.Expr.EvalWithY(x, TheGraph.State.Time, 0, y)
}

// Compile compiles evalexpr and sets changes
func (pr *Param) Compile() {
	pr.Expr.Compile()
	expr := pr.Expr.Expr
	for _, d := range basicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if CheckIfChanges(expr) || strings.Contains(expr, "x") || strings.Contains(expr, "y") {
		pr.Changes = true
	} else {
		pr.BaseVal = pr.Expr.Eval(0, 0, 0)
	}
}

// ExprComplete finds the possible completions for the expr in text field
func ExprComplete(data interface{}, text string, posLn, posCh int) (md complete.Matches) {
	seedStart := 0
	for i := len(text) - 1; i >= 0; i-- {
		r := rune(text[i])
		if !unicode.IsLetter(r) || r == []rune("x")[0] || r == []rune("X")[0] {
			seedStart = i + 1
			break
		}
	}
	md.Seed = text[seedStart:]
	possibles := complete.MatchSeedString(completeWords, md.Seed)
	for _, p := range possibles {
		m := complete.Completion{Text: p, Icon: ""}
		md.Matches = append(md.Matches, m)
	}
	return md
}

// ExprCompleteEdit is the editing function called when using complete
func ExprCompleteEdit(data interface{}, text string, cursorPos int, completion complete.Completion, seed string) (ed complete.Edit) {
	ed = complete.EditWord(text, cursorPos, completion.Text, seed)
	return ed
}

// SetCompleteWords sets the words used for complete in the expressions
func SetCompleteWords(functions Functions) {
	completeWords = []string{}
	for k := range functions {
		completeWords = append(completeWords, k)
	}
	completeWords = append(completeWords, "true", "false", "pi", "a", "t")
}
