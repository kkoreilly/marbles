// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"image/color"
	"sort"
	"strings"
	"unicode"

	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/mat32/v2"
	"goki.dev/pi/v2/complete"
	"goki.dev/svg"
)

// Graph contains the lines and parameters of a graph
type Graph struct { //gti:add

	// the parameters for updating the marbles
	Params Params `view:"-"`

	// the lines of the graph -- can have any number
	Lines Lines `view:"-"`

	Marbles []*Marble `view:"-" json:"-"`

	State State `view:"-" json:"-"`

	Functions Functions `view:"-" json:"-"`

	Vectors Vectors `view:"-" json:"-"`

	Objects Objects `view:"-" json:"-"`
}

// State has the state of the graph
type State struct {
	Running        bool
	Time           float64
	PrevTime       float64
	Step           int
	Error          error
	SelectedMarble int
	File           string
}

// Line represents one line with an equation etc
type Line struct {

	// Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)
	Expr Expr

	// Graph this line if this condition is true. Ex: x>3
	GraphIf Expr

	// how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all
	Bounce Expr `min:"0" max:"2" step:".05"`

	// Line color and colorswitch
	Colors LineColors ` view:"no-inline"`

	TimesHit int `view:"-" json:"-"`

	Changes bool `view:"-" json:"-"`
}

// Params are the parameters of the graph
type Params struct { //gti:add

	// Number of marbles
	NMarbles int `min:"1" max:"10000" step:"10" label:"Number of marbles"`

	// Marble horizontal start position
	MarbleStartX Expr

	// Marble vertical start position
	MarbleStartY Expr

	// Starting horizontal velocity of the marbles
	StartVelY Param `label:"Starting velocity y"`

	// Starting vertical velocity of the marbles
	StartVelX Param `label:"Starting velocity x"`

	// how fast to move along velocity vector -- lower = smoother, more slow-mo
	UpdateRate Param

	// how fast time increases
	TimeStep Param

	// how fast it accelerates down
	YForce Param `label:"Y force (Gravity)"`

	// how fast the marbles move side to side without collisions, set to 0 for no movement
	XForce Param `label:"X force (Wind)"`

	// the center point of the graph, x
	CenterX Param `label:"Graph center x"`

	// the center point of the graph, y
	CenterY Param `label:"Graph center y"`

	TrackingSettings TrackingSettings
}

// Param is the type of certain parameters that can change over time and x
type Param struct {
	Expr Expr `label:""`

	Changes bool `view:"-"`

	BaseVal float64 `view:"-"`
}

