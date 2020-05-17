// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"math/rand"

	"github.com/Knetic/govaluate"
	"github.com/chewxy/math32"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

var colors = []string{"black", "red", "blue", "green", "purple", "brown", "orange"}

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
		{"OpenJSON", ki.Props{
			"label": "Open...",
			"desc":  "Opens line equations and params from a .json file.",
			"icon":  "file-open",
			"Args": ki.PropSlice{
				{"File Name", ki.Props{
					"ext": ".json",
				}},
			},
		}},
		{"SaveJSON", ki.Props{
			"label": "Save As...",
			"desc":  "Saves line equations and params to a .json file.",
			"icon":  "file-save",
			"Args": ki.PropSlice{
				{"File Name", ki.Props{
					"ext": ".json",
				}},
			},
		}},
		{"sep-ctrl", ki.BlankProp{}},
		{"Graph", ki.Props{
			"desc": "updates graph for current equations",
			"icon": "file-image",
		}},
		{"Run", ki.Props{
			"desc":            "runs the marbles for NSteps",
			"icon":            "run",
			"no-update-after": true,
		}},
		{"Stop", ki.Props{
			"desc":            "runs the marbles for NSteps",
			"icon":            "stop",
			"no-update-after": true,
		}},
		{"Step", ki.Props{
			"desc":            "steps the marbles for one step",
			"icon":            "step-fwd",
			"no-update-after": true,
		}},
	},
}

func (gr *Graph) Defaults() {
	gr.Params.Defaults()
	gr.Lines.Defaults()
}

// OpenJSON open from JSON file
func (gr *Graph) OpenJSON(filename gi.FileName) error {
	b, err := ioutil.ReadFile(string(filename))
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	err = json.Unmarshal(b, gr)
	gr.Graph()
	return err
}

