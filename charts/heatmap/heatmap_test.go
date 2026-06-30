package heatmap_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func pf(v float64) *float64 { return &v }

// sampleData is a deterministic 3-serie × 3-column dataset.
func sampleData() []heatmap.HeatMapSerie {
	return []heatmap.HeatMapSerie{
		{ID: "Japan", Data: []heatmap.HeatMapDatum{{X: "Train", Y: pf(10)}, {X: "Car", Y: pf(20)}, {X: "Bike", Y: pf(30)}}},
		{ID: "France", Data: []heatmap.HeatMapDatum{{X: "Train", Y: pf(40)}, {X: "Car", Y: pf(50)}, {X: "Bike", Y: pf(60)}}},
		{ID: "USA", Data: []heatmap.HeatMapDatum{{X: "Train", Y: pf(70)}, {X: "Car", Y: nil}, {X: "Bike", Y: pf(90)}}},
	}
}

func render(t *testing.T, props heatmap.HeatMapProps) string {
	t.Helper()
	var b strings.Builder
	if err := heatmap.HeatMap(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("HeatMap.Render: %v", err)
	}
	return b.String()
}

func TestHeatMap_RendersSVG(t *testing.T) {
	out := render(t, heatmap.HeatMapProps{
		Width: 500, Height: 360,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 60},
		Data:   sampleData(),
	})
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestHeatMap_CellCount(t *testing.T) {
	// 3 series × 3 columns = 9 cells; each cell emits one <rect>. Plus the
	// SvgWrapper background rect = 10 total <rect>.
	out := render(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
	})
	if got := strings.Count(out, "<rect"); got != 10 {
		t.Errorf("rect count = %d, want 10 (9 cells + 1 background)", got)
	}
}

func TestHeatMap_EmptyCellUsesEmptyColor(t *testing.T) {
	// The nil value (USA/Car) must render with the empty color, not a scale color.
	out := render(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		EmptyColor: "#123456",
	})
	if !strings.Contains(out, "#123456") {
		t.Errorf("expected empty color #123456 in output for the nil cell")
	}
}

func TestHeatMap_Golden(t *testing.T) {
	out := render(t, heatmap.HeatMapProps{
		Width: 500, Height: 360,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 60},
		Data:   sampleData(),
	})
	golden.Assert(t, "heatmap-basic", out)
}
