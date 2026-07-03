package boxplot_test

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// sampleData generates a deterministic distribution per group.
func sampleData() []boxplot.BoxPlotDatum {
	var out []boxplot.BoxPlotDatum
	for gi, g := range []string{"Alpha", "Beta", "Gamma"} {
		for i := 0; i < 20; i++ {
			// Deterministic spread, shifted per group.
			v := float64((i*7+gi*11)%40) + float64(gi)*10 + 10
			out = append(out, boxplot.BoxPlotDatum{Group: g, Value: v})
		}
	}
	return out
}

func baseProps() boxplot.BoxPlotProps {
	return boxplot.BoxPlotProps{
		Width: 500, Height: 400,
		Margin:  core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 60},
		Data:    sampleData(),
		ColorBy: "group",
	}
}

func renderChart(t *testing.T, props boxplot.BoxPlotProps) string {
	t.Helper()
	var b strings.Builder
	if err := boxplot.BoxPlot(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("BoxPlot.Render: %v", err)
	}
	return b.String()
}

func TestBoxPlot_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestBoxPlot_BoxCount(t *testing.T) {
	// 3 groups → 3 box rects (one <rect> each). Plus the SvgWrapper background
	// rect = 4 total.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<rect"); got != 4 {
		t.Errorf("rect count = %d, want 4 (3 boxes + 1 background)", got)
	}
}

func TestBoxPlot_QuantileSummary(t *testing.T) {
	// Verify the median quantile matches d3-array semantics for a known set.
	data := []boxplot.BoxPlotDatum{}
	for _, v := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		data = append(data, boxplot.BoxPlotDatum{Group: "g", Value: v})
	}
	res := boxplot.UseBoxPlot(boxplot.BoxPlotProps{
		Width: 400, Height: 300, Data: data,
		Quantiles: []float64{0.1, 0.25, 0.5, 0.75, 0.9}, Padding: 0.1, InnerPadding: 6,
		Layout: boxplot.BoxPlotLayoutVertical,
	})
	if len(res.Boxes) != 1 {
		t.Fatalf("box count = %d, want 1", len(res.Boxes))
	}
	vals := res.Boxes[0].Summary.Values
	// d3 quantileSorted of 1..10: q10=1.9, q25=3.25, q50=5.5, q75=7.75, q90=9.1.
	want := []float64{1.9, 3.25, 5.5, 7.75, 9.1}
	for i, w := range want {
		if math.Abs(vals[i]-w) > 1e-9 {
			t.Errorf("quantile[%d] = %v, want %v", i, vals[i], w)
		}
	}
}

func TestBoxPlot_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "boxplot-basic", out)
}

func TestBoxPlot_Golden_Horizontal(t *testing.T) {
	p := baseProps()
	p.Layout = boxplot.BoxPlotLayoutHorizontal
	out := renderChart(t, p)
	golden.Assert(t, "boxplot-horizontal", out)
}

func TestBoxPlot_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive boxplot should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive boxplot must not emit data-tc-tooltip")
	}
}

func TestBoxPlot_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.01
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("animated boxplot should emit <animate>")
	}
	p.Animate = false
	if strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("non-animated boxplot must not emit <animate>")
	}
}
