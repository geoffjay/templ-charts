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
		EnableGridX: true,
		EnableGridY: true,
	}
}

func render(t *testing.T, props scatterplot.ScatterPlotProps) string {
	t.Helper()
	var b strings.Builder
	if err := scatterplot.ScatterPlot(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ScatterPlot.Render: %v", err)
	}
	return b.String()
}

func TestScatterPlot_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestScatterPlot_NodeCount(t *testing.T) {
	// 2 series × 3 nodes = 6 dots, each a <circle>.
	out := render(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 6 {
		t.Errorf("circle count = %d, want 6", got)
	}
}

func TestScatterPlot_GridDisabled(t *testing.T) {
	p := baseProps()
	p.EnableGridX = false
	p.EnableGridY = false
	withGrid := strings.Count(render(t, baseProps()), "<line")
	noGrid := strings.Count(render(t, p), "<line")
	if noGrid >= withGrid {
		t.Errorf("disabling grid should reduce line count: with=%d without=%d", withGrid, noGrid)
	}
}

func TestScatterPlot_Golden(t *testing.T) {
	out := render(t, baseProps())
	golden.Assert(t, "scatterplot-basic", out)
}

func TestScatterPlot_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	out := render(t, p)
	if !strings.Contains(out, "data-tc-tooltip") {
		t.Errorf("interactive scatterplot should emit data-tc-tooltip")
	}
	// default (non-interactive) must not.
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive scatterplot must not emit data-tc-tooltip")
	}
}
