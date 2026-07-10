package swarmplot

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

// applyDefaults fills zero-valued SwarmPlotProps fields from Defaults.
func applyDefaults(p SwarmPlotProps) SwarmPlotProps {
	if p.ValueScale == nil {
		p.ValueScale = Defaults.ValueScale
	}
	if p.Size == 0 {
		p.Size = Defaults.Size
	}
	if p.Spacing == 0 {
		p.Spacing = Defaults.Spacing
	}
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.ForceStrength == 0 {
		p.ForceStrength = Defaults.ForceStrength
	}
	if p.SimulationIterations == 0 {
		p.SimulationIterations = Defaults.SimulationIterations
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.ColorBy == "" {
		p.ColorBy = Defaults.ColorBy
	}
	if p.BorderColor == "" {
		p.BorderColor = Defaults.BorderColor
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

// renderLayers renders the enabled layers as an inner SVG string. Annotations
// are deferred (no swarmplot annotation specs).
func renderLayers(props SwarmPlotProps, result SwarmPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case SwarmPlotLayerGrid:
			b.WriteString(renderGridLayer(props, result, dims, theme))
		case SwarmPlotLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case SwarmPlotLayerCircles:
			b.WriteString(renderCirclesLayer(props, result))
		case SwarmPlotLayerMesh:
			b.WriteString(renderMeshLayer(props, result, dims))
		case SwarmPlotLayerAnnotations:
			// annotations: deferred.
		}
	}
	return b.String()
}

func renderGridLayer(props SwarmPlotProps, result SwarmPlotResult, dims core.Dimensions, theme *theming.Theme) string {
	var s strings.Builder
	if props.GridXEnabled() {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "x", Scale: result.XScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridXValues, Theme: theme,
		})))
	}
	if props.GridYEnabled() {
		s.WriteString(renderComponent(axes.Grid(axes.GridProps{
			Axis: "y", Scale: result.YScale,
			Width: dims.InnerWidth, Height: dims.InnerHeight,
			TickValues: props.GridYValues, Theme: theme,
		})))
	}
	return s.String()
}

func renderAxesLayer(props SwarmPlotProps, result SwarmPlotResult, dims core.Dimensions, theme *theming.Theme) string {
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

// resolveAxis builds an axes.AxisProps positioned at (originX, originY),
// substituting default tick size/padding when unset.
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

func renderCirclesLayer(props SwarmPlotProps, result SwarmPlotResult) string {
	var b strings.Builder
	drawBorder := props.BorderWidth > 0 && !isTransparent(props.BorderColor)
	for i, n := range result.Nodes {
		b.WriteString(`<circle cx="`)
		b.WriteString(fmtF(n.X))
		b.WriteString(`" cy="`)
		b.WriteString(fmtF(n.Y))
		b.WriteString(`" r="`)
		b.WriteString(fmtF(n.Size / 2))
		b.WriteString(`" fill="`)
		b.WriteString(n.Color)
		b.WriteString(`"`)
		if drawBorder {
			b.WriteString(` stroke="`)
			b.WriteString(props.BorderColor)
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.BorderWidth))
			b.WriteString(`"`)
		}
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(n.Color, n.ID, n.FormattedValue)))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILAnimate("r", "0", fmtF(n.Size/2), core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</circle>`)
	}
	return b.String()
}

// renderMeshLayer emits the accurate voronoi-mesh hover overlay over the nodes
// (charts/interact, backed by internal/d3/delaunay). Active only when
// Interactive && UseMesh.
func renderMeshLayer(props SwarmPlotProps, result SwarmPlotResult, dims core.Dimensions) string {
	if !props.Interactive || !props.UseMesh {
		return ""
	}
	pts := make([]interact.MeshPoint, 0, len(result.Nodes))
	for _, n := range result.Nodes {
		pts = append(pts, interact.MeshPoint{
			X:    n.X,
			Y:    n.Y,
			HTML: interact.TooltipHTML(n.Color, n.ID, n.FormattedValue),
		})
	}
	return interact.MeshOverlay(pts, dims.InnerWidth, dims.InnerHeight, props.DebugMesh, props.DetectionRadius)
}

func isTransparent(c string) bool {
	return strings.ReplaceAll(c, " ", "") == "rgba(0,0,0,0)" || c == "transparent"
}

func fmtF(v float64) string {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		r = 0 // normalize -0 to +0 for cross-platform-stable output
	}
	return strconv.FormatFloat(r, 'g', -1, 64)
}
