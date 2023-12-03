package main

import (
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/goosi/events"
	"goki.dev/icons"
	"goki.dev/mat32/v2"
	"goki.dev/svg"
)

func (gr *Graph) TopAppBar(tb *gi.TopAppBar) {
	gi.DefaultTopAppBar(tb)
	giv.NewFuncButton(tb, gr.Graph).SetIcon(icons.ShowChart)
	giv.NewFuncButton(tb, gr.Run).SetIcon(icons.PlayArrow)
	giv.NewFuncButton(tb, gr.Stop)
	giv.NewFuncButton(tb, gr.Step)

	gi.NewSeparator(tb)
	giv.NewFuncButton(tb, gr.SelectNextMarble).SetText("Next marble").SetIcon(icons.ArrowForward)
	giv.NewFuncButton(tb, gr.StopSelecting).SetText("Unselect").SetIcon(icons.Close)
	giv.NewFuncButton(tb, gr.TrackSelectedMarble).SetText("Track").SetIcon(icons.PinDrop)
}

func (gr *Graph) MakeBasicElements(b *gi.Body) {
	sp := gi.NewSplits(b)

	lsp := gi.NewSplits(sp).SetDim(mat32.Y)

	lns := giv.NewTableView(lsp).SetSlice(&gr.Lines)
	lns.OnChange(func(e events.Event) {
		gr.Graph()
	})

	params := giv.NewStructView(lsp).SetStruct(&gr.Params)
	params.OnChange(func(e events.Event) {
		gr.Graph()
	})

	lsp.SetSplits(0.7, 0.3)

	gr.Objects.Graph = gi.NewSVG(sp)

	gr.Objects.SVG = gr.Objects.Graph.SVG
	gr.Objects.SVG.InvertY = true

	gr.Vectors.Min = mat32.Vec2{X: -graphViewBoxSize, Y: -graphViewBoxSize}
	gr.Vectors.Max = mat32.Vec2{X: graphViewBoxSize, Y: graphViewBoxSize}
	gr.Vectors.Size = gr.Vectors.Max.Sub(gr.Vectors.Min)
	var n float32 = 1.0 / float32(TheSettings.GraphInc)
	gr.Vectors.Inc = mat32.Vec2{X: n, Y: n}

	gr.Objects.Root = &gr.Objects.SVG.Root
	gr.Objects.Root.ViewBox.Min = gr.Vectors.Min
	gr.Objects.Root.ViewBox.Size = gr.Vectors.Size
	gr.Objects.Root.SetProp("stroke-width", "0.1dp")
	gr.Objects.Root.SetProp("fill", colors.Scheme.Surface)

	svg.NewCircle(gr.Objects.Root).SetRadius(50)

	sp.SetSplits(0.5, 0.5)

	gr.Objects.Lines = svg.NewGroup(gr.Objects.Root, "lines")
	gr.Objects.Marbles = svg.NewGroup(gr.Objects.Root, "marbles")
	gr.Objects.Coords = svg.NewGroup(gr.Objects.Root, "coords")
	gr.Objects.TrackingLines = svg.NewGroup(gr.Objects.Root, "tracking-lines")

	gr.Objects.Coords.SetProp("stroke-width", "0.05dp")
	gr.Objects.TrackingLines.SetProp("stroke-width", "0.05dp")
}
