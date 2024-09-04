package main

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/paint"
)

// draw renders the graph.
func (gr *Graph) draw(pc *paint.Context) {
	gr.drawAxes(pc)
}

// canvasCoord converts the given coordinate to a normalized 0-1 canvas coordinate.
func (gr *Graph) canvasCoord(v math32.Vector2) math32.Vector2 {
	res := math32.Vector2{}
	res.X = v.X / (gr.Vectors.Max.X - gr.Vectors.Min.X)
	res.Y = 1 - v.Y/(gr.Vectors.Max.Y-gr.Vectors.Min.Y)
	return res
}

func (gr *Graph) drawAxes(pc *paint.Context) {
	pc.StrokeStyle.Color = colors.Scheme.OutlineVariant

	start := gr.canvasCoord(math32.Vec2(gr.Vectors.Min.X, 0))
	end := gr.canvasCoord(math32.Vec2(gr.Vectors.Max.X, 0))
	pc.MoveTo(start.X, start.Y)
	pc.LineTo(end.X, end.Y)
	pc.Stroke()

	start = gr.canvasCoord(math32.Vec2(0, gr.Vectors.Min.Y))
	end = gr.canvasCoord(math32.Vec2(0, gr.Vectors.Max.Y))
	pc.MoveTo(start.X, start.Y)
	pc.LineTo(end.X, end.Y)
	pc.Stroke()
}
