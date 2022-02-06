// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"math/rand"

	"github.com/Knetic/govaluate"
	"github.com/goki/gi/gist"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

// Type Graph contains the lines and parameters of a graph
type Graph struct {
	Params Params `view:"inline" desc:"the parameters for updating the marbles"`
	Lines  Lines  `view:"-" desc:"the lines of the graph -- can have any number"`
}

// Line represents one line with an equation etc
type Line struct {
	Expr       Expr       `width:"60" label:"y=" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	MinX       Expr       `width:"30" step:"1" desc:"Minimum x value for this line."`
	MaxX       Expr       `step:"1" desc:"Maximum x value for this line."`
	MinY       Expr       `step:"1" desc:"Minimum y value for this line."`
	MaxY       Expr       `step:"1" desc:"Maximum y value for this line."`
	Bounce     Expr       `min:"0" max:"2" step:".05" desc:"how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all"`
	LineColors LineColors `desc:"Line color and colorswitch" view:"no-inline"`
	TimesHit   int        `view:"-" json:"-"`
}

// Params is the parameters of the graph
type Params struct {
	NMarbles   int     `min:"1" max:"10000" step:"10" desc:"number of marbles"`
	NSteps     int     `min:"100" max:"10000" step:"10" desc:"number of steps to take when running"`
	StartSpeed float32 `min:"0" max:"2" step:".05" desc:"Coordinates per unit of time"`
	UpdtRate   float32 `min:"0.001" max:"1" step:".01" desc:"how fast to move along velocity vector -- lower = smoother, more slow-mo"`
	Gravity    float32 `min:"0" max:"2" step:".01" desc:"how fast it accelerates down"`
	Width      float32 `min:"0" max:"10" step:"1" desc:"length of spawning zone for marbles, set to 0 for all spawn in a column"`
	TimeStep   float32 `min:"0.001" max:"100" step:".01" desc:"how fast time increases"`
	Time       float32 `view:"-" json:"-" inactive:"+" desc:"time in msecs since starting"`
	MinSize    mat32.Vec2
	MaxSize    mat32.Vec2
}
type LineColors struct {
	Color       gist.Color `desc:"color to draw the line in" view:"no-inline"`
	ColorSwitch gist.Color `desc:"Switch the color of the marble that hits this line" view:"no-inline"`
}

// Lines is a collection of lines
type Lines []*Line

// colors is all of the colors that are used for marbles and default lines
var colors = []string{"black", "red", "blue", "green", "purple", "brown", "orange"}

// Last Saved file is the last saved or opened file, used for the save button
var LastSavedFile string

// Stop is used to tell RunMarbles to stop
var Stop = false

// Gr is current graph
var Gr Graph

// KiT_Graph is there to have the toolbar
var KiT_Graph = kit.Types.AddType(&Graph{}, GraphProps)

// X and Y axis are the x and y axis
var xAxis, yAxis *svg.Line

// GraphProps define the ToolBar for overall app
var GraphProps = ki.Props{
	"ToolBar": ki.PropSlice{
		{Name: "OpenJSON", Value: ki.Props{
			"label": "Open",
			"desc":  "Opens line equations and params from a .json file.",
			"icon":  "file-open",
			"Args": ki.PropSlice{
				{Name: "File Name", Value: ki.Props{
					"ext":     ".json",
					"default": "savedGraphs/",
				}},
			},
		}},
		{Name: "SaveLast", Value: ki.Props{
			"label": "Save",
			"desc":  "Save line equations and params to the last opened / saved file.",
			"icon":  "file-save",
		}},
		{Name: "SaveJSON", Value: ki.Props{
			"label": "Save As...",
			"desc":  "Saves line equations and params to a .json file.",
			"icon":  "file-save",
			"Args": ki.PropSlice{
				{Name: "File Name", Value: ki.Props{
					"ext":     ".json",
					"default": "savedGraphs/",
				}},
			},
		}},
		{Name: "OpenAutoSave", Value: ki.Props{
			"label": "Open Autosaved",
			"desc":  "Opens the most recently graphed set of equations and parameters.",
			"icon":  "file-open",
		}},
		{Name: "sep-ctrl", Value: ki.BlankProp{}},
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
		{Name: "Upload", Value: ki.Props{
			"label": "Upload Graph",
			"desc":  "Allows other people to download your graph. Enter a name for your graph",
			"icon":  "file-upload",
			"Args": ki.PropSlice{
				{Name: "Name", Value: ki.Props{}},
			},
		}},
		{Name: "Download", Value: ki.Props{
			"label": "Download Graph",
			"desc":  "Download a graph from the global database",
			"icon":  "file-download",
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
	gr.CompileExprs()
	ResetMarbles()
	gr.Params.Time = 0
	problemWithEval = false
	errorText.SetText("")
	gr.Lines.Graph()
	gr.AutoSave()
}

// Run runs the marbles for NSteps
func (gr *Graph) Run() {
	go RunMarbles()
}

// Stop stops the marbles
func (gr *Graph) Stop() {
	Stop = true
}

// Step does one step update of marbles
func (gr *Graph) Step() {
	UpdateMarbles()
}

func (gr *Graph) AddLine() {
	newLine := &Line{Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, Expr{"", nil, nil}, LineColors{gist.NilColor, gist.NilColor}, 0}
	// newLine.Defaults(rand.Intn(10))
	Gr.Lines = append(Gr.Lines, newLine)
}

// Gets the lines of the graph ready for graphing
func (gr *Graph) CompileExprs() {
	for _, ln := range gr.Lines {
		if ln.Expr.Expr == "" {
			ln.Expr.Expr = TheSettings.LineDefaults.Expr
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
		if ln.MinX.Expr == "" {
			ln.MinX.Expr = TheSettings.LineDefaults.MinX
		}
		if ln.MaxX.Expr == "" {
			ln.MaxX.Expr = TheSettings.LineDefaults.MaxX
		}
		if ln.MinY.Expr == "" {
			ln.MinY.Expr = TheSettings.LineDefaults.MinY
		}
		if ln.MaxY.Expr == "" {
			ln.MaxY.Expr = TheSettings.LineDefaults.MaxY
		}
		ln.TimesHit = 0
		ln.LoopEquationChangeSlice()
		ln.Compile()
	}
}

// Compiles all of the expressions in a line
func (ln *Line) Compile() {
	ln.Expr.Compile()
	ln.Bounce.Compile()
	ln.MinX.Compile()
	ln.MaxX.Compile()
	ln.MinY.Compile()
	ln.MaxY.Compile()
}

// Sets the line to the defaults specified in settings
func (ln *Line) Defaults(lidx int) {
	ln.Expr.Expr = TheSettings.LineDefaults.Expr
	white, _ := gist.ColorFromName("white")
	if TheSettings.LineDefaults.LineColors.Color == white {
		color, _ := gist.ColorFromName(colors[lidx%len(colors)])
		ln.LineColors.Color = color
	} else {
		ln.LineColors.Color = TheSettings.LineDefaults.LineColors.Color
	}
	ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
	ln.MinX.Expr = TheSettings.LineDefaults.MinX
	ln.MaxX.Expr = TheSettings.LineDefaults.MaxX
	ln.MinY.Expr = TheSettings.LineDefaults.MinY
	ln.MaxY.Expr = TheSettings.LineDefaults.MaxY
	ln.LineColors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
}

// Makes the lines and then defaults them
func (ls *Lines) Defaults() {
	*ls = make(Lines, 1, 10)
	ln := Line{}
	(*ls)[0] = &ln
	ln.Defaults(0)

}

// Graphs the lines
func (ls *Lines) Graph() {
	updt := SvgGraph.UpdateStart()
	SvgGraph.ViewBox.Min = Gr.Params.MinSize
	SvgGraph.ViewBox.Size = Gr.Params.MaxSize.Sub(Gr.Params.MinSize)
	gmin = Gr.Params.MinSize
	gmax = Gr.Params.MaxSize
	gsz = Gr.Params.MaxSize.Sub(Gr.Params.MinSize)
	nln := len(*ls)
	if SvgLines.NumChildren() != nln {
		SvgLines.SetNChildren(nln, svg.KiT_Path, "line")
	}
	for i, ln := range *ls {
		ln.Graph(i)
	}
	SvgGraph.UpdateEnd(updt)
}

// Graphs a single line
func (ln *Line) Graph(lidx int) {
	if ln.Expr.Expr == "" {
		ln.Defaults(lidx)
	}
	if ln.LineColors.Color == gist.NilColor {
		white, _ := gist.ColorFromName("white")
		if TheSettings.LineDefaults.LineColors.Color == white {
			color, _ := gist.ColorFromName(colors[lidx%len(colors)])
			ln.LineColors.Color = color
		} else {
			ln.LineColors.Color = TheSettings.LineDefaults.LineColors.Color
		}
	}
	if ln.LineColors.ColorSwitch == gist.NilColor {
		ln.LineColors.ColorSwitch = TheSettings.LineDefaults.LineColors.ColorSwitch
	}
	if ln.Bounce.Expr == "" {
		ln.Bounce.Expr = TheSettings.LineDefaults.Bounce
	}
	if ln.MinX.Expr == "" {
		ln.MinX.Expr = TheSettings.LineDefaults.MinX
	}
	if ln.MaxX.Expr == "0" {
		ln.MaxX.Expr = TheSettings.LineDefaults.MaxX
	}
	if ln.MinY.Expr == "0" {
		ln.MinY.Expr = TheSettings.LineDefaults.MinY
	}
	if ln.MaxY.Expr == "0" {
		ln.MaxY.Expr = TheSettings.LineDefaults.MaxY
	}
	path := SvgLines.Child(lidx).(*svg.Path)
	path.SetProp("fill", "none")
	clr := ln.LineColors.Color
	// clr = EvalColorIf(clr, ln.TimesHit)
	path.SetProp("stroke", clr)
	var err error
	ln.Expr.Val, err = govaluate.NewEvaluableExpressionWithFunctions(ln.Expr.Expr, functions)
	if HandleError(err) {
		ln.Expr.Val = nil
		return
	}

	ps := ""
	start := true
	for x := gmin.X; x < gmax.X; x += ginc.X {
		if problemWithEval {
			return
		}
		MinX := ln.MinX.Eval(x, Gr.Params.Time, ln.TimesHit)
		MaxX := ln.MaxX.Eval(x, Gr.Params.Time, ln.TimesHit)
		MinY := ln.MinY.Eval(x, Gr.Params.Time, ln.TimesHit)
		MaxY := ln.MaxY.Eval(x, Gr.Params.Time, ln.TimesHit)
		y := ln.Expr.Eval(x, Gr.Params.Time, ln.TimesHit)
		if x > MinX && x < MaxX && y > MinY && y < MaxY {

			if start {
				ps += fmt.Sprintf("M %v %v ", x, y)
				start = false
			} else {
				ps += fmt.Sprintf("L %v %v ", x, y)
			}
		}
	}
	path.SetData(ps)
}

// Makes the x and y axis
func InitCoords() {
	updt := SvgGraph.UpdateStart()
	SvgCoords.DeleteChildren(true)

	xAxis = svg.AddNewLine(SvgCoords, "xAxis", -1000, 0, 1000, 0)
	xAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	yAxis = svg.AddNewLine(SvgCoords, "yAxis", 0, -1000, 0, 1000)
	yAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)

	SvgGraph.UpdateEnd(updt)
}

func (mb *Marble) Init(diff float32) {
	randNum := (rand.Float32() * 2) - 1
	xPos := randNum * Gr.Params.Width
	mb.Pos = mat32.Vec2{X: xPos, Y: Gr.Params.MaxSize.Y - diff}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	mb.Vel = mat32.Vec2{X: 0, Y: float32(-Gr.Params.StartSpeed)}
	mb.PrvPos = mb.Pos
}

func (pr *Params) Defaults() {
	pr.NMarbles = TheSettings.GraphDefaults.NMarbles
	pr.NSteps = TheSettings.GraphDefaults.NSteps
	pr.StartSpeed = TheSettings.GraphDefaults.StartSpeed
	pr.UpdtRate = TheSettings.GraphDefaults.UpdtRate
	pr.Gravity = TheSettings.GraphDefaults.Gravity
	pr.TimeStep = TheSettings.GraphDefaults.TimeStep
	pr.MinSize = TheSettings.GraphDefaults.MinSize
	pr.MaxSize = TheSettings.GraphDefaults.MaxSize
	pr.Width = TheSettings.GraphDefaults.Width
}

func (pr *Params) BasicDefaults() {
	pr.NMarbles = 10
	pr.NSteps = 10000
	pr.StartSpeed = 0
	pr.UpdtRate = .02
	pr.Gravity = 0.1
	pr.TimeStep = 0.01
	pr.MinSize = mat32.Vec2{X: -10, Y: -10}
	pr.MaxSize = mat32.Vec2{X: 10, Y: 10}
	pr.Width = 0
}
