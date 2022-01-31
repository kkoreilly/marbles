// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"math/rand"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/chewxy/math32"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

type EquationChange struct {
	Old string
	New string
}

// Lines is a collection of lines
type Lines []*Line

// Line represents one line with an equation etc
type Line struct {
	Expr        Expr   `width:"60" label:"y=" desc:"Equation: use x for the x value, t for the time passed since the marbles were ran (incremented by TimeStep), and a for 10*sin(t) (swinging back and forth version of t)"`
	MinX        Expr   `width:"30" step:"1" desc:"Minimum x value for this line."`
	MaxX        Expr   `step:"1" desc:"Maximum x value for this line."`
	MinY        Expr   `step:"1" desc:"Minimum y value for this line."`
	MaxY        Expr   `step:"1" desc:"Maximum y value for this line."`
	Bounce      Expr   `min:"0" max:"2" step:".05" desc:"how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all"`
	Color       string `desc:"color to draw the line in"`
	ColorSwitch string `desc:"Switch the color of the marble that hits this line"`
	TimesHit    int    `view:"-"`
}

var EquationChangeSlice = []EquationChange{
	{"^", "**"},
	{"x(", "x*("},
	{"t(", "t*("},
	{"a(", "a*("},
	{"ax", "a*x"},
	{"xa", "x*a"},
	{"tx", "t*x"},
	{"xt", "x*t"},
	{"at", "a*t"},
	{"ta", "t*a"},
	{"aa", "a*a"},
	{"tt", "t*t"},
	{"xx", "x*x"},
	{"h(", "h*("},
	{"hx", "h*x"},
	{"xh", "x*h"},
	{"ah", "a*h"},
	{"ha", "h*a"},
	{"th", "t*h"},
	{"ht", "h*t"},
	{"+.", "+0."},
	{"-.", "-0."},
	{"*.", "*0."},
	{"/.", "/0."},
	{")(", ")*("},
	{"sqrt*(", "sqrt("},
	{"cot*(", "cot("},
	{"t*an(", "tan("},
	{"a*tan(", "atan("},
}

var colors = []string{"black", "red", "blue", "green", "purple", "brown", "orange"}
var lineColors = []string{"black", "red", "blue", "green", "purple", "brown", "orange", "yellow", "white"}
var MarbleRadius = .1
var Stop = false

// Graph represents the overall graph parameters -- lines and params
type Graph struct {
	Params Params `view:"inline" desc:"the parameters for updating the marbles"`
	Lines  Lines  `view:"-" desc:"the lines of the graph -- can have any number"`
}

// Gr is current graph
var Gr Graph

var KiT_Graph = kit.Types.AddType(&Graph{}, GraphProps)

