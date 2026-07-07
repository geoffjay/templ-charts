package polaraxes_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	polaraxes "github.com/geoffjay/templ-charts/charts/polar-axes"
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

func angleScale() scales.Scale {
	return scales.NewLinearScaleWithRange(0, 360, 0, 360)
}

func radiusScale() scales.Scale {
	return scales.NewLinearScaleWithRange(0, 100, 0, 100)
}

// The templ components emit the Render* SVG fragments as escaped text; the
// numeric/transform content survives escaping, so assertions target those.

func TestCircularAxisComponent(t *testing.T) {
	out := render(t, polaraxes.CircularAxis(polaraxes.CircularAxisProps{
		Type:       polaraxes.CircularAxisOuter,
		Center:     [2]float64{150, 150},
		Radius:     80,
		StartAngle: 0,
		EndAngle:   360,
		Scale:      angleScale(),
		Theme:      &theming.DefaultTheme,
	}))
	if !strings.Contains(out, "translate(150,150)") {
		t.Errorf("missing center translate: %s", out)
	}
	if !strings.Contains(out, "text") {
		t.Errorf("missing tick labels: %s", out)
	}
}

func TestRadialAxisComponent(t *testing.T) {
	out := render(t, polaraxes.RadialAxis(polaraxes.RadialAxisProps{
		Center:        [2]float64{100, 100},
		Angle:         0,
		Scale:         radiusScale(),
		TicksPosition: polaraxes.TicksAfter,
		Theme:         &theming.DefaultTheme,
	}))
	if !strings.Contains(out, "translate(100,100)") {
		t.Errorf("missing center translate: %s", out)
	}
	if !strings.Contains(out, "rotate(-90)") {
		t.Errorf("missing rotation to angle 0: %s", out)
	}
}

func TestPolarGridComponent(t *testing.T) {
	out := render(t, polaraxes.PolarGrid(polaraxes.PolarGridProps{
		Center:             [2]float64{200, 200},
		EnableRadialGrid:   true,
		AngleScale:         angleScale(),
		EnableCircularGrid: true,
		RadiusScale:        radiusScale(),
		StartAngle:         0,
		EndAngle:           360,
		OuterRadius:        100,
		Theme:              &theming.DefaultTheme,
	}))
	if !strings.Contains(out, "translate(200,200)") {
		t.Errorf("missing center translate: %s", out)
	}
}

func TestRadialGridComponent(t *testing.T) {
	out := render(t, polaraxes.RadialGrid(polaraxes.RadialGridProps{
		Scale:       angleScale(),
		InnerRadius: 10,
		OuterRadius: 90,
		Theme:       &theming.DefaultTheme,
	}))
	if !strings.Contains(out, "rotate(") {
		t.Errorf("missing rotated grid rays: %s", out)
	}
}

func TestCircularGridComponent(t *testing.T) {
	out := render(t, polaraxes.CircularGrid(polaraxes.CircularGridProps{
		Scale:      radiusScale(),
		StartAngle: 0,
		EndAngle:   270,
		Theme:      &theming.DefaultTheme,
	}))
	if !strings.Contains(out, "path") {
		t.Errorf("missing arc paths: %s", out)
	}
}

func TestRenderRadialAxis_BranchesByPositionAndAngle(t *testing.T) {
	tests := []struct {
		name       string
		pos        polaraxes.TicksPosition
		angle      float64
		wantAnchor string
		wantRotate string
	}{
		{"before-angle-45", polaraxes.TicksBefore, 45, `text-anchor="end"`, `rotate(90)`},
		{"before-angle-180", polaraxes.TicksBefore, 180, `text-anchor="start"`, `rotate(-90)`},
		{"before-angle-300", polaraxes.TicksBefore, 300, `text-anchor="end"`, `rotate(90)`},
		{"after-angle-45", polaraxes.TicksAfter, 45, `text-anchor="start"`, `rotate(90)`},
		{"after-angle-180", polaraxes.TicksAfter, 180, `text-anchor="end"`, `rotate(-90)`},
		{"after-angle-300", polaraxes.TicksAfter, 300, `text-anchor="start"`, `rotate(90)`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svg := polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
				Center:        [2]float64{0, 0},
				Angle:         tc.angle,
				Scale:         radiusScale(),
				TicksPosition: tc.pos,
				Theme:         &theming.DefaultTheme,
			})
			if !strings.Contains(svg, tc.wantAnchor) {
				t.Errorf("missing %s: %s", tc.wantAnchor, svg)
			}
			if !strings.Contains(svg, tc.wantRotate) {
				t.Errorf("missing tick %s: %s", tc.wantRotate, svg)
			}
		})
	}
}

