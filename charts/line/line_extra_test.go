package line_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
)

// renderLineComp renders any templ.Component to a string.
func renderLineComp(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

// TestLine_PointsBorderAndLabels enables point borders and point labels:
// circles get stroke/stroke-width and each point renders its y value.
func TestLine_PointsBorderAndLabels(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:              sampleData(),
		EnablePoints:      true,
		PointBorderWidth:  2,
		EnablePointLabel:  true,
		PointLabelYOffset: -12,
	}
	out := renderChart(t, props)
	// The line paths also use stroke-width="2" (default LineWidth), so count
	// only the point circles.
	if got := strings.Count(out, `stroke-width="2"></circle>`); got != 6 {
		t.Errorf("bordered point count = %d, want 6", got)
	}
	if got := strings.Count(out, `y="-12"`); got != 6 {
		t.Errorf("point label count = %d, want 6", got)
	}
	for _, want := range []string{">10</text>", ">25</text>"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected point label %q", want)
		}
	}
}

// TestLine_AreaBlendModeAndAnimate covers the area blend-mode style attribute
// and the SMIL area enter animation.
func TestLine_AreaBlendModeAndAnimate(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:          sampleData(),
		EnableArea:    true,
		AreaBlendMode: core.MixBlendMultiply,
	}
	props.Animate = true
	out := renderChart(t, props)
	if !strings.Contains(out, "mix-blend-mode: multiply") {
		t.Errorf("expected area mix-blend-mode style")
	}
	// 2 line paths + 2 area paths, each animating "d".
	if got := strings.Count(out, `<animate attributeName="d"`); got != 4 {
		t.Errorf("path d animations = %d, want 4", got)
	}
}

// TestLine_SlicesY renders y-axis hover slices: one full-width rect per
// distinct y position (6 points at distinct pixel ys → 6 slices).
func TestLine_SlicesY(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesY,
	}
	out := renderChart(t, props)
	if got := strings.Count(out, `data-ref="`); got != 6 {
		t.Errorf("y-slice count = %d, want 6", got)
	}
	// Client hover is the default: slices carry their tooltip payloads.
	if !strings.Contains(out, "data-tc-tooltip") {
		t.Errorf("expected client tooltip payloads on y slices")
	}
}

// TestLine_CrosshairOnHover verifies the crosshair layer renders the hover
// crosshair lines when the hover state is set.
func TestLine_CrosshairOnHover(t *testing.T) {
	base := line.LineProps{
		Width: 500, Height: 300,
		Data:            sampleData(),
		Interactive:     true,
		UseMesh:         true,
		EnableCrosshair: true,
	}
	without := renderChart(t, base)
	withHover := base
	withHover.HasHover = true
	withHover.HoverX = 120
	withHover.HoverY = 80
	out := renderChart(t, withHover)
	delta := strings.Count(out, "<line") - strings.Count(without, "<line")
	if delta != 2 {
		t.Errorf("crosshair added %d <line> elements, want 2", delta)
	}
	if !strings.Contains(out, `"120"`) || !strings.Contains(out, `"80"`) {
		t.Errorf("expected crosshair lines positioned at the hover point (120, 80)")
	}
}

// TestLine_CrosshairDebugCentre covers the debug-mesh centre crosshair (no
// active hover, DebugMesh on → crosshair drawn at the chart centre).
func TestLine_CrosshairDebugCentre(t *testing.T) {
	base := line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
		DebugMesh:   true,
	}
	without := renderChart(t, base)
	withCross := base
	withCross.EnableCrosshair = true
	out := renderChart(t, withCross)
	delta := strings.Count(out, "<line") - strings.Count(without, "<line")
	if delta != 2 {
		t.Errorf("debug crosshair added %d <line> elements, want 2", delta)
	}
	// Centre of a 500x300 chart with no margin.
	if !strings.Contains(out, `"250"`) || !strings.Contains(out, `"150"`) {
		t.Errorf("expected debug crosshair at the chart centre (250, 150)")
	}
}

// TestLine_ServerHoverDebugMesh covers the mesh debug rectangle used when the
// legacy htmx path is active (no client voronoi cells available).
func TestLine_ServerHoverDebugMesh(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
		ServerHover: true,
		DebugMesh:   true,
		ChartID:     "m1",
	}
	out := renderChart(t, props)
	if !strings.Contains(out, `stroke-opacity="0.5"`) {
		t.Errorf("expected the debug mesh capture rect to have a visible stroke")
	}
	if !strings.Contains(out, "hover?mesh=1") {
		t.Errorf("expected the htmx mesh round-trip in ServerHover mode")
	}
	if strings.Contains(out, "data-tc-mesh") {
		t.Errorf("ServerHover must not emit the client mesh payload")
	}
}

