package main

import (
	"fmt"
	"math"
	"time"

	"github.com/chewxy/math32"
	"github.com/goki/gi/gist"
	"github.com/goki/gi/svg"
	"github.com/goki/mat32"
)

// Marble contains the information of a marble
type Marble struct {
	Pos    mat32.Vec2
	Vel    mat32.Vec2
	PrvPos mat32.Vec2
}

// Marbles contains all of the marbles
var Marbles []*Marble

// GraphMarblesInit initializes the graph drawing of the marbles
func GraphMarblesInit() {
	updt := SvgGraph.UpdateStart()

	SvgMarbles.DeleteChildren(true)
	for i, m := range Marbles {
		size := float32(TheSettings.MarbleSettings.MarbleSize) * gsz.Y / 20
		// fmt.Printf("size: %v \n", size)
		circle := svg.AddNewCircle(SvgMarbles, "circle", m.Pos.X, m.Pos.Y, size)
		circle.SetProp("stroke", "none")
		if TheSettings.MarbleSettings.MarbleColor == "default" {
			circle.SetProp("fill", colors[i%len(colors)])
		} else {
			circle.SetProp("fill", TheSettings.MarbleSettings.MarbleColor)
		}
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

// Update the marbles for one step
func UpdateMarbles() {
	wupdt := SvgGraph.TopUpdateStart()
	defer SvgGraph.TopUpdateEnd(wupdt)

	updt := SvgGraph.UpdateStart()
	defer SvgGraph.UpdateEnd(updt)
	SvgGraph.SetNeedsFullRender()

	Gr.Lines.Graph()
	white, _ := gist.ColorFromName("white")

	for i, m := range Marbles {
		var setColor = white

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
				// setColor = EvalColorIf(ln.ColorSwitch, ln.TimesHit)
				setColor = ln.LineColors.ColorSwitch

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

		if setColor != white {
			circle.SetProp("fill", setColor)
		}

	}
}

// Run the marbles for NSteps
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

// Turn radians to degrees
func RadToDeg(rad float32) float32 {
	return rad * 180 / math.Pi
}
