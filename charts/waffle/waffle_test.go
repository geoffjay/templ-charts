package waffle_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
	"github.com/geoffjay/templ-charts/charts/waffle"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []waffle.WaffleDatum {
	return []waffle.WaffleDatum{
		{ID: "men", Label: "Men", Value: 30},
		{ID: "women", Label: "Women", Value: 45},
		{ID: "children", Label: "Children", Value: 25},
	}
}

func renderChart(t *testing.T, props waffle.WaffleProps) string {
	t.Helper()
	var b strings.Builder
	if err := waffle.Waffle(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Waffle.Render: %v", err)
	}
	return b.String()
}

func TestWaffle_RendersSVG(t *testing.T) {
	out := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Total:  100, Rows: 10, Columns: 10,
		Data: sampleData(),
	})
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("expected a full <svg>…</svg>")
	}
}

func TestWaffle_CellCount(t *testing.T) {
	// 10×10 grid → 100 cells, each a <rect>, plus the background rect = 101.
	out := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400, Total: 100, Rows: 10, Columns: 10, Data: sampleData(),
	})
	if got := strings.Count(out, "<rect"); got != 101 {
		t.Errorf("rect count = %d, want 101 (100 cells + 1 background)", got)
	}
}

func TestWaffle_FilledCellCount(t *testing.T) {
	// total=100, 100 cells → unit=1. Values 30/45/25 fill 30+45+25=100 cells,
	// none empty. The empty color (#cccccc) should therefore be absent.
	out := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400, Total: 100, Rows: 10, Columns: 10, Data: sampleData(),
	})
	if strings.Contains(out, "#cccccc") {
		t.Errorf("expected no empty cells (sum of values fills the grid)")
	}

	// total=200 → unit=2 → 15+23+13 = 51 filled, 49 empty → empty color present.
	out2 := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400, Total: 200, Rows: 10, Columns: 10, Data: sampleData(),
	})
	if !strings.Contains(out2, "#cccccc") {
		t.Errorf("expected empty cells when total exceeds the data sum")
	}
}

func TestWaffle_Golden(t *testing.T) {
	out := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Total:  100, Rows: 10, Columns: 10,
		Data: sampleData(),
	})
	golden.Assert(t, "waffle-basic", out)
}

func TestWaffle_GoldenFillDirection(t *testing.T) {
	out := renderChart(t, waffle.WaffleProps{
		Width: 400, Height: 400,
		Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Total:  100, Rows: 10, Columns: 10,
		FillDirection: grid.GridFillBottom,
		Data:          sampleData(),
	})
	golden.Assert(t, "waffle-fill-bottom", out)
}

func TestWaffle_InteractiveEmitsTooltip(t *testing.T) {
	p := waffle.WaffleProps{
		Width: 400, Height: 400, Total: 100, Rows: 10, Columns: 10, Data: sampleData(),
		Interactive: true,
	}
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive waffle should emit data-tc-tooltip")
	}
	p.Interactive = false
	if strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("non-interactive waffle must not emit data-tc-tooltip")
	}
}
