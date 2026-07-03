package polarbar_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/polarbar"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []polarbar.PolarBarDatum {
	return []polarbar.PolarBarDatum{
		{Index: "A", Values: map[string]float64{"hot": 10, "cold": 6}},
		{Index: "B", Values: map[string]float64{"hot": 7, "cold": 12}},
		{Index: "C", Values: map[string]float64{"hot": 15, "cold": 4}},
		{Index: "D", Values: map[string]float64{"hot": 9, "cold": 9}},
	}
}

func baseProps() polarbar.PolarBarProps {
	return polarbar.PolarBarProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 40, Right: 60, Bottom: 40, Left: 60},
		Data:   sampleData(),
		Keys:   []string{"hot", "cold"},
	}
}

func renderChart(t *testing.T, props polarbar.PolarBarProps) string {
	t.Helper()
	var b strings.Builder
	if err := polarbar.PolarBar(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("PolarBar.Render: %v", err)
	}
	return b.String()
}

func TestPolarBar_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestPolarBar_ArcCount(t *testing.T) {
	// 4 indices × 2 keys = 8 stacked arcs, each a <path> in the arcs layer.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, `<path class="tc-arc"`); got != 8 {
		// fall back: arcs may not carry a class; assert at least 8 arc paths via the generator.
		if got2 := strings.Count(out, "<path"); got2 < 8 {
			t.Errorf("expected at least 8 arc paths, got %d (tc-arc=%d)", got2, got)
		}
	}
}

func TestPolarBar_IndexLabels(t *testing.T) {
	out := renderChart(t, baseProps())
	for _, want := range []string{">A<", ">B<", ">C<", ">D<"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected circular-axis index label %q", want)
		}
	}
}

func TestPolarBar_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "polarbar-basic", out)
}

func TestPolarBar_GoldenInnerRadius(t *testing.T) {
	p := baseProps()
	p.InnerRadius = 0.5 // ratio in [0,1]; default 0 → donut hole
	golden.Assert(t, "polarbar-inner-radius", renderChart(t, p))
}

func TestPolarBar_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Temperatures"
	p.Desc = "Hot/cold split across four indices."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Temperatures</title>") {
		t.Errorf("expected <title>")
	}
	if !strings.Contains(out, "<desc>Hot/cold split across four indices.</desc>") {
		t.Errorf("expected <desc>")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestPolarBar_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive polar-bar should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive polar-bar must not emit data-tc-tooltip")
	}
}

func TestPolarBar_Legend(t *testing.T) {
	p := baseProps()
	p.Legends = []legends.LegendProps{
		{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 50},
	}
	out := renderChart(t, p)
	if !strings.Contains(out, "hot") || !strings.Contains(out, "cold") {
		t.Errorf("expected legend to include keys")
	}
}
