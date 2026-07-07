package axes_test

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

func linearScale(min, max, size float64) scales.Scale {
	return scales.NewLinearScaleWithRange(min, max, 0, size)
}

func bandScale(domain []string, size float64) scales.Scale {
	return scales.NewBandScaleWithRange(domain, 0, size, 0, false)
}

func timeScale(t *testing.T) scales.Scale {
	t.Helper()
	t0, _ := time.Parse("2006-01-02", "2024-01-01")
	t1, _ := time.Parse("2006-01-02", "2024-12-31")
	data := scales.ComputedSerieAxis{All: []any{t0, t1}, Min: t0, Max: t1}
	return scales.ComputeScale(scales.ScaleTimeSpec{}, data, 100, scales.ScaleAxisX)
}

func baseXAxisProps() axes.AxisProps {
	return axes.AxisProps{
		Axis:        "x",
		Scale:       linearScale(0, 100, 100),
		Length:      100,
		TickSize:    5,
		TickPadding: 5,
		TickValues:  []any{0.0, 50.0, 100.0},
	}
}

func TestAxis_XLinearDefaultTheme(t *testing.T) {
	out := render(t, axes.Axis(baseXAxisProps(), &theming.DefaultTheme))

	// Root group translated to the axis origin.
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("missing root translate: %s", out)
	}
	// Domain line spans the axis length horizontally, styled from the theme.
	if !strings.Contains(out, `x2="100"`) {
		t.Errorf("domain line should end at x2=100: %s", out)
	}
	if !strings.Contains(out, `stroke="transparent"`) {
		t.Errorf("domain line should use theme stroke: %s", out)
	}
	// Tick lines use the theme's tick line style.
	if got := strings.Count(out, `stroke="#777777"`); got != 3 {
		t.Errorf("tick line stroke count = %d, want 3", got)
	}
	// One label per tick with the theme's text style.
	if got := strings.Count(out, "<text"); got != 3 {
		t.Errorf("text count = %d, want 3", got)
	}
	for _, label := range []string{">0</text>", ">50</text>", ">100</text>"} {
		if !strings.Contains(out, label) {
			t.Errorf("missing tick label %q", label)
		}
	}
	if got := strings.Count(out, `fill="#333333"`); got != 3 {
		t.Errorf("tick text fill count = %d, want 3", got)
	}
	if !strings.Contains(out, `font-size="11"`) {
		t.Errorf("tick text should inherit font-size 11: %s", out)
	}
	if !strings.Contains(out, `font-family="sans-serif"`) {
		t.Errorf("tick text should inherit font-family: %s", out)
	}
	// x axis ticks (after) point down: y2 = tickSize.
	if got := strings.Count(out, `y2="5"`); got != 3 {
		t.Errorf("tick y2=5 count = %d, want 3", got)
	}
}

func TestAxis_YOrientation(t *testing.T) {
	props := baseXAxisProps()
	props.Axis = "y"
	props.X = 10
	props.Y = 20
	out := render(t, axes.Axis(props, &theming.DefaultTheme))

	if !strings.Contains(out, `transform="translate(10,20)"`) {
		t.Errorf("missing origin translate: %s", out)
	}
	// y-axis domain line is vertical: x2=0, y2=length.
	if !strings.Contains(out, `x2="0"`) || !strings.Contains(out, `y2="100"`) {
		t.Errorf("y domain line endpoints wrong: %s", out)
	}
	// Ticks protrude along x: x2 - x1 = tickSize.
	if got := strings.Count(out, `x2="5"`); got != 3 {
		t.Errorf("y tick x2=5 count = %d, want 3", got)
	}
}

func TestAxis_TicksPositionBefore(t *testing.T) {
	props := baseXAxisProps()
	props.TicksPosition = "before"
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	// Ticks point up: y2 = -tickSize.
	if got := strings.Count(out, `y2="-5"`); got != 3 {
		t.Errorf("before tick y2=-5 count = %d, want 3", got)
	}
}

func TestAxis_NilTheme(t *testing.T) {
	out := render(t, axes.Axis(baseXAxisProps(), nil))
	// Without a theme no stroke/fill/font styling is emitted.
	if strings.Contains(out, "stroke=") {
		t.Errorf("nil theme should not emit stroke: %s", out)
	}
	if strings.Contains(out, "font-family=") || strings.Contains(out, "font-size=") || strings.Contains(out, "fill=") {
		t.Errorf("nil theme should not emit text styling: %s", out)
	}
	// Ticks + labels still render (ComputeCartesianTicks falls back to the
	// default theme for formatting).
	if got := strings.Count(out, "<text"); got != 3 {
		t.Errorf("text count = %d, want 3", got)
	}
}