// GraphProps define the ToolBar for overall app
var GraphProps = ki.Props{
	"ToolBar": ki.PropSlice{
		{Name: "OpenJSON", Value: ki.Props{
			"label": "Open...",
			"desc":  "Opens line equations and params from a .json file.",
			"icon":  "file-open",
			"Args": ki.PropSlice{
				{Name: "File Name", Value: ki.Props{
					"ext": ".json",
				}},
			},
		}},
		{Name: "SaveJSON", Value: ki.Props{
			"label": "Save As...",
			"desc":  "Saves line equations and params to a .json file.",
			"icon":  "file-save",
			"Args": ki.PropSlice{
				{Name: "File Name", Value: ki.Props{
					"ext": ".json",
				}},
			},
		}},
		{Name: "OpenAutoSave", Value: ki.Props{
			"label": "Open Autosaved...",
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
	},
}

func InitEquationChangeSlice() {
	for i := 0; i < 10; i++ {
		is := strconv.Itoa(i)
		EquationChangeSlice = append(EquationChangeSlice, EquationChange{is + "(", is + "*("}, EquationChange{is + "x", is + "*x"}, EquationChange{is + "t", is + "*t"}, EquationChange{is + "a", is + "*a"}, EquationChange{is + "h", is + "*h"})
	}
}
func (gr *Graph) Defaults() {
	gr.Params.Defaults()
	gr.Lines.Defaults()
}

// OpenJSON open from JSON file
func (gr *Graph) OpenJSON(filename gi.FileName) error {
	b, err := ioutil.ReadFile(string(filename))
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	gr.Graph()
	return err
}
func (gr *Graph) OpenAutoSave() error {
	b, err := ioutil.ReadFile("autosave.json")
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	gr.Graph()
	return err
}

// SaveJSON save to JSON file
func (gr *Graph) SaveJSON(filename gi.FileName) error {
	b, err := json.MarshalIndent(gr, "", "  ")
	if HandleError(err) {
		return err
	}
	err = ioutil.WriteFile(string(filename), b, 0644)
	HandleError(err)
	return err
}
func (gr *Graph) AutoSave() error {
	b, err := json.MarshalIndent(gr, "", "  ")
	if HandleError(err) {
		return err
	}
	err = ioutil.WriteFile("autosave.json", b, 0644)
	HandleError(err)
	return err
}

// Graph updates graph for current equations, and resets marbles too
func (gr *Graph) Graph() {
	gr.CompileExprs()
	ResetMarbles()
	gr.Params.Time = 0
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

func (gr *Graph) CompileExprs() {
	for _, ln := range gr.Lines {
		if ln.Expr.Expr == "" {
			ln.Expr.Expr = "x"
		}
		if ln.Color == "" {
			ln.Color = "black"
		}
		if ln.Bounce.Expr == "" {
			ln.Bounce.Expr = "0.95"
		}
		if ln.MinX.Expr == "" {
			ln.MinX.Expr = "-10"
		}
		if ln.MaxX.Expr == "" {
			ln.MaxX.Expr = "10"
		}
		if ln.MinY.Expr == "" {
			ln.MinY.Expr = "-10"
		}
		if ln.MaxY.Expr == "" {
			ln.MaxY.Expr = "10"
		}
		ln.TimesHit = 0
		ln.MakeGraphable()
		ln.Compile()
	}
}

func (ln *Line) MakeGraphable() { // Change things like 9x to 9*x to make the line graphable
	for _, d := range EquationChangeSlice {
		ln.Expr.Expr = strings.Replace(ln.Expr.Expr, d.Old, d.New, -1)
	}
}

func (ln *Line) Compile() {
	ln.Expr.Compile()
	ln.Bounce.Compile()
	ln.MinX.Compile()
	ln.MaxX.Compile()
	ln.MinY.Compile()
	ln.MaxY.Compile()
}

func (ln *Line) Defaults(lidx int) {
	ln.Expr.Expr = "x"
	ln.Color = colors[lidx%len(colors)]
	ln.Bounce.Expr = "0.95"
	ln.MinX.Expr = "-10"
	ln.MaxX.Expr = "10"
	ln.MinY.Expr = "-10"
	ln.MaxY.Expr = "10"
	ln.ColorSwitch = "none"
}

func (ls *Lines) Defaults() {
	*ls = make(Lines, 1, 10)
	ln := Line{}
	(*ls)[0] = &ln
	ln.Defaults(0)

}

// OpenJSON open from JSON file
func (ls *Lines) OpenJSON(filename gi.FileName) error {
	b, err := ioutil.ReadFile(string(filename))
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, ls)
	HandleError(err)
	return err
}

// SaveJSON save to JSON file
func (ls *Lines) SaveJSON(filename gi.FileName) error {
	b, err := json.MarshalIndent(ls, "", "  ")
	if HandleError(err) {
		return err
	}
	err = ioutil.WriteFile(string(filename), b, 0644)
	HandleError(err)
	return err
}

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

// Graph graphs this line in the SvgLines group
func (ln *Line) Graph(lidx int) {
	if ln.Expr.Expr == "" {
		ln.Defaults(lidx)
	}
	if ln.Color == "" || ln.Color == "black" {
		ln.Color = colors[lidx%len(colors)]
	}
	if ln.ColorSwitch == "" {
		ln.ColorSwitch = "none"
	}
	if ln.Bounce.Expr == "" {
		ln.Bounce.Expr = "0.95"
	}
	if ln.MinX.Expr == "" {
		ln.MinX.Expr = "-10"
	}
	if ln.MaxX.Expr == "0" {
		ln.MaxX.Expr = "10"
	}
	if ln.MinY.Expr == "0" {
		ln.MinY.Expr = "-10"
	}
	if ln.MaxY.Expr == "0" {
		ln.MaxY.Expr = "10"
	}
	path := SvgLines.Child(lidx).(*svg.Path)
	path.SetProp("fill", "none")
	clr := ln.Color
	clr = EvalColorIf(clr, ln.TimesHit)
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
func EvalColorIf(clr string, h int) string {
	if strings.Contains(clr, "if") {
		var err error

		expr := Expr{clr, nil, nil}
		expr.Val, err = govaluate.NewEvaluableExpressionWithFunctions(expr.Expr, functions)
		HandleError(err)
		expr.Params = make(map[string]interface{}, 2)
		expr.Params["h"] = float64(h)
		expr.Params["t"] = float64(Gr.Params.Time)
		expr.Params["a"] = float64(10 * math.Sin(float64(Gr.Params.Time)))
		for k, d := range lineColors {
			expr.Params[d] = float64(k)
		}
		yi, err := expr.Val.Evaluate(expr.Params)
		HandleError(err)
		yf := yi.(float64)
		for k, d := range lineColors {
			if yf == float64(k) {
				return d
			}
		}
	}
	return clr
}

func InitCoords() {
	updt := SvgGraph.UpdateStart()
	SvgCoords.DeleteChildren(true)

	xAxis := svg.AddNewLine(SvgCoords, "xAxis", -1000, 0, 1000, 0)
	xAxis.SetProp("stroke", "#888")

	yAxis := svg.AddNewLine(SvgCoords, "yAxis", 0, -1000, 0, 1000)
	yAxis.SetProp("stroke", "#888")

	SvgGraph.UpdateEnd(updt)
}

/////////////////////////////////////////////////////////////////////////

//  Marbles

type Marble struct {
	Pos    mat32.Vec2
	Vel    mat32.Vec2
	PrvPos mat32.Vec2
}

func (mb *Marble) Init(diff float32) {
	randNum := (rand.Float32() * 2) - 1
	xPos := randNum * Gr.Params.Width
	mb.Pos = mat32.Vec2{X: xPos, Y: Gr.Params.MaxSize.Y - diff}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	mb.Vel = mat32.Vec2{X: 0, Y: float32(-Gr.Params.StartSpeed)}
	mb.PrvPos = mb.Pos
}

var Marbles []*Marble

// Kai: put all these in a struct, and add a StructInlineView to edit them.
// see your other code for how to do it..

// Params holds our parameters
type Params struct {
	NMarbles   int     `min:"1" max:"10000" step:"10" desc:"number of marbles"`
	NSteps     int     `min:"100" max:"10000" step:"10" desc:"number of steps to take when running"`
	StartSpeed float32 `min:"0" max:"2" step:".05" desc:"Coordinates per unit of time"`
	UpdtRate   float32 `min:"0.001" max:"1" step:".01" desc:"how fast to move along velocity vector -- lower = smoother, more slow-mo"`
	Gravity    float32 `min:"0" max:"2" step:".01" desc:"how fast it accelerates down"`
	Width      float32 `min:"0" max:"10" step:"1" desc:"length of spawning zone for marbles, set to 0 for all spawn in a column"`
	TimeStep   float32 `min:"0.001" max:"100" step:".01" desc:"how fast time increases"`
	Time       float32 `view:"-" inactive:"+" desc:"time in msecs since starting"`
	MinSize    mat32.Vec2
	MaxSize    mat32.Vec2
}

func (pr *Params) Defaults() {
	pr.NMarbles = 10
	pr.NSteps = 10000
	pr.StartSpeed = 0
	pr.UpdtRate = .02
	pr.Gravity = 0.1
	pr.TimeStep = 0.01
	pr.MinSize = mat32.Vec2{X: -10, Y: -10}
	pr.MaxSize = mat32.Vec2{X: 10, Y: 10}
}

func RadToDeg(rad float32) float32 {
	return rad * 180 / math.Pi
}

// GraphMarblesInit initializes the graph drawing of the marbles
func GraphMarblesInit() {
	updt := SvgGraph.UpdateStart()

	SvgMarbles.DeleteChildren(true)
	for i, m := range Marbles {
		size := float32(MarbleRadius) * gsz.Y / 20
		// fmt.Printf("size: %v \n", size)
		circle := svg.AddNewCircle(SvgMarbles, "circle", m.Pos.X, m.Pos.Y, size)
		circle.SetProp("stroke", "none")
		circle.SetProp("fill", colors[i%len(colors)])
	}
	SvgGraph.UpdateEnd(updt)
}

// InitMarbles creates the marbles and puts them at their initial positions
func InitMarbles() {
	Marbles = make([]*Marble, 0)
	for n := 0; n < Gr.Params.NMarbles; n++ {
		diff := (gsz.Y / 20) * 2 * float32(n) / float32(Gr.Params.NMarbles)
		m := Marble{}
		m.Init(diff)
		Marbles = append(Marbles, &m)
	}
}

// ResetMarbles just calls InitMarbles and GraphMarblesInit
func ResetMarbles() {
	InitMarbles()
	GraphMarblesInit()
}

func UpdateMarbles() {
	wupdt := SvgGraph.TopUpdateStart()
	defer SvgGraph.TopUpdateEnd(wupdt)

	updt := SvgGraph.UpdateStart()
	defer SvgGraph.UpdateEnd(updt)
	SvgGraph.SetNeedsFullRender()

	Gr.Lines.Graph()

	for i, m := range Marbles {
		var setColor = "none"

		m.Vel.Y -= Gr.Params.Gravity * ((gsz.Y * gsz.X) / 400)
		updtrate := Gr.Params.UpdtRate
		npos := m.Pos.Add(m.Vel.MulScalar(updtrate))
		ppos := m.Pos

		for _, ln := range Gr.Lines {
			if ln.Expr.Val == nil {
				continue
			}

			yp := ln.Expr.Eval(m.Pos.X, Gr.Params.Time, ln.TimesHit)
			yn := ln.Expr.Eval(npos.X, Gr.Params.Time, ln.TimesHit)

			// fmt.Printf("y: %v npos: %v pos: %v\n", y, npos.Y, m.Pos.Y)
			MinX := ln.MinX.Eval(npos.X, Gr.Params.Time, ln.TimesHit)
			MaxX := ln.MaxX.Eval(npos.X, Gr.Params.Time, ln.TimesHit)
			MinY := ln.MinY.Eval(npos.X, Gr.Params.Time, ln.TimesHit)
			MaxY := ln.MaxY.Eval(npos.X, Gr.Params.Time, ln.TimesHit)
			if ((npos.Y < yn && m.Pos.Y >= yp) || (npos.Y > yn && m.Pos.Y <= yp)) && (npos.X < MaxX && npos.X > MinX) && (npos.Y < MaxY && npos.Y > MinY) {
				// fmt.Printf("Collided! Equation is: %v \n", ln.Eq)
				ln.TimesHit++
				setColor = EvalColorIf(ln.ColorSwitch, ln.TimesHit)

				dly := yn - yp // change in the lines y
				dx := npos.X - m.Pos.X

				var yi, xi float32

				if dx == 0 {

					xi = npos.X
					yi = yn

				} else {

					ml := dly / dx
					dmy := npos.Y - m.Pos.Y
					mm := dmy / dx

					xi = (npos.X*(ml-mm) + npos.Y - yn) / (ml - mm)
					yi = ln.Expr.Eval(xi, Gr.Params.Time, ln.TimesHit)
					//		fmt.Printf("xi: %v, yi: %v \n", xi, yi)
				}

				yl := ln.Expr.Eval(xi-.01, Gr.Params.Time, ln.TimesHit) // point to the left of x
				yr := ln.Expr.Eval(xi+.01, Gr.Params.Time, ln.TimesHit) // point to the right of x

				//slp := (yr - yl) / .02
				angLn := math32.Atan2(yr-yl, 0.02)
				angN := angLn + math.Pi/2 // + 90 deg

				angI := math32.Atan2(m.Vel.Y, m.Vel.X)
				angII := angI + math.Pi

				angNII := angN - angII
				angR := math.Pi + 2*angNII

				// fmt.Printf("angLn: %v  angN: %v  angI: %v  angII: %v  angNII: %v  angR: %v\n",
				// 	RadToDeg(angLn), RadToDeg(angN), RadToDeg(angI), RadToDeg(angII), RadToDeg(angNII), RadToDeg(angR))

				Bounce := ln.Bounce.Eval(npos.X, Gr.Params.Time, ln.TimesHit)

				nvx := Bounce * (m.Vel.X*math32.Cos(angR) - m.Vel.Y*math32.Sin(angR))
				nvy := Bounce * (m.Vel.X*math32.Sin(angR) + m.Vel.Y*math32.Cos(angR))

				m.Vel = mat32.Vec2{X: nvx, Y: nvy}

				m.Pos = mat32.Vec2{X: xi, Y: yi}

			}
		}

		m.PrvPos = ppos
		m.Pos = m.Pos.Add(m.Vel.MulScalar(Gr.Params.UpdtRate))

		circle := SvgMarbles.Child(i).(*svg.Circle)
		circle.Pos = m.Pos
		if setColor != "none" {
			circle.SetProp("fill", setColor)
		}

	}
}

func RunMarbles() {
	Stop = false
	startFrames := 0
	start := time.Now()
	for i := 0; i < Gr.Params.NSteps; i++ {
		UpdateMarbles()
		if time.Since(start).Milliseconds() >= 1000 {
			fpsText.SetText(fmt.Sprintf("FPS: %v", i-startFrames))
			start = time.Now()
			startFrames = i
		}

		Gr.Params.Time += Gr.Params.TimeStep
		if Stop {
			break
		}
	}
}