// TestLine_MarkersLayer renders x + y cartesian markers with legends.
func TestLine_MarkersLayer(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: sampleData(),
		Markers: []core.CartesianMarker{
			{Axis: "y", Value: float64(12), Legend: "threshold", LineColor: "#e25c3b"},
			{Axis: "x", Value: "two", Legend: "release"},
		},
	}
	out := renderChart(t, props)
	for _, want := range []string{"threshold", "release", "#e25c3b"} {
		if !strings.Contains(out, want) {
			t.Errorf("markers output missing %q", want)
		}
	}
}

// TestLine_LegendsCustomItemsAndToggle covers legends with pre-filled Items
// (bypassing the legendData fill) plus the htmx series-toggle wiring that a
// ChartID enables.
func TestLine_LegendsCustomItemsAndToggle(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:    sampleData(),
		ChartID: "lc1",
		Legends: []legends.LegendProps{
			{
				Anchor: legends.LegendAnchorBottomRight, Direction: legends.LegendDirectionRow,
				ItemWidth: 80, ItemHeight: 20,
				Items: []legends.Datum{{ID: "A", Label: "Alpha", Color: "#123456"}},
			},
			{
				Anchor: legends.LegendAnchorTopLeft, Direction: legends.LegendDirectionColumn,
				ItemWidth: 80, ItemHeight: 20,
			},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, ">Alpha</text>") {
		t.Errorf("expected the custom legend item label Alpha")
	}
	// The second legend fills from legendData (reversed → B first).
	if !strings.Contains(out, ">B</text>") {
		t.Errorf("expected the auto-filled legend item B")
	}
	if !strings.Contains(out, "hx-post") {
		t.Errorf("expected the series-toggle hx-post attr when ChartID is set")
	}
}

// TestLine_HiddenSeries verifies InitialHiddenIDs removes the series from the
// plot but keeps it (flagged hidden) in the legend data.
func TestLine_HiddenSeries(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:             sampleData(),
		InitialHiddenIDs: []string{"B"},
		Legends: []legends.LegendProps{
			{Anchor: legends.LegendAnchorTopRight, Direction: legends.LegendDirectionColumn, ItemWidth: 80, ItemHeight: 20},
		},
	}
	out := renderChart(t, props)
	if got := strings.Count(out, `fill="none"`); got != 1 {
		t.Errorf("visible line paths = %d, want 1", got)
	}
	if !strings.Contains(out, ">B</text>") {
		t.Errorf("hidden series B must stay in the legend")
	}
}

// TestLine_CustomFormats applies XFormat/YFormat funcs; the formatted values
// flow into the point labels.
func TestLine_CustomFormats(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:             sampleData(),
		EnablePoints:     true,
		EnablePointLabel: true,
		XFormat:          func(v any) string { return fmt.Sprintf("x=%v", v) },
		YFormat:          func(v any) string { return fmt.Sprintf("y=%v", v) },
	}
	out := renderChart(t, props)
	if !strings.Contains(out, ">y=10</text>") {
		t.Errorf("expected custom-formatted point label y=10")
	}
}

// TestLine_UseSlicesDirect exercises UseSlices' edge branches: nil-data points
// are skipped, unsorted positions are sorted, and empty mode returns nil.
func TestLine_UseSlicesDirect(t *testing.T) {
	points := []line.Point{
		{X: 40, Y: 25, Data: line.PointDatum{X: "c", Y: float64(3)}},
		{X: 10, Y: 20, Data: line.PointDatum{X: "a", Y: float64(1)}},
		{X: 30, Y: 5, Data: line.PointDatum{X: nil, Y: float64(2)}}, // skipped
		{X: 20, Y: 40, Data: line.PointDatum{X: "b", Y: float64(4)}},
	}
	xs := line.UseSlices("t", line.EnableSlicesX, points, 100, 50)
	if len(xs) != 3 {
		t.Fatalf("x slices = %d, want 3 (nil-data point skipped)", len(xs))
	}
	if xs[0].X != 10 || xs[1].X != 20 || xs[2].X != 40 {
		t.Errorf("x slices not sorted: %v %v %v", xs[0].X, xs[1].X, xs[2].X)
	}
	// First slice starts at its own x, last extends to the chart width.
	if xs[0].X0 != 10 {
		t.Errorf("first slice X0 = %v, want 10", xs[0].X0)
	}
	if got := xs[2].X0 + xs[2].Width; got != 100 {
		t.Errorf("last slice must extend to width 100, ends at %v", got)
	}

	ys := line.UseSlices("t", line.EnableSlicesY, points, 100, 50)
	if len(ys) != 3 {
		t.Fatalf("y slices = %d, want 3", len(ys))
	}
	if ys[0].Y != 20 || ys[1].Y != 25 || ys[2].Y != 40 {
		t.Errorf("y slices not sorted: %v %v %v", ys[0].Y, ys[1].Y, ys[2].Y)
	}
	if got := ys[2].Y0 + ys[2].Height; got != 50 {
		t.Errorf("last y slice must extend to height 50, ends at %v", got)
	}

	if got := line.UseSlices("t", line.EnableSlicesFalse, points, 100, 50); got != nil {
		t.Errorf("EnableSlicesFalse must return nil, got %v", got)
	}
}

