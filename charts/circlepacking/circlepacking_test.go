package circlepacking_test

import (
	"context"
	"strings"
	"testing"

	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() cp.CirclePackingNode {
	return cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
		{ID: "A", Children: []cp.CirclePackingNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
		{ID: "B", Children: []cp.CirclePackingNode{{ID: "b1", Value: 6}}},
		{ID: "C", Value: 10},
	}}
}

func baseProps() cp.CirclePackingProps {
	return cp.CirclePackingProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props cp.CirclePackingProps) string {
	t.Helper()
	var b strings.Builder
	if err := cp.CirclePacking(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("CirclePacking.Render: %v", err)
	}
	return b.String()
}

func TestCirclePacking_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestCirclePacking_CircleCount(t *testing.T) {
	// root + A,B,C + a1,a2,b1 = 7 circles.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 7 {
		t.Errorf("circle count = %d, want 7", got)
	}
}

func TestCirclePacking_Golden(t *testing.T) {
	golden.Assert(t, "circlepacking-basic", renderChart(t, baseProps()))
}

func TestCirclePacking_GoldenPadding(t *testing.T) {
	p := baseProps()
	p.Padding = 12
	golden.Assert(t, "circlepacking-padding", renderChart(t, p))
}

func TestCirclePacking_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Packed"
	p.Desc = "Nested packed circles."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Packed</title>") || !strings.Contains(out, "<desc>Nested packed circles.</desc>") {
		t.Errorf("expected title/desc")
	}
}

func TestCirclePacking_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive circle-packing should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive circle-packing must not emit data-tc-tooltip")
	}
}

func TestCirclePacking_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.05
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("Animate=true should emit an <animate> enter transition")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("Animate=false should not emit any <animate>")
	}
}

func TestCirclePacking_ZoomOffNoZoomAttrs(t *testing.T) {
	if strings.Contains(renderChart(t, baseProps()), "/zoom?node=") {
		t.Errorf("default circlepacking must not emit /zoom hx-get")
	}
	p := baseProps()
	p.EnableZooming = true // no ChartID → still off
	if strings.Contains(renderChart(t, p), "/zoom?node=") {
		t.Errorf("circlepacking with EnableZooming but no ChartID must not emit /zoom hx-get")
	}
}

func TestCirclePacking_ZoomOnEmitsZoomTargets(t *testing.T) {
	p := baseProps()
	p.EnableZooming = true
	p.ChartID = "cp1"
	out := renderChart(t, p)
	for _, id := range []string{"A", "B", "C", "a1", "a2", "b1"} {
		if !strings.Contains(out, "/charts/cp1/zoom?node="+id+`"`) {
			t.Errorf("zoomable circlepacking should emit zoom target for node %q", id)
		}
	}
	if !strings.Contains(out, `hx-target="#chart-cp1"`) {
		t.Errorf("zoom targets should swap into #chart-cp1")
	}
}

func TestCirclePacking_ZoomedShowsBreadcrumb(t *testing.T) {
	p := baseProps()
	p.EnableZooming = true
	p.ChartID = "cp1"
	p.FocusID = "A"
	out := renderChart(t, p)
	if !strings.Contains(out, "/charts/cp1/zoom?node=root") {
		t.Errorf("focused circlepacking breadcrumb should link back to root")
	}
	if strings.Contains(out, "/charts/cp1/zoom?node=b1") {
		t.Errorf("focused circlepacking should not render nodes outside the focus subtree")
	}
}

func TestCirclePacking_GoldenZoomable(t *testing.T) {
	p := baseProps()
	p.EnableZooming = true
	p.ChartID = "cp1"
	golden.Assert(t, "circlepacking-zoomable", renderChart(t, p))
}

func TestCirclePacking_GoldenZoomed(t *testing.T) {
	p := baseProps()
	p.EnableZooming = true
	p.ChartID = "cp1"
	p.FocusID = "A"
	golden.Assert(t, "circlepacking-zoomed", renderChart(t, p))
}
