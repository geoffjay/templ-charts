package core_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/core"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}
	return b.String()
}

func TestSvgWrapper_AllOptions(t *testing.T) {
	out := renderWrapper(t, core.SvgWrapperProps{
		Width: 200, Height: 100,
		Margin:          core.Margin{Top: 10, Left: 20},
		Defs:            []core.Def{core.PatternDotsDef("dots", nil)},
		Background:      "#eee",
		Role:            "img",
		AriaLabel:       "chart",
		AriaLabelledBy:  "chart-title",
		AriaDescribedBy: "chart-desc",
		Responsive:      true,
	})
	for _, want := range []string{
		`width="200"`, `height="100"`, `viewBox="0 0 200 100"`,
		`aria-label="chart"`, `aria-labelledby="chart-title"`, `aria-describedby="chart-desc"`,
		`style="width:100%;height:auto;display:block"`,
		`<rect width="200" height="100" fill="#eee">`,
		`<defs>`, `<pattern id="dots"`,
		`transform="translate(20,10)"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("SvgWrapper missing %q in %q", want, out)
		}
	}
}

func TestSvgWrapper_Minimal(t *testing.T) {
	out := renderWrapper(t, core.SvgWrapperProps{Width: 10, Height: 10})
	for _, absent := range []string{"role=", "aria-", "tabindex", "style=", "<rect", "<defs>"} {
		if strings.Contains(out, absent) {
			t.Errorf("minimal SvgWrapper should not contain %q; got %q", absent, out)
		}
	}
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("expected zero-margin translate; got %q", out)
	}
}

func TestDefs_Empty(t *testing.T) {
	if out := render(t, core.Defs(nil)); out != "" {
		t.Errorf("empty defs should render nothing, got %q", out)
	}
}

func TestDefs_UnknownTypeSkipped(t *testing.T) {
	out := render(t, core.Defs([]core.Def{{ID: "x", Type: "bogus"}}))
	if !strings.Contains(out, "<defs>") || strings.Contains(out, "<pattern") || strings.Contains(out, "<linearGradient") {
		t.Errorf("unknown def type should render an empty defs block, got %q", out)
	}
}

func TestDefs_LinearGradient(t *testing.T) {
	out := render(t, core.Defs([]core.Def{core.LinearGradientDef("g1", []core.GradientStop{
		{Offset: 0, Color: "#f00"},
		{Offset: 100, Color: "#00f", Opacity: 0.5},
	}, nil)}))
	for _, want := range []string{
		`<linearGradient id="g1" x1="0" x2="0" y1="0" y2="1">`,
		`offset="0%"`, `offset="100%"`,
		`stop-color="#f00"`, `stop-color="#00f"`,
		// Zero opacity defaults to 1.
		`stop-opacity="1"`, `stop-opacity="0.5"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("gradient missing %q in %q", want, out)
		}
	}
	if got := strings.Count(out, "<stop "); got != 2 {
		t.Errorf("stop count = %d, want 2", got)
	}
}

func TestDefs_PatternDots(t *testing.T) {
	// Defaults: size 4, padding 4 → fullSize 8, radius 2, first dot at (4,4).
	out := render(t, core.Defs([]core.Def{core.PatternDotsDef("d1", nil)}))
	for _, want := range []string{
		`<pattern id="d1" width="8" height="8" patternUnits="userSpaceOnUse">`,
		`<rect width="8" height="8" fill="#ffffff">`,
		`<circle cx="4" cy="4" r="2" fill="#000000">`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("patternDots missing %q in %q", want, out)
		}
	}
	if got := strings.Count(out, "<circle"); got != 1 {
		t.Errorf("circle count = %d, want 1", got)
	}
}

