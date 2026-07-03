package sankey

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
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

// applyDefaults fills zero-valued SankeyProps fields from Defaults.
func applyDefaults(p SankeyProps) SankeyProps {
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Align == "" {
		p.Align = Defaults.Align
	}
	if p.Sort == "" {
		p.Sort = Defaults.Sort
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.NodeOpacity == 0 {
		p.NodeOpacity = Defaults.NodeOpacity
	}
	if p.NodeThickness == 0 {
		p.NodeThickness = Defaults.NodeThickness
	}
	if p.NodeSpacing == 0 {
		p.NodeSpacing = Defaults.NodeSpacing
	}
	if p.NodeBorderWidth == 0 {
		p.NodeBorderWidth = Defaults.NodeBorderWidth
	}
	if isZeroInherited(p.NodeBorderColor) {
		p.NodeBorderColor = Defaults.NodeBorderColor
	}
	if p.LinkOpacity == 0 {
		p.LinkOpacity = Defaults.LinkOpacity
	}
	if p.LinkBlendMode == "" {
		p.LinkBlendMode = Defaults.LinkBlendMode
	}
	if p.EnableLabels == nil {
		p.EnableLabels = Defaults.EnableLabels
	}
	if p.LabelPosition == "" {
		p.LabelPosition = Defaults.LabelPosition
	}
	if p.LabelPadding == 0 {
		p.LabelPadding = Defaults.LabelPadding
	}
	if p.LabelOrientation == "" {
		p.LabelOrientation = Defaults.LabelOrientation
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
// are the inner dimensions (used for label side selection and legends).
func renderLayers(props SankeyProps, result SankeyResult, width, height float64, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case SankeyLayerLinks:
			b.WriteString(renderLinksLayer(props, result))
		case SankeyLayerNodes:
			b.WriteString(renderNodesLayer(props, result, theme))
		case SankeyLayerLabels:
			b.WriteString(renderLabelsLayer(props, result, width, height, theme))
		case SankeyLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, width, height))
		}
	}
	return b.String()
}

func renderLinksLayer(props SankeyProps, result SankeyResult) string {
	var b strings.Builder
	for i, l := range result.Links {
		b.WriteString(`<path d="`)
		b.WriteString(l.Path)
		b.WriteString(`" fill="`)
		b.WriteString(l.Color)
		b.WriteString(`" fill-opacity="`)
		b.WriteString(fmtF(props.LinkOpacity))
		b.WriteString(`" style="mix-blend-mode:`)
		b.WriteString(props.LinkBlendMode)
		b.WriteString(`"`)
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(l.Color, l.Source+" > "+l.Target, l.FormattedValue)))
			b.WriteString(`" pointer-events="auto"`)
		} else {
			b.WriteString(` pointer-events="none"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</path>`)
	}
	return b.String()
}

func renderNodesLayer(props SankeyProps, result SankeyResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.NodeBorderColor, theme)
	var b strings.Builder
	for i, n := range result.Nodes {
		b.WriteString(`<rect x="`)
		b.WriteString(fmtF(n.X))
		b.WriteString(`" y="`)
		b.WriteString(fmtF(n.Y))
		b.WriteString(`" width="`)
		b.WriteString(fmtF(math.Max(n.Width, 0)))
		b.WriteString(`" height="`)
		b.WriteString(fmtF(math.Max(n.Height, 0)))
		b.WriteString(`"`)
		if props.NodeBorderRadius > 0 {
			b.WriteString(` rx="`)
			b.WriteString(fmtF(props.NodeBorderRadius))
			b.WriteString(`" ry="`)
			b.WriteString(fmtF(props.NodeBorderRadius))
			b.WriteString(`"`)
		}
		b.WriteString(` fill="`)
		b.WriteString(n.Color)
		b.WriteString(`" fill-opacity="`)
		b.WriteString(fmtF(props.NodeOpacity))
		b.WriteString(`"`)
		if props.NodeBorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(getBorderColor(map[string]any{"color": n.Color}))
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.NodeBorderWidth))
			b.WriteString(`" stroke-opacity="`)
			b.WriteString(fmtF(props.NodeOpacity))
			b.WriteString(`"`)
		}
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(n.Color, n.Label, n.FormattedValue)))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</rect>`)
	}
	return b.String()
}

// renderLabelsLayer positions each node label inside/outside the node, choosing
// the side based on which half of the chart the node sits in (matching
// @nivo/sankey SankeyLabels).
func renderLabelsLayer(props SankeyProps, result SankeyResult, width, height float64, theme *theming.Theme) string {
	if !props.LabelsEnabled() {
		return ""
	}
	getFill := colors.GetInheritedColorGenerator(props.LabelTextColor, theme)
	_, fontSize, fontFamily := labelsTextStyle(theme)
	rotation := 0.0
	if props.LabelOrientation == SankeyLabelVertical {
		rotation = -90
	}
	horizontal := props.Layout != SankeyLayoutVertical

	var b strings.Builder
	for _, n := range result.Nodes {
		var x, y float64
		var anchor string
		if horizontal {
			y = n.Y + n.Height/2
			inside := props.LabelPosition == SankeyLabelInside
			if n.X < width/2 {
				if inside {
					x = n.X1 + props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "middle", "start")
				} else {
					x = n.X - props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "middle", "end")
				}
			} else {
				if inside {
					x = n.X - props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "middle", "end")
				} else {
					x = n.X1 + props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "middle", "start")
				}
			}
		} else {
			x = n.X + n.Width/2
			inside := props.LabelPosition == SankeyLabelInside
			if n.Y < height/2 {
				if inside {
					y = n.Y1 + props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "end", "middle")
				} else {
					y = n.Y - props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "start", "middle")
				}
			} else {
				if inside {
					y = n.Y - props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "start", "middle")
				} else {
					y = n.Y1 + props.LabelPadding
					anchor = ternary(props.LabelOrientation == SankeyLabelVertical, "end", "middle")
				}
			}
		}
		fill := getFill(map[string]any{"color": n.Color})
		b.WriteString(labelText(x, y, rotation, anchor, n.Label, fill, fontSize, fontFamily))
	}
	return b.String()
}

func renderLegendsLayer(props SankeyProps, result SankeyResult, width, height float64) string {
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

// labelText emits a <text> translated to (x,y) and rotated, with central
// baseline — matching @nivo/sankey's label transform.
func labelText(x, y, rotation float64, anchor, s, fill string, fontSize float64, fontFamily string) string {
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
	b.WriteString(`" dominant-baseline="central" style="pointer-events:none;fill:`)
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

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer).
func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
