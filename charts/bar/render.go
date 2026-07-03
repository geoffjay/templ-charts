package bar

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/core"
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

// renderGrid renders an axes.Grid component to a string.
func renderGrid(props axes.GridProps) string {
	return renderComponent(axes.Grid(props))
}

// renderAxes renders an axes.Axes component to a string.
func renderAxes(props axes.AxesProps) string {
	return renderComponent(axes.Axes(props))
}

// renderBarItem renders a BarItem component to a string.
func renderBarItem(props BarItemProps) string {
	return renderComponent(BarItem(props))
}

// renderBarTotals renders a BarTotals component to a string.
func renderBarTotals(props BarTotalsProps) string {
	return renderComponent(BarTotals(props))
}

// renderCartesianMarkers renders core.CartesianMarkers to a string.
func renderCartesianMarkers(props core.CartesianMarkersProps) string {
	return renderComponent(core.CartesianMarkers(props))
}

// renderBarLegends renders BarLegends to a string.
func renderBarLegends(props BarLegendsProps) string {
	return renderComponent(BarLegends(props))
}

// renderBarAnnotations renders BarAnnotations to a string.
func renderBarAnnotations(props BarAnnotationsProps) string {
	return renderComponent(BarAnnotations(props))
}

// resolveAxis builds an axes.AxisProps for the given axis, merging defaults.
// originX/originY position the axis <g>; length is the axis extent.
func resolveAxis(props *axes.AxisProps, axis string, scale scales.Scale, length, originX, originY float64) *axes.AxisProps {
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
		out.TicksPosition = "after"
	}
	if out.LegendPosition == "" {
		out.LegendPosition = axes.AxisLegendEnd
	}
	return &out
}

// renderBarLayers renders the enabled layers as an SVG string. Exported so the
// htmx handler can re-render just the inner SVG (for series-toggle swaps).
func renderBarLayers(layers []BarLayerId, props BarProps, result BarResult, dims core.Dimensions, theme *theming.Theme, bound core.SvgDefsAndFill) string {
	var b strings.Builder
	for _, layer := range layers {
		switch layer {
		case BarLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case BarLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case BarLayerBars:
			b.WriteString(renderBarsLayer(props, result, bound, theme))
		case BarLayerTotals:
			if props.EnableTotals {
				b.WriteString(renderTotalsLayer(props, result, theme))
			}
		case BarLayerMarkers:
			b.WriteString(renderMarkersLayer(props, dims, result))
		case BarLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		case BarLayerAnnotations:
			b.WriteString(renderAnnotationsLayer(props, result))
		}
	}
	return b.String()
}

func renderGridLayer(props BarProps, result BarResult, dims core.Dimensions, theme *theming.Theme) string {
	var s strings.Builder
	if props.EnableGridX {
		s.WriteString(renderGrid(axes.GridProps{
			Axis: "x", Scale: result.XScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridXValues, Theme: theme,
		}))
	}
	if props.EnableGridY {
		s.WriteString(renderGrid(axes.GridProps{
			Axis: "y", Scale: result.YScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridYValues, Theme: theme,
		}))
	}
	return s.String()
}

func renderAxesLayer(props BarProps, result BarResult, dims core.Dimensions, theme *theming.Theme) string {
	xAxis := resolveAxis(props.AxisBottom, "x", result.XScale, dims.InnerWidth, 0, dims.InnerHeight)
	yAxis := resolveAxis(props.AxisLeft, "y", result.YScale, dims.InnerHeight, 0, 0)
	return renderAxes(axes.AxesProps{
		XAxis: xAxis, YAxis: yAxis,
		Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
	})
}

func renderBarsLayer(props BarProps, result BarResult, bound core.SvgDefsAndFill, theme *theming.Theme) string {
	var s strings.Builder
	horizontal := props.Layout == LayoutHorizontal
	for i, barDatum := range result.BarsWithValue {
		ll := result.ComputeLabelLayout(barDatum.Width, barDatum.Height)
		fill := barDatum.Color
		if override, ok := bound.FillByNodeIndex[i]; ok {
			fill = override
		}
		borderColor := result.GetBorderColor(barDatumWithColor(barDatum))
		labelColor := result.GetLabelColor(barDatumWithColor(barDatum))
		label := result.GetLabel(barDatum.Data)
		shouldLabel := result.ShouldRenderLabel(barDatum.Width, barDatum.Height)
		bip := BarItemProps{
			Bar:               barDatum,
			Label:             label,
			ShouldRenderLabel: shouldLabel,
			BorderRadius:      props.BorderRadius,
			BorderWidth:       props.BorderWidth,
			BorderColor:       borderColor,
			Fill:              fill,
			LabelColor:        labelColor,
			LabelLayout:       ll,
			Animate:           props.Animate,
			Horizontal:        horizontal,
			IsActive:          props.HoveredKey != "" && props.HoveredKey == barDatum.Key,
			Dimmed:            props.HoveredKey != "" && props.HoveredKey != barDatum.Key,
		}
		if props.Interactive && props.ChartID != "" {
			bip.HxGet = fmt.Sprintf("/charts/%s/hover?bar=%s", props.ChartID, barDatum.Key)
			bip.HxTrigger = "mouseenter"
			bip.HxSwap = "innerHTML"
			bip.HxTarget = fmt.Sprintf("#tooltip-%s", props.ChartID)
		}
		s.WriteString(renderBarItem(bip))
	}
	return s.String()
}

