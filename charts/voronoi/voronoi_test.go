package voronoi_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/voronoi"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []voronoi.VoronoiDatum {
	return []voronoi.VoronoiDatum{
		{ID: "A", X: 0.2, Y: 0.2},
		{ID: "B", X: 0.8, Y: 0.3},
		{ID: "C", X: 0.5, Y: 0.7},
		{ID: "D", X: 0.9, Y: 0.9},
		{ID: "E", X: 0.15, Y: 0.85},
		{ID: "F", X: 0.55, Y: 0.25},
		{ID: "G", X: 0.35, Y: 0.5},
		{ID: "H", X: 0.75, Y: 0.6},
	}
}

func baseProps() voronoi.VoronoiProps {
	return voronoi.VoronoiProps{
		Width: 500, Height: 500,
		Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props voronoi.VoronoiProps) string {
	t.Helper()
	var b strings.Builder
	if err := voronoi.Voronoi(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Voronoi.Render: %v", err)
	}
	return b.String()
}

func TestVoronoi_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestVoronoi_Animate(t *testing.T) {
	if out := renderChart(t, baseProps()); strings.Contains(out, "<animate") {
		t.Errorf("no <animate> expected when Animate=false")
	}
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.05
	out := renderChart(t, p)
	if !strings.Contains(out, `<animate attributeName="r" from="0"`) {
		t.Errorf("expected point radius <animate> when Animate=true")
	}
}

func TestVoronoi_PointCount(t *testing.T) {
	// One <circle> per input datum (points layer, enabled by default).
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 8 {
		t.Errorf("circle count = %d, want 8", got)
	}
}

func TestVoronoi_CellsRenderedByDefault(t *testing.T) {
	// Default: cells on (single path), links off.
	out := renderChart(t, baseProps())
	if !strings.Contains(out, `stroke="#000000"`) {
		t.Errorf("expected cell paths with default cell color")
	}
	if strings.Contains(out, `stroke="#bbbbbb"`) {
		t.Errorf("links should be off by default")
	}
}

func TestVoronoi_LinksToggle(t *testing.T) {
	p := baseProps()
	p.EnableLinks = voronoi.BoolPtr(true)
	out := renderChart(t, p)
	if !strings.Contains(out, `stroke="#bbbbbb"`) {
		t.Errorf("expected delaunay links when enabled")
	}
}

func TestVoronoi_PointsToggle(t *testing.T) {
	p := baseProps()
	p.EnablePoints = voronoi.BoolPtr(false)
	out := renderChart(t, p)
	if strings.Contains(out, "<circle") {
		t.Errorf("points disabled but circles rendered")
	}
}

func TestVoronoi_Bounds(t *testing.T) {
	// Inner area is 500-40 = 460 square; bounds path closes the rectangle.
	out := renderChart(t, baseProps())
	if !strings.Contains(out, "M0,0L460,0L460,460L0,460Z") {
		t.Errorf("expected bounds rectangle path for the inner area")
	}
}

func TestVoronoi_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "voronoi-basic", out)
}

func TestVoronoi_GoldenCellFill(t *testing.T) {
	// Enabling cell fill exercises the per-cell path rendering path (each cell
	// filled with its site color) instead of the default hollow single-path cells.
	p := baseProps()
	p.EnableCellFill = voronoi.BoolPtr(true)
	golden.Assert(t, "voronoi-cell-fill", renderChart(t, p))
}

func TestVoronoi_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Voronoi diagram"
	p.Desc = "Eight sites partitioned into cells."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Voronoi diagram</title>") {
		t.Errorf("expected <title> threaded to SvgWrapper")
	}
	if !strings.Contains(out, "<desc>Eight sites partitioned into cells.</desc>") {
		t.Errorf("expected <desc> threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestVoronoi_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	out := renderChart(t, p)
	if !strings.Contains(out, "data-tc-tooltip") {
		t.Errorf("interactive voronoi should emit per-cell data-tc-tooltip")
	}
	// Cells are unfilled, so the whole area must be hoverable via pointer-events,
	// otherwise the tooltip only triggers on the cell edges (the stroke).
	if !strings.Contains(out, `pointer-events="all"`) {
		t.Errorf("interactive cells need pointer-events=all so the interior is hoverable")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive voronoi must not emit data-tc-tooltip")
	}
}