func TestDefs_PatternDotsStagger(t *testing.T) {
	out := render(t, core.Defs([]core.Def{{
		ID: "d2", Type: core.DefTypePatternDots,
		Size: 4, Padding: 4, Stagger: true, Color: "#111", Background: "#222",
	}}))
	// Stagger doubles the tile: fullSize 16, second dot at (10,10).
	for _, want := range []string{
		`width="16" height="16"`,
		`fill="#222"`, `fill="#111"`,
		`<circle cx="10" cy="10" r="2"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("staggered dots missing %q in %q", want, out)
		}
	}
	if got := strings.Count(out, "<circle"); got != 2 {
		t.Errorf("circle count = %d, want 2", got)
	}
}

func TestDefs_PatternSquares(t *testing.T) {
	out := render(t, core.Defs([]core.Def{core.PatternSquaresDef("s1", nil)}))
	for _, want := range []string{
		`<pattern id="s1" width="8" height="8"`,
		`<rect x="2" y="2" width="4" height="4" fill="#000000">`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("patternSquares missing %q in %q", want, out)
		}
	}
	if got := strings.Count(out, "<rect"); got != 2 {
		t.Errorf("rect count = %d, want 2 (background + square)", got)
	}

	stag := render(t, core.Defs([]core.Def{{ID: "s2", Type: core.DefTypePatternSquares, Stagger: true}}))
	if !strings.Contains(stag, `width="16" height="16"`) || strings.Count(stag, "<rect") != 3 {
		t.Errorf("staggered squares wrong tile/rect count: %q", stag)
	}
	if !strings.Contains(stag, `<rect x="10" y="10" width="4" height="4"`) {
		t.Errorf("staggered square offset missing in %q", stag)
	}
}

func TestDefs_PatternLinesRotations(t *testing.T) {
	cases := []struct {
		name     string
		rotation float64
		wants    []string
	}{
		// Defaults: spacing 5, lineWidth 2.
		{"rotation 0", 0, []string{`width="5" height="5"`, `d="M 0 0 L 5 0 M 0 5 L 5 5"`}},
		{"rotation 90", 90, []string{`d="M 0 0 L 0 5 M 5 0 L 5 5"`}},
		{"rotation 45", 45, []string{`width="7.071" height="7.071"`, `d="M 0 -7.071`}},
		{"rotation -45", -45, []string{`d="M -7.071 7.071`}},
		// Normalization: 100→-80, 200→-160, -100→80, -200→160.
		{"rotation 100", 100, []string{`<pattern id="l"`}},
		{"rotation 200", 200, []string{`<pattern id="l"`}},
		{"rotation -100", -100, []string{`<pattern id="l"`}},
		{"rotation -200", -200, []string{`<pattern id="l"`}},
	}
	for _, c := range cases {
		out := render(t, core.Defs([]core.Def{{ID: "l", Type: core.DefTypePatternLines, Rotation: c.rotation}}))
		for _, want := range append(c.wants, `stroke="#ffffff"`, `strokeWidth="2"`, `fill="#000000"`) {
			if !strings.Contains(out, want) {
				t.Errorf("%s: missing %q in %q", c.name, want, out)
			}
		}
	}
}

func TestDotsItem_Basic(t *testing.T) {
	out := render(t, core.DotsItem(core.DotsItemProps{X: 10, Y: 20, Size: 8, Color: "#f00"}))
	for _, want := range []string{
		`transform="translate(10,20)"`,
		`style="pointer-events: none"`,
		`r="4"`, `fill="#f00"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("DotsItem missing %q in %q", want, out)
		}
	}
	for _, absent := range []string{"stroke", "<text", "data-tc-tooltip", "<animate"} {
		if strings.Contains(out, absent) {
			t.Errorf("basic DotsItem should not contain %q; got %q", absent, out)
		}
	}
}

func TestDotsItem_FullOptions(t *testing.T) {
	out := render(t, core.DotsItem(core.DotsItemProps{
		X: 1, Y: 2, Size: 6, Color: "#0f0",
		BorderWidth: 1.5, BorderColor: "#333",
		Label: "12", LabelTextAnchor: "start", LabelYOffset: -8,
		LabelFill: "#444", LabelFontSize: 11, LabelFontFamily: "sans-serif",
		Tooltip: "point tip",
		Animate: true, AnimateBegin: "0.2s",
	}))
	for _, want := range []string{
		`data-tc-tooltip="point tip"`,
		`style="pointer-events: auto; cursor: pointer"`,
		`stroke="#333"`, `stroke-width="1.5"`,
		`text-anchor="start"`, `y="-8"`,
		`fill="#444"`, `font-size="11"`, `font-family="sans-serif"`,
		`>12</text>`,
		`<animate attributeName="r" from="0" to="3" begin="0.2s"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("DotsItem missing %q in %q", want, out)
		}
	}
}