// LineColors contains the color and colorswitch for a line
type LineColors struct {

	// color to draw the line in
	Color color.RGBA

	// Switch the color of the marble that hits this line
	ColorSwitch color.RGBA
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
	Graph         *gi.SVG
	SVG           *svg.SVG
	Root          *svg.SVGNode
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

var basicFunctionList = []string{}

var completeWords = []string{}

// functionNames has all of the supported function names, in order
var functionNames = []string{"f", "g", "b", "c", "j", "k", "l", "m", "o", "p", "q", "r", "s", "u", "v", "w"}

// TheGraph is current graph
var TheGraph Graph

/*
// KiTGraph is there to have the toolbar
// var KiTGraph = kit.Types.AddType(&Graph{}, GraphProps)

// // GraphProps define the ToolBar for overall app
// var GraphProps = ki.Props{
// 	"ToolBar": ki.PropSlice{
// 		{Name: "Graph", Value: ki.Props{
// 			"desc": "updates graph for current equations",
// 			"icon": "file-image",
// 		}},
// 		{Name: "Run", Value: ki.Props{
// 			"desc":            "runs the marbles for NSteps",
// 			"icon":            "run",
// 			"no-update-after": true,
// 		}},
// 		{Name: "Stop", Value: ki.Props{
// 			"desc":            "runs the marbles for NSteps",
// 			"icon":            "stop",
// 			"no-update-after": true,
// 		}},
// 		{Name: "Step", Value: ki.Props{
// 			"desc":            "steps the marbles for one step",
// 			"icon":            "step-fwd",
// 			"no-update-after": true,
// 		}},
// 		{Name: "sep-ctrl", Value: ki.BlankProp{}},
// 		{Name: "SelectNextMarble", Value: ki.Props{
// 			"label":           "Next Marble",
// 			"desc":            "selects the next marble",
// 			"icon":            "forward",
// 			"no-update-after": true,
// 			"shortcut":        gi.KeyFunFocusNext,
// 		}},
// 		{Name: "StopSelecting", Value: ki.Props{
// 			"label":           "Unselect",
// 			"desc":            "stops selecting the marble",
// 			"icon":            "stop",
// 			"no-update-after": true,
// 		}},
// 		{Name: "TrackSelectedMarble", Value: ki.Props{
// 			"label":           "Track",
// 			"desc":            "toggles track for the currently selected marble",
// 			"icon":            "edit",
// 			"no-update-after": true,
// 			"shortcut":        gi.KeyFunTranspose,
// 		}},
// 		{Name: "sep-ctrl", Value: ki.BlankProp{}},
// 		{Name: "AddLine", Value: ki.Props{
// 			"label":    "Add New Line",
// 			"desc":     "Adds a new line",
// 			"icon":     "plus",
// 			"shortcut": "Command+M",
// 		}},
// 	},
// }
*/

// Init sets up the graph for the given body. It should only be called once.
func (gr *Graph) Init(b *gi.Body) {
	gr.Defaults()
	gr.MakeBasicElements(b)
	gr.SetFunctionsTo(DefaultFunctions)
	gr.InitCoords()
	gr.CompileExprs()
	gr.Lines.Graph()
	gr.ResetMarbles()
}

// Defaults sets the default parameters and lines for the graph, specified in settings
func (gr *Graph) Defaults() {
	gr.Params.Defaults()
	gr.Lines.Defaults()
}

// Graph updates graph for current equations, and resets marbles too
func (gr *Graph) Graph() { //gti:add
	updt := gr.Objects.Graph.UpdateStart()
	if gr.State.Running {
		gr.Stop()
	}
	gr.State.Error = nil
	gr.SetFunctionsTo(DefaultFunctions)
	gr.AddLineFunctions()
	gr.CompileExprs()
	if gr.State.Error != nil {
		gr.Objects.Graph.UpdateEndRender(updt)
		return
	}
	gr.ResetMarbles()
	gr.State.Time = 0
	SetRandNum()
	if gr.State.Error != nil {
		gr.Objects.Graph.UpdateEndRender(updt)
		return
	}
	gr.Lines.Graph()
	SetCompleteWords(TheGraph.Functions)
	// if gr.State.Error == nil {
	// 	errorText.SetText("Graphed successfully")
	// }
	gr.Objects.Graph.UpdateEndRender(updt)
}

// Run runs the marbles for NSteps
func (gr *Graph) Run() { //gti:add
	// gr.AutoSave()
	go gr.RunMarbles()
}

// Stop stops the marbles
func (gr *Graph) Stop() { //gti:add
	gr.State.Running = false
}

// Step does one step update of marbles
func (gr *Graph) Step() { //gti:add
	if gr.State.Running {
		return
	}
	gr.UpdateMarbles()
	gr.State.Time += gr.Params.TimeStep.Eval(0, 0)
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
		gr.Objects.Graph.UpdateEndRender(updt)
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
	var color color.RGBA
	if TheSettings.LineDefaults.LineColors.Color == colors.White {
		color = colors.AccentVariantList(k)[k-1]
	} else {
		color = TheSettings.LineDefaults.LineColors.Color
	}
	newLine := &Line{Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, LineColors{color, TheSettings.LineDefaults.LineColors.ColorSwitch}, 0, false}
	gr.Lines = append(gr.Lines, newLine)
}

// Reset resets the graph to its starting position (one default line and default params)
func (gr *Graph) Reset() {
	gr.State.File = ""
	// UpdateCurrentFileText()
	gr.Lines = nil
	gr.Lines.Defaults()
	gr.Params.Defaults()
	gr.Graph()
}

// CompileExprs gets the lines of the graph ready for graphing
func (gr *Graph) CompileExprs() {
	for k, ln := range gr.Lines {
		ln.Changes = false
		if ln.Expr.Expr == "" {
			ln.Expr.Expr = TheSettings.LineDefaults.Expr
		}
		if colors.IsNil(ln.Colors.Color) {
			if TheSettings.LineDefaults.LineColors.Color == colors.White {
				ln.Colors.Color = colors.AccentVariantList(len(gr.Lines))[k]
			} else {
				ln.Colors.Color = TheSettings.LineDefaults.LineColors.Color
			}
		}
		if colors.IsNil(ln.Colors.ColorSwitch) {
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
	gr.Params.UpdateRate.Compile()
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
	if TheSettings.LineDefaults.LineColors.Color == colors.White {
		ln.Colors.Color = colors.AccentVariantList(len(TheGraph.Lines))[lidx]
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
			TheGraph.Objects.Lines.SetNChildren(nln, svg.PathType, "line")
		}
	}
	if !TheGraph.State.Running || TheGraph.Params.CenterX.Changes || TheGraph.Params.CenterY.Changes {
		sizeFromCenter := mat32.Vec2{X: graphViewBoxSize, Y: graphViewBoxSize}
		center := mat32.Vec2{X: float32(TheGraph.Params.CenterX.Eval(0, 0)), Y: float32(TheGraph.Params.CenterY.Eval(0, 0))}
		TheGraph.Vectors.Min = center.Sub(sizeFromCenter)
		TheGraph.Vectors.Max = center.Add(sizeFromCenter)
		TheGraph.Vectors.Size = sizeFromCenter.MulScalar(2)
		TheGraph.Objects.Root.ViewBox.Min = mat32.Vec2{X: TheGraph.Vectors.Min.X, Y: -TheGraph.Vectors.Min.Y - 2*graphViewBoxSize}
		TheGraph.Objects.Root.ViewBox.Size = TheGraph.Vectors.Size
		TheGraph.UpdateCoords()
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
func (gr *Graph) InitCoords() {
	updt := gr.Objects.Graph.UpdateStart()
	gr.Objects.Coords.DeleteChildren(true)

	gr.Objects.XAxis = svg.NewLine(gr.Objects.Coords, "x-axis")
	gr.Objects.XAxis.Start.X = gr.Vectors.Min.X
	gr.Objects.XAxis.End.X = gr.Vectors.Max.X
	gr.Objects.XAxis.SetProp("stroke", colors.Scheme.Outline)

	gr.Objects.YAxis = svg.NewLine(gr.Objects.Coords, "y-axis")
	gr.Objects.YAxis.Start.Y = gr.Vectors.Min.Y
	gr.Objects.YAxis.End.Y = gr.Vectors.Max.Y
	gr.Objects.YAxis.SetProp("stroke", colors.Scheme.Outline)

	gr.Objects.Graph.UpdateEnd(updt)
}

// UpdateCoords updates the x and y axis
func (gr *Graph) UpdateCoords() {
	updt := gr.Objects.Graph.UpdateStart()

	gr.Objects.XAxis.SetProp("stroke", colors.Scheme.Outline)
	gr.Objects.XAxis.Start, gr.Objects.XAxis.End = mat32.Vec2{X: gr.Vectors.Min.X, Y: 0}, mat32.Vec2{X: gr.Vectors.Max.X, Y: 0}

	gr.Objects.YAxis.SetProp("stroke", colors.Scheme.Outline)
	gr.Objects.YAxis.Start, gr.Objects.YAxis.End = mat32.Vec2{X: 0, Y: gr.Vectors.Min.Y}, mat32.Vec2{X: 0, Y: gr.Vectors.Max.Y}

	gr.Objects.Graph.UpdateEnd(updt)
}

// Defaults sets the graph parameters to the default settings
func (pr *Params) Defaults() {
	pr.NMarbles = TheSettings.GraphDefaults.NMarbles
	pr.MarbleStartX = TheSettings.GraphDefaults.MarbleStartX
	pr.MarbleStartY = TheSettings.GraphDefaults.MarbleStartY
	pr.StartVelY = TheSettings.GraphDefaults.StartVelY
	pr.StartVelX = TheSettings.GraphDefaults.StartVelX
	pr.UpdateRate = TheSettings.GraphDefaults.UpdateRate
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
	pr.MarbleStartX.Expr = "rand(1)-0.5"
	pr.MarbleStartY.Expr = "10-2n/nmarbles()"
	pr.StartVelY.Expr.Expr = "0"
	pr.StartVelX.Expr.Expr = "0"
	pr.UpdateRate.Expr.Expr = ".02"
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
