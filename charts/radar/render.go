package radar

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued RadarProps fields from Defaults.
func applyDefaults(p RadarProps) RadarProps {
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Curve == "" {
		p.Curve = Defaults.Curve
	}
	if p.BorderWidth == 0 {
		p.BorderWidth = Defaults.BorderWidth
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if p.GridLevels == 0 {
		p.GridLevels = Defaults.GridLevels
	}
	if p.GridShape == "" {
		p.GridShape = Defaults.GridShape
	}
	if p.GridLabelOffset == 0 {
		p.GridLabelOffset = Defaults.GridLabelOffset
	}
	if p.DotSize == 0 {
		p.DotSize = Defaults.DotSize
	}
	if isZeroColor(p.DotColor) {
		p.DotColor = Defaults.DotColor
	}
	if isZeroColor(p.DotBorderColor) {
		p.DotBorderColor = Defaults.DotBorderColor
	}
	if p.DotLabelYOffset == 0 {
		p.DotLabelYOffset = Defaults.DotLabelYOffset
	}
	if p.FillOpacity == 0 {
		p.FillOpacity = Defaults.FillOpacity
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// isZeroColor reports whether an InheritedColorConfig is unset.
func isZeroColor(c colors.InheritedColorConfig) bool {
	return c.Type == 0 && c.Static == "" && c.ThemePath == "" && c.FromPath == "" && c.Func == nil
}

// isZeroOrdinal reports whether an OrdinalColorScaleConfig is unset.
func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

// renderLayers renders the enabled layers as an inner SVG string. The
// interactive "slices" layer is omitted in the v2 static pipeline.
func renderLayers(props RadarProps, result RadarResult, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case RadarLayerGrid:
			b.WriteString(renderGridLayer(props, result, theme))
		case RadarLayerLayers:
			b.WriteString(renderShapeLayer(props, result))
		case RadarLayerDots:
			if props.DotsEnabled() {
				b.WriteString(renderDotsLayer(props, result))
			}
		case RadarLayerLegends:
			b.WriteString(renderLegendsLayer(props, result))
		case RadarLayerSlices:
			// Interactive-only; arrives with the Phase 5 client layer.
		}
	}
	return b.String()
}

// renderGridLayer renders the axis rays, concentric levels, and index labels
// inside a <g translate(center)>. Mirrors @nivo/radar RadarGrid.
func renderGridLayer(props RadarProps, result RadarResult, theme *theming.Theme) string {
	stroke, strokeWidth := gridLineStroke(theme)
	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtR(result.CenterX), fmtR(result.CenterY))

	// Axis rays: one per index, from center to the outer radius.
	for i := range result.Indices {
		angle := result.Rotation + float64(i)*result.AngleStep - math.Pi/2
		p := arcs.PositionFromAngle(angle, result.Radius)
		fmt.Fprintf(&b, `<line x1="0" y1="0" x2="%s" y2="%s"`, fmtR(p.X), fmtR(p.Y))
		writeStroke(&b, stroke, strokeWidth)
		b.WriteString("/>")
	}

	// Levels: concentric circles (circular) or polygons (linear), outermost
	// first (nivo reverses; ordering is cosmetic for static output).
	for j := 0; j < props.GridLevels; j++ {
		levelRadius := (result.Radius / float64(props.GridLevels)) * float64(j+1)
		if props.GridShape == GridShapeLinear {
			b.WriteString(`<path fill="none" d="`)
			b.WriteString(levelPolygonPath(levelRadius, result.Rotation, result.AngleStep, len(result.Indices)))
			b.WriteString(`"`)
			writeStroke(&b, stroke, strokeWidth)
			b.WriteString("/>")
		} else {
			fmt.Fprintf(&b, `<circle fill="none" r="%s"`, fmtR(math.Max(levelRadius, 0)))
			writeStroke(&b, stroke, strokeWidth)
			b.WriteString("/>")
		}
	}

	// Index labels.
	b.WriteString(renderGridLabels(props, result, theme))

	b.WriteString("</g>")
	return b.String()
}

// levelPolygonPath builds the closed polygon path for a linear grid level at
// the given radius. Vertices align with the axis rays.
func levelPolygonPath(radius, rotation, angleStep float64, n int) string {
	if n == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		angle := rotation + float64(i)*angleStep - math.Pi/2
		p := arcs.PositionFromAngle(angle, radius)
		if i == 0 {
			fmt.Fprintf(&b, "M%s,%s", fmtR(p.X), fmtR(p.Y))
		} else {
			fmt.Fprintf(&b, "L%s,%s", fmtR(p.X), fmtR(p.Y))
		}
	}
	b.WriteString("Z")
	return b.String()
}

// renderGridLabels renders one index label per axis at radius+labelOffset.
// Mirrors @nivo/radar RadarGridLabels + the default RadarGridLabel.
func renderGridLabels(props RadarProps, result RadarResult, theme *theming.Theme) string {
	fill, fontSize, fontFamily := axisTickText(theme)
	var b strings.Builder
	for i, index := range result.Indices {
		angle := result.Rotation + float64(i)*result.AngleStep - math.Pi/2
		p := arcs.PositionFromAngle(angle, result.Radius+props.GridLabelOffset)
		anchor := textAnchorFromAngle(angle)
		fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtR(p.X), fmtR(p.Y))
		fmt.Fprintf(&b, `<text dominant-baseline="central" text-anchor="%s"`, anchor)
		if fill != "" {
			fmt.Fprintf(&b, ` fill="%s"`, fill)
		}
		if fontSize != "" {
			fmt.Fprintf(&b, ` font-size="%s"`, fontSize)
		}
		if fontFamily != "" {
			fmt.Fprintf(&b, ` font-family="%s"`, fontFamily)
		}
		fmt.Fprintf(&b, ">%s</text></g>", escapeText(index))
	}
	return b.String()
}