func TestDotsItem_DefaultLabelAnchor(t *testing.T) {
	out := render(t, core.DotsItem(core.DotsItemProps{Size: 4, Label: "x"}))
	if !strings.Contains(out, `text-anchor="middle"`) {
		t.Errorf("default label anchor should be middle; got %q", out)
	}
}

func TestCartesianMarkers_Empty(t *testing.T) {
	out := render(t, core.CartesianMarkers(core.CartesianMarkersProps{Width: 100, Height: 50}))
	// An empty marker list must render nothing at all — not even the stray
	// literal "return " that a bare `return` inside a templ block used to emit.
	if out != "" {
		t.Errorf("no markers should render empty output, got %q", out)
	}
}

func TestCartesianMarkers_Axes(t *testing.T) {
	props := core.CartesianMarkersProps{
		Width: 100, Height: 50,
		XScale: func(v any) float64 { return 30 },
		YScale: func(v any) float64 { return 40 },
		Markers: []core.CartesianMarker{
			{Axis: "x", Value: 3, LineColor: "#f00", LineStrokeWidth: 2},
			{Axis: "y", Value: 4},
		},
	}
	out := render(t, core.CartesianMarkers(props))
	// x marker: vertical line at xScale(v), spanning the height.
	if !strings.Contains(out, `transform="translate(30,0)"`) || !strings.Contains(out, `y2="50"`) {
		t.Errorf("x marker line wrong: %q", out)
	}
	// y marker: horizontal line at yScale(v), spanning the width.
	if !strings.Contains(out, `transform="translate(0,40)"`) || !strings.Contains(out, `x2="100"`) {
		t.Errorf("y marker line wrong: %q", out)
	}
	if !strings.Contains(out, `stroke="#f00"`) || !strings.Contains(out, `stroke-width="2"`) {
		t.Errorf("marker line style missing: %q", out)
	}
	if strings.Contains(out, "<text") {
		t.Errorf("markers without legends should not emit text: %q", out)
	}
	if got := strings.Count(out, "<line"); got != 2 {
		t.Errorf("line count = %d, want 2", got)
	}
}

func TestCartesianMarkers_Legend(t *testing.T) {
	base := core.CartesianMarkersProps{
		Width: 100, Height: 80,
		XScale: func(v any) float64 { return 10 },
		YScale: func(v any) float64 { return 20 },
	}
	cases := []struct {
		name   string
		marker core.CartesianMarker
		wants  []string
	}{
		{
			// Defaults: top-right, horizontal, offsets 14.
			"y axis defaults",
			core.CartesianMarker{Axis: "y", Value: 1, Legend: "limit"},
			[]string{`>limit</text>`, `text-anchor="end"`, `translate(86,-14) rotate(0)`},
		},
		{
			"x axis vertical top",
			core.CartesianMarker{Axis: "x", Value: 1, Legend: "v", LegendPosition: "top", LegendOrientation: "vertical"},
			[]string{`rotate(-90)`, `text-anchor="start"`},
		},
		{
			"x axis left horizontal",
			core.CartesianMarker{Axis: "x", Value: 1, Legend: "h", LegendPosition: "left", LegendOffsetX: 5, LegendOffsetY: 5},
			[]string{`text-anchor="end"`, `translate(-5,40) rotate(0)`},
		},
		{
			"y axis bottom horizontal",
			core.CartesianMarker{Axis: "y", Value: 1, Legend: "b", LegendPosition: "bottom", LegendOffsetX: 5, LegendOffsetY: 5},
			[]string{`text-anchor="middle"`, `translate(50,5) rotate(0)`},
		},
	}
	for _, c := range cases {
		props := base
		props.Markers = []core.CartesianMarker{c.marker}
		out := render(t, core.CartesianMarkers(props))
		for _, want := range c.wants {
			if !strings.Contains(out, want) {
				t.Errorf("%s: missing %q in %q", c.name, want, out)
			}
		}
	}
}
