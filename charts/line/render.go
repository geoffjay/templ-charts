package line

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/tooltip"
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

// renderLines renders a Lines component to a string.
func renderLines(props LinesProps) string {
	return renderComponent(Lines(props))
}

// renderAreas renders an Areas component to a string.
func renderAreas(props AreasProps) string {
	return renderComponent(Areas(props))
}

// renderPoints renders a Points component to a string.
func renderPoints(props PointsProps) string {
	return renderComponent(Points(props))
}

// renderSlices renders a Slices component to a string.
func renderSlices(props SlicesProps) string {
	return renderComponent(Slices(props))
}

// renderMesh renders a Mesh component to a string.
func renderMesh(props MeshProps) string {
	return renderComponent(Mesh(props))
}

// renderCartesianMarkers renders core.CartesianMarkers to a string.
func renderCartesianMarkers(props core.CartesianMarkersProps) string {
	return renderComponent(core.CartesianMarkers(props))
}

// renderCrosshair renders a tooltip.Crosshair to a string.
func renderCrosshair(props tooltip.CrosshairProps) string {
	return renderComponent(tooltip.Crosshair(props))
}

// renderBoxLegend renders a legends.BoxLegendSvg to a string.
func renderBoxLegend(props legends.BoxLegendSvgProps) string {
	return renderComponent(legends.BoxLegendSvg(props))
}

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color (or "" for transparent).
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}

