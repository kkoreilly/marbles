// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
	Vectors   Vectors   `view:"-" json:"-"`
	Objects   Objects   `view:"-" json:"-"`
}

// State has the state of the graph
type State struct {
	Running        bool
	Time           float64
	Step           int
	Error          error
	SelectedMarble int
	File           string
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
	NMarbles         int              `min:"1" max:"10000" step:"10" desc:"number of marbles"`
	MarbleStartX     Expr             `width:"100" desc:"Marble start position, x"`
	MarbleStartY     Expr             `width:"100" desc:"Marble start position, y"`
	StartVelY        Param            `label:"Starting Velocity Y" desc:"Starting velocity of the marbles, y"`
	StartVelX        Param            `label:"Starting Velocity X" desc:"Starting velocity of the marbles, x"`
	UpdtRate         Param            `desc:"how fast to move along velocity vector -- lower = smoother, more slow-mo"`
	TimeStep         Param            `desc:"how fast time increases"`
	YForce           Param            `label:"Y Force (Gravity)" desc:"how fast it accelerates down"`
	XForce           Param            `label:"X Force (Wind)" desc:"how fast the marbles move side to side without collisions, set to 0 for no movement"`
	CenterX          Param            `label:"Graph Center X" desc:"the center point of the graph, x"`
	CenterY          Param            `label:"Graph Center Y" desc:"the center point of the graph, y"`
	TrackingSettings TrackingSettings `view:"inline"`
}

// Param is the type of certain parameters that can change over time and x
type Param struct {
	Expr    Expr    `width:"100" label:""`
	Changes bool    `view:"-"`
	BaseVal float64 `view:"-"`
}

// LineColors contains the color and colorswitch for a line
type LineColors struct {
	Color       gist.Color `desc:"color to draw the line in" view:"no-inline"`
	ColorSwitch gist.Color `desc:"Switch the color of the marble that hits this line" view:"no-inline"`
}

// Vectors contains the size and increment of the graph
type Vectors struct {
	Min  mat32.Vec2
	Max  mat32.Vec2
	Size mat32.Vec2
	Inc  mat32.Vec2
}

// Objects contains the svg graph and the svg groups, plus the axes
type Objects struct {
	Graph         *svg.SVG
	Lines         *svg.Group
	Marbles       *svg.Group
	Coords        *svg.Group
	TrackingLines *svg.Group
	XAxis         *svg.Line
	YAxis         *svg.Line
}

// Lines is a collection of lines
type Lines []*Line

const graphViewBoxSize = 10

// colors is all of the colors that are used for marbles and default lines
var colors = []string{"black", "red", "blue", "green", "purple", "brown", "orange"}

var basicFunctionList = []string{}

var completeWords = []string{}

// functionNames has all of the supported function names, in order
var functionNames = []string{"f", "g", "b", "c", "j", "k", "l", "m", "o", "p", "q", "r", "s", "u", "v", "w"}

// TheGraph is current graph
var TheGraph Graph

// KiTGraph is there to have the toolbar
var KiTGraph = kit.Types.AddType(&Graph{}, GraphProps)

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
		gr.Stop()
	}
	gr.State.Error = nil
	gr.SetFunctionsTo(DefaultFunctions)
	gr.AddLineFunctions()
	gr.CompileExprs()
	if gr.State.Error != nil {
		return
	}
	ResetMarbles()
	gr.State.Time = 0
	SetRandNum()
	if gr.State.Error != nil {
		return
	}
	gr.Lines.Graph()
	SetCompleteWords(TheGraph.Functions)
	if gr.State.Error == nil {
		errorText.SetText("Graphed successfully")
	}
}

// AutoGraph is used to graph the function when something is changed
func (gr *Graph) AutoGraph() {
	updt := gr.Objects.Graph.UpdateStart()
	gr.Graph()
	gr.Objects.Graph.SetNeedsFullRender()
	gr.Objects.Graph.UpdateEnd(updt)
}

// AutoGraphAndUpdate calls autograph, and updates lns and params
func (gr *Graph) AutoGraphAndUpdate() {
	gr.AutoGraph()
	lns.Update()
	params.UpdateFields()
}

// Run runs the marbles for NSteps
func (gr *Graph) Run() {
	gr.AutoSave()
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
	gr.State.Time += gr.Params.TimeStep.Eval(0, 0)
}

