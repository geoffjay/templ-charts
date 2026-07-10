package radialbar_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/radialbar"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// sampleData is a deterministic 2-serie × 2-category dataset.
func sampleData() []radialbar.RadialBarSerie {
	return []radialbar.RadialBarSerie{
		{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18},
		}},
		{ID: "Combini", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 12}, {X: "Fruits", Y: 9},
		}},
	}
}

func baseProps() radialbar.RadialBarProps {
	return radialbar.RadialBarProps{
		Width: 460, Height: 460,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 40},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props radialbar.RadialBarProps) string {
	t.Helper()
	var b strings.Builder
	if err := radialbar.RadialBar(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("RadialBar.Render: %v", err)
	}
	return b.String()
}

func TestRadialBar_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestRadialBar_BarArcCount(t *testing.T) {
	// 2 series × 2 categories = 4 bar arcs. Tracks add 2 arcs (one per serie
	// band). The grid's circular axis + circular grid also emit arc paths, so
	// just assert the bar arcs are present (>= 4 <path with fill that isn't
	// "none"). Simpler: count the bars+tracks group paths by checking the
	// number of arc <path d="M ...">. Assert at least 6 paths total.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<path"); got < 6 {
		t.Errorf("path count = %d, want >= 6 (4 bars + 2 tracks + grid)", got)
	}
}

func TestRadialBar_NoTracksWhenDisabled(t *testing.T) {
	withTracks := strings.Count(renderChart(t, baseProps()), "<path")
	p := baseProps()
	p.EnableTracks = core.BoolPtr(false)
	withoutTracks := strings.Count(renderChart(t, p), "<path")
	if withoutTracks >= withTracks {
		t.Errorf("disabling tracks should reduce path count: with=%d without=%d", withTracks, withoutTracks)
	}
}

func TestRadialBar_LabelsWhenEnabled(t *testing.T) {
	p := baseProps()
	p.EnableLabels = core.BoolPtr(true)
	out := renderChart(t, p)
	// formattedValue labels for the 4 bars (all spans > 10° skip angle here).
	if !strings.Contains(out, ">25<") {
		t.Errorf("expected a value label (25) when labels enabled")
	}
}

func TestRadialBar_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "radialbar-basic", out)
}

func TestRadialBar_Golden_Labels(t *testing.T) {
	p := baseProps()
	p.EnableLabels = core.BoolPtr(true)
	p.CornerRadius = 4
	out := renderChart(t, p)
	golden.Assert(t, "radialbar-labels", out)
}

func TestRadialBar_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	if !strings.Contains(renderChart(t, p), `<animate attributeName="opacity"`) {
		t.Errorf("animated radial bar should emit an opacity fade-in <animate>")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("non-animated radial bar must not emit any <animate>")
	}
}

func TestRadialBar_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive radial-bar should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive radial-bar must not emit data-tc-tooltip")
	}
}