// SaveJSON save to JSON file
func (gr *Graph) SaveJSON(filename gi.FileName) error {
	b, err := json.MarshalIndent(gr, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(string(filename), b, 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

// Graph updates graph for current equations, and resets marbles too
func (gr *Graph) Graph() {
	gr.CompileExprs()
	ResetMarbles()
	gr.Params.Time = 0
	gr.Lines.Graph()
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
		if &ln.Expr == nil {
			ln.Expr.Expr = "x"
		}
		if ln.Color == "" {
			ln.Color = "black"
		}
		if &ln.Bounce == nil || ln.Bounce.Expr == "" {
			ln.Bounce.Expr = "0.95"
		}
		if &ln.MinX == nil || ln.MinX.Expr == "" {
			ln.MinX.Expr = "-10"
		}
		if &ln.MaxX == nil || ln.MaxX.Expr == "" {
			ln.MaxX.Expr = "10"
		}
		if &ln.MinY == nil || ln.MinY.Expr == "" {
			ln.MinY.Expr = "-10"
		}
		if &ln.MaxY == nil || ln.MaxY.Expr == "" {
			ln.MaxY.Expr = "10"
		}
		ln.Compile()
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

///////////////////////////////////////////////////////////////////////////
//  Lines

// Line represents one line with an equation etc
type Line struct {
	Expr   Expr   `width:"60" label:"y=" desc:"equation: use 'x' for the x value, and must use * for multiplication, and start with 0 for decimal numbers (0.01 instead of .01)"`
	MinX   Expr   `width:"30" step:"1" desc:"Minimum x value for this line."`
	MaxX   Expr   `step:"1" desc:"Maximum x value for this line."`
	MinY   Expr   `step:"1" desc:"Minimum y value for this line."`
	MaxY   Expr   `step:"1" desc:"Maximum y value for this line."`
	Color  string `desc:"color to draw the line in"`
	Bounce Expr   `min:"0" max:"2" step:".05" desc:"how bouncy the line is -- 1 = perfectly bouncy, 0 = no bounce at all"`
}

func (ln *Line) Defaults(lidx int) {
	ln.Expr.Expr = "x"
	ln.Color = colors[lidx%len(colors)]
	ln.Bounce.Expr = "0.95"
	ln.MinX.Expr = "-10"
	ln.MaxX.Expr = "10"
	ln.MinY.Expr = "-10"
	ln.MaxY.Expr = "10"
}

// Eval gives the y value of the function for given x value

// Lines is a collection of lines
type Lines []*Line

var KiT_Lines = kit.Types.AddType(&Lines{}, LinesProps)

// LinesProps define the ToolBar for lines
var LinesProps = ki.Props{
	// "ToolBar": ki.PropSlice{
	// 	{"OpenJSON", ki.Props{
	// 		"label": "Open...",
	// 		"desc":  "opens equations from a .json file.",
	// 		"icon":  "file-open",
	// 		"Args": ki.PropSlice{
	// 			{"File Name", ki.Props{
	// 				"ext": ".json",
	// 			}},
	// 		},
	// 	}},
	// 	{"SaveJSON", ki.Props{
	// 		"label": "Save As...",
	// 		"desc":  "Saves equations from a .json file.",
	// 		"icon":  "file-save",
	// 		"Args": ki.PropSlice{
	// 			{"File Name", ki.Props{
	// 				"ext": ".json",
	// 			}},
	// 		},
	// 	}},
	// },
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
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, ls)
	return err
}

// SaveJSON save to JSON file
func (ls *Lines) SaveJSON(filename gi.FileName) error {
	b, err := json.MarshalIndent(ls, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(string(filename), b, 0644)
	if err != nil {
		log.Println(err)
	}
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
	if ln.Color == "" {
		ln.Color = colors[lidx%len(colors)]
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
	path.SetProp("stroke", clr)

	var err error
	ln.Expr.Val, err = govaluate.NewEvaluableExpressionWithFunctions(ln.Expr.Expr, functions)
	if err != nil {
		ln.Expr.Val = nil
		log.Println(err)
		return
	}

	ps := ""
	start := true
	for x := gmin.X; x < gmax.X; x += ginc.X {
		MinX := ln.MinX.Eval(x, Gr.Params.Time)
		MaxX := ln.MaxX.Eval(x, Gr.Params.Time)
		MinY := ln.MinY.Eval(x, Gr.Params.Time)
		MaxY := ln.MaxY.Eval(x, Gr.Params.Time)
		y := ln.Expr.Eval(x, Gr.Params.Time)
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
	mb.Pos = mat32.Vec2{xPos, Gr.Params.MaxSize.Y - diff}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	mb.Vel = mat32.Vec2{0, float32(-Gr.Params.StartSpeed)}
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
	Width      float32 `length of spawning zone for marbles, set to 0 for all spawn in a column`
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
	pr.MinSize = mat32.Vec2{-10, -10}
	pr.MaxSize = mat32.Vec2{10, 10}
}

var MarbleRadius = .1

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

		m.Vel.Y -= Gr.Params.Gravity * ((gsz.Y * gsz.X) / 400)
		updtrate := Gr.Params.UpdtRate
		npos := m.Pos.Add(m.Vel.MulScalar(updtrate))
		ppos := m.Pos

		for _, ln := range Gr.Lines {
			if ln.Expr.Val == nil {
				continue
			}

			yp := ln.Expr.Eval(m.Pos.X, Gr.Params.Time)
			yn := ln.Expr.Eval(npos.X, Gr.Params.Time)

			// fmt.Printf("y: %v npos: %v pos: %v\n", y, npos.Y, m.Pos.Y)
			MinX := ln.MinX.Eval(npos.X, Gr.Params.Time)
			MaxX := ln.MaxX.Eval(npos.X, Gr.Params.Time)
			MinY := ln.MinY.Eval(npos.X, Gr.Params.Time)
			MaxY := ln.MaxY.Eval(npos.X, Gr.Params.Time)
			if ((npos.Y < yn && m.Pos.Y >= yp) || (npos.Y > yn && m.Pos.Y <= yp)) && (npos.X < MaxX && npos.X > MinX) && (npos.Y < MaxY && npos.Y > MinY) {
				// fmt.Printf("Collided! Equation is: %v \n", ln.Eq)

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
					yi = ln.Expr.Eval(xi, Gr.Params.Time)
					//		fmt.Printf("xi: %v, yi: %v \n", xi, yi)
				}

				yl := ln.Expr.Eval(xi-.01, Gr.Params.Time) // point to the left of x
				yr := ln.Expr.Eval(xi+.01, Gr.Params.Time) // point to the right of x

				//slp := (yr - yl) / .02
				angLn := math32.Atan2(yr-yl, 0.02)
				angN := angLn + math.Pi/2 // + 90 deg

				angI := math32.Atan2(m.Vel.Y, m.Vel.X)
				angII := angI + math.Pi

				angNII := angN - angII
				angR := math.Pi + 2*angNII

				// fmt.Printf("angLn: %v  angN: %v  angI: %v  angII: %v  angNII: %v  angR: %v\n",
				// 	RadToDeg(angLn), RadToDeg(angN), RadToDeg(angI), RadToDeg(angII), RadToDeg(angNII), RadToDeg(angR))

				Bounce := ln.Bounce.Eval(npos.X, Gr.Params.Time)

				nvx := Bounce * (m.Vel.X*math32.Cos(angR) - m.Vel.Y*math32.Sin(angR))
				nvy := Bounce * (m.Vel.X*math32.Sin(angR) + m.Vel.Y*math32.Cos(angR))

				m.Vel = mat32.Vec2{nvx, nvy}

				m.Pos = mat32.Vec2{xi, yi}

			}
		}

		m.PrvPos = ppos
		m.Pos = m.Pos.Add(m.Vel.MulScalar(Gr.Params.UpdtRate))

		circle := SvgMarbles.Child(i).(*svg.Circle)
		circle.Pos = m.Pos
	}
}

var Stop = false

func RunMarbles() {
	Stop = false
	for i := 0; i < Gr.Params.NSteps; i++ {
		UpdateMarbles()
		Gr.Params.Time += Gr.Params.TimeStep
		if Stop {
			break
		}
	}
}