func TestRenderRadialAxis_TickSizeAndPadding(t *testing.T) {
	svg := polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
		Angle:         0,
		Scale:         radiusScale(),
		TicksPosition: polaraxes.TicksAfter,
		Theme:         &theming.DefaultTheme,
		RadialAxisConfig: polaraxes.RadialAxisConfig{
			TickSize:    10,
			TickPadding: 4,
		},
	})
	// angle 0, after → lineX = tickSize, textX = tickSize + padding.
	if !strings.Contains(svg, `<line x2="10"`) {
		t.Errorf("tick line should use tickSize 10: %s", svg)
	}
	if !strings.Contains(svg, `dx="14"`) {
		t.Errorf("label offset should be 14: %s", svg)
	}
}

func TestRenderRadialAxis_CustomTickComponent(t *testing.T) {
	var custom polaraxes.RadialTickRenderer = func(props polaraxes.RadialAxisTickProps) string {
		return "<!--radial-tick:" + props.Label + "-->"
	}
	svg := polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
		Angle:         0,
		Scale:         radiusScale(),
		TicksPosition: polaraxes.TicksAfter,
		Theme:         &theming.DefaultTheme,
		RadialAxisConfig: polaraxes.RadialAxisConfig{
			TickComponent: &custom,
			Ticks:         scales.TicksSpec{Count: 5, HasCount: true},
		},
	})
	if !strings.Contains(svg, "<!--radial-tick:") {
		t.Errorf("custom tick renderer not used: %s", svg)
	}
	if strings.Contains(svg, "<line") {
		t.Errorf("default tick renderer should be bypassed: %s", svg)
	}
}

func TestRenderRadialAxis_BandScaleCentersTicks(t *testing.T) {
	scale := scales.NewBandScaleWithRange([]string{"a", "b"}, 0, 100, 0, false)
	svg := polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
		Angle:         0,
		Scale:         scale,
		TicksPosition: polaraxes.TicksAfter,
		Theme:         &theming.DefaultTheme,
	})
	// Bands of width 50 → tick groups at 25 and 75.
	if !strings.Contains(svg, `translate(25,0)`) || !strings.Contains(svg, `translate(75,0)`) {
		t.Errorf("band ticks not centered: %s", svg)
	}
	if !strings.Contains(svg, ">a</text>") || !strings.Contains(svg, ">b</text>") {
		t.Errorf("missing band labels: %s", svg)
	}
}

func TestRenderRadialAxis_Animate(t *testing.T) {
	svg := polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
		Angle:         0,
		Scale:         radiusScale(),
		TicksPosition: polaraxes.TicksAfter,
		Theme:         &theming.DefaultTheme,
		Animate:       true,
	})
	if !strings.Contains(svg, `opacity="0"`) || !strings.Contains(svg, "<animate") {
		t.Errorf("animate should emit opacity + <animate>: %s", svg)
	}
}

func TestRenderCircularAxis_InnerType(t *testing.T) {
	props := polaraxes.CircularAxisProps{
		Center:     [2]float64{0, 0},
		Radius:     80,
		StartAngle: 0,
		EndAngle:   360,
		Scale:      angleScale(),
		Theme:      &theming.DefaultTheme,
	}
	props.Type = polaraxes.CircularAxisInner
	inner := polaraxes.RenderCircularAxis(props)
	props.Type = polaraxes.CircularAxisOuter
	outer := polaraxes.RenderCircularAxis(props)

	// Inner: text at radius - tickSize - padding = 80-5-12 = 63; the tick for
	// value 0 sits at the top (angle -90) → dy = -63. Outer: 80+5+12 = 97.
	if !strings.Contains(inner, `dy="-63"`) {
		t.Errorf("inner axis label radius wrong: %s", inner)
	}
	if !strings.Contains(outer, `dy="-97"`) {
		t.Errorf("outer axis label radius wrong: %s", outer)
	}
}

