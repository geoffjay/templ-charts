package chord

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued ChordProps fields from Defaults.
func applyDefaults(p ChordProps) ChordProps {
	if p.InnerRadiusRatio == 0 {
		p.InnerRadiusRatio = Defaults.InnerRadiusRatio
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.ArcOpacity == 0 {
		p.ArcOpacity = Defaults.ArcOpacity
	}
	if p.ArcBorderWidth == 0 {
		p.ArcBorderWidth = Defaults.ArcBorderWidth
	}
	if isZeroInherited(p.ArcBorderColor) {
		p.ArcBorderColor = Defaults.ArcBorderColor
	}
	if p.RibbonOpacity == 0 {
		p.RibbonOpacity = Defaults.RibbonOpacity
	}
	if p.RibbonBorderWidth == 0 {
		p.RibbonBorderWidth = Defaults.RibbonBorderWidth
	}
	if isZeroInherited(p.RibbonBorderColor) {
		p.RibbonBorderColor = Defaults.RibbonBorderColor
	}
	if p.RibbonBlendMode == "" {
		p.RibbonBlendMode = Defaults.RibbonBlendMode
	}
	if p.EnableLabel == nil {
		p.EnableLabel = Defaults.EnableLabel
	}
	if p.Label == "" {
		p.Label = Defaults.Label
	}
	if p.LabelOffset == 0 {
		p.LabelOffset = Defaults.LabelOffset
	}
	if isZeroInherited(p.LabelTextColor) {
		p.LabelTextColor = Defaults.LabelTextColor
	}
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

func isZeroInherited(c colors.InheritedColorConfig) bool {
	return c.Type == colors.InheritedColorTypeStatic && c.Static == "" && c.Func == nil &&
		c.ThemePath == "" && c.FromPath == "" && len(c.Modifiers) == 0
}

// renderLayers renders the enabled layers as an inner SVG string. width/height
// are the inner dimensions (used for legends).
func renderLayers(props ChordProps, result ChordResult, width, height float64, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case ChordLayerRibbons:
			b.WriteString(centered(result.Center, renderRibbonsLayer(props, result, theme)))
		case ChordLayerArcs:
			b.WriteString(centered(result.Center, renderArcsLayer(props, result, theme)))
		case ChordLayerLabels:
			if props.LabelEnabled() {
				b.WriteString(centered(result.Center, renderLabelsLayer(props, result, theme)))
			}
		case ChordLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, width, height))
		}
	}
	return b.String()
}

// centered wraps inner SVG in a group translated to the circle center, matching
// @nivo/chord which translates each of the ribbons/arcs/labels layers.
func centered(center [2]float64, inner string) string {
	if inner == "" {
		return ""
	}
	return `<g transform="translate(` + fmtF(center[0]) + `,` + fmtF(center[1]) + `)">` + inner + `</g>`
}

func renderRibbonsLayer(props ChordProps, result ChordResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.RibbonBorderColor, theme)
	var b strings.Builder
	for _, r := range result.Ribbons {
		b.WriteString(`<path d="`)
		b.WriteString(r.Path)
		b.WriteString(`" fill="`)
		b.WriteString(r.Color)
		b.WriteString(`" fill-opacity="`)
		b.WriteString(fmtF(props.RibbonOpacity))
		b.WriteString(`"`)
		if props.RibbonBorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(getBorderColor(map[string]any{"color": r.Color}))
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.RibbonBorderWidth))
			b.WriteString(`"`)
		}
		b.WriteString(` style="mix-blend-mode:`)
		b.WriteString(props.RibbonBlendMode)
		if props.Interactive {
			b.WriteString(`;pointer-events:auto`)
		}
		b.WriteString(`"`)
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(r.Color, r.Source.ID+" > "+r.Target.ID, r.Source.FormattedValue)))
			b.WriteString(`"`)
		} else {
			b.WriteString(` pointer-events="none"`)
		}
		b.WriteString(`></path>`)
	}
	return b.String()
}