// SelectNextMarble calls select next marble
func (gr *Graph) SelectNextMarble() {
	SelectNextMarble()
}

// StopSelecting stops selecting current marble
func (gr *Graph) StopSelecting() {
	var updt bool
	if !gr.State.Running {
		updt = gr.Objects.Graph.UpdateStart()
	}
	if gr.State.SelectedMarble != -1 {
		gr.Objects.Marbles.Child(gr.State.SelectedMarble).SetProp("stroke", "none")
		gr.State.SelectedMarble = -1
	}
	if !gr.State.Running {
		gr.Objects.Graph.UpdateEnd(updt)
		gr.Objects.Graph.SetNeedsFullRender()
	}
}

// TrackSelectedMarble toggles track for the currently selected marble
func (gr *Graph) TrackSelectedMarble() {
	if gr.State.SelectedMarble == -1 {
		return
	}
	gr.Marbles[gr.State.SelectedMarble].ToggleTrack(gr.State.SelectedMarble)
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
	newLine := &Line{Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, LineColors{color, TheSettings.LineDefaults.LineColors.ColorSwitch}, 0, false}
	gr.Lines = append(gr.Lines, newLine)
}

// Reset resets the graph to its starting position (one default line and default params)
func (gr *Graph) Reset() {
	gr.State.File = ""
	UpdateCurrentFileText()
	gr.Lines = nil
	gr.Lines.Defaults()
	gr.Params.Defaults()
	gr.AutoGraphAndUpdate()
}

// AddLineFunctions adds all of the line functions
func (gr *Graph) AddLineFunctions() {
	for k, ln := range gr.Lines {
		ln.SetFunctionName(k)
	}
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
		if CheckIfChanges(ln.Expr.Expr) || CheckIfChanges(ln.GraphIf.Expr) || CheckIfChanges(ln.Bounce.Expr) {
			ln.Changes = true
		}
		ln.TimesHit = 0
		ln.Compile()
	}
	gr.CompileParams()
}