func TestRenderCircularAxis_CustomTickComponentAndConfig(t *testing.T) {
	var custom polaraxes.CircularTickRenderer = func(props polaraxes.CircularAxisTickProps) string {
		return "<!--circ-tick:" + props.Label + "-->"
	}
	svg := polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
		Type:       polaraxes.CircularAxisOuter,
		Radius:     50,
		StartAngle: 0,
		EndAngle:   270,
		Scale:      angleScale(),
		Theme:      &theming.DefaultTheme,
		CircularAxisConfig: polaraxes.CircularAxisConfig{
			TickSize:      8,
			TickPadding:   6,
			TickComponent: &custom,
		},
	})
	if !strings.Contains(svg, "<!--circ-tick:") {
		t.Errorf("custom tick renderer not used: %s", svg)
	}
	if strings.Contains(svg, "<text") {
		t.Errorf("default tick renderer should be bypassed: %s", svg)
	}
	// Domain arc still rendered by the axis itself.
	if !strings.Contains(svg, "<path") {
		t.Errorf("missing domain arc: %s", svg)
	}
}

func TestRenderCircularAxis_AnimateAndFormat(t *testing.T) {
	svg := polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
		Type:       polaraxes.CircularAxisOuter,
		Radius:     50,
		StartAngle: 0,
		EndAngle:   360,
		Scale:      angleScale(),
		Theme:      &theming.DefaultTheme,
		Animate:    true,
		CircularAxisConfig: polaraxes.CircularAxisConfig{
			Format: func(v any) string { return "F!" },
		},
	})
	if !strings.Contains(svg, "<animate") {
		t.Errorf("animate missing: %s", svg)
	}
	if !strings.Contains(svg, ">F!</text>") {
		t.Errorf("custom formatter not applied: %s", svg)
	}
}

func TestRenderCircularAxisTick_EmptyThemeAndEscaping(t *testing.T) {
	svg := polaraxes.RenderCircularAxisTick(polaraxes.CircularAxisTickProps{
		Label: "a<b&c>d",
		X1:    1, Y1: 2, X2: 3, Y2: 4,
		TextX: 5, TextY: 6,
	})
	if !strings.Contains(svg, ">a&lt;b&amp;c&gt;d</text>") {
		t.Errorf("label not escaped: %s", svg)
	}
	// Empty theme → no stroke / fill / font attributes.
	for _, attr := range []string{"stroke=", "fill=", "font-size=", "font-family=", "opacity="} {
		if strings.Contains(svg, attr) {
			t.Errorf("empty theme should not emit %s: %s", attr, svg)
		}
	}
}

func TestRenderRadialAxisTick_EmptyTheme(t *testing.T) {
	svg := polaraxes.RenderRadialAxisTick(polaraxes.RadialAxisTickProps{
		Label:      "v",
		TextAnchor: "start",
		Y:          10,
		Length:     5,
		TextX:      8,
		Rotation:   90,
	})
	if !strings.Contains(svg, `transform="translate(10,0) rotate(90)"`) {
		t.Errorf("tick transform wrong: %s", svg)
	}
	for _, attr := range []string{"stroke=", "fill=", "font-size=", "font-family="} {
		if strings.Contains(svg, attr) {
			t.Errorf("empty theme should not emit %s: %s", attr, svg)
		}
	}
}

func TestRenderPolarGrid_Toggles(t *testing.T) {
	base := polaraxes.PolarGridProps{
		Center:      [2]float64{50, 50},
		AngleScale:  angleScale(),
		RadiusScale: radiusScale(),
		StartAngle:  0,
		EndAngle:    360,
		OuterRadius: 100,
		Theme:       &theming.DefaultTheme,
	}

	neither := polaraxes.RenderPolarGrid(base)
	if strings.Contains(neither, "<line") || strings.Contains(neither, "<path") {
		t.Errorf("disabled grids should render no lines/paths: %s", neither)
	}
	if !strings.Contains(neither, `translate(50,50)`) {
		t.Errorf("missing container translate: %s", neither)
	}

	radialOnly := base
	radialOnly.EnableRadialGrid = true
	svg := polaraxes.RenderPolarGrid(radialOnly)
	if !strings.Contains(svg, "<line") || strings.Contains(svg, "<path") {
		t.Errorf("radial-only grid wrong: %s", svg)
	}

	circularOnly := base
	circularOnly.EnableCircularGrid = true
	svg = polaraxes.RenderPolarGrid(circularOnly)
	if strings.Contains(svg, "<line") || !strings.Contains(svg, "<path") {
		t.Errorf("circular-only grid wrong: %s", svg)
	}
}

