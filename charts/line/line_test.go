package line_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []line.LineSeries {
	return []line.LineSeries{
		{ID: "A", Data: []line.LinePointData{
			{X: "one", Y: float64(10)},
			{X: "two", Y: float64(20)},
			{X: "three", Y: float64(30)},
		}},
		{ID: "B", Data: []line.LinePointData{
			{X: "one", Y: float64(5)},
			{X: "two", Y: float64(15)},
			{X: "three", Y: float64(25)},
		}},
	}
}

func renderChart(t *testing.T, props line.LineProps) string {
	t.Helper()
	var b strings.Builder
	if err := line.Line(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Line.Render: %v", err)
	}
	return b.String()
}

func TestLine_RendersSVG(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: sampleData(),
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got: %q", out[:minLen(out)])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("output missing closing </svg>")
	}
}

func TestLine_PathCount(t *testing.T) {
	// 2 series → 2 line paths. Disable grid so fill="none" counts only lines.
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		EnableGridX: core.BoolPtr(false),
		EnableGridY: core.BoolPtr(false),
	}
	out := renderChart(t, props)
	// Each line is a <path d="M…" fill="none" stroke=…>.
	pathCount := strings.Count(out, `fill="none"`)
	if pathCount != 2 {
		t.Errorf("expected 2 line paths, got %d", pathCount)
	}
}

func TestLine_AreaEnabled(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:       sampleData(),
		EnableArea: core.BoolPtr(true),
	}
	out := renderChart(t, props)
	// Areas emit <path fill="…" fill-opacity="…">. The fill-opacity attr is
	// unique to area paths (line paths have fill="none").
	if !strings.Contains(out, `fill-opacity`) {
		t.Errorf("expected area fill-opacity attribute, not found")
	}
}

func TestLine_PointsEnabled(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		EnablePoints: core.BoolPtr(true),
	}
	out := renderChart(t, props)
	// 6 points (2 series × 3 x-values) → 6 <circle> elements.
	circleCount := strings.Count(out, "<circle")
	if circleCount != 6 {
		t.Errorf("expected 6 <circle> points, got %d", circleCount)
	}
}

func TestLine_PointsDisabled(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		EnablePoints: core.BoolPtr(false),
	}
	out := renderChart(t, props)
	if strings.Contains(out, "<circle") {
		t.Errorf("expected no <circle> when EnablePoints=core.BoolPtr(false), found one")
	}
}

func TestLine_GridLines(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		EnableGridX: core.BoolPtr(true),
		EnableGridY: core.BoolPtr(true),
	}
	out := renderChart(t, props)
	lineCount := strings.Count(out, "<line")
	if lineCount == 0 {
		t.Errorf("expected grid <line> elements, found 0")
	}
}

func TestLine_AnimationsEmitSMIL(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: sampleData(),
	}
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

func TestLine_CustomColors(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: sampleData(),
		Colors: colors.OrdinalColorScaleConfig{
			Type:   colors.OrdinalTypeColors,
			Colors: []string{"#ff0000", "#00ff00"},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "#ff0000") {
		t.Errorf("expected custom color #ff0000 in output")
	}
}

func TestLine_Curve(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:  sampleData(),
		Curve: core.CurveMonotoneX,
	}
	out := renderChart(t, props)
	// A monotone curve should produce a path with C (cubic bezier) commands.
	if !strings.Contains(out, "C") {
		t.Errorf("expected cubic-bezier path commands for monotoneX curve")
	}
}

func TestLine_MultipleSeries(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: []line.LineSeries{
			{ID: "A", Data: []line.LinePointData{{X: "a", Y: float64(1)}}},
			{ID: "B", Data: []line.LinePointData{{X: "a", Y: float64(2)}}},
			{ID: "C", Data: []line.LinePointData{{X: "a", Y: float64(3)}}},
		},
		EnableGridX: core.BoolPtr(false),
		EnableGridY: core.BoolPtr(false),
	}
	out := renderChart(t, props)
	pathCount := strings.Count(out, `fill="none"`)
	if pathCount != 3 {
		t.Errorf("expected 3 line paths, got %d", pathCount)
	}
}