// textAnchorFromAngle mirrors @nivo/radar RadarGridLabels.textAnchorFromAngle:
// degrees = radiansToDegrees(angle) + 90.
func textAnchorFromAngle(angleRad float64) string {
	angle := math.Mod(arcs.RadToDeg(angleRad)+90, 360)
	if angle < 0 {
		angle += 360
	}
	if angle <= 10 || angle >= 350 || (angle >= 170 && angle <= 190) {
		return "middle"
	}
	if angle > 180 {
		return "end"
	}
	return "start"
}

// renderShapeLayer renders one closed polygon per key inside a
// <g translate(center)>. Mirrors @nivo/radar RadarLayer.
func renderShapeLayer(props RadarProps, result RadarResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtR(result.CenterX), fmtR(result.CenterY))
	for i, serie := range result.Series {
		if serie.Path == "" {
			continue
		}
		fmt.Fprintf(&b, `<path d="%s" fill="%s" fill-opacity="%s"`, serie.Path, serie.Color, fmtR(props.FillOpacity))
		if props.BorderWidth > 0 && serie.Stroke != "" {
			fmt.Fprintf(&b, ` stroke="%s" stroke-width="%s"`, serie.Stroke, fmtR(props.BorderWidth))
		}
		b.WriteString(">")
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString("</path>")
	}
	b.WriteString("</g>")
	return b.String()
}

// renderDotsLayer renders the per-point dots inside a <g translate(center)>.
// Mirrors @nivo/radar RadarDots.
func renderDotsLayer(props RadarProps, result RadarResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtR(result.CenterX), fmtR(result.CenterY))
	for _, p := range result.Points {
		label := ""
		if props.EnableDotLabel {
			label = p.FormattedValue
		}
		dp := core.DotsItemProps{
			X:               p.X,
			Y:               p.Y,
			Size:            props.DotSize,
			Color:           p.Color,
			BorderWidth:     props.DotBorderWidth,
			BorderColor:     p.BorderColor,
			Label:           label,
			LabelTextAnchor: "middle",
			LabelYOffset:    props.DotLabelYOffset,
		}
		if props.Interactive {
			dp.Tooltip = interact.TooltipHTML(p.Color, p.Index+" - "+p.Key, p.FormattedValue)
		}
		b.WriteString(renderComponent(core.DotsItem(dp)))
	}
	b.WriteString("</g>")
	return b.String()
}

// renderLegendsLayer renders any box legends.
func renderLegendsLayer(props RadarProps, result RadarResult) string {
	if len(props.Legends) == 0 {
		return ""
	}
	var b strings.Builder
	for _, lg := range props.Legends {
		legend := lg
		if len(legend.Items) == 0 {
			legend.Items = result.LegendData
		}
		b.WriteString(renderComponent(legends.BoxLegendSvg(legends.BoxLegendSvgProps{
			Props:       legend,
			ChartWidth:  props.Width,
			ChartHeight: props.Height,
		})))
	}
	return b.String()
}

// --- theme helpers ---

func gridLineStroke(theme *theming.Theme) (string, float64) {
	if theme == nil {
		return "", 0
	}
	return strokeFromExtra(theme.Grid.Line.Extra)
}

func axisTickText(theme *theming.Theme) (fill, fontSize, fontFamily string) {
	if theme == nil {
		return "", "", ""
	}
	t := theme.Axis.Ticks.Text
	fill = t.Fill
	if fill == "" {
		fill = theme.Text.Fill
	}
	fs := t.FontSize
	if fs == nil {
		fs = theme.Text.FontSize
	}
	if fs != nil {
		fontSize = fmt.Sprintf("%v", fs)
	}
	fontFamily = t.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	return fill, fontSize, fontFamily
}

func strokeFromExtra(m map[string]any) (string, float64) {
	if m == nil {
		return "", 0
	}
	var stroke string
	if v, ok := m["stroke"]; ok {
		if s, ok := v.(string); ok {
			stroke = s
		} else {
			stroke = fmt.Sprintf("%v", v)
		}
	}
	var sw float64
	if v, ok := m["strokeWidth"]; ok {
		switch x := v.(type) {
		case float64:
			sw = x
		case int:
			sw = float64(x)
		}
	}
	return stroke, sw
}

func writeStroke(b *strings.Builder, stroke string, strokeWidth float64) {
	if stroke != "" {
		fmt.Fprintf(b, ` stroke="%s"`, stroke)
	}
	if strokeWidth > 0 {
		fmt.Fprintf(b, ` stroke-width="%s"`, fmtR(strokeWidth))
	}
}

// escapeText minimal XML escaping for labels.
func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// fmtR formats a float for SVG output (3 dp, trimmed).
func fmtR(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := fmt.Sprintf("%.3f", v)
	for len(s) > 1 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 1 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
