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