func TestAxis_EmptyThemeOmitsLineStyling(t *testing.T) {
	// A non-nil theme with no axis styling: stroke/width branches are skipped,
	// text falls back to (empty) root text style.
	out := render(t, axes.Axis(baseXAxisProps(), &theming.Theme{}))
	if strings.Contains(out, "stroke=") || strings.Contains(out, "stroke-width=") {
		t.Errorf("empty theme should not emit stroke attrs: %s", out)
	}
}

func TestAxis_StyleOverrides(t *testing.T) {
	props := baseXAxisProps()
	props.LineColor = "#ff0000"
	props.LineWidth = 2
	props.TickLineColor = "#00ff00"
	props.TickLineWidth = 3
	props.TickTextColor = "#0000ff"
	props.TickTextFontSize = "13px" // string font-size branch
	props.TickTextFontFamily = "Arial"
	props.Legend = "value"
	props.LegendTextColor = "#123456"
	props.LegendTextFontSize = 14.0 // float64 font-size branch
	props.LegendTextFontFamily = "Georgia"
	out := render(t, axes.Axis(props, &theming.DefaultTheme))

	for _, want := range []string{
		`stroke="#ff0000"`, `stroke-width="2"`,
		`stroke="#00ff00"`, `stroke-width="3"`,
		`fill="#0000ff"`, `font-size="13px"`, `font-family="Arial"`,
		`fill="#123456"`, `font-size="14"`, `font-family="Georgia"`,
		">value</text>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output: %s", want, out)
		}
	}
}

func TestAxis_BandScaleCenteredTicks(t *testing.T) {
	props := axes.AxisProps{
		Axis:     "x",
		Scale:    bandScale([]string{"a", "b"}, 100),
		Length:   100,
		TickSize: 5,
	}
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	// Two bands of width 50 → centers at 25 and 75.
	if !strings.Contains(out, `translate(25,0)`) || !strings.Contains(out, `translate(75,0)`) {
		t.Errorf("band ticks not centered at 25/75: %s", out)
	}
	if !strings.Contains(out, ">a</text>") || !strings.Contains(out, ">b</text>") {
		t.Errorf("missing band labels: %s", out)
	}
}

func TestAxis_TickRotationTransform(t *testing.T) {
	tests := []struct {
		rotation     float64
		wantAnchor   string
		wantBaseline string
	}{
		{45, `text-anchor="end"`, `dominant-baseline="auto"`},
		{90, `text-anchor="middle"`, `dominant-baseline="hanging"`},
		{135, `text-anchor="start"`, `dominant-baseline="auto"`},
		{270, `text-anchor="middle"`, `dominant-baseline="auto"`},
		{300, `text-anchor="start"`, `dominant-baseline="auto"`},
		{-45, `text-anchor="start"`, `dominant-baseline="auto"`}, // normalized to 315
	}
	for _, tc := range tests {
		props := baseXAxisProps()
		props.TickRotation = tc.rotation
		out := render(t, axes.Axis(props, &theming.DefaultTheme))
		if !strings.Contains(out, "rotate(") {
			t.Errorf("rotation %v: missing rotate transform", tc.rotation)
		}
		if !strings.Contains(out, tc.wantAnchor) {
			t.Errorf("rotation %v: missing %s in %s", tc.rotation, tc.wantAnchor, out)
		}
		if !strings.Contains(out, tc.wantBaseline) {
			t.Errorf("rotation %v: missing %s in %s", tc.rotation, tc.wantBaseline, out)
		}
	}
}

func TestAxis_DisableTickLabel(t *testing.T) {
	props := baseXAxisProps()
	props.DisableTickLabel = true
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	if strings.Contains(out, "<text") {
		t.Errorf("DisableTickLabel should drop labels: %s", out)
	}
	// Tick lines still drawn.
	if got := strings.Count(out, "<line"); got != 4 { // domain + 3 ticks
		t.Errorf("line count = %d, want 4", got)
	}
}

