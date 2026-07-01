package sunburst

import (
	"context"
	"math"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

func applyDefaults(p SunburstProps) SunburstProps {
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.BorderColor == "" {
		p.BorderColor = Defaults.BorderColor
	}
	if p.EnableArcLabels == nil {
		p.EnableArcLabels = Defaults.EnableArcLabels
	}
	if p.ArcLabelsRadiusOffset == 0 {
		p.ArcLabelsRadiusOffset = Defaults.ArcLabelsRadiusOffset
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

func isZeroOrdinal(c colors.OrdinalColorScaleConfig) bool {
	return c.Type == 0 && c.Scheme == "" && c.Static == "" && len(c.Colors) == 0 && c.Func == nil && c.DatumPath == ""
}

// renderArcs draws the sunburst arcs (and optional labels) via charts/arcs.
func renderArcs(props SunburstProps, result SunburstResult, theme *theming.Theme) string {
	arcGen := arcs.CreateArcGenerator(props.CornerRadius, 0)
	items := make([]arcs.ArcLayerItem, 0, len(result.Arcs))
	for _, a := range result.Arcs {
		sp := arcs.ArcShapeProps{
			Path: arcGen.GenerateSvgArc(a.Arc),
			Fill: a.Color,
		}
		if props.BorderWidth > 0 {
			sp.Stroke = props.BorderColor
			sp.StrokeWidth = props.BorderWidth
		}
		if props.Interactive {
			sp.DataTooltip = interact.TooltipHTML(a.Color, a.ID, a.FormattedValue)
		}
		items = append(items, arcs.ArcLayerItem{Arc: a.Arc, Props: sp})
	}
	out := renderComponent(arcs.ArcsLayer(arcs.ArcsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
	if props.ArcLabelsEnabled() {
		out += renderArcLabels(props, result, theme)
	}
	return out
}

func renderArcLabels(props SunburstProps, result SunburstResult, theme *theming.Theme) string {
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	items := make([]arcs.ArcLabelItem, 0, len(result.Arcs))
	for _, a := range result.Arcs {
		spanDeg := arcs.RadToDeg(math.Abs(a.Arc.EndAngle - a.Arc.StartAngle))
		if props.ArcLabelsSkipAngle > 0 && spanDeg < props.ArcLabelsSkipAngle {
			continue
		}
		x, y := arcs.ComputeArcCenter(a.Arc, props.ArcLabelsRadiusOffset)
		items = append(items, arcs.ArcLabelItem{Props: arcs.ArcLabelProps{
			X: x, Y: y, Label: a.FormattedValue, Fill: fill, FontSize: fontSize, FontFamily: fontFamily,
		}})
	}
	return renderComponent(arcs.ArcLabelsLayer(arcs.ArcLabelsLayerProps{
		Items:   items,
		CenterX: result.Center[0],
		CenterY: result.Center[1],
	}))
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
