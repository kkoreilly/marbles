package main

/*
// Marble contains the information of a marble
type Marble struct {
	Pos          mat32.Vec2
	Vel          mat32.Vec2
	PrvPos       mat32.Vec2
	Color        gist.Color
	TrackingInfo TrackingInfo
}

// TrackingInfo contains all of the tracking info for a marbles
type TrackingInfo struct {
	Track                 bool
	LastPos               mat32.Vec2
	FramesSinceLastUpdate int
	StartedTrackingAt     int
}

// GraphMarblesInit initializes the graph drawing of the marbles
func GraphMarblesInit() {
	updt := TheGraph.Objects.Graph.UpdateStart()

	TheGraph.Objects.Marbles.DeleteChildren(true)
	TheGraph.Objects.TrackingLines.DeleteChildren(true)
	for i, m := range TheGraph.Marbles {
		svg.AddNewGroup(TheGraph.Objects.TrackingLines, "tlm"+strconv.Itoa(i))
		size := float32(TheSettings.MarbleSettings.MarbleSize) * TheGraph.Vectors.Size.Y / 20
		// fmt.Printf("size: %v \n", size)
		circle := svg.AddNewCircle(TheGraph.Objects.Marbles, "circle", m.Pos.X, m.Pos.Y, size)
		circle.SetProp("stroke", "none")
		circle.SetProp("stroke-width", 4*TheSettings.MarbleSettings.MarbleSize)
		if TheSettings.MarbleSettings.MarbleColor == "default" {
			circle.SetProp("fill", colors[i%len(colors)])
			m.Color, _ = gist.ColorFromName(colors[i%len(colors)])
		} else {
			circle.SetProp("fill", TheSettings.MarbleSettings.MarbleColor)
			m.Color, _ = gist.ColorFromName(TheSettings.MarbleSettings.MarbleColor)
		}
		m.TrackingInfo.LastPos = mat32.Vec2{X: m.Pos.X, Y: m.Pos.Y}
		m.TrackingInfo.StartedTrackingAt = 0
	}
	TheGraph.Objects.Graph.UpdateEnd(updt)
}

// Init makes a marble
func (m *Marble) Init(n int) {
	// diff := (TheGraph.Vectors.Size.Y / 20) * 2 * float32(n) / float32(TheGraph.Params.NMarbles)
	SetRandNum()
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

	m.Pos = mat32.Vec2{X: float32(xPos), Y: float32(yPos)}
	// fmt.Printf("mb.Pos: %v \n", mb.Pos)
	startY := TheGraph.Params.StartVelY.Eval(float64(m.Pos.X), float64(m.Pos.Y))
	startX := TheGraph.Params.StartVelX.Eval(float64(m.Pos.X), float64(m.Pos.Y))
	m.Vel = mat32.Vec2{X: float32(startX), Y: float32(startY)}
	m.PrvPos = m.Pos
	tls := TheGraph.Params.TrackingSettings
	m.TrackingInfo.Track = tls.TrackByDefault
}

// InitMarbles creates the marbles and puts them at their initial positions
func InitMarbles() {
	TheGraph.Marbles = make([]*Marble, 0)
	for n := 0; n < TheGraph.Params.NMarbles; n++ {
		m := Marble{}
		m.Init(n)
		TheGraph.Marbles = append(TheGraph.Marbles, &m)
	}
	TheGraph.State.SelectedMarble = -1
}

// ResetMarbles just calls InitMarbles and GraphMarblesInit
func ResetMarbles() {
	InitMarbles()
	GraphMarblesInit()
}

// UpdateMarbles calls update marbles graph and update marbles data
func UpdateMarbles() bool {
	if !UpdateMarblesGraph() {
		UpdateMarblesData()
	} else {
		return true
	}
	return false
}

// UpdateMarblesGraph updates the graph of the marbles
func UpdateMarblesGraph() bool {
	if TheGraph.Objects.Graph.IsRendering() || TheGraph.Objects.Graph.IsUpdating() || vp.IsUpdatingNode() {
		return true
	}
	wupdt := TheGraph.Objects.Graph.TopUpdateStart()
	defer TheGraph.Objects.Graph.TopUpdateEnd(wupdt)

	if vp.IsUpdatingNode() {
		return true
	}

	updt := TheGraph.Objects.Graph.UpdateStart()
	defer TheGraph.Objects.Graph.UpdateEnd(updt)

	if vp.IsUpdatingNode() {
		return true
	}

	TheGraph.Objects.Graph.SetNeedsFullRender()

	TheGraph.Lines.Graph()
	for i, m := range TheGraph.Marbles {
		circle := TheGraph.Objects.Marbles.Child(i).(*svg.Circle)
		circle.Pos = m.Pos
		circle.SetProp("fill", m.Color)
		m.UpdateTrackingLines(circle, i)

	}
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
			svgGroup := TheGraph.Objects.TrackingLines.Child(idx)
			lpos := m.TrackingInfo.LastPos
			m.TrackingInfo.FramesSinceLastUpdate = 0
			m.TrackingInfo.LastPos = m.Pos
			if TheGraph.State.Step-m.TrackingInfo.StartedTrackingAt >= tls.NTrackingFrames {
				TheGraph.Objects.TrackingLines.Child(idx).DeleteChildAtIndex(0, true)
			}
			line := svg.AddNewLine(svgGroup, "line", lpos.X, lpos.Y, m.Pos.X, m.Pos.Y)
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
	for _, m := range TheGraph.Marbles {

		m.Vel.Y += float32(TheGraph.Params.YForce.Eval(float64(m.Pos.X), float64(m.Pos.Y))) * ((TheGraph.Vectors.Size.Y * TheGraph.Vectors.Size.X) / 400)
		m.Vel.X += float32(TheGraph.Params.XForce.Eval(float64(m.Pos.X), float64(m.Pos.Y))) * ((TheGraph.Vectors.Size.Y * TheGraph.Vectors.Size.X) / 400)
		updtrate := float32(TheGraph.Params.UpdtRate.Eval(float64(m.Pos.X), float64(m.Pos.Y)))
		npos := m.Pos.Add(m.Vel.MulScalar(updtrate))
		ppos := m.Pos
		setColor := gist.White
		for _, ln := range TheGraph.Lines {
			if ln.Expr.Val == nil {
				continue
			}

			yp := ln.Expr.Eval(float64(m.Pos.X), TheGraph.State.Time, ln.TimesHit)
			yn := ln.Expr.Eval(float64(npos.X), TheGraph.State.Time, ln.TimesHit)

			if m.Collided(ln, npos, yp, yn) {
				ln.TimesHit++
				setColor = ln.Colors.ColorSwitch
				m.Pos, m.Vel = m.CalcCollide(ln, npos, yp, yn)
				break
			}
		}

		m.PrvPos = ppos
		m.Pos = m.Pos.Add(m.Vel.MulScalar(float32(TheGraph.Params.UpdtRate.Eval(float64(m.Pos.X), float64(m.Pos.Y)))))
		if setColor != gist.White {
			m.Color = setColor
		}

	}
}

// Collided returns true if the marble has collided with the line, and false if the marble has not.
func (m *Marble) Collided(ln *Line, npos mat32.Vec2, yp, yn float64) bool {
	graphIf := ln.GraphIf.EvalBool(float64(npos.X), yn, TheGraph.State.Time, ln.TimesHit)
	inBounds := InBounds(npos)
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
		yi = float32(ln.Expr.Eval(float64(xi), TheGraph.State.Time, ln.TimesHit))
		//		fmt.Printf("xi: %v, yi: %v \n", xi, yi)
	}

	yl := ln.Expr.Eval(float64(xi)-.01, TheGraph.State.Time, ln.TimesHit) // point to the left of x
	yr := ln.Expr.Eval(float64(xi)+.01, TheGraph.State.Time, ln.TimesHit) // point to the right of x

	//slp := (yr - yl) / .02
	angLn := math32.Atan2(float32(yr-yl), 0.02)
	angN := angLn + math.Pi/2 // + 90 deg

	angI := math32.Atan2(m.Vel.Y, m.Vel.X)
	angII := angI + math.Pi

	angNII := angN - angII
	angR := math.Pi + 2*angNII

	Bounce := ln.Bounce.EvalWithY(float64(npos.X), TheGraph.State.Time, ln.TimesHit, float64(yi))

	nvx := float32(Bounce) * (m.Vel.X*math32.Cos(angR) - m.Vel.Y*math32.Sin(angR))
	nvy := float32(Bounce) * (m.Vel.X*math32.Sin(angR) + m.Vel.Y*math32.Cos(angR))

	vel := mat32.Vec2{X: nvx, Y: nvy}
	pos := mat32.Vec2{X: xi, Y: yi}

	return pos, vel
}

// InBounds checks whether a point is in the bounds of the graph
func InBounds(pos mat32.Vec2) bool {
	if pos.Y > TheGraph.Vectors.Min.Y && pos.Y < TheGraph.Vectors.Max.Y && pos.X > TheGraph.Vectors.Min.X && pos.X < TheGraph.Vectors.Max.X {
		return true
	}
	return false
}

// RunMarbles runs the marbles for NSteps
func RunMarbles() {
	if TheGraph.State.Running {
		return
	}
	TheGraph.State.Running = true
	startFrames := 0
	start := time.Now()
	for TheGraph.State.Step = 0; TheGraph.State.Running; TheGraph.State.Step++ {
		if TheGraph.State.Error != nil {
			TheGraph.State.Running = false
		}
		for j := 0; j < TheSettings.NFramesPer-1; j++ {
			UpdateMarblesData()
			TheGraph.State.Time += TheGraph.Params.TimeStep.Eval(0, 0)
		}
		if UpdateMarbles() {
			TheGraph.State.Step--
			continue
		}
		if time.Since(start).Milliseconds() >= 3000 {
			fpsText.SetText(fmt.Sprintf("FPS: %v", (TheGraph.State.Step-startFrames)/3))
			start = time.Now()
			startFrames = TheGraph.State.Step
		}
		TheGraph.State.Time += TheGraph.Params.TimeStep.Eval(0, 0)
	}
}

// ToggleTrack toogles tracking setting for a certain marble
func (m *Marble) ToggleTrack(idx int) {
	m.TrackingInfo.Track = !m.TrackingInfo.Track
	TheGraph.Objects.TrackingLines.Child(idx).DeleteChildren(true)
	m.TrackingInfo.FramesSinceLastUpdate = 0
	m.TrackingInfo.LastPos = mat32.Vec2{X: m.Pos.X, Y: m.Pos.Y}
	m.TrackingInfo.StartedTrackingAt = TheGraph.State.Step
}

// SelectNextMarble selects the next marble in the viewbox
func SelectNextMarble() {
	if !TheGraph.State.Running {
		updt := TheGraph.Objects.Graph.UpdateStart()
		defer TheGraph.Objects.Graph.UpdateEnd(updt)
	}
	if TheGraph.State.SelectedMarble != -1 {
		TheGraph.Objects.Marbles.Child(TheGraph.State.SelectedMarble).SetProp("stroke", "none")
	}
	TheGraph.State.SelectedMarble++
	if TheGraph.State.SelectedMarble >= len(TheGraph.Marbles) {
		TheGraph.State.SelectedMarble = 0
	}
	newMarble := TheGraph.Marbles[TheGraph.State.SelectedMarble]
	if !InBounds(newMarble.Pos) { // If the marble isn't in bounds, don't select it
		for _, m := range TheGraph.Marbles { // If all marbles aren't in bounds, do nothing
			if InBounds(m.Pos) {
				SelectNextMarble()
				return
			}
		}
		return

	}
	TheGraph.Objects.Marbles.Child(TheGraph.State.SelectedMarble).SetProp("stroke", "yellow")
	TheGraph.Objects.Graph.SetNeedsFullRender()
}
*/
