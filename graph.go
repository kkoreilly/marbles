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
	"sync"
	"unicode"

	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/mat32/v2"
	"goki.dev/pi/v2/complete"
	"goki.dev/svg"
)

// Graph contains the lines and parameters of a graph
type Graph struct { //gti:add

	// the parameters for updating the marbles
	Params Params

	// the lines of the graph -- can have any number
	Lines Lines

	Marbles []*Marble `json:"-"`

	State State `json:"-"`

	Functions Functions `json:"-"`

	Vectors Vectors `json:"-"`

	Objects Objects `json:"-"`

	EvalMu sync.Mutex `json:"-"`
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
	Colors LineColors

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
	StartVelY Param `view:"inline" label:"Starting velocity y"`

	// Starting vertical velocity of the marbles
	StartVelX Param `view:"inline" label:"Starting velocity x"`

	// how fast to move along velocity vector -- lower = smoother, more slow-mo
	UpdateRate Param `view:"inline"`

	// how fast time increases
	TimeStep Param `view:"inline"`

	// how fast it accelerates down
	YForce Param `view:"inline" label:"Y force (Gravity)"`

	// how fast the marbles move side to side without collisions, set to 0 for no movement
	XForce Param `view:"inline" label:"X force (Wind)"`

	// the center point of the graph, x
	CenterX Param `view:"inline" label:"Graph center x"`

	// the center point of the graph, y
	CenterY Param `view:"inline" label:"Graph center y"`

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

	LinesView *giv.TableView
}

// Lines is a collection of lines
type Lines []*Line

const GraphViewBoxSize = 10

var BasicFunctionList = []string{}

var CompleteWords = []string{}

// FunctionNames has all of the supported function names, in order
var FunctionNames = []string{"f", "g", "b", "c", "j", "k", "l", "m", "o", "p", "q", "r", "s", "u", "v", "w"}

