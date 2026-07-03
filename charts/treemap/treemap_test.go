package treemap_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/treemap"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() treemap.TreemapNode {
	return treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
		{ID: "A", Children: []treemap.TreemapNode{
			{ID: "a1", Value: 12}, {ID: "a2", Value: 8}, {ID: "a3", Value: 5},
		}},
		{ID: "B", Children: []treemap.TreemapNode{
			{ID: "b1", Value: 10}, {ID: "b2", Value: 6},
		}},
		{ID: "C", Value: 15},
	}}
}

func baseProps() treemap.TreemapProps {
	return treemap.TreemapProps{
		Width: 500, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props treemap.TreemapProps) string {
	t.Helper()
	var b strings.Builder
	if err := treemap.Treemap(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Treemap.Render: %v", err)
	}
	return b.String()
}

func TestTreemap_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestTreemap_NodeCount(t *testing.T) {
	// root + 3 top-level (A,B,C) + 5 leaves under A,B = 9 node rects (the extra
	// <rect> is the SvgWrapper background, which carries no fill-opacity).
	out := renderChart(t, baseProps())
	if got := strings.Count(out, `fill-opacity=`); got != 9 {
		t.Errorf("node rect count = %d, want 9", got)
	}
}

func TestTreemap_LeafLabels(t *testing.T) {
	out := renderChart(t, baseProps())
	// leaf formattedValue labels present.
	for _, want := range []string{">12<", ">15<", ">6<"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected leaf value label %q", want)
		}
	}
}

func TestTreemap_ParentLabels(t *testing.T) {
	out := renderChart(t, baseProps())
	for _, want := range []string{">A<", ">B<"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected parent label %q", want)
		}
	}
}

func TestTreemap_Golden(t *testing.T) {
	golden.Assert(t, "treemap-basic", renderChart(t, baseProps()))
}

func TestTreemap_GoldenBinaryTile(t *testing.T) {
	p := baseProps()
	p.Tile = treemap.TileBinary
	golden.Assert(t, "treemap-binary", renderChart(t, p))
}

func TestTreemap_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Sales"
	p.Desc = "Sales by region and product."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Sales</title>") || !strings.Contains(out, "<desc>Sales by region and product.</desc>") {
		t.Errorf("expected title/desc")
	}
}

func TestTreemap_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive treemap should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive treemap must not emit data-tc-tooltip")
	}
}

func TestTreemap_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.01
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("animated treemap should emit <animate>")
	}
	p.Animate = false
	if strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("non-animated treemap must not emit <animate>")
	}
}