func TestAxis_LegendPositions(t *testing.T) {
	tests := []struct {
		name   string
		axis   string
		pos    axes.AxisLegendPosition
		offset float64
		want   []string
	}{
		{"x-end-default-offset", "x", "", 0, []string{`translate(100,32) rotate(0)`, `text-anchor="end"`}},
		{"x-start", "x", axes.AxisLegendStart, 20, []string{`translate(0,20) rotate(0)`, `text-anchor="start"`}},
		{"x-middle", "x", axes.AxisLegendMiddle, 20, []string{`translate(50,20) rotate(0)`, `text-anchor="middle"`}},
		{"y-end", "y", axes.AxisLegendEnd, 40, []string{`translate(-40,100) rotate(-90)`, `text-anchor="end"`}},
		{"y-start", "y", axes.AxisLegendStart, 40, []string{`translate(-40,0) rotate(-90)`}},
		{"y-middle", "y", axes.AxisLegendMiddle, 40, []string{`translate(-40,50) rotate(-90)`}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			props := baseXAxisProps()
			props.Axis = tc.axis
			props.Legend = "Legend"
			props.LegendPosition = tc.pos
			props.LegendOffset = tc.offset
			out := render(t, axes.Axis(props, &theming.DefaultTheme))
			if !strings.Contains(out, ">Legend</text>") {
				t.Fatalf("missing legend text: %s", out)
			}
			// Default theme legend font-size is 12.
			if !strings.Contains(out, `font-size="12"`) {
				t.Errorf("legend should use theme font-size 12: %s", out)
			}
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Errorf("missing %q in %s", want, out)
				}
			}
		})
	}
}

func TestAxisLegend_NilTheme(t *testing.T) {
	props := baseXAxisProps()
	props.Legend = "L"
	out := render(t, axes.AxisLegend(props, nil))
	if !strings.Contains(out, ">L</text>") {
		t.Errorf("legend text missing: %s", out)
	}
	if strings.Contains(out, "fill=") || strings.Contains(out, "font-size=") {
		t.Errorf("nil theme should not style legend: %s", out)
	}
}

func TestAxes_BothAxes(t *testing.T) {
	x := baseXAxisProps()
	y := baseXAxisProps()
	y.Axis = "y"
	y.Legend = "Y legend"
	out := render(t, axes.Axes(axes.AxesProps{
		XAxis: &x,
		YAxis: &y,
		Theme: &theming.DefaultTheme,
	}))
	// One horizontal domain line (x axis) and one vertical (y axis).
	if !strings.Contains(out, `x2="100" y2="0"`) {
		t.Errorf("missing x-axis domain line: %s", out)
	}
	if !strings.Contains(out, `x2="0" y2="100"`) {
		t.Errorf("missing y-axis domain line: %s", out)
	}
	if !strings.Contains(out, ">Y legend</text>") {
		t.Errorf("missing y legend: %s", out)
	}
}

func TestAxes_HiddenAndNil(t *testing.T) {
	x := baseXAxisProps()
	x.Hidden = true
	out := render(t, axes.Axes(axes.AxesProps{XAxis: &x, YAxis: nil, Theme: &theming.DefaultTheme}))
	if out != "" {
		t.Errorf("hidden x + nil y should render nothing, got %q", out)
	}
	out = render(t, axes.Axes(axes.AxesProps{}))
	if out != "" {
		t.Errorf("empty AxesProps should render nothing, got %q", out)
	}
}

func TestAxisTick_DirectRender(t *testing.T) {
	props := baseXAxisProps()
	computed := axes.ComputeCartesianTicks(props, &theming.DefaultTheme)
	out := render(t, axes.AxisTick(computed, props, &theming.DefaultTheme))
	if got := strings.Count(out, "<line"); got != 3 {
		t.Errorf("tick line count = %d, want 3", got)
	}
	if got := strings.Count(out, "<text"); got != 3 {
		t.Errorf("tick text count = %d, want 3", got)
	}
}

func TestAxis_TimeScaleLabels(t *testing.T) {
	jan1, _ := time.Parse("2006-01-02", "2024-01-01")
	feb1, _ := time.Parse("2006-01-02", "2024-02-01")
	mar15, _ := time.Parse("2006-01-02", "2024-03-15")
	props := axes.AxisProps{
		Axis:       "x",
		Scale:      timeScale(t),
		Length:     100,
		TickSize:   5,
		TickValues: []any{jan1, feb1, mar15, "not-a-time"},
	}
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	for _, want := range []string{">2024</text>", ">Feb</text>", ">Mar 15</text>", ">not-a-time</text>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing time label %q in %s", want, out)
		}
	}
}

func TestAxis_NumberFormatterVariants(t *testing.T) {
	d, _ := time.Parse("2006-01-02", "2024-06-01")
	props := baseXAxisProps()
	props.TickValues = []any{2.5, 3, "x", d}
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	for _, want := range []string{">2.5</text>", ">3</text>", ">x</text>", ">2024-06-01</text>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing formatted label %q in %s", want, out)
		}
	}
}