func renderArcsLayer(props ChordProps, result ChordResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.ArcBorderColor, theme)
	var b strings.Builder
	for _, a := range result.Arcs {
		b.WriteString(`<path d="`)
		b.WriteString(a.Path)
		b.WriteString(`" fill="`)
		b.WriteString(a.Color)
		b.WriteString(`" fill-opacity="`)
		b.WriteString(fmtF(props.ArcOpacity))
		b.WriteString(`"`)
		if props.ArcBorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(getBorderColor(map[string]any{"color": a.Color}))
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.ArcBorderWidth))
			b.WriteString(`"`)
		}
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(a.Color, a.Label, a.FormattedValue)))
			b.WriteString(`" style="pointer-events:auto"`)
		} else {
			b.WriteString(` pointer-events="none"`)
		}
		b.WriteString(`></path>`)
	}
	return b.String()
}

// renderLabelsLayer positions each arc's label radially, mirroring @nivo/chord
// ChordLabels (getPolarLabelProps at radius + labelOffset with the mid angle).
func renderLabelsLayer(props ChordProps, result ChordResult, theme *theming.Theme) string {
	getFill := colors.GetInheritedColorGenerator(props.LabelTextColor, theme)
	_, fontSize, fontFamily := labelsTextStyle(theme)
	labelRadius := result.Radius + props.LabelOffset

	var b strings.Builder
	for _, a := range result.Arcs {
		mid := a.StartAngle + (a.EndAngle-a.StartAngle)/2
		x, y, rotate, anchor, baseline := polarLabelProps(labelRadius, mid, props.LabelRotation)
		fill := getFill(map[string]any{"color": a.Color})
		b.WriteString(labelText(x, y, rotate, anchor, baseline, a.Label, fill, fontSize, fontFamily))
	}
	return b.String()
}

// polarLabelProps ports @nivo/core getPolarLabelProps (svg engine): the label
// anchor point at (angle - π/2, radius), the text rotation, and the SVG
// text-anchor / dominant-baseline.
func polarLabelProps(radius, angle, rotation float64) (x, y, rotate float64, align, baseline string) {
	a := angle - math.Pi/2
	x = math.Cos(a) * radius
	y = math.Sin(a) * radius
	rotate = angle * 180 / math.Pi
	align = "middle"
	baseline = "alphabetic"
	if rotation > 0 {
		align = "end"
		baseline = "central"
	} else if rotation < 0 {
		align = "start"
		baseline = "central"
	}
	if rotation != 0 && rotate > 180 {
		rotate -= 180
		if align == "end" {
			align = "start"
		} else {
			align = "end"
		}
	}
	rotate += rotation
	return x, y, rotate, align, baseline
}

func renderLegendsLayer(props ChordProps, result ChordResult, width, height float64) string {
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
			ChartWidth:  width,
			ChartHeight: height,
		})))
	}
	return b.String()
}

// labelText emits a <text> translated to (x,y) and rotated, with the given
// text-anchor and dominant-baseline — matching @nivo/chord's label transform.
func labelText(x, y, rotation float64, anchor, baseline, s, fill string, fontSize float64, fontFamily string) string {
	var b strings.Builder
	b.WriteString(`<text transform="translate(`)
	b.WriteString(fmtF(x))
	b.WriteString(`,`)
	b.WriteString(fmtF(y))
	b.WriteString(`)`)
	if rotation != 0 {
		b.WriteString(` rotate(`)
		b.WriteString(fmtF(rotation))
		b.WriteString(`)`)
	}
	b.WriteString(`" text-anchor="`)
	b.WriteString(anchor)
	b.WriteString(`" dominant-baseline="`)
	b.WriteString(baseline)
	b.WriteString(`" style="pointer-events:none;fill:`)
	b.WriteString(fill)
	if fontSize > 0 {
		b.WriteString(`;font-size:`)
		b.WriteString(fmtF(fontSize))
		b.WriteString(`px`)
	}
	if fontFamily != "" {
		b.WriteString(`;font-family:`)
		b.WriteString(fontFamily)
	}
	b.WriteString(`">`)
	b.WriteString(templ.EscapeString(s))
	b.WriteString(`</text>`)
	return b.String()
}

func labelsTextStyle(theme *theming.Theme) (fill string, fontSize float64, fontFamily string) {
	if theme == nil {
		return "", 0, ""
	}
	t := theme.Labels.Text
	fill = t.Fill
	if fill == "" {
		fill = theme.Text.Fill
	}
	fs := t.FontSize
	if fs == nil {
		fs = theme.Text.FontSize
	}
	fontFamily = t.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	return fill, toNum(fs), fontFamily
}

func toNum(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	}
	return 0
}

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer).
func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
