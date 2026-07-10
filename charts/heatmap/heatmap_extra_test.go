package heatmap_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// renderComp renders any templ.Component to a string.
func renderComp(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

// TestHeatMap_ForceSquareLayout checks the forceSquare layout squares the
// cells and centers the grid (offsets + laid-out size).
func TestHeatMap_ForceSquareLayout(t *testing.T) {
	res := heatmap.UseHeatMap(heatmap.HeatMapProps{
		Width: 300, Height: 150,
		ForceSquare: true,
		Data:        sampleData(),
	})
	// 3 cols × 3 rows in 300×150 → cell size min(100, 50) = 50.
	if res.Width != 150 || res.Height != 150 {
		t.Errorf("layout size = %vx%v, want 150x150", res.Width, res.Height)
	}
	if res.OffsetX != 75 || res.OffsetY != 0 {
		t.Errorf("layout offset = (%v,%v), want (75,0)", res.OffsetX, res.OffsetY)
	}
	if len(res.Cells) != 9 {
		t.Fatalf("cells = %d, want 9", len(res.Cells))
	}
	if res.Cells[0].Width != 50 || res.Cells[0].Height != 50 {
		t.Errorf("cell size = %vx%v, want 50x50", res.Cells[0].Width, res.Cells[0].Height)
	}
	// Empty data leaves the layout untouched (cols/rows guard).
	empty := heatmap.UseHeatMap(heatmap.HeatMapProps{Width: 300, Height: 150, ForceSquare: true})
	if empty.OffsetX != 0 || empty.Width != 300 {
		t.Errorf("empty forceSquare layout must keep the full area, got offset %v width %v", empty.OffsetX, empty.Width)
	}
	if empty.MinValue != 0 || empty.MaxValue != 0 {
		t.Errorf("empty data domain = [%v,%v], want [0,0]", empty.MinValue, empty.MaxValue)
	}
}

// TestHeatMap_BorderRadiusAndWidth covers the cell border branches.
func TestHeatMap_BorderRadiusAndWidth(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		BorderWidth:  2,
		BorderRadius: 3,
	})
	if got := strings.Count(out, `rx="3" ry="3"`); got != 9 {
		t.Errorf("rounded cell count = %d, want 9", got)
	}
	if got := strings.Count(out, `stroke-width="2"`); got != 9 {
		t.Errorf("bordered cell count = %d, want 9", got)
	}
}

// TestHeatMap_InteractiveTooltips: cells with data emit client tooltips; the
// nil cell must not be hoverable.
func TestHeatMap_InteractiveTooltips(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Interactive: true,
		ChartID:     "hm1",
	})
	if got := strings.Count(out, "data-tc-tooltip"); got != 8 {
		t.Errorf("tooltip cell count = %d, want 8 (nil cell excluded)", got)
	}
	if !strings.Contains(out, "cursor: pointer") {
		t.Errorf("expected hoverable cells to set cursor: pointer")
	}
}

// TestHeatMap_HoveredKeyDimsOthers: the hovered cell keeps the active opacity
// while every other cell dims to the inactive opacity.
func TestHeatMap_HoveredKeyDimsOthers(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		HoveredKey: "Japan.Car",
	})
	if got := strings.Count(out, `fill-opacity="0.15"`); got != 8 {
		t.Errorf("dimmed cell count = %d, want 8", got)
	}
	if got := strings.Count(out, `fill-opacity="1"`); got != 1 {
		t.Errorf("active cell count = %d, want 1", got)
	}
}

// TestHeatMap_DivergingColors selects the diverging color scale.
func TestHeatMap_DivergingColors(t *testing.T) {
	divergeAt := 0.3
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Colors: heatmap.HeatMapColorConfig{Type: "diverging", DivergeAt: &divergeAt},
	})
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
	if got := strings.Count(out, `fill="#`); got < 8 {
		t.Errorf("colored cell count = %d, want >= 8", got)
	}
}

// TestHeatMap_SchemeOnlyColors leaves Type empty with a Scheme set, so
// applyDefaults fills only the type.
func TestHeatMap_SchemeOnlyColors(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Colors: heatmap.HeatMapColorConfig{Scheme: "brown_blueGreen"},
	})
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
}

// TestHeatMap_ValueFormat applies a d3-format spec to the cell labels.
func TestHeatMap_ValueFormat(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		ValueFormat: ".1f",
	})
	if !strings.Contains(out, ">10.0</text>") {
		t.Errorf("expected d3-formatted label 10.0")
	}
}

// TestHeatMap_LabelsDisabled turns cell labels off.
func TestHeatMap_LabelsDisabled(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		EnableLabels: core.BoolPtr(false),
	})
	if strings.Contains(out, ">10</text>") {
		t.Errorf("expected no cell value labels when EnableLabels=false")
	}
}

