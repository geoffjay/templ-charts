package heatmap

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/canvas"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func cf(v float64) *float64 { return &v }

func canvasSampleProps() HeatMapProps {
	return HeatMapProps{
		Width: 500, Height: 360,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 60},
		Data: []HeatMapSerie{
			{ID: "Japan", Data: []HeatMapDatum{{X: "Train", Y: cf(10)}, {X: "Car", Y: cf(20)}, {X: "Bike", Y: cf(30)}}},
			{ID: "France", Data: []HeatMapDatum{{X: "Train", Y: cf(40)}, {X: "Car", Y: cf(50)}, {X: "Bike", Y: cf(60)}}},
			{ID: "USA", Data: []HeatMapDatum{{X: "Train", Y: cf(70)}, {X: "Car", Y: nil}, {X: "Bike", Y: cf(90)}}},
		},
	}
}

func computeModel(p HeatMapProps) (HeatMapProps, HeatMapResult, core.Dimensions) {
	p = applyDefaults(p)
	dims := core.UseDimensions(p.Width, p.Height, p.Margin)
	p.Width = dims.InnerWidth
	p.Height = dims.InnerHeight
	return p, UseHeatMap(p), dims
}

func cellOps(p HeatMapProps, result HeatMapResult) []canvas.Op {
	rec := canvas.NewRecorder()
	recordCells(rec, p, result)
	return rec.Ops()
}

// TestCanvas_DrawListGolden pins the recorded cell draw-list.
func TestCanvas_DrawListGolden(t *testing.T) {
	p, result, _ := computeModel(canvasSampleProps())
	golden.Assert(t, "heatmap-canvas-drawlist", canvas.EncodeJSON(cellOps(p, result)))
}

// TestCanvas_OneRectPerCell asserts one FillRect op per computed cell, matching
// the SVG cells layer's one <rect> per cell (position/size within rounding).
func TestCanvas_OneRectPerCell(t *testing.T) {
	p, result, _ := computeModel(canvasSampleProps())
	var rects []canvas.Op
	for _, op := range cellOps(p, result) {
		if op.Kind == canvas.OpFillRect {
			rects = append(rects, op)
		}
	}
	if len(rects) != len(result.Cells) {
		t.Fatalf("FillRect count = %d, want %d (one per cell)", len(rects), len(result.Cells))
	}
	for i, cell := range result.Cells {
		r := rects[i]
		wantX, wantY := cell.XPos-cell.Width/2, cell.YPos-cell.Height/2
		if r.A != wantX || r.B != wantY {
			t.Errorf("cell %d: rect origin (%v,%v), want (%v,%v)", i, r.A, r.B, wantX, wantY)
		}
	}
}

// TestCanvas_CorrespondsToSVGCellCount cross-checks the draw-list against the
// SVG render: FillRect ops equal the SVG cell <rect> count (total rects minus
// the SvgWrapper background rect).
func TestCanvas_CorrespondsToSVGCellCount(t *testing.T) {
	p, result, _ := computeModel(canvasSampleProps())
	fr := 0
	for _, op := range cellOps(p, result) {
		if op.Kind == canvas.OpFillRect {
			fr++
		}
	}

	var b strings.Builder
	svgProps := canvasSampleProps()
	svgProps.Render = theming.EngineSVG
	if err := HeatMap(svgProps).Render(context.Background(), &b); err != nil {
		t.Fatalf("render SVG: %v", err)
	}
	svgCells := strings.Count(b.String(), "<rect") - 1 // minus background rect
	if svgCells != fr {
		t.Errorf("SVG cell rects = %d, canvas FillRect = %d (must correspond)", svgCells, fr)
	}
}

// TestCanvas_RendersComposition checks the Canvas engine emits the layered
// wrapper: SVG panes, a <canvas>, and the ops script.
func TestCanvas_RendersComposition(t *testing.T) {
	p := canvasSampleProps()
	p.Render = theming.EngineCanvas
	p.ChartID = "hm1"
	var b strings.Builder
	if err := HeatMap(p).Render(context.Background(), &b); err != nil {
		t.Fatalf("render canvas: %v", err)
	}
	out := b.String()
	for _, want := range []string{
		`<div class="tc-canvas-chart"`,
		`class="tc-canvas-pane"`,
		`<canvas class="tc-canvas" id="hm1"`,
		`id="hm1-ops"`,
		`["fr",`,
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
	if err := HeatMap(canvasSampleProps()).Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.HasPrefix(b.String(), "<svg") {
		t.Error("default (zero Render) engine should emit <svg>")
	}
}
