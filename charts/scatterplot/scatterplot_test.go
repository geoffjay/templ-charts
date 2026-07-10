package scatterplot_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []scatterplot.ScatterPlotSerie {
	return []scatterplot.ScatterPlotSerie{
		{ID: "group A", Data: []scatterplot.ScatterPlotDatum{
			{X: 10.0, Y: 20.0}, {X: 30.0, Y: 40.0}, {X: 55.0, Y: 12.0},
		}},
		{ID: "group B", Data: []scatterplot.ScatterPlotDatum{
			{X: 15.0, Y: 60.0}, {X: 42.0, Y: 33.0}, {X: 70.0, Y: 80.0},
		}},
	}
}

func baseProps() scatterplot.ScatterPlotProps {
	return scatterplot.ScatterPlotProps{
		Width: 500, Height: 400,
		Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data:        sampleData(),
		EnableGridX: core.BoolPtr(true),
		EnableGridY: core.BoolPtr(true),
	}
}

func renderChart(t *testing.T, props scatterplot.ScatterPlotProps) string {
	t.Helper()
	var b strings.Builder
	if err := scatterplot.ScatterPlot(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ScatterPlot.Render: %v", err)
	}
	return b.String()
}

func TestScatterPlot_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestScatterPlot_NodeCount(t *testing.T) {
	// 2 series × 3 nodes = 6 dots, each a <circle>.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 6 {
		t.Errorf("circle count = %d, want 6", got)
	}
}

func TestScatterPlot_GridDisabled(t *testing.T) {
	p := baseProps()
	p.EnableGridX = core.BoolPtr(false)
	p.EnableGridY = core.BoolPtr(false)
	withGrid := strings.Count(renderChart(t, baseProps()), "<line")
	noGrid := strings.Count(renderChart(t, p), "<line")
	if noGrid >= withGrid {
		t.Errorf("disabling grid should reduce line count: with=%d without=%d", withGrid, noGrid)
	}
}

func TestScatterPlot_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "scatterplot-basic", out)
}

func TestScatterPlot_GoldenNoGrid(t *testing.T) {
	p := baseProps()
	p.EnableGridX = core.BoolPtr(false)
	p.EnableGridY = core.BoolPtr(false)
	golden.Assert(t, "scatterplot-no-grid", renderChart(t, p))
}

func TestScatterPlot_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Groups A and B"
	p.Desc = "Scatter of x/y values for two groups."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Groups A and B</title>") {
		t.Errorf("expected <title> threaded through to SvgWrapper")
	}
	if !strings.Contains(out, "<desc>Scatter of x/y values for two groups.</desc>") {
		t.Errorf("expected <desc> threaded through to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestScatterPlot_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	out := renderChart(t, p)
	if !strings.Contains(out, "data-tc-tooltip") {
		t.Errorf("interactive scatterplot should emit data-tc-tooltip")
	}
	// default (non-interactive) must not.
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive scatterplot must not emit data-tc-tooltip")
	}
}

func TestScatterPlot_VoronoiMesh(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	p.UseMesh = true
	out := renderChart(t, p)
	if !strings.Contains(out, "data-tc-mesh") {
		t.Errorf("Interactive+UseMesh scatterplot should emit a voronoi-mesh overlay")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-mesh") {
		t.Errorf("scatterplot without UseMesh must not emit data-tc-mesh")
	}
	p.DetectionRadius = 25
	if !strings.Contains(renderChart(t, p), "data-tc-mesh-radius") {
		t.Errorf("DetectionRadius should emit data-tc-mesh-radius")
	}
}

func TestScatterPlot_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.05
	out := renderChart(t, p)
	if !strings.Contains(out, `<animate attributeName="r" from="0"`) {
		t.Errorf("Animate=true should emit an r enter animation")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("Animate=false should not emit any <animate>")
	}
	golden.Assert(t, "scatterplot-animated", out)
}
