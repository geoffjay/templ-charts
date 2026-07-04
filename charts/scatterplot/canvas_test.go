package scatterplot

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/canvas"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func canvasSampleProps() ScatterPlotProps {
	return ScatterPlotProps{
		Width: 500, Height: 400,
		Margin: core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data: []ScatterPlotSerie{
			{ID: "group A", Data: []ScatterPlotDatum{{X: 10.0, Y: 20.0}, {X: 30.0, Y: 40.0}, {X: 55.0, Y: 12.0}}},
			{ID: "group B", Data: []ScatterPlotDatum{{X: 15.0, Y: 60.0}, {X: 42.0, Y: 33.0}, {X: 70.0, Y: 80.0}}},
		},
		EnableGridX: true,
		EnableGridY: true,
	}
}

// computeModel runs the same pipeline the templ does, returning the resolved
// props, computed result, and dimensions.
func computeModel(p ScatterPlotProps) (ScatterPlotProps, ScatterPlotResult, core.Dimensions) {
	p = applyDefaults(p)
	dims := core.UseDimensions(p.Width, p.Height, p.Margin)
	p.Width = dims.InnerWidth
	p.Height = dims.InnerHeight
	return p, UseScatterPlot(p), dims
}

func nodeOps(result ScatterPlotResult) []canvas.Op {
	rec := canvas.NewRecorder()
	recordNodes(rec, result)
	return rec.Ops()
}

// TestCanvas_DrawListGolden pins the recorded node draw-list.
func TestCanvas_DrawListGolden(t *testing.T) {
	_, result, _ := computeModel(canvasSampleProps())
	golden.Assert(t, "scatterplot-canvas-drawlist", canvas.EncodeJSON(nodeOps(result)))
}

// TestCanvas_OneCirclePerNode asserts a strict correspondence: exactly one
// FillCircle op per computed node, matching the SVG nodes layer's one <circle>
// per node, with the same centre and radius.
func TestCanvas_OneCirclePerNode(t *testing.T) {
	_, result, _ := computeModel(canvasSampleProps())
	ops := nodeOps(result)

	var circles []canvas.Op
	for _, op := range ops {
		if op.Kind == canvas.OpFillCircle {
			circles = append(circles, op)
		}
	}
	if len(circles) != len(result.Nodes) {
		t.Fatalf("FillCircle count = %d, want %d (one per node)", len(circles), len(result.Nodes))
	}
	for i, n := range result.Nodes {
		c := circles[i]
		if c.A != n.X || c.B != n.Y || c.C != n.Size/2 {
			t.Errorf("node %d: circle (%v,%v,r=%v), want (%v,%v,r=%v)", i, c.A, c.B, c.C, n.X, n.Y, n.Size/2)
		}
	}
}

// TestCanvas_CorrespondsToSVGCircleCount cross-checks the draw-list against the
// actual SVG render: the number of FillCircle ops equals the number of <circle>
// elements the SVG backend emits for the same data.
func TestCanvas_CorrespondsToSVGCircleCount(t *testing.T) {
	p, result, _ := computeModel(canvasSampleProps())
	fc := 0
	for _, op := range nodeOps(result) {
		if op.Kind == canvas.OpFillCircle {
			fc++
		}
	}

	var b strings.Builder
	svgProps := canvasSampleProps()
	svgProps.Render = theming.EngineSVG
	if err := ScatterPlot(svgProps).Render(context.Background(), &b); err != nil {
		t.Fatalf("render SVG: %v", err)
	}
	if svgCircles := strings.Count(b.String(), "<circle"); svgCircles != fc {
		t.Errorf("SVG circles = %d, canvas FillCircle = %d (must correspond)", svgCircles, fc)
	}
	_ = p
}

// TestCanvas_RendersComposition checks the Canvas engine emits the layered
// wrapper: an underlay/overlay SVG pane, a <canvas>, and the ops script.
func TestCanvas_RendersComposition(t *testing.T) {
	p := canvasSampleProps()
	p.Render = theming.EngineCanvas
	p.ChartID = "sp1"
	var b strings.Builder
	if err := ScatterPlot(p).Render(context.Background(), &b); err != nil {
		t.Fatalf("render canvas: %v", err)
	}
	out := b.String()
	for _, want := range []string{
		`<div class="tc-canvas-chart"`,
		`class="tc-canvas-pane"`, // grid / axes panes
		`<canvas class="tc-canvas" id="sp1"`,
		`id="sp1-ops"`,
		`["fc",`, // draw-list embedded
	} {
		if !strings.Contains(out, want) {
			t.Errorf("canvas render missing %q", want)
		}
	}
	if strings.HasPrefix(out, "<svg") {
		t.Error("canvas engine should not emit a top-level <svg> wrapper")
	}
}

// TestCanvas_SVGDefaultUnchanged asserts the default engine still renders SVG.
func TestCanvas_SVGDefaultUnchanged(t *testing.T) {
	var b strings.Builder
	if err := ScatterPlot(canvasSampleProps()).Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.HasPrefix(b.String(), "<svg") {
		t.Error("default (zero Render) engine should emit <svg>")
	}
}
