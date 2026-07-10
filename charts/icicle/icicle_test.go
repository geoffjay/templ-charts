package icicle_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() icicle.IcicleNode {
	return icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
		{ID: "A", Children: []icicle.IcicleNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
		{ID: "B", Children: []icicle.IcicleNode{{ID: "b1", Value: 6}}},
		{ID: "C", Value: 10},
	}}
}

func baseProps() icicle.IcicleProps {
	return icicle.IcicleProps{
		Width: 500, Height: 300,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props icicle.IcicleProps) string {
	t.Helper()
	var b strings.Builder
	if err := icicle.Icicle(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Icicle.Render: %v", err)
	}
	return b.String()
}

func TestIcicle_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestIcicle_RectCount(t *testing.T) {
	// root + A,B,C + a1,a2,b1 = 7 node rects (plus the SvgWrapper background).
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<rect x="); got != 7 {
		t.Errorf("node rect count = %d, want 7", got)
	}
}

func TestIcicle_Orientations(t *testing.T) {
	for _, o := range []icicle.Orientation{icicle.OrientationBottom, icicle.OrientationTop, icicle.OrientationLeft, icicle.OrientationRight} {
		p := baseProps()
		p.Orientation = o
		out := renderChart(t, p)
		if !strings.HasPrefix(out, "<svg") {
			t.Errorf("orientation %q failed to render", o)
		}
	}
}

func TestIcicle_Golden(t *testing.T) {
	golden.Assert(t, "icicle-basic", renderChart(t, baseProps()))
}

func TestIcicle_GoldenRight(t *testing.T) {
	p := baseProps()
	p.Orientation = icicle.OrientationRight
	golden.Assert(t, "icicle-right", renderChart(t, p))
}

func TestIcicle_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Icicle"
	p.Desc = "Depth-banded hierarchy."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Icicle</title>") || !strings.Contains(out, "<desc>Depth-banded hierarchy.</desc>") {
		t.Errorf("expected title/desc")
	}
}

func TestIcicle_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive icicle should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive icicle must not emit data-tc-tooltip")
	}
}

func TestIcicle_ZoomOffNoZoomAttrs(t *testing.T) {
	// Zoom off (default), and zoom on but without a ChartID, must not emit any
	// /zoom hx-get — output stays byte-stable to the pre-zoom behavior.
	if strings.Contains(renderChart(t, baseProps()), "/zoom?node=") {
		t.Errorf("default icicle must not emit /zoom hx-get")
	}
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true) // no ChartID → still off
	if strings.Contains(renderChart(t, p), "/zoom?node=") {
		t.Errorf("icicle with EnableZooming but no ChartID must not emit /zoom hx-get")
	}
}

func TestIcicle_ZoomOnEmitsZoomTargets(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "ic1"
	out := renderChart(t, p)
	for _, id := range []string{"A", "B", "C", "a1", "a2", "b1"} {
		if !strings.Contains(out, "/charts/ic1/zoom?node="+id+`"`) {
			t.Errorf("zoomable icicle should emit zoom target for node %q", id)
		}
	}
	if !strings.Contains(out, `hx-target="#chart-ic1"`) || !strings.Contains(out, `hx-swap="innerHTML"`) {
		t.Errorf("zoom targets should swap innerHTML into #chart-ic1")
	}
}

func TestIcicle_ZoomedShowsBreadcrumb(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "ic1"
	p.FocusID = "A"
	out := renderChart(t, p)
	// Breadcrumb links to the root (zoom-out) and to A.
	if !strings.Contains(out, "/charts/ic1/zoom?node=root") {
		t.Errorf("focused icicle breadcrumb should link back to root")
	}
	// Sibling B/C collapse out of the focused subtree.
	if strings.Contains(out, "/charts/ic1/zoom?node=b1") {
		t.Errorf("focused icicle should not render nodes outside the focus subtree")
	}
}

func TestIcicle_GoldenZoomable(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "ic1"
	golden.Assert(t, "icicle-zoomable", renderChart(t, p))
}

func TestIcicle_GoldenZoomed(t *testing.T) {
	p := baseProps()
	p.EnableZooming = core.BoolPtr(true)
	p.ChartID = "ic1"
	p.FocusID = "A"
	golden.Assert(t, "icicle-zoomed", renderChart(t, p))
}

func TestIcicle_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.01
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("animated icicle should emit <animate>")
	}
	p.Animate = false
	if strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("non-animated icicle must not emit <animate>")
	}
}
