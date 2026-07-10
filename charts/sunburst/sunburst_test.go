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
	p.EnableArcLabels = core.BoolPtr(true)
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

func TestSunburst_Animate(t *testing.T) {
	on := renderChart(t, func() sunburst.SunburstProps {
		p := baseProps()
		p.Animate = true
		return p
	}())
	if !strings.Contains(on, `<animate attributeName="opacity"`) {
		t.Errorf("animated sunburst should emit an opacity fade-in <animate>")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("non-animated sunburst must not emit any <animate>")
	}
}

func TestSunburst_GoldenAnimated(t *testing.T) {
	p := baseProps()
	p.Animate = true
	golden.Assert(t, "sunburst-animated", renderChart(t, p))
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

func TestSunburst_ZoomOffNoZoomAttrs(t *testing.T) {
	if strings.Contains(renderChart(t, baseProps()), "/zoom?node=") {
		t.Errorf("default sunburst must not emit /zoom hx-get")
	}
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true) // no ChartID → still off
	if strings.Contains(renderChart(t, p), "/zoom?node=") {
		t.Errorf("sunburst with EnableZooming but no ChartID must not emit /zoom hx-get")
	}
}

func TestSunburst_ZoomOnEmitsZoomTargets(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "sb1"
	out := renderChart(t, p)
	for _, id := range []string{"A", "B", "C", "a1", "a2", "b1"} {
		if !strings.Contains(out, "/charts/sb1/zoom?node="+id+`"`) {
			t.Errorf("zoomable sunburst should emit zoom target for node %q", id)
		}
	}
	if !strings.Contains(out, `hx-target="#chart-sb1"`) {
		t.Errorf("zoom targets should swap into #chart-sb1")
	}
}

func TestSunburst_ZoomedShowsBreadcrumb(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "sb1"
	p.FocusID = "A"
	out := renderChart(t, p)
	if !strings.Contains(out, "/charts/sb1/zoom?node=root") {
		t.Errorf("focused sunburst breadcrumb should link back to root")
	}
	if strings.Contains(out, "/charts/sb1/zoom?node=b1") {
		t.Errorf("focused sunburst should not render nodes outside the focus subtree")
	}
}

func TestSunburst_GoldenZoomable(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "sb1"
	golden.Assert(t, "sunburst-zoomable", renderChart(t, p))
}

func TestSunburst_GoldenZoomed(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "sb1"
	p.FocusID = "A"
	golden.Assert(t, "sunburst-zoomed", renderChart(t, p))
}
