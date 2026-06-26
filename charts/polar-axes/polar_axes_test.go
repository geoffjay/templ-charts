package polaraxes

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func mustLinearScale(t *testing.T, min, max float64, size float64) scales.Scale {
	t.Helper()
	data := scales.ComputedSerieAxis{All: []any{min, max}, Min: min, Max: max}
	s := scales.ComputeScale(scales.ScaleLinearSpec{Min: scales.FloatVal(min), Max: scales.FloatVal(max)}, data, size, scales.ScaleAxisY)
	return s
}

func TestRenderCircularGrid_EmitsArcPaths(t *testing.T) {
	radiusScale := mustLinearScale(t, 0, 100, 100)
	svg := RenderCircularGrid(CircularGridProps{
		Scale:      radiusScale,
		StartAngle: 0,
		EndAngle:   270,
		Theme:      &theming.DefaultTheme,
	})
	if !strings.Contains(svg, "<path") {
		t.Fatalf("circular grid missing <path>: %s", svg)
	}
	if !strings.Contains(svg, "fill=\"none\"") {
		t.Fatalf("circular grid should not fill: %s", svg)
	}
}

func TestRenderRadialGrid_EmitsLines(t *testing.T) {
	angleScale := mustLinearScale(t, 0, 360, 360)
	svg := RenderRadialGrid(RadialGridProps{
		Scale:       angleScale,
		InnerRadius: 20,
		OuterRadius: 100,
		Theme:       &theming.DefaultTheme,
	})
	if !strings.Contains(svg, "<line") {
		t.Fatalf("radial grid missing <line>: %s", svg)
	}
}

func TestRenderPolarGrid_ComposesBoth(t *testing.T) {
	angleScale := mustLinearScale(t, 0, 360, 360)
	radiusScale := mustLinearScale(t, 0, 100, 100)
	svg := RenderPolarGrid(PolarGridProps{
		Center:             [2]float64{200, 200},
		EnableRadialGrid:   true,
		AngleScale:         angleScale,
		StartAngle:         0,
		EndAngle:           270,
		EnableCircularGrid: true,
		RadiusScale:        radiusScale,
		InnerRadius:        0,
		OuterRadius:        100,
		Theme:              &theming.DefaultTheme,
	})
	if !strings.Contains(svg, "translate(200,200)") {
		t.Fatalf("polar grid missing center translate: %s", svg)
	}
	if !strings.Contains(svg, "<line") || !strings.Contains(svg, "<path") {
		t.Fatalf("polar grid should contain both lines and paths: %s", svg)
	}
}

func TestRenderCircularAxis_IncludesDomainAndTicks(t *testing.T) {
	angleScale := mustLinearScale(t, 0, 360, 360)
	svg := RenderCircularAxis(CircularAxisProps{
		Type:       CircularAxisOuter,
		Center:     [2]float64{100, 100},
		Radius:     80,
		StartAngle: 0,
		EndAngle:   360,
		Scale:      angleScale,
		Theme:      &theming.DefaultTheme,
	})
	if !strings.Contains(svg, "translate(100,100)") {
		t.Fatalf("circular axis missing translate: %s", svg)
	}
	if !strings.Contains(svg, "<path") {
		t.Fatalf("circular axis missing domain arc path: %s", svg)
	}
	if !strings.Contains(svg, "<text") {
		t.Fatalf("circular axis missing tick labels: %s", svg)
	}
}

func TestRenderRadialAxis_RotatesToAngle(t *testing.T) {
	scale := mustLinearScale(t, 0, 10, 100)
	svg := RenderRadialAxis(RadialAxisProps{
		Center:        [2]float64{100, 100},
		Angle:         90,
		Scale:         scale,
		TicksPosition: TicksAfter,
		Theme:         &theming.DefaultTheme,
	})
	if !strings.Contains(svg, "translate(100,100)") {
		t.Fatalf("radial axis missing translate: %s", svg)
	}
	if !strings.Contains(svg, "rotate(") {
		t.Fatalf("radial axis missing rotate transform: %s", svg)
	}
}

func TestRenderPolarGrid_NilThemeReturnsEmpty(t *testing.T) {
	svg := RenderPolarGrid(PolarGridProps{Theme: nil})
	if svg != "" {
		t.Fatalf("nil theme should yield empty svg, got %q", svg)
	}
}

func TestAnimateOpacity(t *testing.T) {
	if got := animateOpacity(false); got != "" {
		t.Fatalf("animate=false should emit nothing, got %q", got)
	}
	if got := animateOpacity(true); !strings.Contains(got, "<animate") {
		t.Fatalf("animate=true should emit <animate>, got %q", got)
	}
}

func TestArcPath_Semicircle(t *testing.T) {
	p := arcPath(50, 0, 180)
	if !strings.HasPrefix(p, "M") {
		t.Fatalf("arc path should start with M: %s", p)
	}
	if !strings.Contains(p, "A50,50") {
		t.Fatalf("arc path should contain arc command: %s", p)
	}
}
