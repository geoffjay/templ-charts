package bar_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/core"
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

// TestBarItem_RoundedRectFullOptions exercises the <path> variant of BarItem
// with every optional attribute enabled: border, aria-label, htmx hover
// attrs, dimmed opacity, and the SMIL enter animation.
func TestBarItem_RoundedRectFullOptions(t *testing.T) {
	props := bar.BarItemProps{
		Bar: bar.ComputedBarDatum{
			Key: "value.one", X: 10, Y: 20, Width: 40, Height: 60,
			Data: bar.ComputedDatum{ID: "value", IndexValue: "one", Value: float64(10)},
		},
		Label:             "10",
		ShouldRenderLabel: true,
		BorderRadius:      4,
		BorderWidth:       2,
		BorderColor:       "#112233",
		Fill:              "#abcdef",
		LabelColor:        "#000000",
		LabelLayout:       bar.BarLabelLayout{LabelX: 20, LabelY: 30, TextAnchor: "middle"},
		Animate:           true,
		Dimmed:            true,
		HxGet:             "/charts/c1/hover?bar=value.one",
		HxTrigger:         "mouseenter",
		HxSwap:            "innerHTML",
		HxTarget:          "#tooltip-c1",
		AriaLabel:         "bar one",
	}
	out := renderComp(t, bar.BarItem(props))
	for _, want := range []string{
		`<path d="M`,
		`fill="#abcdef"`,
		`fill-opacity="0.5"`,
		`stroke="#112233"`,
		`stroke-width="2"`,
		`aria-label="bar one"`,
		`hx-get="/charts/c1/hover?bar=value.one"`,
		`hx-trigger="mouseenter"`,
		`hx-swap="innerHTML"`,
		`hx-target="#tooltip-c1"`,
		`<animate attributeName="height"`,
		`text-anchor="middle"`,
		`>10</text>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("rounded BarItem missing %q in output:\n%s", want, out)
		}
	}
}

// TestBarItem_RectFullOptions exercises the <rect> variant with borders,
// aria-label, htmx attrs, dim, label, and the horizontal enter animation.
func TestBarItem_RectFullOptions(t *testing.T) {
	props := bar.BarItemProps{
		Bar: bar.ComputedBarDatum{
			Key: "value.two", X: 0, Y: 0, Width: 80, Height: 20,
			Data: bar.ComputedDatum{ID: "value", IndexValue: "two", Value: float64(20)},
		},
		Label:             "20",
		ShouldRenderLabel: true,
		BorderWidth:       1.5,
		BorderColor:       "#445566",
		Fill:              "#fedcba",
		LabelColor:        "#111111",
		LabelLayout:       bar.BarLabelLayout{LabelX: 40, LabelY: 10, TextAnchor: "start"},
		Animate:           true,
		Horizontal:        true,
		Dimmed:            true,
		HxGet:             "/charts/c2/hover?bar=value.two",
		HxTrigger:         "mouseenter",
		HxSwap:            "innerHTML",
		HxTarget:          "#tooltip-c2",
		AriaLabel:         "bar two",
	}
	out := renderComp(t, bar.BarItem(props))
	for _, want := range []string{
		`<rect width="80" height="20"`,
		`fill-opacity="0.5"`,
		`stroke="#445566"`,
		`stroke-width="1.5"`,
		`aria-label="bar two"`,
		`hx-get="/charts/c2/hover?bar=value.two"`,
		`hx-trigger="mouseenter"`,
		`hx-swap="innerHTML"`,
		`hx-target="#tooltip-c2"`,
		`<animate attributeName="width" from="0" to="80"`,
		`>20</text>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("rect BarItem missing %q in output:\n%s", want, out)
		}
	}
}