func TestLine_Legends(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: sampleData(),
		Legends: []legends.LegendProps{
			{
				Anchor: legends.LegendAnchorTopRight, Direction: legends.LegendDirectionColumn,
				ItemWidth: 80, ItemHeight: 20,
			},
		},
	}
	out := renderChart(t, props)
	// Legend emits a <g> with a symbol + <text>.
	if !strings.Contains(out, "<text") {
		t.Errorf("expected legend <text>, not found")
	}
}

func TestLine_SlicesX(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesX,
	}
	out := renderChart(t, props)
	// Slices emit transparent <rect> with data-ref="slice:…".
	if !strings.Contains(out, `data-ref="slice:`) {
		t.Errorf("expected slice rects with data-ref, not found")
	}
}

func TestLine_ClientHoverMesh(t *testing.T) {
	out := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
		ClientHover: true,
	})
	if !strings.Contains(out, "data-tc-mesh") {
		t.Errorf("client-hover mesh should emit data-tc-mesh")
	}
	if strings.Contains(out, "hover?mesh=1") {
		t.Errorf("client-hover mesh should not emit the htmx mesh round-trip")
	}
}

func TestLine_ClientHoverSlices(t *testing.T) {
	out := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesX,
		ClientHover:  true,
	})
	if !strings.Contains(out, "data-tc-tooltip") {
		t.Errorf("client-hover slices should emit data-tc-tooltip")
	}
	if strings.Contains(out, "/slice?axis=") {
		t.Errorf("client-hover slices should not emit the htmx slice round-trip")
	}
}

func TestLine_ClientHoverIsDefault(t *testing.T) {
	// the per-mousemove server round-trip is retired: with a ChartID set but no
	// ServerHover, mesh/slice hover defaults to the client path (data-tc-*), not
	// the htmx round-trip.
	mesh := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
		ChartID:     "abc",
	})
	if !strings.Contains(mesh, "data-tc-mesh") {
		t.Errorf("mesh should default to client data-tc-mesh")
	}
	if strings.Contains(mesh, "hover?mesh=1") {
		t.Errorf("mesh should not emit the htmx round-trip by default (retired)")
	}
	slices := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesX,
		ChartID:      "abc",
	})
	if !strings.Contains(slices, "data-tc-tooltip") {
		t.Errorf("slices should default to client data-tc-tooltip")
	}
	if strings.Contains(slices, "/slice?axis=") {
		t.Errorf("slices should not emit the htmx round-trip by default (retired)")
	}
}

func TestLine_ServerHoverOptIn(t *testing.T) {
	// ServerHover opts back into the legacy htmx per-mousemove round-trip.
	mesh := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
		ChartID:     "abc",
		ServerHover: true,
	})
	if !strings.Contains(mesh, "hover?mesh=1") {
		t.Errorf("ServerHover mesh should emit the htmx round-trip")
	}
	if strings.Contains(mesh, "data-tc-mesh=") {
		t.Errorf("ServerHover mesh should not emit client data-tc-mesh")
	}
	slices := renderChart(t, line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesX,
		ChartID:      "abc",
		ServerHover:  true,
	})
	if !strings.Contains(slices, "/slice?axis=") {
		t.Errorf("ServerHover slices should emit the htmx round-trip")
	}
}

func TestLine_SlicesDebug(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		Interactive:  true,
		EnableSlices: line.EnableSlicesX,
		DebugSlices:  true,
	}
	out := renderChart(t, props)
	// Debug slices get a red stroke.
	if !strings.Contains(out, `stroke="red"`) {
		t.Errorf("expected debug slice red stroke, not found")
	}
}

func TestLine_Mesh(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		Interactive: true,
		UseMesh:     true,
	}
	out := renderChart(t, props)
	// Mesh emits an overlay <rect> with hx-get when ChartID is set; without
	// ChartID it still emits the rect (transparent). We just check a rect
	// exists with fill-opacity="0".
	if !strings.Contains(out, `fill-opacity="0"`) && !strings.Contains(out, `fill-opacity="0.0"`) {
		t.Errorf("expected mesh overlay rect with fill-opacity=0, not found")
	}
}

