package heatmap

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/canvas"
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
			// Annotations deferred (no heatmap annotation specs).
		}
	}
	return b.String()
}

// renderCanvas renders the heatmap with the Canvas backend: each cell becomes a
// FillRect (plus an optional StrokeRect border and FillText label) in a
// <canvas> draw-list, while grid sits in an SVG pane behind the canvas and
// axes/legends in an SVG pane in front (charts/canvas Compose). FillStyle /
// GlobalAlpha are emitted only when they change, so the draw-list stays compact
// across many cells. The draw-list is margin-translated to share the SVG panes'
// coordinate space.
func renderCanvas(props HeatMapProps, result HeatMapResult, dims core.Dimensions, theme *theming.Theme) string {
	rec := canvas.NewRecorder()
	rec.Translate(dims.Margin.Left, dims.Margin.Top)
	recordCells(rec, props, result)

	id := props.ChartID
	if id == "" {
		id = "tc-heatmap"
	}
	grid := renderGridLayer(props, result, dims, theme)
	overlay := renderCanvasOverlay(props, result, dims, theme)
	return canvas.Compose(id, dims.OuterWidth, dims.OuterHeight,
		dims.Margin.Left, dims.Margin.Top, themeBackground(theme),
		rec.Ops(), grid, overlay)
}

// recordCells writes one FillRect per cell into rec (matching the SVG cell
// <rect>), plus an optional border StrokeRect and value-label FillText. FillStyle
// and GlobalAlpha are tracked and emitted only on change, so the draw-list stays
// compact; the FillRect count corresponds one-to-one with the SVG cells layer.
func recordCells(rec *canvas.Recorder, props HeatMapProps, result HeatMapResult) {
	enableLabel := props.LabelsEnabled()
	if enableLabel {
		rec.Font("11px sans-serif")
		rec.TextAlign("center")
		rec.TextBaseline("middle")
	}

	curFill, curAlpha := "", 1.0
	setFill := func(c string) {
		if c != curFill {
			rec.FillStyle(c)
			curFill = c
		}
	}
	setAlpha := func(a float64) {
		if a != curAlpha {
			rec.GlobalAlpha(a)
			curAlpha = a
		}
	}

	for _, cell := range result.Cells {
		x := cell.XPos - cell.Width/2
		y := cell.YPos - cell.Height/2
		w := maxH(cell.Width, 0)
		h := maxH(cell.Height, 0)
		setAlpha(cell.Opacity)
		setFill(cell.Color)
		rec.FillRect(x, y, w, h)
		if props.BorderWidth > 0 && cell.BorderColor != "" {
			setAlpha(1)
			rec.StrokeStyle(cell.BorderColor)
			rec.LineWidth(props.BorderWidth)
			rec.StrokeRect(x, y, w, h)
		}
		if enableLabel && cell.Label != "" {
			setAlpha(1)
			setFill(cell.LabelTextColor)
			rec.FillText(cell.Label, cell.XPos, cell.YPos)
		}
	}
}

// renderCanvasOverlay renders the SVG panes drawn on top of the canvas marks:
// the axes and legends layers (grid goes behind via renderGridLayer; cells
// became the canvas draw-list).
func renderCanvasOverlay(props HeatMapProps, result HeatMapResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case HeatMapLayerAxes:
			b.WriteString(renderAxesLayer(props, result, theme))
		case HeatMapLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
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
	// Orientation-aware label anchor: a vertical axis should anchor its labels
	// toward the axis (end for a left axis, start for a right one) so they sit
	// beside the grid instead of centering on the offset point and overlapping
	// the cells. Top/bottom axes keep the centered default.
	if out.TextAlign == "" && axis == "y" {
		if out.TicksPosition == "before" {
			out.TextAlign = "end"
		} else {
			out.TextAlign = "start"
		}
	}
	return &out
}

func renderCellsLayer(props HeatMapProps, result HeatMapResult) string {
	var s strings.Builder
	interactive := props.Interactive && props.ChartID != ""
	for i, cell := range result.Cells {
		cp := HeatMapCellProps{
			Cell:         cell,
			BorderWidth:  props.BorderWidth,
			BorderRadius: props.BorderRadius,
			EnableLabel:  props.LabelsEnabled(),
			Animate:      props.Animate,
		}
		if props.Animate {
			cp.AnimateBegin = core.StaggerBegin(i, props.MotionStagger)
		}
		// Only cells with data are hoverable. hover is routed through the
		// client interactivity layer (charts/interact) — no server round-trip.
		if interactive && cell.Value != nil {
			cp.Tooltip = interact.TooltipHTML(cell.Color, cell.SerieID+" - "+cell.X, cell.FormattedValue)
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
