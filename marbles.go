package main

import (
	"image/color"
	"math"
	"time"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/svg"
)

// Marble contains the information of a marble
type Marble struct {
	Pos          math32.Vector2
	Velocity     math32.Vector2
	PrevPos      math32.Vector2
	Color        color.RGBA
	TrackingInfo TrackingInfo
}

// TrackingInfo contains all of the tracking info for a marbles
type TrackingInfo struct {
	Track                 bool
	LastPos               math32.Vector2
	FramesSinceLastUpdate int
	StartedTrackingAt     int
}

// GraphMarblesInit initializes the graph drawing of the marbles
func (gr *Graph) GraphMarblesInit() {
	// gr.Objects.Marbles.DeleteChildren()
	// gr.Objects.TrackingLines.DeleteChildren()
	for i, m := range gr.Marbles {
		// svg.NewGroup(gr.Objects.TrackingLines)
		// size := float32(TheSettings.MarbleSettings.MarbleSize) * gr.Vectors.Size.Y / 20
		// fmt.Printf("size: %v \n", size)
		// circle := svg.NewCircle(gr.Objects.Marbles).SetPos(m.Pos).SetRadius(size)
		// circle.SetProperty("stroke", "none")
		// circle.SetProperty("stroke-width", 4*TheSettings.MarbleSettings.MarbleSize)
		if TheSettings.MarbleSettings.MarbleColor == "default" {
			m.Color = colors.Spaced(i)
			// circle.SetProperty("fill", m.Color)
		} else {
			m.Color = errors.Log1(colors.FromName(TheSettings.MarbleSettings.MarbleColor))
			// circle.SetProperty("fill", TheSettings.MarbleSettings.MarbleColor)
		}
		m.TrackingInfo.LastPos = math32.Vector2{X: m.Pos.X, Y: m.Pos.Y}
		m.TrackingInfo.StartedTrackingAt = 0
	}
}

// Init makes a marble
func (m *Marble) Init(n int) {
	// diff := (TheGraph.Vectors.Size.Y / 20) * 2 * float32(n) / float32(TheGraph.Params.NMarbles)
	if TheGraph.Params.MarbleStartX.Compile() != nil {
		return
	}
	TheGraph.Params.MarbleStartX.Params["n"] = n
	xPos := TheGraph.Params.MarbleStartX.Eval(0, 0, 0)

	if TheGraph.Params.MarbleStartY.Compile() != nil {
		return
	}
	TheGraph.Params.MarbleStartY.Params["n"] = n
	yPos := TheGraph.Params.MarbleStartY.Eval(xPos, 0, 0)

	m.Pos = math32.Vector2{X: float32(xPos), Y: float32(yPos)}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	startY := TheGraph.Params.StartVelocityY.Eval(float64(m.Pos.X), float64(m.Pos.Y))
	startX := TheGraph.Params.StartVelocityX.Eval(float64(m.Pos.X), float64(m.Pos.Y))
	m.Velocity = math32.Vector2{X: float32(startX), Y: float32(startY)}
	m.PrevPos = m.Pos
	tls := TheGraph.Params.TrackingSettings
	m.TrackingInfo.Track = tls.TrackByDefault
}

// InitMarbles creates the marbles and puts them at their initial positions
func (gr *Graph) InitMarbles() {
	gr.Marbles = make([]*Marble, 0)
	for n := 0; n < gr.Params.NMarbles; n++ {
		m := Marble{}
		m.Init(n)
		gr.Marbles = append(gr.Marbles, &m)
	}
	gr.State.SelectedMarble = -1
}

// ResetMarbles just calls InitMarbles and GraphMarblesInit
func (gr *Graph) ResetMarbles() {
	gr.InitMarbles()
	gr.GraphMarblesInit()
}

// UpdateMarbles calls update marbles graph and update marbles data
func (gr *Graph) UpdateMarbles() bool {
	if !gr.UpdateMarblesGraph() {
		gr.UpdateMarblesData()
		return false
	}
	return true
}

// UpdateMarblesGraph updates the graph of the marbles
func (gr *Graph) UpdateMarblesGraph() bool {
	// for i, m := range gr.Marbles {
	// 	circle := gr.Objects.Marbles.Child(i).(*svg.Circle)
	// 	circle.Pos = m.Pos
	// 	circle.SetProperty("fill", m.Color)
	// 	m.UpdateTrackingLines(circle, i)
	// }

	gr.Objects.Graph.NeedsRender()
	return false
}

