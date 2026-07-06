package marimekko

import (
	"context"
	"math"
	"strconv"
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

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued MarimekkoProps fields from Defaults.
func applyDefaults(p MarimekkoProps) MarimekkoProps {
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Offset == "" {
		p.Offset = Defaults.Offset
	}
	if p.InnerPadding == 0 {
		p.InnerPadding = Defaults.InnerPadding
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
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

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

func renderLayers(props MarimekkoProps, result MarimekkoResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case MarimekkoLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case MarimekkoLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case MarimekkoLayerBars:
			b.WriteString(renderBarsLayer(props, result, theme))
		case MarimekkoLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		}
	}
	return b.String()
}

// scaleAxes returns the (x, y) scales for the current layout.
func scaleAxes(props MarimekkoProps, result MarimekkoResult) (xScale, yScale scales.Scale) {
	if props.Layout == MarimekkoLayoutHorizontal {
		return result.StackScale, result.ThicknessScale
	}
	return result.ThicknessScale, result.StackScale
}

func renderGridLayer(props MarimekkoProps, result MarimekkoResult, dims core.Dimensions, theme *theming.Theme) string {
	xScale, yScale := scaleAxes(props, result)
	var s strings.Builder
	if props.EnableGridX {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "x", Scale: xScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	if props.EnableGridY {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "y", Scale: yScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	return s.String()
}

func renderAxesLayer(props MarimekkoProps, result MarimekkoResult, dims core.Dimensions, theme *theming.Theme) string {
	xScale, yScale := scaleAxes(props, result)
	var b strings.Builder
	var xAxis, yAxis *axes.AxisProps
	if props.AxisBottom != nil {
		xAxis = resolveAxis(props.AxisBottom, "x", xScale, dims.InnerWidth, 0, dims.InnerHeight, "after")
	}
	if props.AxisLeft != nil {
		yAxis = resolveAxis(props.AxisLeft, "y", yScale, dims.InnerHeight, 0, 0, "after")
	}
	if xAxis != nil || yAxis != nil {
		b.WriteString(renderComponent(axes.Axes(axes.AxesProps{
			XAxis: xAxis, YAxis: yAxis, Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
		})))
	}
	var xTop, yRight *axes.AxisProps
	if props.AxisTop != nil {
		xTop = resolveAxis(props.AxisTop, "x", xScale, dims.InnerWidth, 0, 0, "before")
	}
	if props.AxisRight != nil {
		yRight = resolveAxis(props.AxisRight, "y", yScale, dims.InnerHeight, dims.InnerWidth, 0, "before")
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

// renderBarsLayer emits one <rect> per stacked segment.
func renderBarsLayer(props MarimekkoProps, result MarimekkoResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	var b strings.Builder
	for i, bar := range result.Bars {
		if bar.Width <= 0 || bar.Height <= 0 {
			continue
		}
		b.WriteString(`<rect x="`)
		b.WriteString(fmtF(bar.X))
		b.WriteString(`" y="`)
		b.WriteString(fmtF(bar.Y))
		b.WriteString(`" width="`)
		b.WriteString(fmtF(bar.Width))
		b.WriteString(`" height="`)
		b.WriteString(fmtF(bar.Height))
		b.WriteString(`" fill="`)
		b.WriteString(bar.Color)
		b.WriteString(`"`)
		if props.BorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(getBorderColor(map[string]any{"color": bar.Color}))
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.BorderWidth))
			b.WriteString(`"`)
		}
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(bar.Color, bar.Index+" - "+bar.DimensionID, bar.FormattedValue)))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		if props.Animate {
			b.WriteString(`>`)
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
			b.WriteString(`</rect>`)
		} else {
			b.WriteString(`></rect>`)
		}
	}
	return b.String()
}

func renderLegendsLayer(props MarimekkoProps, result MarimekkoResult, dims core.Dimensions) string {
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

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer).
func fmtF(v float64) string {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		r = 0 // normalize -0 to +0 for cross-platform-stable output
	}
	return strconv.FormatFloat(r, 'g', -1, 64)
}
