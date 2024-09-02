package main

import (
	"strconv"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/svg"
	"cogentcore.org/core/tree"
)

func (gr *Graph) MakeToolbar(p *tree.Plan) {
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.Graph).SetIcon(icons.ShowChart)
	})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.Run).SetIcon(icons.PlayArrow)
	})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.Stop).SetIcon(icons.Stop)
	})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.Step).SetIcon(icons.Step)
	})

	tree.Add(p, func(w *core.Separator) {})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.AddLine).SetIcon(icons.Add)
	})

	tree.Add(p, func(w *core.Separator) {})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.SelectNextMarble).SetText("Next marble").SetIcon(icons.ArrowForward)
	})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.StopSelecting).SetText("Unselect").SetIcon(icons.Close)
	})
	tree.Add(p, func(w *core.FuncButton) {
		w.SetFunc(gr.TrackSelectedMarble).SetText("Track").SetIcon(icons.PinDrop)
	})
}

func (gr *Graph) MakeBasicElements(b *core.Body) {
	sp := core.NewSplits(b)

	lsp := core.NewSplits(sp)
	lsp.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	gr.Objects.LinesTable = core.NewTable(lsp).SetSlice(&gr.Lines)
	gr.Objects.LinesTable.OnChange(func(e events.Event) {
		gr.Graph()
	})

	params := core.NewForm(lsp).SetStruct(&gr.Params)
	params.OnChange(func(e events.Event) {
		gr.Graph()
	})

	lsp.SetSplits(0.7, 0.3)

	gr.Objects.Graph = core.NewSVG(sp)
	gr.Objects.Graph.SetReadOnly(false)

	gr.Objects.SVG = gr.Objects.Graph.SVG
	gr.Objects.SVG.InvertY = true

	gr.Vectors.Min = math32.Vector2{X: -GraphViewBoxSize, Y: -GraphViewBoxSize}
	gr.Vectors.Max = math32.Vector2{X: GraphViewBoxSize, Y: GraphViewBoxSize}
	gr.Vectors.Size = gr.Vectors.Max.Sub(gr.Vectors.Min)
	var n float32 = 1.0 / float32(TheSettings.GraphInc)
	gr.Vectors.Inc = math32.Vector2{X: n, Y: n}

	gr.Objects.Root = gr.Objects.SVG.Root
	gr.Objects.Root.ViewBox.Min = gr.Vectors.Min
	gr.Objects.Root.ViewBox.Size = gr.Vectors.Size
	gr.Objects.Root.SetProperty("stroke-width", "0.1dp")
	gr.Objects.Root.SetProperty("fill", colors.Scheme.Surface)

	svg.NewCircle(gr.Objects.Root).SetRadius(50)

	sp.SetSplits(0.5, 0.5)

	gr.Objects.Lines = svg.NewGroup(gr.Objects.Root)
	gr.Objects.Lines.Maker(func(p *tree.Plan) {
		for i := range len(gr.Lines) {
			tree.AddAt(p, strconv.Itoa(i), func(w *svg.Path) {})
		}
	})
	gr.Objects.Marbles = svg.NewGroup(gr.Objects.Root)
	gr.Objects.Coords = svg.NewGroup(gr.Objects.Root)
	gr.Objects.TrackingLines = svg.NewGroup(gr.Objects.Root)

	gr.Objects.Coords.SetProperty("stroke-width", "0.05dp")
	gr.Objects.TrackingLines.SetProperty("stroke-width", "0.05dp")
}
