package stream

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
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

// applyDefaults fills zero-valued StreamProps fields from Defaults.
func applyDefaults(p StreamProps) StreamProps {
	if p.OffsetType == "" {
		p.OffsetType = Defaults.OffsetType
	}
	if p.Order == "" {
		p.Order = Defaults.Order
	}
	if p.Curve == "" {
		p.Curve = Defaults.Curve
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.FillOpacity == 0 {
		p.FillOpacity = Defaults.FillOpacity
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	if p.AxisTop == nil && p.AxisBottom == nil && p.AxisLeft == nil && p.AxisRight == nil {
		p.AxisBottom = &axes.AxisProps{}
		p.AxisLeft = &axes.AxisProps{}
	}
	return p
}

func isZeroColor(c colors.InheritedColorConfig) bool {
	return c.Type == 0 && c.Static == "" && c.ThemePath == "" && c.FromPath == "" && c.Func == nil
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

// renderLayers renders the enabled layers as an inner SVG string. The dots and
// slices layers are omitted in the v2 static pipeline (dots default to off;
// slices are interactive — Phase 5).
func renderLayers(props StreamProps, result StreamResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case StreamLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case StreamLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case StreamLayerLayers:
			b.WriteString(renderAreasLayer(props, result))
		case StreamLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		case StreamLayerDots, StreamLayerSlices:
			// dots: default off; slices: interactive (Phase 5).
		}
	}
	return b.String()
}

func renderGridLayer(props StreamProps, result StreamResult, dims core.Dimensions, theme *theming.Theme) string {
	var s strings.Builder
	if props.EnableGridX {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "x", Scale: result.XScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridXValues, Theme: theme,
		})))
	}
	if props.EnableGridY {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "y", Scale: result.YScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridYValues, Theme: theme,
		})))
	}
	return s.String()
}

func renderAxesLayer(props StreamProps, result StreamResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	var xAxis, yAxis *axes.AxisProps
	if props.AxisBottom != nil {
		xAxis = resolveAxis(props.AxisBottom, "x", result.XScale, dims.InnerWidth, 0, dims.InnerHeight, "after")
	}
	if props.AxisLeft != nil {
		yAxis = resolveAxis(props.AxisLeft, "y", result.YScale, dims.InnerHeight, 0, 0, "after")
	}
	if xAxis != nil || yAxis != nil {
		b.WriteString(renderComponent(axes.Axes(axes.AxesProps{
			XAxis: xAxis, YAxis: yAxis, Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	var xTop, yRight *axes.AxisProps
	if props.AxisTop != nil {
		xTop = resolveAxis(props.AxisTop, "x", result.XScale, dims.InnerWidth, 0, 0, "before")
	}
	if props.AxisRight != nil {
		yRight = resolveAxis(props.AxisRight, "y", result.YScale, dims.InnerHeight, dims.InnerWidth, 0, "before")
	}
	if xTop != nil || yRight != nil {
		b.WriteString(renderComponent(axes.Axes(axes.AxesProps{
			XAxis: xTop, YAxis: yRight, Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	return b.String()
}

func resolveAxis(props *axes.AxisProps, axis string, scale scales.Scale, length, originX, originY float64, ticksPos string) *axes.AxisProps {
	out := *props
	out.Axis = axis
	out.Scale = scale
	out.Length = length
	out.X = originX
	out.Y = originY
	if out.TicksPosition == "" {
		out.TicksPosition = ticksPos
	}
	if out.TickSize == 0 && out.TickPadding == 0 {
		out.TickSize = 5
		out.TickPadding = 5
	}
	if out.TextAlign == "" {
		out.TextAlign = "center"
	}
	if out.TextBaseline == "" {
		out.TextBaseline = "middle"
	}
	if out.LegendPosition == "" {
		out.LegendPosition = axes.AxisLegendEnd
	}
	return &out
}

func renderAreasLayer(props StreamProps, result StreamResult) string {
	var b strings.Builder
	for i, layer := range result.Layers {
		if layer.Path == "" {
			continue
		}
		fmt.Fprintf(&b, `<path d="%s" fill="%s" fill-opacity="%s"`, layer.Path, layer.Color, fmtN(props.FillOpacity))
		if props.BorderWidth > 0 && layer.BorderColor != "" {
			fmt.Fprintf(&b, ` stroke="%s" stroke-width="%s"`, layer.BorderColor, fmtN(props.BorderWidth))
		}
		if props.Interactive {
			fmt.Fprintf(&b, ` style="cursor:pointer" data-tc-tooltip="%s"`, interact.EscapeAttr(interact.TooltipHTML(layer.Color, layer.Label, "")))
		}
		b.WriteString(">")
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString("</path>")
	}
	return b.String()
}

func renderLegendsLayer(props StreamProps, result StreamResult, dims core.Dimensions) string {
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
			ChartWidth:  dims.InnerWidth,
			ChartHeight: dims.InnerHeight,
		})))
	}
	return b.String()
}

// fmtN formats a float for SVG output (3 dp, trimmed).
func fmtN(v float64) string {
	s := fmt.Sprintf("%.3f", v)
	for len(s) > 1 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 1 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