// TestHeatMap_ContinuousLegend renders the continuous color legend with a
// title, plus a second legend relying on the length/thickness defaults.
func TestHeatMap_ContinuousLegend(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Legends: []heatmap.HeatMapLegend{
			{Anchor: legends.LegendAnchorBottom, Title: "value scale", TranslateY: 30, Length: 150, Thickness: 10},
			{Anchor: legends.LegendAnchorTopRight},
		},
	})
	if got := strings.Count(out, "<linearGradient"); got != 2 {
		t.Errorf("legend gradient count = %d, want 2", got)
	}
	if !strings.Contains(out, ">value scale</text>") {
		t.Errorf("expected the legend title text")
	}
	// Min/max domain labels.
	if !strings.Contains(out, ">10</text>") || !strings.Contains(out, ">90</text>") {
		t.Errorf("expected legend min/max labels 10 and 90")
	}
}

// TestHeatMap_AxesBottomRight covers the bottom/right axis fallbacks (the
// default is top + left).
func TestHeatMap_AxesBottomRight(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		AxisBottom: &axes.AxisProps{Legend: "transport"},
		AxisRight:  &axes.AxisProps{Legend: "country"},
	})
	for _, want := range []string{"transport", "country", ">Train</text>", ">Japan</text>"} {
		if !strings.Contains(out, want) {
			t.Errorf("axes output missing %q", want)
		}
	}
	// Right axis labels anchor toward the start (right of the grid).
	if !strings.Contains(out, `text-anchor="start"`) {
		t.Errorf("expected right-axis labels anchored start")
	}
}

// TestHeatMap_GridLayers enables both grid directions.
func TestHeatMap_GridLayers(t *testing.T) {
	base := heatmap.HeatMapProps{Width: 500, Height: 360, Data: sampleData()}
	without := renderChart(t, base)
	withGrid := base
	withGrid.EnableGridX = core.BoolPtr(true)
	withGrid.EnableGridY = core.BoolPtr(true)
	out := renderChart(t, withGrid)
	if strings.Count(out, "<line") <= strings.Count(without, "<line") {
		t.Errorf("expected grid to add <line> elements")
	}
}

// TestHeatMap_ThemeBackground renders with a custom theme background.
func TestHeatMap_ThemeBackground(t *testing.T) {
	th := theming.DefaultTheme
	th.Background = "#0b0b0b"
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Theme: &th,
	})
	if !strings.Contains(out, `fill="#0b0b0b"`) {
		t.Errorf("expected the themed background rect")
	}
}

// TestHeatMapCell_DirectHxAttrs renders the cell component directly with the
// htmx hover attributes, tooltip, borders, animation, and label.
func TestHeatMapCell_DirectHxAttrs(t *testing.T) {
	out := renderComp(t, heatmap.HeatMapCell(heatmap.HeatMapCellProps{
		Cell: heatmap.ComputedCell{
			XPos: 50, YPos: 40, Width: 20, Height: 20,
			Color: "#336699", Opacity: 1, BorderColor: "#111111",
			Label: "42", LabelTextColor: "#eeeeee",
		},
		BorderWidth:  1,
		BorderRadius: 2,
		EnableLabel:  true,
		HxGet:        "/charts/hm1/hover?cell=a.b",
		HxTrigger:    "mouseenter",
		HxSwap:       "innerHTML",
		HxTarget:     "#tooltip-hm1",
		Tooltip:      "<b>tip</b>",
		Animate:      true,
		AnimateBegin: "0.1s",
	}))
	for _, want := range []string{
		`hx-get="/charts/hm1/hover?cell=a.b"`,
		`hx-trigger="mouseenter"`,
		`hx-swap="innerHTML"`,
		`hx-target="#tooltip-hm1"`,
		`data-tc-tooltip=`,
		`rx="2" ry="2"`,
		`stroke="#111111"`,
		`<animate`,
		`begin="0.1s"`,
		`>42</text>`,
		`fill="#eeeeee"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("HeatMapCell missing %q in output:\n%s", want, out)
		}
	}
}

// TestHeatMap_CanvasBordersAndLabels drives the canvas backend with borders so
// the draw-list contains StrokeRect ("sr") and FillText ("fx") ops.
func TestHeatMap_CanvasBordersAndLabels(t *testing.T) {
	out := renderChart(t, heatmap.HeatMapProps{
		Width: 500, Height: 360, Data: sampleData(),
		Render:      theming.EngineCanvas,
		BorderWidth: 1,
		HoveredKey:  "Japan.Car", // varying opacities exercise the alpha tracking
	})
	for _, want := range []string{`["sr",`, `["fx",`, `["fr",`, `id="tc-heatmap"`} {
		if !strings.Contains(out, want) {
			t.Errorf("canvas draw-list missing %q", want)
		}
	}
}
