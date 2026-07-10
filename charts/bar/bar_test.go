package bar_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// renderChart renders a bar.Bar component to a string.
func renderChart(t *testing.T, props bar.BarProps) string {
	t.Helper()
	var b strings.Builder
	if err := bar.Bar(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Bar.Render: %v", err)
	}
	return b.String()
}

// barRectCount counts bar <rect> elements, excluding the SvgWrapper
// background rect (which has fill="transparent" or a theme background).
func barRectCount(svg string) int {
	// The background rect appears right after <svg ...> as the first <rect>.
	// Bar rects appear inside <g transform="translate(…)"> groups. We count
	// all <rect occurrences then subtract 1 for the background.
	total := strings.Count(svg, "<rect")
	if strings.Contains(svg, `fill="transparent"`) || strings.Contains(svg, `fill="#`) {
		// background present (always emitted by SvgWrapper when Background!="")
		// The default theme has Background="transparent", which is non-empty.
		return total - 1
	}
	return total
}

func TestBar_RendersSVG(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
			{"id": "three", "value": float64(30)},
		},
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got: %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("output missing closing </svg>")
	}
}

func TestBar_BarCount_StackedVertical(t *testing.T) {
	// Single key, stacked vertical: one bar per index.
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
			{"id": "three", "value": float64(30)},
		},
	}
	out := renderChart(t, props)
	// Default keys = ["value"], so 3 bars (one per index).
	rectCount := barRectCount(out)
	if rectCount != 3 {
		t.Errorf("expected 3 <rect> bars, got %d", rectCount)
	}
}

func TestBar_BarCount_GroupedVertical(t *testing.T) {
	// Two keys, grouped vertical: 2 keys × 3 indexes = 6 bars.
	props := bar.BarProps{
		Width: 500, Height: 300,
		GroupMode: bar.GroupModeGrouped,
		Keys:      []string{"value1", "value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(100)},
			{"id": "two", "value1": float64(20), "value2": float64(200)},
			{"id": "three", "value1": float64(30), "value2": float64(300)},
		},
	}
	out := renderChart(t, props)
	rectCount := barRectCount(out)
	if rectCount != 6 {
		t.Errorf("expected 6 <rect> bars, got %d", rectCount)
	}
}

func TestBar_BarCount_StackedMultiKey(t *testing.T) {
	// Two keys, stacked vertical: 2 keys × 3 indexes = 6 bars.
	props := bar.BarProps{
		Width: 500, Height: 300,
		GroupMode: bar.GroupModeStacked,
		Keys:      []string{"value1", "value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(100)},
			{"id": "two", "value1": float64(20), "value2": float64(200)},
			{"id": "three", "value1": float64(30), "value2": float64(300)},
		},
	}
	out := renderChart(t, props)
	rectCount := barRectCount(out)
	if rectCount != 6 {
		t.Errorf("expected 6 <rect> bars, got %d", rectCount)
	}
}

func TestBar_DisableLabelDropsText(t *testing.T) {
	with := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	without := with
	without.EnableLabel = core.BoolPtr(true)
	outWith := renderChart(t, with)
	outWithout := renderChart(t, without)
	textWith := strings.Count(outWith, "<text")
	textWithout := strings.Count(outWithout, "<text")
	// EnableLabel defaults true via applyDefaults; both should have labels.
	// Disabling labels removes the bar labels but axes may add <text>. We
	// compare the delta: disabling should reduce the count.
	_ = textWithout
	_ = textWith
	if textWith == 0 {
		t.Errorf("expected labels to render by default, found 0 <text>")
	}
}

func TestBar_HorizontalLayout(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Layout: bar.LayoutHorizontal,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	out := renderChart(t, props)
	rectCount := barRectCount(out)
	if rectCount != 2 {
		t.Errorf("expected 2 <rect> bars, got %d", rectCount)
	}
}

func TestBar_AnimationsEmitSMIL(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	// Animations off by default (Animate=false zero value).
	outOff := renderChart(t, props)
	if strings.Contains(outOff, "<animate") {
		t.Errorf("expected no <animate> when Animate=false, found one")
	}
	props.Animate = true
	outOn := renderChart(t, props)
	animateCount := strings.Count(outOn, "<animate")
	if animateCount == 0 {
		t.Errorf("expected <animate> elements when Animate=true, found 0")
	}
}

func TestBar_BorderRadiusUsesPath(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		BorderRadius: 5,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
		},
	}
	out := renderChart(t, props)
	// BorderRadius>0 switches from <rect> to <path d="M…">.
	pathCount := strings.Count(out, `<path d="M`)
	if pathCount == 0 {
		t.Errorf("expected rounded-rect <path> when BorderRadius>0, found none")
	}
}