// resolveAxis builds an axes.AxisProps for the given axis, substituting the
// default axis config when props is nil. originX/originY position the axis
// <g>; length is the axis extent.
func resolveAxis(props *axes.AxisProps, defaultProps axes.AxisProps, axis string, scale scales.Scale, length, originX, originY float64) *axes.AxisProps {
	var out axes.AxisProps
	if props != nil {
		out = *props
	} else {
		out = defaultProps
	}
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

// renderLineLayers renders the enabled layers as an SVG string. Exported so
// the htmx handler can re-render just the inner SVG (for series-toggle swaps).
func renderLineLayers(layers []LineLayerId, props LineProps, result LineResult, dims core.Dimensions, theme *theming.Theme, bound core.SvgDefsAndFill) string {
	var b strings.Builder
	for _, layer := range layers {
		switch layer {
		case LineLayerGrid:
			if props.EnableGridX || props.EnableGridY {
				b.WriteString(renderGridLayer(props, result, dims, theme))
			}
		case LineLayerMarkers:
			b.WriteString(renderMarkersLayer(props, dims, result))
		case LineLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case LineLayerAreas:
			if props.EnableArea {
				b.WriteString(renderAreasLayer(props, result))
			}
		case LineLayerCrosshair:
			b.WriteString(renderCrosshairLayer(props, result, dims, theme))
		case LineLayerLines:
			b.WriteString(renderLinesLayer(props, result))
		case LineLayerPoints:
			if props.EnablePoints {
				b.WriteString(renderPointsLayer(props, result))
			}
		case LineLayerSlices:
			if props.IsInteractive && props.EnableSlices != "" {
				b.WriteString(renderSlicesLayer(props, result))
			}
		case LineLayerMesh:
			if props.IsInteractive && props.UseMesh && props.EnableSlices == "" {
				b.WriteString(renderMeshLayer(props, dims))
			}
		case LineLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		}
	}
	return b.String()
}

func renderGridLayer(props LineProps, result LineResult, dims core.Dimensions, theme *theming.Theme) string {
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

func renderAxesLayer(props LineProps, result LineResult, dims core.Dimensions, theme *theming.Theme) string {
	xAxis := resolveAxis(props.AxisBottom, defaultAxisBottomProps, "x", result.XScale, dims.InnerWidth, 0, dims.InnerHeight)
	yAxis := resolveAxis(props.AxisLeft, defaultAxisLeftProps, "y", result.YScale, dims.InnerHeight, 0, 0)
	return renderAxes(axes.AxesProps{
		XAxis: xAxis, YAxis: yAxis,
		Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
	})
}

func renderAreasLayer(props LineProps, result LineResult) string {
	return renderAreas(AreasProps{
		Series:        result.Series,
		AreaGenerator: result.AreaGenerator,
		AreaOpacity:   props.AreaOpacity,
		AreaBlendMode: string(props.AreaBlendMode),
		Animate:       props.Animate,
	})
}

func renderLinesLayer(props LineProps, result LineResult) string {
	return renderLines(LinesProps{
		Series:        result.Series,
		LineGenerator: result.LineGenerator,
		LineWidth:     props.LineWidth,
		Animate:       props.Animate,
	})
}

func renderPointsLayer(props LineProps, result LineResult) string {
	return renderPoints(PointsProps{
		Points:       result.Points,
		Size:         props.PointSize,
		BorderWidth:  props.PointBorderWidth,
		EnableLabel:  props.EnablePointLabel,
		LabelYOffset: props.PointLabelYOffset,
		ChartID:      props.ChartID,
	})
}

func renderSlicesLayer(props LineProps, result LineResult) string {
	return renderSlices(SlicesProps{
		Slices:  result.Slices,
		Axis:    props.EnableSlices,
		Debug:   props.DebugSlices,
		ChartID: props.ChartID,
	})
}

func renderMeshLayer(props LineProps, dims core.Dimensions) string {
	return renderMesh(MeshProps{
		Width:   dims.InnerWidth,
		Height:  dims.InnerHeight,
		Debug:   props.DebugMesh,
		ChartID: props.ChartID,
	})
}

func renderCrosshairLayer(props LineProps, result LineResult, dims core.Dimensions, theme *theming.Theme) string {
	if !props.EnableCrosshair {
		return ""
	}
	// Render the crosshair when the htmx hover endpoint has flagged an active
	// hover (mesh mode) or when DebugMesh is on (faint full-area guide).
	if !props.HasHover && !props.DebugMesh {
		return ""
	}
	x, y := props.HoverX, props.HoverY
	if props.DebugMesh && !props.HasHover {
		// Debug: draw a crosshair at the centre to visualise the layer.
		x = dims.InnerWidth / 2
		y = dims.InnerHeight / 2
	}
	cl := theme.Crosshair.Line
	return renderCrosshair(tooltip.CrosshairProps{
		Type:   props.CrosshairType,
		Width:  dims.InnerWidth,
		Height: dims.InnerHeight,
		X:      x,
		Y:      y,
		Theme: tooltip.CrosshairTheme{
			Stroke:          cl.Stroke,
			StrokeWidth:     cl.StrokeWidth,
			StrokeOpacity:   cl.StrokeOpacity,
			StrokeDasharray: cl.StrokeDasharray,
		},
	})
}

func renderMarkersLayer(props LineProps, dims core.Dimensions, result LineResult) string {
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

func renderLegendsLayer(props LineProps, result LineResult, dims core.Dimensions) string {
	if len(props.Legends) == 0 {
		return ""
	}
	var s strings.Builder
	for _, legend := range props.Legends {
		lp := legend
		if len(lp.Items) == 0 {
			// fill from legendData
			items := make([]legends.Datum, len(result.LegendData))
			for i, d := range result.LegendData {
				items[i] = legends.Datum{ID: d.ID, Label: d.Label, Color: d.Color, Hidden: d.Hidden}
			}
			lp.Items = items
		}
		if lp.ChartID == "" {
			lp.ChartID = props.ChartID
		}
		if lp.Toggle == "" && props.ChartID != "" {
			lp.Toggle = "1"
		}
		s.WriteString(renderBoxLegend(legends.BoxLegendSvgProps{
			Props:       lp,
			ChartWidth:  dims.InnerWidth,
			ChartHeight: dims.InnerHeight,
		}))
	}
	return s.String()
}

// seriesAsAny converts []ComputedSeries to []any for BindDefs.
func seriesAsAny(ss []ComputedSeries) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

// applyDefaults overlays Defaults onto a partial LineProps.
func applyDefaults(p LineProps) LineProps {
	if p.XScale == nil {
		p.XScale = Defaults.XScale
	}
	if p.YScale == (scales.ScaleLinearSpec{}) {
		p.YScale = Defaults.YScale
	}
	if p.Curve == "" {
		p.Curve = Defaults.Curve
	}
	if p.Colors.Type == 0 && p.Colors.Scheme == "" && p.Colors.Static == "" && len(p.Colors.Colors) == 0 && p.Colors.Func == nil && p.Colors.DatumPath == "" {
		p.Colors = Defaults.Colors
	}
	if p.PointColor.Type == 0 && p.PointColor.Static == "" && p.PointColor.ThemePath == "" && p.PointColor.FromPath == "" && p.PointColor.Func == nil {
		p.PointColor = Defaults.PointColor
	}
	if p.PointBorderColor.Type == 0 && p.PointBorderColor.Static == "" && p.PointBorderColor.ThemePath == "" && p.PointBorderColor.FromPath == "" && p.PointBorderColor.Func == nil {
		p.PointBorderColor = Defaults.PointBorderColor
	}
	if p.LineWidth == 0 {
		p.LineWidth = Defaults.LineWidth
	}
	if p.PointSize == 0 {
		p.PointSize = Defaults.PointSize
	}
	if p.AreaOpacity == 0 {
		p.AreaOpacity = Defaults.AreaOpacity
	}
	if p.AreaBlendMode == "" {
		p.AreaBlendMode = Defaults.AreaBlendMode
	}
	if p.CrosshairType == "" {
		p.CrosshairType = Defaults.CrosshairType
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	if p.PointLabel == nil {
		p.PointLabel = Defaults.PointLabel
	}
	if len(p.Layers) == 0 {
		p.Layers = DefaultLayers
	}
	if p.InitialHiddenIDs == nil {
		p.InitialHiddenIDs = Defaults.InitialHiddenIDs
	}
	return p
}

// guard against unused imports during refactoring
var _ = fmt.Sprintf
