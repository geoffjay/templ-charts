package parallelcoordinates

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
	"github.com/geoffjay/templ-charts/charts/theming"
)

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued PCProps fields from Defaults.
func applyDefaults(p PCProps) PCProps {
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Curve == "" {
		p.Curve = Defaults.Curve
	}
	if p.LineWidth == 0 {
		p.LineWidth = Defaults.LineWidth
	}
	if p.LineOpacity == 0 {
		p.LineOpacity = Defaults.LineOpacity
	}
	if p.AxesTicksPosition == "" {
		p.AxesTicksPosition = Defaults.AxesTicksPosition
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
	return p
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

func renderLayers(props PCProps, result PCResult, dims core.Dimensions, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case PCLayerLines:
			b.WriteString(renderLinesLayer(props, result))
		case PCLayerAxes:
			b.WriteString(renderAxesLayer(props, result, dims, theme))
		case PCLayerMesh:
			b.WriteString(renderMeshLayer(props, result, dims))
		case PCLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		}
	}
	return b.String()
}

// renderLinesLayer emits one <path> per datum polyline.
func renderLinesLayer(props PCProps, result PCResult) string {
	var b strings.Builder
	for i, ln := range result.Lines {
		if ln.Path == "" {
			continue
		}
		b.WriteString(`<path d="`)
		b.WriteString(ln.Path)
		b.WriteString(`" fill="none" stroke="`)
		b.WriteString(ln.Color)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(props.LineWidth))
		b.WriteString(`" stroke-linecap="round" style="opacity:`)
		b.WriteString(fmtF(props.LineOpacity))
		b.WriteString(`"`)
		// With UseMesh, the mesh overlay handles hover; skip per-line tooltips.
		if props.Interactive && !props.UseMesh && ln.Tooltip != "" {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(ln.Tooltip))
			b.WriteString(`" style="pointer-events:auto;opacity:`)
			b.WriteString(fmtF(props.LineOpacity))
			b.WriteString(`"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</path>`)
	}
	return b.String()
}

// renderMeshLayer emits the accurate voronoi-mesh hover overlay over the
// per-axis line vertices (charts/interact, backed by internal/d3/delaunay):
// each datum contributes one mesh point per variable axis it crosses, so
// hovering anywhere resolves to the nearest datum. Active only when
// Interactive && UseMesh; DebugMesh draws the cells; DetectionRadius bounds
// hit-testing when > 0.
func renderMeshLayer(props PCProps, result PCResult, dims core.Dimensions) string {
	if !props.Interactive || !props.UseMesh {
		return ""
	}
	var pts []interact.MeshPoint
	for _, ln := range result.Lines {
		if ln.Tooltip == "" {
			continue
		}
		for _, p := range ln.Points {
			pts = append(pts, interact.MeshPoint{X: p[0], Y: p[1], HTML: ln.Tooltip})
		}
	}
	return interact.MeshOverlay(pts, dims.InnerWidth, dims.InnerHeight, props.DebugMesh, props.DetectionRadius)
}

// renderAxesLayer emits one axes.Axis per variable, oriented per the layout.
func renderAxesLayer(props PCProps, result PCResult, dims core.Dimensions, theme *theming.Theme) string {
	horizontal := props.Layout != PCLayoutVertical
	var b strings.Builder
	for _, v := range result.Variables {
		ap := axes.AxisProps{
			Scale:          v.Scale,
			TicksPosition:  props.AxesTicksPosition,
			TickSize:       5,
			TickPadding:    5,
			Legend:         v.Label,
			LegendPosition: axes.AxisLegendMiddle,
		}
		if horizontal {
			ap.Axis = "y"
			ap.X = v.Pos
			ap.Y = 0
			ap.Length = dims.InnerHeight
			ap.LegendOffset = -20
		} else {
			ap.Axis = "x"
			ap.X = 0
			ap.Y = v.Pos
			ap.Length = dims.InnerWidth
			ap.LegendOffset = -30
		}
		b.WriteString(renderComponent(axes.Axis(ap, theme)))
	}
	return b.String()
}

func renderLegendsLayer(props PCProps, result PCResult, dims core.Dimensions) string {
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