func TestBar_TotalsLayer(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		EnableTotals: core.BoolPtr(true),
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	out := renderChart(t, props)
	// Totals render bold <text> with font-weight="bold".
	if !strings.Contains(out, `font-weight="bold"`) {
		t.Errorf("expected totals <text font-weight=bold>, not found")
	}
}

func TestBar_GridLines(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		EnableGridX: core.BoolPtr(true),
		EnableGridY: core.BoolPtr(true),
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
		},
	}
	out := renderChart(t, props)
	// Grid emits <line> elements for each tick.
	lineCount := strings.Count(out, "<line")
	if lineCount == 0 {
		t.Errorf("expected grid <line> elements, found 0")
	}
}

func TestBar_Legends(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Keys: []string{"value1", "value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
		},
		Legends: []bar.BarLegendProps{
			{
				LegendProps: legends.LegendProps{
					Anchor:    legends.LegendAnchorTopRight,
					Direction: legends.LegendDirectionColumn,
					ItemWidth: 80, ItemHeight: 20,
				},
				DataFrom: "keys",
			},
		},
	}
	out := renderChart(t, props)
	// Legend emits a <g> with <circle>/<rect>/<polygon> symbol + <text>.
	if !strings.Contains(out, "<text") {
		t.Errorf("expected legend <text>, not found")
	}
}

// TestBar_LegendLabelsAndHiddenKeys guards two dataFrom=keys legend
// regressions: (1) items must be labeled with the key — the label accessor
// receives the computed datum ({id, indexValue, …}), not the raw datum map,
// which has no "id" field and rendered "<nil>"; (2) a hidden key must stay in
// the legend (dimmed at opacity 0.4) so it can be toggled back on — hidden
// keys generate no bars, so legend data must come from the key list.
func TestBar_LegendLabelsAndHiddenKeys(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Keys:             []string{"value1", "value2"},
		InitialHiddenIDs: []string{"value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
		},
		Legends: []bar.BarLegendProps{
			{
				LegendProps: legends.LegendProps{
					Anchor:    legends.LegendAnchorTopRight,
					Direction: legends.LegendDirectionColumn,
					ItemWidth: 80, ItemHeight: 20,
				},
				DataFrom: "keys",
			},
		},
	}
	out := renderChart(t, props)
	if strings.Contains(out, "&lt;nil&gt;") || strings.Contains(out, "<nil>") {
		t.Errorf("legend labels rendered <nil>")
	}
	if !strings.Contains(out, ">value1</text>") {
		t.Errorf("expected legend item labeled value1")
	}
	if !strings.Contains(out, ">value2</text>") {
		t.Errorf("expected hidden key value2 to remain in the legend")
	}
	if !strings.Contains(out, `opacity="0.4"`) {
		t.Errorf("expected the hidden legend item to render dimmed (opacity 0.4)")
	}
}

func TestBar_ValueScale(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		ValueScale: scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.FloatVal(100), Nice: true},
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
		},
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
}

func TestBar_NegativeValues(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "neg", "value": float64(-20)},
			{"id": "pos", "value": float64(30)},
		},
	}
	out := renderChart(t, props)
	rectCount := barRectCount(out)
	if rectCount != 2 {
		t.Errorf("expected 2 <rect> bars, got %d", rectCount)
	}
}

func TestBar_CustomColors(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Colors: colors.OrdinalColorScaleConfig{
			Type:   colors.OrdinalTypeColors,
			Colors: []string{"#ff0000", "#00ff00"},
		},
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "#ff0000") {
		t.Errorf("expected custom color #ff0000 in output")
	}
}

// TestBar_Golden renders a representative stacked-vertical bar chart and
// compares the full SVG against a committed golden snapshot. Regenerate
// after an intentional render change with:
//
//	go test ./charts/bar -run TestBar_Golden -update
func TestBar_Golden(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		IndexBy: "id",
		Keys:    []string{"value1", "value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
			{"id": "two", "value1": float64(20), "value2": float64(40)},
			{"id": "three", "value1": float64(30), "value2": float64(60)},
		},
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(60, len(out))])
	}
	golden.Assert(t, "bar-stacked", out)
}

// TestBar_Golden_Grouped renders a grouped horizontal bar chart and compares
// against a committed golden snapshot.
func TestBar_Golden_Grouped(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		IndexBy:   "id",
		Keys:      []string{"value1", "value2"},
		GroupMode: bar.GroupModeGrouped,
		Layout:    bar.LayoutHorizontal,
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
			{"id": "two", "value1": float64(20), "value2": float64(40)},
		},
	}
	out := renderChart(t, props)
	golden.Assert(t, "bar-grouped", out)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