func renderTotalsLayer(props BarProps, result BarResult, theme *theming.Theme) string {
	return renderBarTotals(BarTotalsProps{
		Data:   result.BarTotals,
		Layout: string(props.Layout),
		Theme: BarTotalsTheme{
			FontSize:   themeLabelFontSize(theme),
			FontFamily: theme.Labels.Text.FontFamily,
			Fill:       theme.Labels.Text.Fill,
		},
	})
}

func renderMarkersLayer(props BarProps, dims core.Dimensions, result BarResult) string {
	if len(props.Markers) == 0 {
		return ""
	}
	return renderCartesianMarkers(core.CartesianMarkersProps{
		Markers: props.Markers,
		Width:   dims.InnerWidth,
		Height:  dims.InnerHeight,
		XScale:  result.XScale.Call,
		YScale:  result.YScale.Call,
	})
}

func renderLegendsLayer(props BarProps, result BarResult, dims core.Dimensions) string {
	return renderBarLegends(BarLegendsProps{
		Width:   dims.InnerWidth,
		Height:  dims.InnerHeight,
		Legends: result.LegendsWithData,
	})
}

func renderAnnotationsLayer(props BarProps, result BarResult) string {
	return renderBarAnnotations(BarAnnotationsProps{
		Bars:        result.Bars,
		Annotations: props.Annotations,
	})
}

// barDatumWithColor wraps a ComputedBarDatum as a map for the inherited-color
// generators (nivo's getBorderColor/getLabelColor receive a datum with a
// `color` field and `data` sub-object).
func barDatumWithColor(b ComputedBarDatum) any {
	return map[string]any{
		"color":  b.Color,
		"width":  b.Width,
		"height": b.Height,
		"x":      b.X,
		"y":      b.Y,
		"data": map[string]any{
			"id":             b.Data.ID,
			"value":          b.Data.Value,
			"formattedValue": b.Data.FormattedValue,
			"indexValue":     b.Data.IndexValue,
			"color":          b.Color,
		},
	}
}

// barsAsAny converts []ComputedBarDatum to []any for BindDefs.
func barsAsAny(bs []ComputedBarDatum) []any {
	out := make([]any, len(bs))
	for i, b := range bs {
		out[i] = b
	}
	return out
}

// applyDefaults overlays Defaults onto a partial BarProps so the templ
// component can accept a partial props struct (mirrors nivo's
// svgDefaultProps merging).
func applyDefaults(p BarProps) BarProps {
	if p.IndexBy == nil {
		p.IndexBy = Defaults.IndexBy
	}
	if len(p.Keys) == 0 {
		p.Keys = Defaults.Keys
	}
	if p.GroupMode == "" {
		p.GroupMode = Defaults.GroupMode
	}
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Padding == 0 {
		p.Padding = Defaults.Padding
	}
	if p.ValueScale == (scales.ScaleLinearSpec{}) {
		p.ValueScale = Defaults.ValueScale
	}
	if p.IndexScale == (scales.ScaleBandSpec{}) {
		p.IndexScale = Defaults.IndexScale
	}
	if p.ColorBy == "" {
		p.ColorBy = Defaults.ColorBy
	}
	if p.Colors.Type == 0 && p.Colors.Scheme == "" && p.Colors.Static == "" && len(p.Colors.Colors) == 0 && p.Colors.Func == nil && p.Colors.DatumPath == "" {
		p.Colors = Defaults.Colors
	}
	if p.BorderColor.Type == 0 && p.BorderColor.Static == "" && p.BorderColor.ThemePath == "" && p.BorderColor.FromPath == "" && p.BorderColor.Func == nil {
		p.BorderColor = Defaults.BorderColor
	}
	if p.LabelTextColor.Type == 0 && p.LabelTextColor.Static == "" && p.LabelTextColor.ThemePath == "" && p.LabelTextColor.FromPath == "" && p.LabelTextColor.Func == nil {
		p.LabelTextColor = Defaults.LabelTextColor
	}
	if p.LabelPosition == "" {
		p.LabelPosition = Defaults.LabelPosition
	}
	if p.TotalsOffset == 0 {
		p.TotalsOffset = Defaults.TotalsOffset
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	if len(p.Layers) == 0 {
		p.Layers = DefaultLayers
	}
	// EnableLabel defaults to true; callers must set false explicitly. Since
	// Go's bool zero value is false, we can't distinguish "unset" from "set
	// to false" — we leave EnableLabel as the caller set it. The Bar templ
	// component documents that EnableLabel=true is the default.
	return p
}

// ensure annotations import is used even if future refactors drop calls.
var _ = annotations.AnnotationTypeCircle