func TestRenderRadialGrid_AnimateAndEmptyTheme(t *testing.T) {
	props := polaraxes.RadialGridProps{
		Scale:       angleScale(),
		Ticks:       scales.TicksSpec{Count: 4, HasCount: true},
		InnerRadius: 10,
		OuterRadius: 90,
		Theme:       &theming.DefaultTheme,
		Animate:     true,
	}
	svg := polaraxes.RenderRadialGrid(props)
	if !strings.Contains(svg, `opacity="0"`) || !strings.Contains(svg, "<animate") {
		t.Errorf("animate branches missing: %s", svg)
	}
	if !strings.Contains(svg, `x1="10"`) || !strings.Contains(svg, `x2="90"`) {
		t.Errorf("ray endpoints wrong: %s", svg)
	}
	if !strings.Contains(svg, `stroke="#dddddd"`) {
		t.Errorf("theme grid stroke missing: %s", svg)
	}

	props.Theme = &theming.Theme{}
	props.Animate = false
	svg = polaraxes.RenderRadialGrid(props)
	if strings.Contains(svg, "stroke=") || strings.Contains(svg, "<animate") {
		t.Errorf("empty theme should not style rays: %s", svg)
	}
}

func TestRenderCircularGrid_BandScaleAndAnimate(t *testing.T) {
	scale := scales.NewBandScaleWithRange([]string{"a", "b"}, 0, 100, 0, false)
	svg := polaraxes.RenderCircularGrid(polaraxes.CircularGridProps{
		Scale:      scale,
		StartAngle: 0,
		EndAngle:   360,
		Theme:      &theming.DefaultTheme,
		Animate:    true,
	})
	// Band radii centered: 25 and 75 → arc radius appears in path.
	if !strings.Contains(svg, "A25,25") || !strings.Contains(svg, "A75,75") {
		t.Errorf("band radii not centered: %s", svg)
	}
	if got := strings.Count(svg, "<animate"); got != 2 {
		t.Errorf("animate count = %d, want 2", got)
	}
	if got := strings.Count(svg, `opacity="0"`); got != 2 {
		t.Errorf("opacity count = %d, want 2", got)
	}
}

func styledTickTheme() theming.AxisTheme {
	return theming.AxisTheme{
		Ticks: theming.AxisTicks{
			Line: theming.AxisTickLine{Extra: map[string]any{"stroke": "#444444", "strokeWidth": float64(2)}},
			Text: theming.TextStyle{Fill: "#eeeeee", FontSize: 9, FontFamily: "TickFont"},
		},
	}
}

func TestRenderCircularAxisTick_StyledText(t *testing.T) {
	svg := polaraxes.RenderCircularAxisTick(polaraxes.CircularAxisTickProps{
		Label: "42",
		Theme: styledTickTheme(),
	})
	for _, want := range []string{
		`stroke="#444444"`, `stroke-width="2"`,
		`fill="#eeeeee"`, `font-size="9"`, `font-family="TickFont"`,
		">42</text>",
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("missing %q: %s", want, svg)
		}
	}
}

func TestRenderRadialAxisTick_StyledText(t *testing.T) {
	svg := polaraxes.RenderRadialAxisTick(polaraxes.RadialAxisTickProps{
		Label:      "7",
		TextAnchor: "end",
		Theme:      styledTickTheme(),
	})
	for _, want := range []string{
		`stroke="#444444"`, `stroke-width="2"`,
		`fill="#eeeeee"`, `font-size="9"`, `font-family="TickFont"`,
		`text-anchor="end"`, ">7</text>",
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("missing %q: %s", want, svg)
		}
	}
}

// plainScale implements scales.Scale but not scales.ScaleWithBandwidth,
// exercising the non-bandwidth angle path.
type plainScale struct{}

func (plainScale) Type() scales.ScaleType { return scales.ScaleTypeLinear }
func (plainScale) Call(v any) float64     { return 0 }

func TestRenderCircularAxis_ScaleWithoutBandwidth(t *testing.T) {
	svg := polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
		Type:       polaraxes.CircularAxisOuter,
		Radius:     50,
		StartAngle: 0,
		EndAngle:   360,
		Scale:      plainScale{},
		Theme:      &theming.DefaultTheme,
	})
	// No ticks can be derived, but the domain arc still renders.
	if !strings.Contains(svg, "<path") {
		t.Errorf("missing domain arc: %s", svg)
	}
	if strings.Contains(svg, "<text") {
		t.Errorf("plain scale should yield no ticks: %s", svg)
	}
}

func TestRadiusScaleRange(t *testing.T) {
	inner, outer := polaraxes.RadiusScaleRange(radiusScale(), 0.0, 100.0)
	if inner != 0 || outer != 100 {
		t.Errorf("RadiusScaleRange = (%v, %v), want (0, 100)", inner, outer)
	}
}
