package heatmap

import (
	"context"
	"fmt"
	"net/url"
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

// applyDefaults fills zero-valued HeatMapProps fields from Defaults.
func applyDefaults(p HeatMapProps) HeatMapProps {
	if p.Colors.Type == "" && p.Colors.Scheme == "" {
		p.Colors = Defaults.Colors
	} else if p.Colors.Type == "" {
		p.Colors.Type = Defaults.Colors.Type
	}
	if p.EmptyColor == "" {
		p.EmptyColor = Defaults.EmptyColor
	}
	if p.Opacity == 0 {
		p.Opacity = Defaults.Opacity
	}
	if p.ActiveOpacity == 0 {
		p.ActiveOpacity = Defaults.ActiveOpacity
	}
	if p.InactiveOpacity == 0 {
		p.InactiveOpacity = Defaults.InactiveOpacity
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if isZeroColor(p.LabelTextColor) {
		p.LabelTextColor = Defaults.LabelTextColor
	}
	if len(p.Layers) == 0 {
		p.Layers = DefaultLayers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	// Default the visible axes (top + left) when no axis at all was provided.
	if p.AxisTop == nil && p.AxisBottom == nil && p.AxisLeft == nil && p.AxisRight == nil {
		p.AxisTop = &axes.AxisProps{}
		p.AxisLeft = &axes.AxisProps{}
	}
	return p
}

// isZeroColor reports whether an InheritedColorConfig is unset (so a default
// should apply). Mirrors the unset check used by charts/bar.
func isZeroColor(c colors.InheritedColorConfig) bool {
	return c.Type == 0 && c.Static == "" && c.ThemePath == "" && c.FromPath == "" && c.Func == nil
}

// renderLayers renders the enabled layers as an inner SVG string.
func renderLayers(props HeatMapProps, result HeatMapResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case HeatMapLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case HeatMapLayerAxes:
			b.WriteString(renderAxesLayer(props, result, theme))
		case HeatMapLayerCells:
			b.WriteString(renderCellsLayer(props, result))
		case HeatMapLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		case HeatMapLayerAnnotations:
			// Annotations deferred (no heatmap annotation specs in v2).
		}
	}
	return b.String()
}

func renderGridLayer(props HeatMapProps, result HeatMapResult, dims core.Dimensions, theme *theming.Theme) string {
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

func renderAxesLayer(props HeatMapProps, result HeatMapResult, theme *theming.Theme) string {
	// Top axis spans the x band scale; left axis spans the y band scale. The
	// forceSquare offsets shift each axis origin so axes track the grid.
	var xAxis, yAxis *axes.AxisProps
	if props.AxisTop != nil {
		xAxis = resolveAxis(props.AxisTop, "x", result.XScale, result.Width, result.OffsetX, result.OffsetY, "before")
	} else if props.AxisBottom != nil {
		xAxis = resolveAxis(props.AxisBottom, "x", result.XScale, result.Width, result.OffsetX, result.OffsetY+result.Height, "after")
	}
	if props.AxisLeft != nil {
		yAxis = resolveAxis(props.AxisLeft, "y", result.YScale, result.Height, result.OffsetX, result.OffsetY, "before")
	} else if props.AxisRight != nil {
		yAxis = resolveAxis(props.AxisRight, "y", result.YScale, result.Height, result.OffsetX+result.Width, result.OffsetY, "after")
	}
	return renderComponent(axes.Axes(axes.AxesProps{
		XAxis: xAxis, YAxis: yAxis,
		Width: result.Width, Height: result.Height, Theme: theme,
	}))
}

// resolveAxis builds an axes.AxisProps positioned at (originX, originY),
// substituting the default tick size/padding so labels sit outside the grid
// (an empty AxisProps{} otherwise leaves TickSize/TickPadding at 0, dropping
// the labels onto the grid edge).
func resolveAxis(props *axes.AxisProps, axis string, scale scales.Scale, length, originX, originY float64, ticksPos string) *axes.AxisProps {
	if props == nil {
		return nil
	}
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
		out.TickSize = axes.DefaultAxisProps.TickSize
		out.TickPadding = axes.DefaultAxisProps.TickPadding
	}
	return &out
}

func renderCellsLayer(props HeatMapProps, result HeatMapResult) string {
	var s strings.Builder
	interactive := props.IsInteractive && props.ChartID != ""
	for _, cell := range result.Cells {
		cp := HeatMapCellProps{
			Cell:         cell,
			BorderWidth:  props.BorderWidth,
			BorderRadius: props.BorderRadius,
			EnableLabel:  props.LabelsEnabled(),
		}
		// Only cells with data are hoverable.
		if interactive && cell.Value != nil {
			cp.HxGet = fmt.Sprintf("/charts/%s/hover?cell=%s", props.ChartID, url.QueryEscape(cell.ID))
			cp.HxTrigger = "mouseenter"
			cp.HxSwap = "innerHTML"
			cp.HxTarget = fmt.Sprintf("#tooltip-%s", props.ChartID)
		}
		s.WriteString(renderComponent(HeatMapCell(cp)))
	}
	return s.String()
}

func renderLegendsLayer(props HeatMapProps, result HeatMapResult, dims core.Dimensions) string {
	if len(props.Legends) == 0 {
		return ""
	}
	var s strings.Builder
	for _, lg := range props.Legends {
		length := lg.Length
		if length == 0 {
			length = 200
		}
		thickness := lg.Thickness
		if thickness == 0 {
			thickness = 16
		}
		s.WriteString(renderComponent(legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
			Scale:       result.ColorScale,
			Min:         result.MinValue,
			Max:         result.MaxValue,
			Width:       length,
			Height:      thickness,
			Anchor:      lg.Anchor,
			TranslateX:  lg.TranslateX,
			TranslateY:  lg.TranslateY,
			Title:       lg.Title,
			Samples:     32,
			ChartWidth:  dims.InnerWidth,
			ChartHeight: dims.InnerHeight,
		})))
	}
	return s.String()
}
