package scatterplot

import (
	"context"
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

// applyDefaults fills zero-valued ScatterPlotProps fields from Defaults.
func applyDefaults(p ScatterPlotProps) ScatterPlotProps {
	if p.XScale == nil {
		p.XScale = Defaults.XScale
	}
	if p.YScale == nil {
		p.YScale = Defaults.YScale
	}
	if p.NodeSize == 0 {
		p.NodeSize = Defaults.NodeSize
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
	// Default the visible axes (bottom + left) when no axis was provided.
	if p.AxisTop == nil && p.AxisBottom == nil && p.AxisLeft == nil && p.AxisRight == nil {
		p.AxisBottom = &axes.AxisProps{}
		p.AxisLeft = &axes.AxisProps{}
	}
	return p
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

// renderLayers renders the enabled layers as an inner SVG string. The
// interactive "mesh" layer is omitted in the v2 static pipeline; annotations
// are deferred (no scatterplot annotation specs in v2).
func renderLayers(props ScatterPlotProps, result ScatterPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case ScatterPlotLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case ScatterPlotLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case ScatterPlotLayerNodes:
			b.WriteString(renderNodesLayer(result))
		case ScatterPlotLayerMarkers:
			b.WriteString(renderMarkersLayer(props, result, dims))
		case ScatterPlotLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		case ScatterPlotLayerMesh, ScatterPlotLayerAnnotations:
			// mesh: interactive (Phase 5); annotations: deferred in v2.
		}
	}
	return b.String()
}

func renderGridLayer(props ScatterPlotProps, result ScatterPlotResult, dims core.Dimensions, theme *theming.Theme) string {
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

func renderAxesLayer(props ScatterPlotProps, result ScatterPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	// Bottom + left in one pass.
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
	// Top + right in a second pass when provided.
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

// resolveAxis builds an axes.AxisProps positioned at (originX, originY),
// substituting the default tick size/padding when unset.
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

func renderNodesLayer(result ScatterPlotResult) string {
	var b strings.Builder
	for _, n := range result.Nodes {
		b.WriteString(renderComponent(core.DotsItem(core.DotsItemProps{
			X: n.X, Y: n.Y, Size: n.Size, Color: n.Color,
		})))
	}
	return b.String()
}

func renderMarkersLayer(props ScatterPlotProps, result ScatterPlotResult, dims core.Dimensions) string {
	if len(props.Markers) == 0 {
		return ""
	}
	return renderComponent(core.CartesianMarkers(core.CartesianMarkersProps{
		Markers: props.Markers,
		Width:   dims.InnerWidth,
		Height:  dims.InnerHeight,
		XScale:  result.XScale.Call,
		YScale:  result.YScale.Call,
	}))
}

func renderLegendsLayer(props ScatterPlotProps, result ScatterPlotResult, dims core.Dimensions) string {
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
