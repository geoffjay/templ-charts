package polaraxes

import (
	"fmt"
	"strings"

	"github.com/geoffjay/templ-charts/charts/scales"
)

// RenderPolarGrid renders the polar grid: a <g translate(center)> containing
// the optional radial grid (rays from inner to outer radius at each angle
// tick) and the optional circular grid (concentric arcs at each radius tick).
// Mirrors @nivo/polar-axes PolarGrid.
func RenderPolarGrid(props PolarGridProps) string {
	if props.Theme == nil {
		return ""
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtN(props.Center[0]), fmtN(props.Center[1]))
	if props.EnableRadialGrid {
		b.WriteString(RenderRadialGrid(RadialGridProps{
			Scale:       props.AngleScale,
			InnerRadius: props.InnerRadius,
			OuterRadius: props.OuterRadius,
			Theme:       props.Theme,
			Animate:     props.Animate,
		}))
	}
	if props.EnableCircularGrid {
		b.WriteString(RenderCircularGrid(CircularGridProps{
			Scale:      props.RadiusScale,
			Ticks:      props.CircularGridTicks,
			StartAngle: props.StartAngle,
			EndAngle:   props.EndAngle,
			Theme:      props.Theme,
			Animate:    props.Animate,
		}))
	}
	b.WriteString("</g>")
	return b.String()
}

// RenderRadialGrid renders the radial grid lines: for each angle tick of
// `scale`, a <g rotate(angle-90)> containing a horizontal <line> from
// innerRadius to outerRadius. Mirrors @nivo/polar-axes RadialGrid.
func RenderRadialGrid(props RadialGridProps) string {
	values := scales.GetScaleTicks(props.Scale, props.Ticks)
	grid := resolveGridTheme(props.Theme)
	stroke := strokeFromExtra(grid.Line.Extra)
	strokeWidth := strokeWidthFromExtra(grid.Line.Extra)

	var b strings.Builder
	for _, value := range values {
		angle := props.Scale.Call(value) - 90
		fmt.Fprintf(&b, `<g transform="rotate(%s)"`, fmtN(angle))
		if props.Animate {
			b.WriteString(` opacity="0"`)
		}
		b.WriteString(">")
		fmt.Fprintf(&b, `<line x1="%s" x2="%s"`, fmtN(props.InnerRadius), fmtN(props.OuterRadius))
		if stroke != "" {
			fmt.Fprintf(&b, ` stroke="%s"`, stroke)
		}
		if strokeWidth > 0 {
			fmt.Fprintf(&b, ` stroke-width="%s"`, fmtN(strokeWidth))
		}
		b.WriteString("/>")
		b.WriteString(animateOpacity(props.Animate))
		b.WriteString("</g>")
	}
	return b.String()
}

// RenderCircularGrid renders the circular grid lines: for each radius tick
// of `scale`, an arc path from startAngle to endAngle at that radius. Mirrors
// @nivo/polar-axes CircularGrid (which uses @nivo/arcs ArcLine).
func RenderCircularGrid(props CircularGridProps) string {
	values := scales.GetScaleTicks(props.Scale, props.Ticks)
	grid := resolveGridTheme(props.Theme)
	stroke := strokeFromExtra(grid.Line.Extra)
	strokeWidth := strokeWidthFromExtra(grid.Line.Extra)

	startAngle := props.StartAngle - 90
	endAngle := props.EndAngle - 90

	var b strings.Builder
	for _, value := range values {
		radius := props.Scale.Call(value)
		if bw, ok := props.Scale.(scales.ScaleWithBandwidth); ok {
			radius += bw.Bandwidth() / 2
		}
		b.WriteString(`<path d="`)
		b.WriteString(arcPath(radius, startAngle, endAngle))
		b.WriteString(`" fill="none"`)
		if stroke != "" {
			fmt.Fprintf(&b, ` stroke="%s"`, stroke)
		}
		if strokeWidth > 0 {
			fmt.Fprintf(&b, ` stroke-width="%s"`, fmtN(strokeWidth))
		}
		if props.Animate {
			b.WriteString(` opacity="0"`)
		}
		b.WriteString(">")
		b.WriteString(animateOpacity(props.Animate))
		b.WriteString("</path>")
	}
	return b.String()
}

// RadiusScaleRange returns (innerRadius, outerRadius) for a polar grid given
// the min/max of the radius domain and the radius scale. nivo derives these
// from radiusScale.range(); our Scale interface doesn't expose the range, so
// callers use this helper (inner = scale.Call(min), outer = scale.Call(max))
// and pass the result into PolarGridProps.InnerRadius/OuterRadius.
func RadiusScaleRange(scale scales.Scale, min, max any) (inner, outer float64) {
	return scale.Call(min), scale.Call(max)
}