// UpdateTrackingLines adds a tracking line for a marble, if needed
func (m *Marble) UpdateTrackingLines(circle *svg.Circle, idx int) {
	tls := TheGraph.Params.TrackingSettings
	if m.TrackingInfo.Track {
		fslu := m.TrackingInfo.FramesSinceLastUpdate
		if fslu <= 100/tls.Accuracy {
			m.TrackingInfo.FramesSinceLastUpdate++
		} else {
			// svgGroup := TheGraph.Objects.TrackingLines.Child(idx)
			// lpos := m.TrackingInfo.LastPos
			// m.TrackingInfo.FramesSinceLastUpdate = 0
			// m.TrackingInfo.LastPos = m.Pos
			// if TheGraph.State.Step-m.TrackingInfo.StartedTrackingAt >= tls.NTrackingFrames {
			// 	TheGraph.Objects.TrackingLines.Child(idx).AsTree().DeleteChildAt(0)
			// }
			// line := svg.NewLine(svgGroup).SetStart(lpos).SetEnd(m.Pos)
			// clr := tls.LineColor
			// if clr == colors.White {
			// 	clr = errors.Log1(colors.FromAny(circle.Property("fill"), colors.White))
			// }
			// line.SetProperty("stroke", clr)
		}
	}
}

// UpdateMarblesData updates marbles data
func (gr *Graph) UpdateMarblesData() {
	gr.EvalMu.Lock()
	defer gr.EvalMu.Unlock()

	for _, m := range gr.Marbles {

		m.Velocity.Y += float32(gr.Params.YForce.Eval(float64(m.Pos.X), float64(m.Pos.Y))) * ((gr.Vectors.Size.Y * gr.Vectors.Size.X) / 400)
		m.Velocity.X += float32(gr.Params.XForce.Eval(float64(m.Pos.X), float64(m.Pos.Y))) * ((gr.Vectors.Size.Y * gr.Vectors.Size.X) / 400)
		updtrate := float32(gr.Params.UpdateRate.Eval(float64(m.Pos.X), float64(m.Pos.Y)))
		npos := m.Pos.Add(m.Velocity.MulScalar(updtrate))
		ppos := m.Pos
		setColor := colors.White
		for _, ln := range gr.Lines {
			if ln.Expr.Val == nil {
				continue
			}

			// previous line y (with old time)
			yp := ln.Expr.Eval(float64(m.Pos.X), gr.State.PrevTime, ln.TimesHit)
			// new line y with old time
			yno := ln.Expr.Eval(float64(npos.X), gr.State.PrevTime, ln.TimesHit)
			// new line y
			yn := ln.Expr.Eval(float64(npos.X), gr.State.Time, ln.TimesHit)

			if m.Collided(ln, npos, yp, yn) {
				ln.TimesHit++
				setColor = ln.Colors.ColorSwitch
				m.Pos, m.Velocity = m.CalcCollide(ln, npos, yp, yn, yno)
				break
			}
		}

		m.PrevPos = ppos
		m.Pos = m.Pos.Add(m.Velocity.MulScalar(float32(gr.Params.UpdateRate.Eval(float64(m.Pos.X), float64(m.Pos.Y)))))
		if setColor != colors.White {
			m.Color = setColor
		}
	}
}

// Collided returns true if the marble has collided with the line, and false if the marble has not.
func (m *Marble) Collided(ln *Line, npos math32.Vector2, yp, yn float64) bool {
	graphIf := ln.GraphIf.EvalBool(float64(npos.X), yn, TheGraph.State.Time, ln.TimesHit)
	inBounds := TheGraph.InBounds(npos)
	collided := (float64(npos.Y) < yn && float64(m.Pos.Y) >= yp) || (float64(npos.Y) > yn && float64(m.Pos.Y) <= yp)
	if collided && graphIf && inBounds {
		return true
	}
	return false
}

