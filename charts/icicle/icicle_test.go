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
