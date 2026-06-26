package polaraxes

import (
	"fmt"
	"strings"

	"github.com/geoffjay/templ-charts/charts/arcs"
)

// RenderRadialAxis renders a radial axis: a <g translate(center)> rotating a
// subgroup to the axis angle, then placing one tick per scale value along
// that rotated x-axis. Mirrors @nivo/polar-axes RadialAxis (minus
// react-spring; SMIL fade-in when Animate).
//
// Returns the SVG fragment string.
func RenderRadialAxis(props RadialAxisProps) string {
	tickSize := props.TickSize
	if tickSize == 0 {
		tickSize = 5
	}
	tickPadding := props.TickPadding
	if tickPadding == 0 {
		tickPadding = 5
	}
	extraRotation := props.TickRotation

	angle := arcs.NormalizeAngleDegrees(props.Angle)

	var textAnchor string
	var lineX, textX, tickRotation float64

	if props.TicksPosition == TicksBefore {
		tickRotation = 90 + extraRotation
		if angle <= 90 {
			lineX = -tickSize
			textX = lineX - tickPadding
			textAnchor = "end"
		} else if angle < 270 {
			lineX = tickSize
			textX = lineX + tickPadding
			textAnchor = "start"
			tickRotation -= 180
		} else {
			lineX = -tickSize
			textX = lineX - tickPadding
			textAnchor = "end"
		}
	} else { // after
		tickRotation = 90 + extraRotation
		if angle < 90 {
			lineX = tickSize
			textX = lineX + tickPadding
			textAnchor = "start"
		} else if angle < 270 {
			lineX = -tickSize
			textX = lineX - tickPadding
			textAnchor = "end"
			tickRotation -= 180
		} else {
			lineX = tickSize
			textX = lineX + tickPadding
			textAnchor = "start"
		}
	}

	axisTheme := resolveAxisTheme(props.Theme, props.Style)
	ticks := radialTickPositions(props.Scale, props.Ticks)

	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)" style="pointer-events:none">`, fmtN(props.Center[0]), fmtN(props.Center[1]))
	fmt.Fprintf(&b, `<g transform="rotate(%s)">`, fmtN(props.Angle-90))

	for _, tick := range ticks {
		label := tickLabel(tick.value, props.Format)
		tp := RadialAxisTickProps{
			Label:      label,
			TextAnchor: textAnchor,
			Theme:      axisTheme,
			Y:          tick.position,
			Length:     lineX,
			TextX:      textX,
			Rotation:   tickRotation,
			Animate:    props.Animate,
		}
		if props.TickComponent != nil {
			b.WriteString((*props.TickComponent)(tp))
		} else {
			b.WriteString(RenderRadialAxisTick(tp))
		}
	}

	b.WriteString("</g></g>")
	return b.String()
}

// RenderRadialAxisTick renders one radial-axis tick (a translated+rotated
// <g> with a tick line and label). Mirrors @nivo/polar-axes RadialAxisTick.
func RenderRadialAxisTick(props RadialAxisTickProps) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,0) rotate(%s)"`, fmtN(props.Y), fmtN(props.Rotation))
	if props.Animate {
		b.WriteString(` opacity="0"`)
	}
	b.WriteString(">")
	tStroke := strokeFromExtra(props.Theme.Ticks.Line.Extra)
	tStrokeWidth := strokeWidthFromExtra(props.Theme.Ticks.Line.Extra)
	fmt.Fprintf(&b, `<line x2="%s"`, fmtN(props.Length))
	if tStroke != "" {
		fmt.Fprintf(&b, ` stroke="%s"`, tStroke)
	}
	if tStrokeWidth > 0 {
		fmt.Fprintf(&b, ` stroke-width="%s"`, fmtN(tStrokeWidth))
	}
	b.WriteString("/>")
	fill := props.Theme.Ticks.Text.Fill
	fontSize := props.Theme.Ticks.Text.FontSize
	fontFamily := props.Theme.Ticks.Text.FontFamily
	fmt.Fprintf(&b, `<text dx="%s" text-anchor="%s" dominant-baseline="central"`, fmtN(props.TextX), props.TextAnchor)
	if fill != "" {
		fmt.Fprintf(&b, ` fill="%s"`, fill)
	}
	if fontSize != nil {
		fmt.Fprintf(&b, ` font-size="%v"`, fontSize)
	}
	if fontFamily != "" {
		fmt.Fprintf(&b, ` font-family="%s"`, fontFamily)
	}
	fmt.Fprintf(&b, `>%s</text>`, escapeText(props.Label))
	b.WriteString(animateOpacity(props.Animate))
	b.WriteString("</g>")
	return b.String()
}