// TestBar_HoverDimsSiblings renders an interactive chart with a hovered key:
// the hovered key's bars stay opaque while sibling keys dim, and each bar
// emits the htmx hover attributes.
func TestBar_HoverDimsSiblings(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Keys:        []string{"value1", "value2"},
		Interactive: true,
		ChartID:     "c1",
		HoveredKey:  "value1.one", // bar keys are "<key>.<indexValue>"
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
			{"id": "two", "value1": float64(20), "value2": float64(40)},
		},
	}
	out := renderChart(t, props)
	// The three non-hovered bars must be dimmed.
	if got := strings.Count(out, `fill-opacity="0.5"`); got != 3 {
		t.Errorf("dimmed bar count = %d, want 3", got)
	}
	if !strings.Contains(out, `hx-get="/charts/c1/hover?bar=value1.one"`) {
		t.Errorf("expected htmx hover attr for the value1.one bar")
	}
	if !strings.Contains(out, `hx-target="#tooltip-c1"`) {
		t.Errorf("expected htmx hover target #tooltip-c1")
	}
}

// TestBar_HorizontalAnimatedRounded combines horizontal layout, animation and
// border radius: the rounded <path> bars must animate width (not height).
func TestBar_HorizontalAnimatedRounded(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Layout:       bar.LayoutHorizontal,
		BorderRadius: 3,
		BorderWidth:  1,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	props.Animate = true
	out := renderChart(t, props)
	if !strings.Contains(out, `<path d="M`) {
		t.Fatalf("expected rounded-rect <path> bars")
	}
	if got := strings.Count(out, `<animate attributeName="width"`); got != 2 {
		t.Errorf("horizontal width animations = %d, want 2", got)
	}
	if strings.Contains(out, `<animate attributeName="height"`) {
		t.Errorf("horizontal bars must not animate height")
	}
	// Default BorderColor inherits from the bar color, so a stroke must appear.
	if !strings.Contains(out, `stroke=`) {
		t.Errorf("expected stroked bars when BorderWidth > 0")
	}
}

// TestBar_MarkersLayer verifies the markers layer renders both an x-axis
// (band) and a y-axis (linear) marker with their legends.
func TestBar_MarkersLayer(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		Markers: []core.CartesianMarker{
			{Axis: "y", Value: float64(15), Legend: "y target", LineColor: "#e25c3b", LineStrokeWidth: 2},
			{Axis: "x", Value: "one", Legend: "x mark"},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "y target") {
		t.Errorf("expected y marker legend text")
	}
	if !strings.Contains(out, "x mark") {
		t.Errorf("expected x marker legend text")
	}
	if !strings.Contains(out, "#e25c3b") {
		t.Errorf("expected marker line color #e25c3b")
	}
}

// TestBar_AnnotationsLayer binds a circle annotation to one bar and checks
// the annotation note + outline render.
func TestBar_AnnotationsLayer(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		Annotations: []annotations.AnnotationSpec[bar.ComputedBarDatum]{
			{
				Type: annotations.AnnotationTypeCircle,
				Circle: &annotations.CircleAnnotationSpec[bar.ComputedBarDatum]{
					Match:  func(b bar.ComputedBarDatum) bool { return b.Data.IndexValue == "one" },
					Radius: 12,
					Note:   "peak note",
					NoteX:  30,
					NoteY:  -20,
				},
			},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "peak note") {
		t.Errorf("expected annotation note text in output")
	}
	if !strings.Contains(out, "<circle") {
		t.Errorf("expected annotation circle outline")
	}
}

// TestBar_AxesRendered configures bottom + left axes (one relying on the
// resolveAxis defaults, one overriding ticks/legend position) and checks tick
// labels + axis legends render.
func TestBar_AxesRendered(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		AxisBottom: &axes.AxisProps{Legend: "bottom legend", TickSize: 5, TickPadding: 5},
		AxisLeft: &axes.AxisProps{
			Legend:         "left legend",
			TicksPosition:  "before",
			LegendPosition: axes.AxisLegendMiddle,
			TickSize:       5,
			TickPadding:    5,
		},
	}
	out := renderChart(t, props)
	for _, want := range []string{"bottom legend", "left legend"} {
		if !strings.Contains(out, want) {
			t.Errorf("axes output missing %q", want)
		}
	}
	// Both axes emit their domain <line>.
	if got := strings.Count(out, "<line"); got < 2 {
		t.Errorf("axis <line> count = %d, want >= 2", got)
	}
}