func TestAxis_NonFiniteOriginNormalized(t *testing.T) {
	props := baseXAxisProps()
	props.X = math.NaN()
	props.Y = math.Inf(1)
	out := render(t, axes.Axis(props, &theming.DefaultTheme))
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("NaN/Inf origin should normalize to 0: %s", out)
	}
}

func TestAxis_CustomThemeExtraCoercion(t *testing.T) {
	// Non-string stroke and int strokeWidth exercise the Extra coercion paths.
	theme := theming.Theme{}
	theme.Axis.Domain.Line.Extra = map[string]any{"stroke": 123, "strokeWidth": int(2)}
	theme.Axis.Ticks.Line.Extra = map[string]any{"other": "x"} // missing keys
	out := render(t, axes.Axis(baseXAxisProps(), &theme))
	if !strings.Contains(out, `stroke="123"`) {
		t.Errorf("non-string stroke should be coerced: %s", out)
	}
	if !strings.Contains(out, `stroke-width="2"`) {
		t.Errorf("int strokeWidth should be coerced: %s", out)
	}
}

func TestGrid_XAxisDefaultTheme(t *testing.T) {
	out := render(t, axes.Grid(axes.GridProps{
		Axis:       "x",
		Scale:      linearScale(0, 100, 100),
		Width:      100,
		Height:     50,
		X:          5,
		Y:          10,
		TickValues: []any{0.0, 50.0, 100.0},
		Theme:      &theming.DefaultTheme,
	}))
	if !strings.Contains(out, `transform="translate(5,10)"`) {
		t.Errorf("missing grid translate: %s", out)
	}
	if got := strings.Count(out, "<line"); got != 3 {
		t.Errorf("grid line count = %d, want 3", got)
	}
	// x grid lines are vertical spanning the height.
	if got := strings.Count(out, `y2="50"`); got != 3 {
		t.Errorf("vertical grid line count = %d, want 3", got)
	}
	if got := strings.Count(out, `stroke="#dddddd"`); got != 3 {
		t.Errorf("theme grid stroke count = %d, want 3", got)
	}
	if got := strings.Count(out, `fill="none"`); got != 3 {
		t.Errorf("fill=none count = %d, want 3", got)
	}
}

func TestGrid_YAxisWithOverrides(t *testing.T) {
	out := render(t, axes.Grid(axes.GridProps{
		Axis:       "y",
		Scale:      linearScale(0, 100, 50),
		Width:      100,
		Height:     50,
		TickValues: []any{0.0, 100.0},
		LineColor:  "#abcdef",
		LineWidth:  2,
		Theme:      &theming.DefaultTheme,
	}))
	// y grid lines are horizontal spanning the width.
	if got := strings.Count(out, `x2="100"`); got != 2 {
		t.Errorf("horizontal grid line count = %d, want 2", got)
	}
	if got := strings.Count(out, `stroke="#abcdef"`); got != 2 {
		t.Errorf("override stroke count = %d, want 2", got)
	}
	if got := strings.Count(out, `stroke-width="2"`); got != 2 {
		t.Errorf("override stroke-width count = %d, want 2", got)
	}
}

func TestGrid_NilThemeNoStyling(t *testing.T) {
	out := render(t, axes.Grid(axes.GridProps{
		Axis:       "x",
		Scale:      linearScale(0, 100, 100),
		Width:      100,
		Height:     50,
		TickValues: []any{50.0},
	}))
	if strings.Contains(out, "stroke=") || strings.Contains(out, "stroke-width=") {
		t.Errorf("nil theme grid should not emit stroke attrs: %s", out)
	}
	if got := strings.Count(out, "<line"); got != 1 {
		t.Errorf("grid line count = %d, want 1", got)
	}
}

func TestGrid_BandScaleCentered(t *testing.T) {
	out := render(t, axes.Grid(axes.GridProps{
		Axis:   "x",
		Scale:  bandScale([]string{"a", "b"}, 100),
		Width:  100,
		Height: 40,
		Theme:  &theming.DefaultTheme,
	}))
	// Band centers at 25/75.
	if !strings.Contains(out, `x1="25"`) || !strings.Contains(out, `x1="75"`) {
		t.Errorf("band grid lines should be centered: %s", out)
	}
}

func TestGrid_DefaultTicksFromScale(t *testing.T) {
	// No TickValues → scale default ticks are used.
	out := render(t, axes.Grid(axes.GridProps{
		Axis:   "y",
		Scale:  linearScale(0, 100, 50),
		Width:  100,
		Height: 50,
		Theme:  &theming.DefaultTheme,
	}))
	if strings.Count(out, "<line") == 0 {
		t.Errorf("expected default scale ticks to produce grid lines: %s", out)
	}
}
