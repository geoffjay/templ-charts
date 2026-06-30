package radialbar

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	polaraxes "github.com/geoffjay/templ-charts/charts/polar-axes"
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

// applyDefaults fills zero-valued RadialBarProps fields from Defaults.
func applyDefaults(p RadialBarProps) RadialBarProps {
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.StartAngle == 0 && p.EndAngle == 0 {
		p.StartAngle = Defaults.StartAngle
		p.EndAngle = Defaults.EndAngle
	}
	if p.InnerRadius == 0 {
		p.InnerRadius = Defaults.InnerRadius
	}
	if p.Padding == 0 {
		p.Padding = Defaults.Padding
	}
	if p.EnableTracks == nil {
		p.EnableTracks = Defaults.EnableTracks
	}
	if p.TracksColor == "" {
		p.TracksColor = Defaults.TracksColor
	}
	if p.EnableRadialGrid == nil {
		p.EnableRadialGrid = Defaults.EnableRadialGrid
	}
	if p.EnableCircularGrid == nil {
		p.EnableCircularGrid = Defaults.EnableCircularGrid
	}
	if p.ShowRadialAxisStart == nil {
		p.ShowRadialAxisStart = Defaults.ShowRadialAxisStart
	}
	if p.ShowRadialAxisEnd == nil {
		p.ShowRadialAxisEnd = Defaults.ShowRadialAxisEnd
	}
	if p.ShowCircularAxisInner == nil {
		p.ShowCircularAxisInner = Defaults.ShowCircularAxisInner
	}
	if p.ShowCircularAxisOuter == nil {
		p.ShowCircularAxisOuter = Defaults.ShowCircularAxisOuter
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if p.EnableLabels == nil {
		p.EnableLabels = Defaults.EnableLabels
	}
	if p.Label == "" {
		p.Label = Defaults.Label
	}
	if p.LabelsSkipAngle == 0 {
		p.LabelsSkipAngle = Defaults.LabelsSkipAngle
	}
	if p.LabelsRadiusOffset == 0 {
		p.LabelsRadiusOffset = Defaults.LabelsRadiusOffset
	}
	if isZeroColor(p.LabelsTextColor) {
		p.LabelsTextColor = Defaults.LabelsTextColor
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

func isZeroColor(c colors.InheritedColorConfig) bool {
	return c.Type == 0 && c.Static == "" && c.ThemePath == "" && c.FromPath == "" && c.Func == nil
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

// renderLayers renders the enabled layers as an inner SVG string.
func renderLayers(props RadialBarProps, result RadialBarResult, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case RadialBarLayerGrid:
			b.WriteString(renderGridLayer(props, result, theme))
		case RadialBarLayerTracks:
			if props.TracksEnabled() {
				b.WriteString(renderTracksLayer(result))
			}
		case RadialBarLayerBars:
			b.WriteString(renderBarsLayer(props, result, theme))
		case RadialBarLayerLabels:
			if props.LabelsEnabled() {
				b.WriteString(renderLabelsLayer(props, result, theme))
			}
		case RadialBarLayerLegends:
			b.WriteString(renderLegendsLayer(props, result))
		}
	}
	return b.String()
}

// renderGridLayer renders the polar grid plus the radial/circular axes, reusing
// charts/polar-axes. Mirrors @nivo/radial-bar's grid layer.
func renderGridLayer(props RadialBarProps, result RadialBarResult, theme *theming.Theme) string {
	var b strings.Builder
	b.WriteString(polaraxes.RenderPolarGrid(polaraxes.PolarGridProps{
		Center:             result.Center,
		EnableRadialGrid:   *props.EnableRadialGrid,
		AngleScale:         result.ValueScale,
		StartAngle:         result.StartAngle,
		EndAngle:           result.EndAngle,
		EnableCircularGrid: *props.EnableCircularGrid,
		RadiusScale:        result.RadiusScale,
		InnerRadius:        result.InnerRadius,
		OuterRadius:        result.OuterRadius,
		Theme:              theme,
	}))
	if *props.ShowRadialAxisStart {
		b.WriteString(polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
			Center:        result.Center,
			Angle:         math.Min(result.StartAngle, result.EndAngle),
			Scale:         result.RadiusScale,
			TicksPosition: polaraxes.TicksBefore,
			Theme:         theme,
		}))
	}
	if *props.ShowRadialAxisEnd {
		b.WriteString(polaraxes.RenderRadialAxis(polaraxes.RadialAxisProps{
			Center:        result.Center,
			Angle:         math.Max(result.StartAngle, result.EndAngle),
			Scale:         result.RadiusScale,
			TicksPosition: polaraxes.TicksAfter,
			Theme:         theme,
		}))
	}
	if *props.ShowCircularAxisInner {
		b.WriteString(polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
			Type:       polaraxes.CircularAxisInner,
			Center:     result.Center,
			Radius:     result.InnerRadius,
			StartAngle: result.StartAngle,
			EndAngle:   result.EndAngle,
			Scale:      result.ValueScale,
			Theme:      theme,
		}))
	}
	if *props.ShowCircularAxisOuter {
		b.WriteString(polaraxes.RenderCircularAxis(polaraxes.CircularAxisProps{
			Type:       polaraxes.CircularAxisOuter,
			Center:     result.Center,
			Radius:     result.OuterRadius,
			StartAngle: result.StartAngle,
			EndAngle:   result.EndAngle,
			Scale:      result.ValueScale,
			Theme:      theme,
		}))
	}
	return b.String()
}