// CalcCollide calculates the new position and velocity of a marble after a collision with the coreen
// line, coreen the previous line y, new line y, and new line y with old time
func (m *Marble) CalcCollide(ln *Line, npos math32.Vector2, yp, yn, yno float64) (math32.Vector2, math32.Vector2) {
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
		yi = float32(ln.Expr.Eval(float64(xi), TheGraph.State.Time, ln.TimesHit))
		//		fmt.Printf("xi: %v, yi: %v \n", xi, yi)
	}

	yl := ln.Expr.Eval(float64(xi)-.01, TheGraph.State.Time, ln.TimesHit) // point to the left of x
	yr := ln.Expr.Eval(float64(xi)+.01, TheGraph.State.Time, ln.TimesHit) // point to the right of x

	//slp := (yr - yl) / .02
	angLn := math32.Atan2(float32(yr-yl), 0.02)
	angN := angLn + math.Pi/2 // + 90 deg

	angI := math32.Atan2(m.Velocity.Y, m.Velocity.X)
	angII := angI + math.Pi

	angNII := angN - angII
	angR := math.Pi + 2*angNII

	Bounce := ln.Bounce.EvalWithY(float64(npos.X), TheGraph.State.Time, ln.TimesHit, float64(yi))

	nvx := float32(Bounce) * (m.Velocity.X*math32.Cos(angR) - m.Velocity.Y*math32.Sin(angR))
	nvy := float32(Bounce) * (m.Velocity.X*math32.Sin(angR) + m.Velocity.Y*math32.Cos(angR))

	vel := math32.Vector2{X: nvx, Y: nvy}
	pos := math32.Vector2{X: xi, Y: yi + float32(yn-yno)} // adding change from prev time to current time in same pos fixes collisions with moving lines

	return pos, vel
}

// InBounds checks whether a point is in the bounds of the graph
func (gr *Graph) InBounds(pos math32.Vector2) bool {
	if pos.Y > gr.Vectors.Min.Y && pos.Y < gr.Vectors.Max.Y && pos.X > gr.Vectors.Min.X && pos.X < gr.Vectors.Max.X {
		return true
	}
	return false
}

// RunMarbles runs the marbles for NSteps
func (gr *Graph) RunMarbles() {
	if gr.State.Running {
		return
	}
	gr.State.Running = true
	gr.State.Step = 0
	startFrames := 0
	start := time.Now()
	ticker := time.NewTicker(time.Second / 60)
	for range ticker.C {
		if !gr.State.Running {
			ticker.Stop()
			return
		}
		gr.State.Step++
		if gr.State.Error != nil {
			gr.State.Running = false
		}
		for j := 0; j < TheSettings.NFramesPer-1; j++ {
			gr.UpdateMarblesData()
			gr.State.PrevTime = gr.State.Time
			gr.State.Time += gr.Params.TimeStep.Eval(0, 0)
		}
		gr.Objects.Graph.AsyncLock()
		ok := gr.UpdateMarbles()
		gr.Objects.Graph.AsyncUnlock()
		if ok {
			gr.State.Step--
			continue
		}
		if time.Since(start).Milliseconds() >= 3000 {
			_ = startFrames
			// fpsText.SetText(fmt.Sprintf("FPS: %v", (gr.State.Step-startFrames)/3))
			start = time.Now()
			startFrames = gr.State.Step
		}
		gr.State.PrevTime = gr.State.Time
		gr.State.Time += gr.Params.TimeStep.Eval(0, 0)
	}
}

// ToggleTrack toogles tracking setting for a certain marble
func (m *Marble) ToggleTrack(idx int) {
	m.TrackingInfo.Track = !m.TrackingInfo.Track
	// TheGraph.Objects.TrackingLines.Child(idx).AsTree().DeleteChildren()
	m.TrackingInfo.FramesSinceLastUpdate = 0
	m.TrackingInfo.LastPos = math32.Vector2{X: m.Pos.X, Y: m.Pos.Y}
	m.TrackingInfo.StartedTrackingAt = TheGraph.State.Step
}

// SelectNextMarble selects the next marble in the viewbox
func (gr *Graph) SelectNextMarble() { //types:add
	if !gr.State.Running {
		defer gr.Objects.Graph.NeedsRender()
	}
	if gr.State.SelectedMarble != -1 {
		// gr.Objects.Marbles.Child(gr.State.SelectedMarble).AsTree().SetProperty("stroke", "none")
	}
	gr.State.SelectedMarble++
	if gr.State.SelectedMarble >= len(gr.Marbles) {
		gr.State.SelectedMarble = 0
	}
	newMarble := gr.Marbles[gr.State.SelectedMarble]
	if !gr.InBounds(newMarble.Pos) { // If the marble isn't in bounds, don't select it
		for _, m := range gr.Marbles { // If all marbles aren't in bounds, do nothing
			if gr.InBounds(m.Pos) {
				gr.SelectNextMarble()
				return
			}
		}
		return

	}
	// gr.Objects.Marbles.Child(gr.State.SelectedMarble).AsTree().SetProperty("stroke", "yellow")
}
