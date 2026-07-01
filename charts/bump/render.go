package bump

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

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued BumpProps fields from Defaults.
func applyDefaults(p BumpProps) BumpProps {
	if p.Interpolation == "" {
		p.Interpolation = Defaults.Interpolation
	}
	if p.XPadding == 0 {
		p.XPadding = Defaults.XPadding
	}
	if p.XOuterPadding == 0 {
		p.XOuterPadding = Defaults.XOuterPadding
	}
	if p.YOuterPadding == 0 {
		p.YOuterPadding = Defaults.YOuterPadding
	}
	if p.LineWidth == 0 {
		p.LineWidth = Defaults.LineWidth
	}
	if p.Opacity == 0 {
		p.Opacity = Defaults.Opacity
	}
	if p.PointSize == 0 {
		p.PointSize = Defaults.PointSize
	}
	if p.StartLabel == nil {
		p.StartLabel = Defaults.StartLabel
	}
	if p.EndLabel == nil {
		p.EndLabel = Defaults.EndLabel
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

// renderLayers renders the enabled layers as an inner SVG string. The mesh
// layer emits an accurate voronoi-mesh hover overlay (charts/interact, backed
// by internal/d3/delaunay) when Interactive && UseMesh.
func renderLayers(props BumpProps, result BumpResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case BumpLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case BumpLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case BumpLayerLines:
			b.WriteString(renderLinesLayer(result))
		case BumpLayerPoints:
			b.WriteString(renderPointsLayer(props, result))
		case BumpLayerLabels:
			b.WriteString(renderLabelsLayer(props, result, theme))
		case BumpLayerMesh:
			b.WriteString(renderMeshLayer(props, result, dims))
		}
	}
	// Legends are a standalone prop (not a layer id) in nivo bump.
	b.WriteString(renderLegendsLayer(props, result, dims))
	return b.String()
}

// renderMeshLayer emits the voronoi-mesh hover overlay over the defined rank
// points. Active only when Interactive && UseMesh; when DebugMesh is set, the
// actual voronoi cells are drawn as a faint guide.
func renderMeshLayer(props BumpProps, result BumpResult, dims core.Dimensions) string {
	if !props.Interactive || !props.UseMesh {
		return ""
	}
	var pts []interact.MeshPoint
	for _, s := range result.Series {
		for _, p := range s.Points {
			if !p.Defined {
				continue
			}
			pts = append(pts, interact.MeshPoint{
				X:    p.X,
				Y:    p.Y,
				HTML: interact.TooltipHTML(p.Color, p.SerieID, p.FormattedY),
			})
		}
	}
	return interact.MeshOverlay(pts, dims.InnerWidth, dims.InnerHeight, props.DebugMesh, 0)
}

func renderGridLayer(props BumpProps, result BumpResult, dims core.Dimensions, theme *theming.Theme) string {
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

func renderAxesLayer(props BumpProps, result BumpResult, dims core.Dimensions, theme *theming.Theme) string {
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

// renderLinesLayer emits one <path> per serie (fill:none stroke:color).
func renderLinesLayer(result BumpResult) string {
	var b strings.Builder
	for _, s := range result.Series {
		if s.LinePath == "" {
			continue
		}
		b.WriteString(`<path d="`)
		b.WriteString(s.LinePath)
		b.WriteString(`" fill="none" stroke="`)
		b.WriteString(s.Color)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(s.LineWidth))
		b.WriteString(`" stroke-linecap="round" style="opacity:`)
		b.WriteString(fmtF(s.Opacity))
		b.WriteString(`"></path>`)
	}
	return b.String()
}

// renderPointsLayer emits a DotsItem per defined rank point.
func renderPointsLayer(props BumpProps, result BumpResult) string {
	var b strings.Builder
	for _, s := range result.Series {
		for _, p := range s.Points {
			if !p.Defined {
				continue
			}
			dp := core.DotsItemProps{
				X: p.X, Y: p.Y, Size: p.Size, Color: "#ffffff",
				BorderWidth: 2, BorderColor: p.Color,
			}
			if props.Interactive {
				dp.Tooltip = interact.TooltipHTML(p.Color, p.SerieID, p.FormattedY)
			}
			b.WriteString(renderComponent(core.DotsItem(dp)))
		}
	}
	return b.String()
}

// renderLabelsLayer emits serie-id labels at the first (start) and/or last
// (end) defined point of each serie.
func renderLabelsLayer(props BumpProps, result BumpResult, theme *theming.Theme) string {
	if !props.StartLabelEnabled() && !props.EndLabelEnabled() {
		return ""
	}
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	var b strings.Builder
	for _, s := range result.Series {
		first, last := firstLastDefined(s.Points)
		if props.StartLabelEnabled() && first != nil {
			b.WriteString(labelText(first.X-props.PointSize-6, first.Y, "end", s.ID, s.Color, fontSize, fontFamily, fill))
		}
		if props.EndLabelEnabled() && last != nil {
			b.WriteString(labelText(last.X+props.PointSize+6, last.Y, "start", s.ID, s.Color, fontSize, fontFamily, fill))
		}
	}
	return b.String()
}

func firstLastDefined(pts []BumpPoint) (first, last *BumpPoint) {
	for i := range pts {
		if pts[i].Defined {
			if first == nil {
				first = &pts[i]
			}
			last = &pts[i]
		}
	}
	return first, last
}

func labelText(x, y float64, anchor, text, color string, fontSize float64, fontFamily, fallbackFill string) string {
	fill := color
	if fill == "" {
		fill = fallbackFill
	}
	var b strings.Builder
	b.WriteString(`<text x="`)
	b.WriteString(fmtF(x))
	b.WriteString(`" y="`)
	b.WriteString(fmtF(y))
	b.WriteString(`" text-anchor="`)
	b.WriteString(anchor)
	b.WriteString(`" dominant-baseline="central" style="fill:`)
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
	b.WriteString(templ.EscapeString(text))
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

func renderLegendsLayer(props BumpProps, result BumpResult, dims core.Dimensions) string {
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

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer),
// so coordinates read cleanly and stay free of float-precision noise.
func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