// TheGraph is current graph
var TheGraph Graph

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
	defer gr.Objects.Graph.UpdateEndRender(updt)

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
	gr.ResetMarbles()
	gr.State.Time = 0
	if gr.State.Error != nil {
		return
	}
	gr.Lines.Graph()
	SetCompleteWords(TheGraph.Functions)
	// if gr.State.Error == nil {
	// 	errorText.SetText("Graphed successfully")
	// }
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
func (gr *Graph) StopSelecting() { //gti:add
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
func (gr *Graph) TrackSelectedMarble() { //gti:add
	if gr.State.SelectedMarble == -1 {
		return
	}
	gr.Marbles[gr.State.SelectedMarble].ToggleTrack(gr.State.SelectedMarble)
}

// AddLine adds a new blank line
func (gr *Graph) AddLine() { //gti:add
	var color color.RGBA
	if TheSettings.LineDefaults.LineColors.Color == colors.White {
		color = colors.BinarySpacedAccentVariant(len(gr.Lines) - 1)
	} else {
		color = TheSettings.LineDefaults.LineColors.Color
	}
	newLine := &Line{Colors: LineColors{color, TheSettings.LineDefaults.LineColors.ColorSwitch}}
	gr.Lines = append(gr.Lines, newLine)
	gr.Objects.LinesView.Update()
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
				ln.Colors.Color = colors.BinarySpacedAccentVariant(k)
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
	for i := range FunctionNames {
		if CheckIfReferences(expr, i) {
			return CheckCircular(TheGraph.Lines[i].Expr.Expr, k)
		}
	}
	return false
}

// CheckIfReferences checks if an expr references a given function
func CheckIfReferences(expr string, k int) bool {
	sort.Slice(BasicFunctionList, func(i, j int) bool {
		return len(BasicFunctionList[i]) > len(BasicFunctionList[j])
	})
	for _, d := range BasicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if k >= len(FunctionNames) || k >= len(TheGraph.Lines) {
		return false
	}
	funcName := FunctionNames[k]
	if strings.Contains(expr, funcName) || strings.Contains(expr, strings.ToUpper(funcName)) {
		return true
	}
	return false
}

// CheckIfChanges checks if an equation changes over time
func CheckIfChanges(expr string) bool {
	for _, d := range BasicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if strings.Contains(expr, "a") || strings.Contains(expr, "h") || strings.Contains(expr, "t") {
		return true
	}
	for k := range FunctionNames {
		if CheckIfReferences(expr, k) {
			return CheckIfChanges(TheGraph.Lines[k].Expr.Expr)
		}
	}
	return false
}

// InitBasicFunctionList adds all of the basic functions to a list
func InitBasicFunctionList() {
	for k := range DefaultFunctions {
		BasicFunctionList = append(BasicFunctionList, k)
	}
	BasicFunctionList = append(BasicFunctionList, "true", "false")
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
		ln.Colors.Color = colors.BinarySpacedAccentVariant(lidx)
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
	TheGraph.EvalMu.Lock()
	defer TheGraph.EvalMu.Unlock()

	if !TheGraph.State.Running {
		updt := TheGraph.Objects.Graph.UpdateStart()
		defer TheGraph.Objects.Graph.UpdateEnd(updt)
		nln := len(*ls)
		if TheGraph.Objects.Lines.NumChildren() != nln {
			TheGraph.Objects.Lines.SetNChildren(nln, svg.PathType, "line")
		}
	}
	if !TheGraph.State.Running || TheGraph.Params.CenterX.Changes || TheGraph.Params.CenterY.Changes {
		sizeFromCenter := mat32.Vec2{X: GraphViewBoxSize, Y: GraphViewBoxSize}
		center := mat32.Vec2{X: float32(TheGraph.Params.CenterX.Eval(0, 0)), Y: float32(TheGraph.Params.CenterY.Eval(0, 0))}
		TheGraph.Vectors.Min = center.Sub(sizeFromCenter)
		TheGraph.Vectors.Max = center.Add(sizeFromCenter)
		TheGraph.Vectors.Size = sizeFromCenter.MulScalar(2)
		TheGraph.Objects.Root.ViewBox.Min = mat32.Vec2{X: TheGraph.Vectors.Min.X, Y: -TheGraph.Vectors.Min.Y - 2*GraphViewBoxSize}
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
	pr.MarbleStartX.Expr = "rand-0.5"
	pr.MarbleStartY.Expr = "10-2n/nmarbles"
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
	for _, d := range BasicFunctionList {
		expr = strings.ReplaceAll(expr, d, "")
	}
	if CheckIfChanges(expr) || strings.Contains(expr, "x") || strings.Contains(expr, "y") {
		pr.Changes = true
	} else {
		pr.BaseVal = pr.Expr.Eval(0, 0, 0)
	}
}

// ExprComplete finds the possible completions for the expr in text field
func ExprComplete(data any, text string, posLn, posCh int) (md complete.Matches) {
	seedStart := 0
	for i := len(text) - 1; i >= 0; i-- {
		r := rune(text[i])
		if !unicode.IsLetter(r) || r == []rune("x")[0] || r == []rune("X")[0] {
			seedStart = i + 1
			break
		}
	}
	md.Seed = text[seedStart:]
	possibles := complete.MatchSeedString(CompleteWords, md.Seed)
	for _, p := range possibles {
		m := complete.Completion{Text: p, Icon: ""}
		md.Matches = append(md.Matches, m)
	}
	return md
}

// ExprCompleteEdit is the editing function called when using complete
func ExprCompleteEdit(data any, text string, cursorPos int, completion complete.Completion, seed string) (ed complete.Edit) {
	ed = complete.EditWord(text, cursorPos, completion.Text, seed)
	return ed
}

// SetCompleteWords sets the words used for complete in the expressions
func SetCompleteWords(functions Functions) {
	CompleteWords = []string{}
	for k := range functions {
		CompleteWords = append(CompleteWords, k)
	}
	CompleteWords = append(CompleteWords, "true", "false", "pi", "a", "t")
}
