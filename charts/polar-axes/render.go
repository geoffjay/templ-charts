package polaraxes

import (
	"fmt"
	"math"
	"strings"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// renderOpts is the resolved styling shared by the polar-axis renderers.
type renderOpts struct {
	axis    theming.AxisTheme
	grid    theming.GridTheme
	animate bool
}

// resolveAxisTheme mirrors @nivo/theming useExtendedAxisTheme: merges optional
// overrides onto the theme's axis block (falling back to the default theme).
func resolveAxisTheme(theme *theming.Theme, overrides *theming.PartialAxisTheme) theming.AxisTheme {
	if theme == nil {
		return theming.AxisTheme{}
	}
	return theming.ExtendAxisTheme(theme.Axis, overrides)
}

func resolveGridTheme(theme *theming.Theme) theming.GridTheme {
	if theme == nil {
		return theming.GridTheme{}
	}
	return theme.Grid
}

// fmtN formats a float for SVG output (3 dp, trimmed), matching the rest of
// the library's formatting.
func fmtN(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.3f", v), "0"), ".")
}

// strokeFromExtra reads a stroke color from a theme line's Extra map.
func strokeFromExtra(m map[string]any) string {
	if m == nil {
		return ""
	}
	if v, ok := m["stroke"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// strokeWidthFromExtra reads a stroke-width from a theme line's Extra map.
func strokeWidthFromExtra(m map[string]any) float64 {
	if m == nil {
		return 0
	}
	if v, ok := m["strokeWidth"]; ok {
		switch x := v.(type) {
		case float64:
			return x
		case int:
			return float64(x)
		}
	}
	return 0
}

// tickValues computes the tick values for a scale, mirroring nivo's use of
// getScaleTicks. When ticksSpec has no count set, nivo defaults to 10.
func tickValues(scale scales.Scale, spec scales.TicksSpec) []any {
	if !spec.HasCount && !spec.HasValues && !spec.HasInterval {
		spec = scales.TicksSpec{Count: 10, HasCount: true}
	}
	return scales.GetScaleTicks(scale, spec)
}

// tickLabel formats a tick value via the configured formatter (or "%v").
func tickLabel(v any, format core.ValueFormat[any]) string {
	f := core.GetValueFormatter[any](format)
	return f(v)
}

// angleScaleToTicks resolves the per-tick angle (degrees, already -90
// offset) for a circular axis. For band/point scales the tick is centered.
func angleScaleToTicks(scale scales.Scale) []tickAngle {
	values := scales.GetScaleTicks(scale, scales.TicksSpec{Count: 10, HasCount: true})
	var angle func(any) float64
	if bw, ok := scale.(scales.ScaleWithBandwidth); ok {
		angle = scales.CenterScale(scale)
		_ = bw
	} else {
		angle = scale.Call
	}
	out := make([]tickAngle, len(values))
	for i, v := range values {
		out[i] = tickAngle{value: v, angle: angle(v) - 90}
	}
	return out
}

type tickAngle struct {
	value any
	angle float64 // degrees, -90 offset already applied
}

// radialTickPositions resolves the per-tick position along the radial axis
// (centered for band/point scales). Mirrors RadialAxis's `ticks` useMemo.
func radialTickPositions(scale scales.Scale, spec scales.TicksSpec) []tickPos {
	values := tickValues(scale, spec)
	out := make([]tickPos, len(values))
	for i, v := range values {
		pos := scale.Call(v)
		if bw, ok := scale.(scales.ScaleWithBandwidth); ok {
			pos += bw.Bandwidth() / 2
		}
		out[i] = tickPos{value: v, position: pos}
	}
	return out
}

type tickPos struct {
	value    any
	position float64
}

// linePositions returns the x1,y1,x2,y2 of a circular-axis tick line at
// `angle` (degrees, -90 offset already applied), spanning innerRadius →
// outerRadius. Mirrors CircularAxis getLinePositions.
func linePositions(angle, innerRadius, outerRadius float64) (x1, y1, x2, y2 float64) {
	start := arcs.PositionFromAngle(arcs.DegToRad(angle), innerRadius)
	end := arcs.PositionFromAngle(arcs.DegToRad(angle), outerRadius)
	return start.X, start.Y, end.X, end.Y
}

// textPosition returns the label x/y for a circular-axis tick at `angle`
// (degrees) and radius. Mirrors CircularAxis getTextPosition.
func textPosition(angle, radius float64) (x, y float64) {
	p := arcs.PositionFromAngle(arcs.DegToRad(angle), radius)
	return p.X, p.Y
}

// arcPath builds an SVG arc path string for a circular grid line / circular
// axis domain, sweeping from startAngle to endAngle (degrees, user-facing
// 0=top convention) at the given radius. This is the ArcLine equivalent
// (nivo uses @nivo/arcs ArcLine). We approximate with a simple circular arc
// since grid lines don't need cornerRadius.
func arcPath(radius, startAngleDeg, endAngleDeg float64) string {
	s := arcs.DegToRad(startAngleDeg)
	e := arcs.DegToRad(endAngleDeg)
	p0 := arcs.PositionFromAngle(s, radius)
	p1 := arcs.PositionFromAngle(e, radius)
	large := 0
	if endAngleDeg-startAngleDeg > 180 {
		large = 1
	}
	return fmt.Sprintf("M%s,%sA%s,%s 0 %d 1 %s,%s",
		fmtN(p0.X), fmtN(p0.Y), fmtN(radius), fmtN(radius), large, fmtN(p1.X), fmtN(p1.Y))
}

// animateOpacity emits a SMIL fade-in <animate> for a container <g> when
// animate is true (mirrors react-spring's opacity 0→1 enter).
func animateOpacity(animate bool) string {
	if !animate {
		return ""
	}
	return `<animate attributeName="opacity" from="0" to="1" begin="0s" dur="0.6s" fill="freeze"/>`
}

// animateRotate emits a SMIL <animate> on the transform attribute growing the
// rotation from rawAngle-90 back to the target... we instead render the final
// transform and fade in (simpler, still matches "enter" intent).
func animateRotate(animate bool, finalRotation float64) string {
	// Reserved for a future faithful rotate-from-0; v1 scaffold uses fade.
	_ = animate
	_ = finalRotation
	return ""
}
