package main

import (
	"fmt"
	"math"
	"math/rand"
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
	Color  gist.Color
	Track  bool
}

// Marbles contains all of the marbles
var Marbles []*Marble

// Whether marbles are currently being ran, used to prevent crashing with double click run marbles
var runningMarbles bool

// GraphMarblesInit initializes the graph drawing of the marbles
func GraphMarblesInit() {
	updt := svgGraph.UpdateStart()

	svgMarbles.DeleteChildren(true)
	for i, m := range Marbles {
		size := float32(TheSettings.MarbleSettings.MarbleSize) * gsz.Y / 20
		// fmt.Printf("size: %v \n", size)
		circle := svg.AddNewCircle(svgMarbles, "circle", m.Pos.X, m.Pos.Y, size)
		circle.SetProp("stroke", "none")
		if TheSettings.MarbleSettings.MarbleColor == "default" {
			circle.SetProp("fill", colors[i%len(colors)])
			m.Color, _ = gist.ColorFromName(colors[i%len(colors)])
		} else {
			circle.SetProp("fill", TheSettings.MarbleSettings.MarbleColor)
			m.Color, _ = gist.ColorFromName(TheSettings.MarbleSettings.MarbleColor)
		}
		circle.SetProp("fslr", 0)
		circle.SetProp("lpos", mat32.Vec2{X: m.Pos.X, Y: m.Pos.Y})
	}
	svgGraph.UpdateEnd(updt)
}

