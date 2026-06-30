package boxplot

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
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

// applyDefaults fills zero-valued BoxPlotProps fields from Defaults.
func applyDefaults(p BoxPlotProps) BoxPlotProps {
	if len(p.Quantiles) == 0 {
		p.Quantiles = Defaults.Quantiles
	}
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Padding == 0 {
		p.Padding = Defaults.Padding
	}
	if p.InnerPadding == 0 {
		p.InnerPadding = Defaults.InnerPadding
	}
	if p.Opacity == 0 {
		p.Opacity = Defaults.Opacity
	}
	if p.ColorBy == "" {
		p.ColorBy = Defaults.ColorBy
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if p.MedianWidth == 0 {
		p.MedianWidth = Defaults.MedianWidth
	}
	if isZeroColor(p.MedianColor) {
		p.MedianColor = Defaults.MedianColor
	}
	if p.WhiskerWidth == 0 {
		p.WhiskerWidth = Defaults.WhiskerWidth
	}
	if isZeroColor(p.WhiskerColor) {
		p.WhiskerColor = Defaults.WhiskerColor
	}
	if p.WhiskerEndSize == 0 {
		p.WhiskerEndSize = Defaults.WhiskerEndSize
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

// renderLayers renders the enabled layers as an inner SVG string. Markers and
// annotations are deferred in the v2 static pipeline.
func renderLayers(props BoxPlotProps, result BoxPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case BoxPlotLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case BoxPlotLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case BoxPlotLayerBoxPlots:
			b.WriteString(renderBoxPlotsLayer(props, result))
		case BoxPlotLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		case BoxPlotLayerMarkers, BoxPlotLayerAnnotations:
			// Deferred in v2.
		}
	}
	return b.String()
}

func renderGridLayer(props BoxPlotProps, result BoxPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var s strings.Builder
	if props.EnableGridX {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "x", Scale: result.XScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	if props.EnableGridY {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "y", Scale: result.YScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	return s.String()
}

func renderAxesLayer(props BoxPlotProps, result BoxPlotResult, dims core.Dimensions, theme *theming.Theme) string {
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

func renderBoxPlotsLayer(props BoxPlotProps, result BoxPlotResult) string {
	var b strings.Builder
	for _, box := range result.Boxes {
		bw := box.Bandwidth
		fmt.Fprintf(&b, `<g transform="%s"`, box.Transform)
		if box.Opacity > 0 && box.Opacity < 1 {
			fmt.Fprintf(&b, ` opacity="%s"`, fmtN(box.Opacity))
		}
		b.WriteString(">")
		// Box rect.
		fmt.Fprintf(&b, `<rect x="%s" y="%s" width="%s" height="%s"`,
			fmtN(-bw/2), fmtN(box.RectY), fmtN(bw), fmtN(box.ValueInterval))
		if props.BorderRadius > 0 {
			fmt.Fprintf(&b, ` rx="%s" ry="%s"`, fmtN(props.BorderRadius), fmtN(props.BorderRadius))
		}
		fmt.Fprintf(&b, ` fill="%s"`, box.Color)
		if props.BorderWidth > 0 && box.BorderColor != "" {
			fmt.Fprintf(&b, ` stroke="%s" stroke-width="%s"`, box.BorderColor, fmtN(props.BorderWidth))
		}
		b.WriteString("></rect>")
		// Median line.
		fmt.Fprintf(&b, `<line x1="%s" x2="%s" y1="0" y2="0" stroke="%s" stroke-width="%s"></line>`,
			fmtN(-bw/2), fmtN(bw/2), box.MedianColor, fmtN(props.MedianWidth))
		// Lower whisker (q25 → q10) + cap.
		writeWhisker(&b, box.VD1, box.VD0, box.WhiskerEnd, box.WhiskerColor, props.WhiskerWidth)
		// Upper whisker (q75 → q90) + cap.
		writeWhisker(&b, box.VD3, box.VD4, box.WhiskerEnd, box.WhiskerColor, props.WhiskerWidth)
		b.WriteString("</g>")
	}
	return b.String()
}

func writeWhisker(b *strings.Builder, distStart, distEnd, end float64, color string, width float64) {
	fmt.Fprintf(b, `<line x1="0" x2="0" y1="%s" y2="%s" stroke="%s" stroke-width="%s"></line>`,
		fmtN(distStart), fmtN(distEnd), color, fmtN(width))
	if end > 0 {
		fmt.Fprintf(b, `<line x1="%s" x2="%s" y1="%s" y2="%s" stroke="%s" stroke-width="%s"></line>`,
			fmtN(-end), fmtN(end), fmtN(distEnd), fmtN(distEnd), color, fmtN(width))
	}
}

func renderLegendsLayer(props BoxPlotProps, result BoxPlotResult, dims core.Dimensions) string {
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

func fmtN(v float64) string {
	if v != v || v > 1e308 || v < -1e308 {
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