// TestBar_TotalsHorizontalWithTheme renders horizontal totals with a themed
// label style so the totals text picks up font-size/family/fill and the
// horizontal anchor/baseline.
func TestBar_TotalsHorizontalWithTheme(t *testing.T) {
	th := theming.DefaultTheme
	th.Background = "#f8f8f8"
	th.Labels.Text.FontSize = 14 // int on purpose: exercises the any → float coercion
	th.Labels.Text.FontFamily = "monospace"
	th.Labels.Text.Fill = "#ff00aa"
	props := bar.BarProps{
		Width: 500, Height: 300,
		Layout:       bar.LayoutHorizontal,
		EnableTotals: core.BoolPtr(true),
		Theme:        &th,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
	}
	out := renderChart(t, props)
	for _, want := range []string{
		`font-weight="bold"`,
		`font-size="14"`,
		`font-family="monospace"`,
		`fill="#ff00aa"`,
		`text-anchor="start"`,        // horizontal totals anchor
		`dominant-baseline="middle"`, // horizontal totals baseline
		`fill="#f8f8f8"`,             // themed background rect
	} {
		if !strings.Contains(out, want) {
			t.Errorf("totals output missing %q", want)
		}
	}
}

// TestBar_LegendsDataFromIndexes checks dataFrom=indexes legends label items
// with the index values instead of the keys.
func TestBar_LegendsDataFromIndexes(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		Legends: []bar.BarLegendProps{
			{
				LegendProps: legends.LegendProps{
					Anchor:    legends.LegendAnchorBottomRight,
					Direction: legends.LegendDirectionRow,
					ItemWidth: 60, ItemHeight: 18,
				},
				DataFrom: "indexes",
			},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, ">one</text>") || !strings.Contains(out, ">two</text>") {
		t.Errorf("expected legend items labeled with index values one/two")
	}
}

// TestBar_LegendsCustomLegendLabel checks a custom LegendLabel accessor is
// applied to dataFrom=keys legend items.
func TestBar_LegendsCustomLegendLabel(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Keys: []string{"value1", "value2"},
		Data: []bar.BarDatum{
			{"id": "one", "value1": float64(10), "value2": float64(20)},
		},
		LegendLabel: func(d map[string]any) string { return "K:" + fmt.Sprint(d["id"]) },
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
	if !strings.Contains(out, ">K:value1</text>") || !strings.Contains(out, ">K:value2</text>") {
		t.Errorf("expected custom legend labels K:value1 / K:value2, got:\n%s", out)
	}
}

// TestBar_UseBarCustomAccessors drives UseBar directly with custom label /
// tooltip / format accessors, grouped horizontal layout, inner padding,
// colorBy=indexValue, label skipping, and totals.
func TestBar_UseBarCustomAccessors(t *testing.T) {
	props := bar.BarProps{
		Width: 400, Height: 300,
		GroupMode:    bar.GroupModeGrouped,
		Layout:       bar.LayoutHorizontal,
		InnerPadding: 2,
		ColorBy:      bar.ColorByIndexValue,
		Keys:         []string{"v1", "v2"},
		Data: []bar.BarDatum{
			{"id": "one", "v1": float64(10), "v2": float64(20)},
			{"id": "two", "v1": float64(30), "v2": float64(40)},
		},
		ValueFormat:     func(v float64) string { return fmt.Sprintf("%.1f!", v) },
		Label:           func(d bar.ComputedDatum) string { return "L:" + d.ID },
		TooltipLabel:    func(d bar.ComputedDatum) string { return "T:" + d.ID },
		EnableLabel:     core.BoolPtr(true),
		LabelSkipWidth:  10,
		LabelSkipHeight: 10,
		EnableTotals:    core.BoolPtr(true),
	}
	res := bar.UseBar(props)
	if got := res.FormatValue(1); got != "1.0!" {
		t.Errorf("FormatValue(1) = %q, want %q", got, "1.0!")
	}
	if got := res.GetLabel(bar.ComputedDatum{ID: "v1"}); got != "L:v1" {
		t.Errorf("GetLabel = %q, want L:v1", got)
	}
	if got := res.GetTooltipLabel(bar.ComputedDatum{ID: "v1"}); got != "T:v1" {
		t.Errorf("GetTooltipLabel = %q, want T:v1", got)
	}
	if res.ShouldRenderLabel(5, 100) {
		t.Errorf("labels must be skipped when width < LabelSkipWidth")
	}
	if res.ShouldRenderLabel(100, 5) {
		t.Errorf("labels must be skipped when height < LabelSkipHeight")
	}
	if !res.ShouldRenderLabel(100, 100) {
		t.Errorf("labels must render when above both skip thresholds")
	}
	if len(res.BarTotals) == 0 {
		t.Errorf("expected computed BarTotals when EnableTotals is set")
	}
	if len(res.BarsWithValue) != 4 {
		t.Errorf("BarsWithValue = %d, want 4", len(res.BarsWithValue))
	}
}

