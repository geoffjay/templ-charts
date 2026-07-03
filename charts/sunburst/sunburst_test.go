package sunburst_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() sunburst.SunburstNode {
	return sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
		{ID: "A", Children: []sunburst.SunburstNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
		{ID: "B", Children: []sunburst.SunburstNode{{ID: "b1", Value: 6}}},
		{ID: "C", Value: 10},
	}}
}

func baseProps() sunburst.SunburstProps {
	return sunburst.SunburstProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props sunburst.SunburstProps) string {
	t.Helper()
	var b strings.Builder
	if err := sunburst.Sunburst(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Sunburst.Render: %v", err)
	}
	return b.String()
}

func TestSunburst_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestSunburst_ArcCount(t *testing.T) {
	// depth>=1 nodes: A,B,C + a1,a2,b1 = 6 arcs.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<path"); got != 6 {
		t.Errorf("arc path count = %d, want 6", got)
	}
}

func TestSunburst_Golden(t *testing.T) {
	golden.Assert(t, "sunburst-basic", renderChart(t, baseProps()))
}

func TestSunburst_GoldenArcLabels(t *testing.T) {
	p := baseProps()
	p.EnableArcLabels = sunburst.BoolPtr(true)
	golden.Assert(t, "sunburst-arc-labels", renderChart(t, p))
}

func TestSunburst_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Breakdown"
	p.Desc = "Nested value breakdown."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Breakdown</title>") || !strings.Contains(out, "<desc>Nested value breakdown.</desc>") {
		t.Errorf("expected title/desc")
	}
}

func TestSunburst_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive sunburst should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive sunburst must not emit data-tc-tooltip")
	}
}