// TestLine_CustomAxes provides explicit axis props (instead of the nil-axis
// defaults) so the custom-axis merge path is exercised.
func TestLine_CustomAxes(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:       sampleData(),
		AxisBottom: &axes.AxisProps{Legend: "bottom legend", TickSize: 5, TickPadding: 5},
		AxisLeft:   &axes.AxisProps{Legend: "left legend", TicksPosition: "before", LegendPosition: axes.AxisLegendMiddle},
	}
	out := renderChart(t, props)
	for _, want := range []string{"bottom legend", "left legend"} {
		if !strings.Contains(out, want) {
			t.Errorf("axes output missing %q", want)
		}
	}
}

// TestLine_UseLineZeroProps calls UseLine directly with minimal props so its
// own default merging (scales/colors/curve/point colors/hidden ids) runs
// without the Line component's applyDefaults.
func TestLine_UseLineZeroProps(t *testing.T) {
	res := line.UseLine(line.LineProps{
		Width: 200, Height: 100,
		Data: sampleData(),
	})
	if len(res.Series) != 2 {
		t.Fatalf("series = %d, want 2", len(res.Series))
	}
	if len(res.Points) != 6 {
		t.Errorf("points = %d, want 6", len(res.Points))
	}
	if res.Series[0].Color == "" {
		t.Errorf("expected a default scheme color for series A")
	}
	if got := res.LineGenerator([]line.PointXY{{X: 0, Y: 0}, {X: 10, Y: 10}}); !strings.HasPrefix(got, "M") {
		t.Errorf("LineGenerator path = %q, want M…", got)
	}
	if got := res.AreaGenerator([]line.PointXY{{X: 0, Y: 0}, {X: 10, Y: 10}}); !strings.HasPrefix(got, "M") {
		t.Errorf("AreaGenerator path = %q, want M…", got)
	}
	// legendData is reversed (B first) and nothing is hidden.
	if len(res.LegendData) != 2 || res.LegendData[0].ID != "B" || res.LegendData[0].Hidden {
		t.Errorf("unexpected legend data: %+v", res.LegendData)
	}
}

// TestLine_LinesItemComponent renders a single LinesItem directly.
func TestLine_LinesItemComponent(t *testing.T) {
	gen := func(points []line.PointXY) string { return "M0,0L10,10" }
	out := renderLineComp(t, line.LinesItem(line.LinesItemProps{
		Points:        []line.PointXY{{X: 0, Y: 0}, {X: 10, Y: 10}},
		LineGenerator: gen,
		Color:         "#c0ffee",
		Thickness:     3,
		Animate:       true,
	}))
	for _, want := range []string{
		`d="M0,0L10,10"`,
		`stroke="#c0ffee"`,
		`stroke-width="3"`,
		`<animate attributeName="d"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("LinesItem missing %q in output:\n%s", want, out)
		}
	}
}

// TestLine_PointTooltipComponent renders the PointTooltip component directly.
func TestLine_PointTooltipComponent(t *testing.T) {
	out := renderLineComp(t, line.PointTooltip(line.PointTooltipProps{
		Point: line.Point{
			SeriesColor: "#ff0000",
			Data:        line.PointDatum{XFormatted: "one", YFormatted: "10"},
		},
	}))
	for _, want := range []string{
		"nivo-tooltip-basic",
		"background: #ff0000",
		"<strong>one</strong>",
		"<strong>10</strong>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("PointTooltip missing %q in output:\n%s", want, out)
		}
	}
}

// TestLine_SliceTooltipComponent renders SliceTooltip for both slice axes:
// x-slices show each point's y value; y-slices show the x value.
func TestLine_SliceTooltipComponent(t *testing.T) {
	slice := line.SliceData{
		Points: []line.Point{
			{SeriesID: "A", SeriesColor: "#00ff00", Data: line.PointDatum{XFormatted: "one", YFormatted: "10"}},
			{SeriesID: "B", SeriesColor: "#0000ff", Data: line.PointDatum{XFormatted: "one", YFormatted: "5"}},
		},
	}
	outX := renderLineComp(t, line.SliceTooltip(line.SliceTooltipProps{Slice: slice, Axis: line.EnableSlicesX}))
	for _, want := range []string{"nivo-tooltip-table", "A", "B", ">10</td>", ">5</td>", "background: #00ff00"} {
		if !strings.Contains(outX, want) {
			t.Errorf("x-axis SliceTooltip missing %q in output:\n%s", want, outX)
		}
	}
	outY := renderLineComp(t, line.SliceTooltip(line.SliceTooltipProps{Slice: slice, Axis: line.EnableSlicesY}))
	if !strings.Contains(outY, ">one</td>") {
		t.Errorf("y-axis SliceTooltip must show the x-formatted value, got:\n%s", outY)
	}
}