// Init makes a marble
func (m *Marble) Init(diff float32) {
	randNum := (rand.Float64() * 2) - 1
	xPos := randNum * Gr.Params.Width
	m.Pos = mat32.Vec2{X: float32(xPos), Y: Gr.Params.MaxSize.Y - diff}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	m.Vel = mat32.Vec2{X: 0, Y: float32(-Gr.Params.StartSpeed)}
	m.PrvPos = m.Pos
	m.Track = TheSettings.TrackingSettings.TrackByDefault
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

// UpdateMarbles calls update marbles graph and update marbles data
func UpdateMarbles() {
	UpdateMarblesData()
	UpdateMarblesGraph()
}

// UpdateMarblesGraph updates the graph of the marbles
func UpdateMarblesGraph() {
	if svgGraph.IsRendering() || svgGraph.IsUpdating() {
		return
	}
	wupdt := svgGraph.TopUpdateStart()
	defer svgGraph.TopUpdateEnd(wupdt)

	updt := svgGraph.UpdateStart()
	defer svgGraph.UpdateEnd(updt)
	svgGraph.SetNeedsFullRender()

	Gr.Lines.Graph(true)
	for i, m := range Marbles {

		circle := svgMarbles.Child(i).(*svg.Circle)
		circle.Pos = m.Pos
		circle.SetProp("fill", m.Color)
		m.UpdateTrackingLines(circle)

	}
}

// UpdateTrackingLines adds a tracking line for a marble, if needed
func (m *Marble) UpdateTrackingLines(circle *svg.Circle) {
	tls := TheSettings.TrackingSettings
	if Gr.Params.TrackingSettings.Override {
		tls = Gr.Params.TrackingSettings.TrackingSettings
	}
	if m.Track {
		fslr := circle.Prop("fslr").(int)
		if fslr <= 100/tls.Accuracy {
			circle.SetProp("fslr", fslr+1)
		} else {
			lpos := circle.Prop("lpos").(mat32.Vec2)
			circle.SetProp("fslr", 0)
			circle.SetProp("lpos", m.Pos)
			line := svg.AddNewLine(svgTrackingLines, "line", lpos.X, lpos.Y, m.Pos.X, m.Pos.Y)
			clr := tls.LineColor
			if clr == gist.White {
				switch circle.Prop("fill").(type) {
				case string:
					clr, _ = gist.ColorFromName(circle.Prop("fill").(string))
				case gist.Color:
					clr = circle.Prop("fill").(gist.Color)
				}

			}
			line.SetProp("stroke", clr)
		}
	}
}

// UpdateMarblesData updates marbles data
func UpdateMarblesData() {
	for _, m := range Marbles {

		m.Vel.Y -= float32(Gr.Params.Gravity) * ((gsz.Y * gsz.X) / 400)
		updtrate := float32(Gr.Params.UpdtRate)
		npos := m.Pos.Add(m.Vel.MulScalar(updtrate))
		ppos := m.Pos
		setColor := gist.White
		for _, ln := range Gr.Lines {
			if ln.Expr.Val == nil {
				continue
			}

			yp := ln.Expr.Eval(float64(m.Pos.X), Gr.Params.Time, ln.TimesHit)
			yn := ln.Expr.Eval(float64(npos.X), Gr.Params.Time, ln.TimesHit)

			if m.Collided(ln, npos, yp, yn) {
				ln.TimesHit++
				setColor = ln.LineColors.ColorSwitch
				m.Pos, m.Vel = m.CalcCollide(ln, npos, yp, yn)
				break
			}
		}

		m.PrvPos = ppos
		m.Pos = m.Pos.Add(m.Vel.MulScalar(float32(Gr.Params.UpdtRate)))
		if setColor != gist.White {
			m.Color = setColor
		}

	}
}

// Collided returns true if the marble has collided with the line, and false if the marble has not.
func (m *Marble) Collided(ln *Line, npos mat32.Vec2, yp, yn float64) bool {
	graphIf := ln.GraphIf.EvalBool(float64(npos.X), yn, Gr.Params.Time, ln.TimesHit)
	inBounds := npos.Y > gmin.Y && npos.Y < gmax.Y && npos.X > gmin.X && npos.X < gmax.X
	collided := (float64(npos.Y) < yn && float64(m.Pos.Y) >= yp) || (float64(npos.Y) > yn && float64(m.Pos.Y) <= yp)
	if collided && graphIf && inBounds {
		return true
	}
	return false
}

// CalcCollide calculates the new position and velocityof a marble after collision
func (m *Marble) CalcCollide(ln *Line, npos mat32.Vec2, yp, yn float64) (mat32.Vec2, mat32.Vec2) {
	dly := yn - yp // change in the lines y
	dx := npos.X - m.Pos.X

	var yi, xi float32

	if dx == 0 {

		xi = npos.X
		yi = float32(yn)

	} else {

		ml := float32(dly) / dx
		dmy := npos.Y - m.Pos.Y
		mm := dmy / dx

		xi = (npos.X*(ml-mm) + npos.Y - float32(yn)) / (ml - mm)
		yi = float32(ln.Expr.Eval(float64(xi), Gr.Params.Time, ln.TimesHit))
		//		fmt.Printf("xi: %v, yi: %v \n", xi, yi)
	}

	yl := ln.Expr.Eval(float64(xi)-.01, Gr.Params.Time, ln.TimesHit) // point to the left of x
	yr := ln.Expr.Eval(float64(xi)+.01, Gr.Params.Time, ln.TimesHit) // point to the right of x

	//slp := (yr - yl) / .02
	angLn := math32.Atan2(float32(yr-yl), 0.02)
	angN := angLn + math.Pi/2 // + 90 deg

	angI := math32.Atan2(m.Vel.Y, m.Vel.X)
	angII := angI + math.Pi

	angNII := angN - angII
	angR := math.Pi + 2*angNII

	Bounce := ln.Bounce.Eval(float64(npos.X), Gr.Params.Time, ln.TimesHit)

	nvx := float32(Bounce) * (m.Vel.X*math32.Cos(angR) - m.Vel.Y*math32.Sin(angR))
	nvy := float32(Bounce) * (m.Vel.X*math32.Sin(angR) + m.Vel.Y*math32.Cos(angR))

	vel := mat32.Vec2{X: nvx, Y: nvy}
	pos := mat32.Vec2{X: xi, Y: yi}

	return pos, vel
}

// RunMarbles runs the marbles for NSteps
func RunMarbles() {
	if runningMarbles {
		return
	}
	runningMarbles = true
	stop = false
	startFrames := 0
	trackingStartFrames := 0
	start := time.Now()
	nsteps := Gr.Params.NSteps
	tls := TheSettings.TrackingSettings
	if Gr.Params.TrackingSettings.Override {
		tls = Gr.Params.TrackingSettings.TrackingSettings
	}
	if nsteps == -1 {
		nsteps = 1000000000000
	}
	for i := 0; i < nsteps; i++ {
		for j := 0; j < TheSettings.NFramesPer-1; j++ {
			UpdateMarblesData()
			Gr.Params.Time += Gr.Params.TimeStep
		}
		UpdateMarbles()
		if time.Since(start).Milliseconds() >= 1000 {
			fpsText.SetText(fmt.Sprintf("FPS: %v", i-startFrames))
			start = time.Now()
			startFrames = i
		}
		usesTrackingLines := false
		for _, m := range Marbles {
			if m.Track {
				usesTrackingLines = true
				break
			}
		}
		if usesTrackingLines && (i-trackingStartFrames > tls.NTrackingFrames) {
			svgTrackingLines.DeleteChildren(true)
			trackingStartFrames = i
		}

		Gr.Params.Time += Gr.Params.TimeStep
		if stop {
			return
		}
	}
}

// Jump jumps n frames forward
func Jump(n int) {
	updt := svgGraph.UpdateStart()
	for i := 0; i < n; i++ {
		UpdateMarblesData()
		Gr.Params.Time += Gr.Params.TimeStep
	}
	Gr.Lines.Graph(true)
	UpdateMarbles()
	svgGraph.UpdateEnd(updt)
}