// TestBar_UseBarZeroProps calls UseBar directly with minimal props so UseBar's
// own default merging (keys/groupMode/layout/labelPosition, empty legend
// DataFrom) is exercised without going through the Bar component's
// applyDefaults.
func TestBar_UseBarZeroProps(t *testing.T) {
	res := bar.UseBar(bar.BarProps{
		Width: 200, Height: 100,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		Legends: []bar.BarLegendProps{
			{LegendProps: legends.LegendProps{ItemWidth: 50, ItemHeight: 20}}, // DataFrom "" → "keys"
		},
	})
	if len(res.Bars) != 2 {
		t.Fatalf("bars = %d, want 2 (default keys [value])", len(res.Bars))
	}
	if len(res.LegendsWithData) != 1 || len(res.LegendsWithData[0].Data) != 1 {
		t.Fatalf("expected 1 legend with 1 item, got %+v", res.LegendsWithData)
	}
	if res.LegendsWithData[0].Data[0].ID != "value" {
		t.Errorf("legend item id = %q, want value", res.LegendsWithData[0].Data[0].ID)
	}
	// Default keys/layout/groupMode produce vertical stacked bars keyed value.<index>.
	if res.Bars[0].Key != "value.one" {
		t.Errorf("bar key = %q, want value.one", res.Bars[0].Key)
	}
}

// TestBar_RectAnnotation binds a rect annotation, which resolves the bar's
// dimensions (not just its center) through bindBarAnnotations.
func TestBar_RectAnnotation(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
		},
		Annotations: []annotations.AnnotationSpec[bar.ComputedBarDatum]{
			{
				Type: annotations.AnnotationTypeRect,
				Rect: &annotations.RectAnnotationSpec[bar.ComputedBarDatum]{
					Match: func(b bar.ComputedBarDatum) bool { return true },
					Note:  "rect note",
					NoteX: 25,
					NoteY: -15,
				},
			},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "rect note") {
		t.Errorf("expected rect annotation note text")
	}
}

// TestBar_DefsFillBars verifies Defs + Fill rules bind to bars: the gradient
// def lands in <defs> and matched bars use the url(#…) fill override.
func TestBar_DefsFillBars(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
		},
		Defs: []core.Def{
			core.LinearGradientDef("barGrad", []core.GradientStop{
				{Offset: 0, Color: "inherit", Opacity: 1},
				{Offset: 100, Color: "inherit", Opacity: 0.1},
			}, nil),
		},
		Fill: []core.DefRule{{ID: "barGrad", Match: "*"}},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "<linearGradient") {
		t.Errorf("expected a <linearGradient> def")
	}
	if !strings.Contains(out, `fill="url(#barGrad`) {
		t.Errorf("expected bar fill url(#barGrad…) override")
	}
}

// TestBar_LabelSkipInRender checks that label skipping removes the bar value
// labels in the rendered SVG (no bar <text> equals the formatted value).
func TestBar_LabelSkipInRender(t *testing.T) {
	props := bar.BarProps{
		Width: 500, Height: 300,
		EnableLabel:    core.BoolPtr(true),
		LabelSkipWidth: 10000, // skip everything
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
		},
	}
	out := renderChart(t, props)
	if strings.Contains(out, `>10</text>`) {
		t.Errorf("expected the bar value label to be skipped")
	}
}