func TestLine_MeshDetectionRadiusAndDebugCells(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:            sampleData(),
		Interactive:     true,
		UseMesh:         true,
		ClientHover:     true,
		DetectionRadius: 40,
		DebugMesh:       true,
	}
	out := renderChart(t, props)
	if !strings.Contains(out, `data-tc-mesh-radius="40"`) {
		t.Errorf("expected data-tc-mesh-radius from DetectionRadius")
	}
	// Debug renders the actual voronoi cells (a closed faint path) rather than a
	// bare debug rectangle.
	if !strings.Contains(out, `stroke-opacity="0.35"`) || !strings.Contains(out, `d="M`) {
		t.Errorf("expected voronoi debug cells path in client-hover debug mesh")
	}
}

func TestLine_NilValues(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data: []line.LineSeries{
			{ID: "A", Data: []line.LinePointData{
				{X: "a", Y: float64(1)},
				{X: "b", Y: nil},
				{X: "c", Y: float64(3)},
			}},
		},
		EnablePoints: core.BoolPtr(true),
		EnableGridX:  core.BoolPtr(false),
		EnableGridY:  core.BoolPtr(false),
	}
	out := renderChart(t, props)
	// The line should still render (with a gap at the nil point).
	pathCount := strings.Count(out, `fill="none"`)
	if pathCount != 1 {
		t.Errorf("expected 1 line path with nil gap, got %d", pathCount)
	}
	// Only 2 points (the nil one is filtered out).
	circleCount := strings.Count(out, "<circle")
	if circleCount != 2 {
		t.Errorf("expected 2 <circle> points (nil filtered), got %d", circleCount)
	}
}

func TestLine_YScaleMax(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Data:   sampleData(),
		YScale: scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.FloatVal(100)},
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
}

// TestLine_Golden renders a multi-series line chart with monotoneX curve +
// points and compares the full SVG against a committed golden snapshot.
// Regenerate after an intentional render change with:
//
//	go test ./charts/line -run TestLine_Golden -update
func TestLine_Golden(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		Curve:        core.CurveMonotoneX,
		EnablePoints: core.BoolPtr(true),
		Data:         sampleData(),
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:minLen(out)])
	}
	golden.Assert(t, "line-multi", out)
}

// TestLine_Golden_Area renders an area + points line chart and compares
// against a committed golden snapshot.
func TestLine_Golden_Area(t *testing.T) {
	props := line.LineProps{
		Width:        500,
		Height:       300,
		Curve:        core.CurveMonotoneX,
		EnableArea:   core.BoolPtr(true),
		AreaOpacity:  0.2,
		EnablePoints: core.BoolPtr(true),
		Data:         sampleData(),
	}
	out := renderChart(t, props)
	golden.Assert(t, "line-area", out)
}

// TestLine_DefsFillAreas verifies that Defs + Fill rules bind to the area
// paths: the gradient def is emitted in <defs> and each area's fill is the
// url(#…) override (with "inherit" stops expanded per series color).
func TestLine_DefsFillAreas(t *testing.T) {
	props := line.LineProps{
		Width: 500, Height: 300,
		EnableArea:  core.BoolPtr(true),
		AreaOpacity: 1,
		Data:        sampleData(),
		Defs: []core.Def{
			core.LinearGradientDef("areaGrad", []core.GradientStop{
				{Offset: 0, Color: "inherit", Opacity: 1},
				{Offset: 100, Color: "inherit", Opacity: 0.05},
			}, nil),
		},
		Fill: []core.DefRule{{ID: "areaGrad", Match: "*"}},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "<linearGradient") {
		t.Errorf("expected a <linearGradient> def, not found")
	}
	if !strings.Contains(out, `fill="url(#areaGrad.`) {
		t.Errorf("expected area fill url(#areaGrad.…) override, not found in output")
	}
	// Line strokes must not pick up the fill override.
	if !strings.Contains(out, `fill="none"`) {
		t.Errorf("expected line paths to keep fill=\"none\"")
	}
}

func minLen(s string) int {
	if len(s) < 50 {
		return len(s)
	}
	return 50
}