// renderTracksLayer renders the background track arcs via arcs.ArcsLayer.
func renderTracksLayer(result RadialBarResult) string {
	items := make([]arcs.ArcLayerItem, 0, len(result.Tracks))
	for _, t := range result.Tracks {
		items = append(items, arcs.ArcLayerItem{
			Arc:   t.Arc,
			Props: arcs.ArcShapeProps{Path: result.ArcGenerator.GenerateSvgArc(t.Arc), Fill: t.Color},
		})
	}
	return renderComponent(arcs.ArcsLayer(arcs.ArcsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
}

// renderBarsLayer renders the stacked bar arcs via arcs.ArcsLayer.
func renderBarsLayer(props RadialBarProps, result RadialBarResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	items := make([]arcs.ArcLayerItem, 0, len(result.Bars))
	for _, bar := range result.Bars {
		sp := arcs.ArcShapeProps{
			Path: result.ArcGenerator.GenerateSvgArc(bar.Arc),
			Fill: bar.Color,
		}
		if props.BorderWidth > 0 {
			sp.Stroke = getBorderColor(map[string]any{"color": bar.Color})
			sp.StrokeWidth = props.BorderWidth
		}
		items = append(items, arcs.ArcLayerItem{Arc: bar.Arc, Props: sp})
	}
	return renderComponent(arcs.ArcsLayer(arcs.ArcsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
}

// renderLabelsLayer renders per-arc labels via arcs.ArcLabelsLayer, skipping
// arcs whose angular span is below LabelsSkipAngle.
func renderLabelsLayer(props RadialBarProps, result RadialBarResult, theme *theming.Theme) string {
	getTextColor := colors.GetInheritedColorGenerator(props.LabelsTextColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	items := make([]arcs.ArcLabelItem, 0, len(result.Bars))
	for _, bar := range result.Bars {
		spanDeg := arcs.RadToDeg(math.Abs(bar.Arc.EndAngle - bar.Arc.StartAngle))
		if props.LabelsSkipAngle > 0 && spanDeg < props.LabelsSkipAngle {
			continue
		}
		x, y := arcs.ComputeArcCenter(bar.Arc, props.LabelsRadiusOffset)
		textColor := getTextColor(map[string]any{"color": bar.Color})
		if textColor == "" {
			textColor = fill
		}
		items = append(items, arcs.ArcLabelItem{Props: arcs.ArcLabelProps{
			X:          x,
			Y:          y,
			Label:      resolveLabel(props.Label, bar),
			Fill:       textColor,
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

// renderLegendsLayer renders any box legends.
func renderLegendsLayer(props RadialBarProps, result RadialBarResult) string {
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

// resolveLabel maps a label path to a bar's value. Supports the common nivo
// paths ("formattedValue", "value", "id"); anything else falls back to the
// formatted value.
func resolveLabel(path string, bar ComputedBar) string {
	switch path {
	case "value":
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", bar.Value), "0"), ".")
	case "id":
		return bar.ID
	default:
		return bar.FormattedValue
	}
}

// labelsTextStyle resolves the labels text style from the theme (labels block,
// falling back to the base text style).
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

// guard against unused scales import if the build trims helpers.
var _ = scales.TicksSpec{}