// CompileParams compiles all of the graph parameter expressions
func (gr *Graph) CompileParams() {
	gr.Params.StartVelY.Compile()
	gr.Params.StartVelX.Compile()
	gr.Params.UpdtRate.Compile()
	gr.Params.YForce.Compile()
	gr.Params.XForce.Compile()
	gr.Params.TimeStep.Compile()
	gr.Params.CenterX.Compile()
	gr.Params.CenterY.Compile()
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
	TheGraph.Functions[functionName+"h"] = func(args ...interface{}) (interface{}, error) {
		err := CheckArgs(functionName+"h", args, "float64")
		if err != nil {
			return 0, err
		}
		return float64(ln.TimesHit) * args[0].(float64), nil
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
func (ls *Lines) Graph() {
	if !TheGraph.State.Running {
		updt := TheGraph.Objects.Graph.UpdateStart()
		defer TheGraph.Objects.Graph.UpdateEnd(updt)
		nln := len(*ls)
		if TheGraph.Objects.Lines.NumChildren() != nln {
			TheGraph.Objects.Lines.SetNChildren(nln, svg.KiT_Path, "line")
		}
	}
	if !TheGraph.State.Running || TheGraph.Params.CenterX.Changes || TheGraph.Params.CenterY.Changes {
		sizeFromCenter := mat32.Vec2{X: graphViewBoxSize, Y: graphViewBoxSize}
		center := mat32.Vec2{X: float32(TheGraph.Params.CenterX.Eval(0, 0)), Y: float32(TheGraph.Params.CenterY.Eval(0, 0))}
		TheGraph.Vectors.Min = center.Sub(sizeFromCenter)
		TheGraph.Vectors.Max = center.Add(sizeFromCenter)
		TheGraph.Vectors.Size = sizeFromCenter.MulScalar(2)
		TheGraph.Objects.Graph.ViewBox.Min = mat32.Vec2{X: TheGraph.Vectors.Min.X, Y: -TheGraph.Vectors.Min.Y - 2*graphViewBoxSize}
		TheGraph.Objects.Graph.ViewBox.Size = TheGraph.Vectors.Size
		UpdateCoords()
	}
	for i, ln := range *ls {
		// If the line doesn't change over time then we don't need to keep graphing it while running marbles
		if !ln.Changes && TheGraph.State.Running && !TheGraph.Params.CenterX.Changes && !TheGraph.Params.CenterY.Changes {
			continue
		}
		ln.Graph(i)
	}
}

// Graph graphs a single line
func (ln *Line) Graph(lidx int) {
	path := TheGraph.Objects.Lines.Child(lidx).(*svg.Path)
	path.SetProp("fill", "none")
	path.SetProp("stroke", ln.Colors.Color)
	ps := ""
	start := true
	skipped := false
	for x := TheGraph.Vectors.Min.X; x < TheGraph.Vectors.Max.X; x += TheGraph.Vectors.Inc.X {
		if TheGraph.State.Error != nil {
			return
		}
		fx := float64(x)
		y := ln.Expr.Eval(fx, TheGraph.State.Time, ln.TimesHit)
		GraphIf := ln.GraphIf.EvalBool(fx, y, TheGraph.State.Time, ln.TimesHit)
		if GraphIf && TheGraph.Vectors.Min.Y < float32(y) && TheGraph.Vectors.Max.Y > float32(y) {
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
	updt := TheGraph.Objects.Graph.UpdateStart()
	TheGraph.Objects.Coords.DeleteChildren(true)

	TheGraph.Objects.XAxis = svg.AddNewLine(TheGraph.Objects.Coords, "TheGraph.Objects.XAxis", TheGraph.Vectors.Min.X, 0, TheGraph.Vectors.Max.X, 0)
	TheGraph.Objects.XAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	TheGraph.Objects.YAxis = svg.AddNewLine(TheGraph.Objects.Coords, "TheGraph.Objects.YAxis", 0, TheGraph.Vectors.Min.Y, 0, TheGraph.Vectors.Max.Y)
	TheGraph.Objects.YAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	TheGraph.Objects.Graph.UpdateEnd(updt)
}

// UpdateCoords updates the x and y axis
func UpdateCoords() {
	updt := TheGraph.Objects.Graph.UpdateStart()

	TheGraph.Objects.XAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	TheGraph.Objects.XAxis.Start, TheGraph.Objects.XAxis.End = mat32.Vec2{X: TheGraph.Vectors.Min.X, Y: 0}, mat32.Vec2{X: TheGraph.Vectors.Max.X, Y: 0}

	TheGraph.Objects.YAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	TheGraph.Objects.YAxis.Start, TheGraph.Objects.YAxis.End = mat32.Vec2{X: 0, Y: TheGraph.Vectors.Min.Y}, mat32.Vec2{X: 0, Y: TheGraph.Vectors.Max.Y}

	TheGraph.Objects.Graph.UpdateEnd(updt)
}

// Defaults sets the graph parameters to the default settings
func (pr *Params) Defaults() {
	pr.NMarbles = TheSettings.GraphDefaults.NMarbles
	pr.MarbleStartX = TheSettings.GraphDefaults.MarbleStartX
	pr.MarbleStartY = TheSettings.GraphDefaults.MarbleStartY
	pr.StartVelY = TheSettings.GraphDefaults.StartVelY
	pr.StartVelX = TheSettings.GraphDefaults.StartVelX
	pr.UpdtRate = TheSettings.GraphDefaults.UpdtRate
	pr.YForce = TheSettings.GraphDefaults.YForce
	pr.XForce = TheSettings.GraphDefaults.XForce
	pr.TimeStep = TheSettings.GraphDefaults.TimeStep
	pr.CenterX = TheSettings.GraphDefaults.CenterX
	pr.CenterY = TheSettings.GraphDefaults.CenterY
	pr.TrackingSettings = TheSettings.GraphDefaults.TrackingSettings
}

// BasicDefaults sets the default defaults for the graph parameters
func (pr *Params) BasicDefaults() {
	pr.NMarbles = 10
	pr.MarbleStartX.Expr = "0(rand(1)-0.5)+0"
	pr.MarbleStartY.Expr = "10-2n/nmarbles()"
	pr.StartVelY.Expr.Expr = "0"
	pr.StartVelX.Expr.Expr = "0"
	pr.UpdtRate.Expr.Expr = ".02"
	pr.TimeStep.Expr.Expr = "0.01"
	pr.YForce.Expr.Expr = "-0.1"
	pr.XForce.Expr.Expr = "0"
	pr.CenterX.Expr.Expr = "0"
	pr.CenterY.Expr.Expr = "0"
	pr.TrackingSettings.Defaults()
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
