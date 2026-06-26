package polaraxes

import (
	"fmt"
	"strings"
)

// RenderCircularAxis renders the full circular-axis SVG: a <g translate(center)>
// containing the domain arc line (an ArcLine) plus one tick per scale value.
// Mirrors @nivo/polar-axes CircularAxis (minus react-spring; SMIL fade-in
// when Animate).
//
// Returns the SVG fragment string.
func RenderCircularAxis(props CircularAxisProps) string {
	tickSize := props.TickSize
	if tickSize == 0 {
		tickSize = 5
	}
	tickPadding := props.TickPadding
	if tickPadding == 0 {
		tickPadding = 12
	}

	startAngle := props.StartAngle - 90
	endAngle := props.EndAngle - 90

	axisTheme := resolveAxisTheme(props.Theme, props.Style)

	outerRadius := props.Radius + tickSize
	textRadius := props.Radius + tickSize + tickPadding
	if props.Type == CircularAxisInner {
		outerRadius = props.Radius - tickSize
		textRadius = outerRadius - tickPadding
	}

	ticks := angleScaleToTicks(props.Scale)

	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)" style="pointer-events:none">`, fmtN(props.Center[0]), fmtN(props.Center[1]))

	// Domain arc line.
	stroke := strokeFromExtra(axisTheme.Domain.Line.Extra)
	strokeWidth := strokeWidthFromExtra(axisTheme.Domain.Line.Extra)
	b.WriteString(`<path d="`)
	b.WriteString(arcPath(props.Radius, startAngle, endAngle))
	b.WriteString(`" fill="none"`)
	if stroke != "" {
		fmt.Fprintf(&b, ` stroke="%s"`, stroke)
	}
	if strokeWidth > 0 {
		fmt.Fprintf(&b, ` stroke-width="%s"`, fmtN(strokeWidth))
	}
	b.WriteString(">")
	b.WriteString(animateOpacity(props.Animate))
	b.WriteString("</path>")

	// Ticks.
	for _, tick := range ticks {
		x1, y1, x2, y2 := linePositions(tick.angle, props.Radius, outerRadius)
		textX, textY := textPosition(tick.angle, textRadius)
		label := tickLabel(tick.value, props.Format)
		tp := CircularAxisTickProps{
			Label:      label,
			TextAnchor: "middle",
			Theme:      axisTheme,
			X1:         x1, Y1: y1, X2: x2, Y2: y2,
			TextX:   textX,
			TextY:   textY,
			Animate: props.Animate,
		}
		if props.TickComponent != nil {
			b.WriteString((*props.TickComponent)(tp))
		} else {
			b.WriteString(RenderCircularAxisTick(tp))
		}
	}

	b.WriteString("</g>")
	return b.String()
}

// RenderCircularAxisTick renders one circular-axis tick (line + label).
// Mirrors @nivo/polar-axes CircularAxisTick.
func RenderCircularAxisTick(props CircularAxisTickProps) string {
	var b strings.Builder
	b.WriteString(`<g`)
	if props.Animate {
		b.WriteString(` opacity="0"`)
	}
	b.WriteString(">")
	tStroke := strokeFromExtra(props.Theme.Ticks.Line.Extra)
	tStrokeWidth := strokeWidthFromExtra(props.Theme.Ticks.Line.Extra)
	fmt.Fprintf(&b, `<line x1="%s" y1="%s" x2="%s" y2="%s"`, fmtN(props.X1), fmtN(props.Y1), fmtN(props.X2), fmtN(props.Y2))
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
	fmt.Fprintf(&b, `<text dx="%s" dy="%s" text-anchor="middle" dominant-baseline="central"`, fmtN(props.TextX), fmtN(props.TextY))
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

// escapeText minimal XML escaping for tick labels.
func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
