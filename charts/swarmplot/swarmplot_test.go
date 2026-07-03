package swarmplot_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []swarmplot.SwarmPlotDatum {
	groups := []string{"A", "B", "C"}
	vals := []float64{12, 45, 23, 67, 34, 89, 5, 56, 78, 30, 41, 62}
	data := make([]swarmplot.SwarmPlotDatum, 0, len(vals))
	for i, v := range vals {
		g := groups[i%len(groups)]
		data = append(data, swarmplot.SwarmPlotDatum{
			ID: g + string(rune('0'+i)), Group: g, Value: v,
		})
	}
	return data
}

func baseProps() swarmplot.SwarmPlotProps {
	return swarmplot.SwarmPlotProps{
		Width: 500, Height: 400,
		Margin: core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data:   sampleData(),
		Groups: []string{"A", "B", "C"},
	}
}

func renderChart(t *testing.T, props swarmplot.SwarmPlotProps) string {
	t.Helper()
	var b strings.Builder
	if err := swarmplot.SwarmPlot(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("SwarmPlot.Render: %v", err)
	}
	return b.String()
}

func TestSwarmPlot_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestSwarmPlot_CircleCount(t *testing.T) {
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 12 {
		t.Errorf("circle count = %d, want 12", got)
	}
}

func TestSwarmPlot_Golden(t *testing.T) {
	golden.Assert(t, "swarmplot-basic", renderChart(t, baseProps()))
}

func TestSwarmPlot_HorizontalGolden(t *testing.T) {
	p := baseProps()
	p.Layout = "horizontal"
	golden.Assert(t, "swarmplot-horizontal", renderChart(t, p))
}

func TestSwarmPlot_Deterministic(t *testing.T) {
	if renderChart(t, baseProps()) != renderChart(t, baseProps()) {
		t.Errorf("swarmplot render is not deterministic")
	}
}

func TestSwarmPlot_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Swarm"
	p.Desc = "Grouped value distribution."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Swarm</title>") || !strings.Contains(out, "<desc>Grouped value distribution.</desc>") {
		t.Errorf("expected title/desc threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestSwarmPlot_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive swarmplot should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive swarmplot must not emit data-tc-tooltip")
	}
}

func TestSwarmPlot_VoronoiMesh(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	p.UseMesh = true
	if !strings.Contains(renderChart(t, p), "data-tc-mesh") {
		t.Errorf("Interactive+UseMesh should emit a voronoi mesh")
	}
}
