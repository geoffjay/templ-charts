package polarbar

import (
	"context"
	"math"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
	polaraxes "github.com/geoffjay/templ-charts/charts/polar-axes"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued PolarBarProps fields from Defaults.
func applyDefaults(p PolarBarProps) PolarBarProps {
	if len(p.Keys) == 0 {
		p.Keys = Defaults.Keys
	}
	if p.IndexBy == "" {
		p.IndexBy = Defaults.IndexBy
	}
	if p.StartAngle == 0 && p.EndAngle == 0 {
		p.StartAngle = Defaults.StartAngle
		p.EndAngle = Defaults.EndAngle
	}
	if p.EnableRadialGrid == nil {
		p.EnableRadialGrid = Defaults.EnableRadialGrid
	}
	if p.EnableCircularGrid == nil {
		p.EnableCircularGrid = Defaults.EnableCircularGrid
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.EnableArcLabels == nil {
		p.EnableArcLabels = Defaults.EnableArcLabels
	}
	if p.ArcLabel == "" {
		p.ArcLabel = Defaults.ArcLabel
	}
	if p.ArcLabelsRadiusOffset == 0 {
		p.ArcLabelsRadiusOffset = Defaults.ArcLabelsRadiusOffset
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

func renderLayers(props PolarBarProps, result PolarBarResult, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case PolarBarLayerGrid:
			b.WriteString(renderGridLayer(props, result, theme))
		case PolarBarLayerArcs:
			b.WriteString(renderArcsLayer(props, result, theme))
		case PolarBarLayerAxes:
			b.WriteString(renderAxesLayer(result, theme))
		case PolarBarLayerLabels:
			if props.ArcLabelsEnabled() {
				b.WriteString(renderLabelsLayer(props, result, theme))
			}
		case PolarBarLayerLegends:
			b.WriteString(renderLegendsLayer(props, result))
		}
	}
	return b.String()
}

// renderGridLayer renders the polar grid (radial rays per index angle +
// concentric rings at radius ticks) via charts/polar-axes.
func renderGridLayer(props PolarBarProps, result PolarBarResult, theme *theming.Theme) string {
	return polaraxes.RenderPolarGrid(polaraxes.PolarGridProps{
		Center:             result.Center,
		EnableRadialGrid:   props.RadialGridEnabled(),
		AngleScale:         result.AngleScale,
		StartAngle:         result.StartAngle,
		EndAngle:           result.EndAngle,
		EnableCircularGrid: props.CircularGridEnabled(),
		RadiusScale:        result.RadiusScale,
		InnerRadius:        result.InnerRadius,
		OuterRadius:        result.OuterRadius,
		Theme:              theme,
	})
}

// renderAxesLayer renders the outer circular axis (index labels) and a radial
// axis (radius/value ticks) via charts/polar-axes.
func renderAxesLayer(result PolarBarResult, theme *theming.Theme) string {
	var b strings.Builder
	b.WriteString(polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
		Type:       polaraxes.CircularAxisOuter,
		Center:     result.Center,
		Radius:     result.OuterRadius,
		StartAngle: result.StartAngle,
		EndAngle:   result.EndAngle,
		Scale:      result.AngleScale,
		Theme:      theme,
	}))
	b.WriteString(polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
		Center:        result.Center,
		Angle:         math.Min(result.StartAngle, result.EndAngle),
		Scale:         result.RadiusScale,
		TicksPosition: polaraxes.TicksBefore,
		Theme:         theme,
	}))
	return b.String()
}

// renderArcsLayer renders the stacked bar arcs via arcs.ArcsLayer.
func renderArcsLayer(props PolarBarProps, result PolarBarResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	items := make([]arcs.ArcLayerItem, 0, len(result.Arcs))
	for i, a := range result.Arcs {
		sp := arcs.ArcShapeProps{
			Path: result.ArcGenerator.GenerateSvgArc(a.Arc),
			Fill: a.Color,
		}
		if props.BorderWidth > 0 {
			sp.Stroke = getBorderColor(map[string]any{"color": a.Color})
			sp.StrokeWidth = props.BorderWidth
		}
		if props.Interactive {
			sp.DataTooltip = interact.TooltipHTML(a.Color, a.Index+" - "+a.Key, a.FormattedValue)
		}
		if props.Animate {
			sp.Animate = true
			sp.AnimateBegin = core.StaggerBegin(i, props.MotionStagger)
		}
		items = append(items, arcs.ArcLayerItem{Arc: a.Arc, Props: sp})
	}
	return renderComponent(arcs.ArcsLayer(arcs.ArcsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
}

// renderLabelsLayer renders per-arc labels via arcs.ArcLabelsLayer, skipping
// arcs whose angular span is below ArcLabelsSkipAngle.
func renderLabelsLayer(props PolarBarProps, result PolarBarResult, theme *theming.Theme) string {
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	items := make([]arcs.ArcLabelItem, 0, len(result.Arcs))
	for _, a := range result.Arcs {
		spanDeg := arcs.RadToDeg(math.Abs(a.Arc.EndAngle - a.Arc.StartAngle))
		if props.ArcLabelsSkipAngle > 0 && spanDeg < props.ArcLabelsSkipAngle {
			continue
		}
		x, y := arcs.ComputeArcCenter(a.Arc, props.ArcLabelsRadiusOffset)
		items = append(items, arcs.ArcLabelItem{Props: arcs.ArcLabelProps{
			X:          x,
			Y:          y,
			Label:      resolveLabel(props.ArcLabel, a),
			Fill:       fill,
			FontSize:   fontSize,
			FontFamily: fontFamily,
		}})
	}
	return renderComponent(arcs.ArcLabelsLayer(arcs.ArcLabelsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
}

func renderLegendsLayer(props PolarBarProps, result PolarBarResult) string {
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
			ChartWidth:  props.Width,
			ChartHeight: props.Height,
		})))
	}
	return b.String()
}

// resolveLabel maps a label path to an arc's value. Supports "formattedValue",
// "value", "id", "key"; anything else falls back to the formatted value.
func resolveLabel(path string, a ComputedArc) string {
	switch path {
	case "value":
		return valueFormatter("")(a.Value)
	case "id", "key":
		return a.Key
	case "index":
		return a.Index
	default:
		return a.FormattedValue
	}
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
