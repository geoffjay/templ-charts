package sunburst

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
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

// zoomEnabled reports whether click-to-zoom wiring should be emitted: gated on
// EnableZooming plus a ChartID (htmx mode). Static renders stay byte-identical.
func zoomEnabled(props SunburstProps) bool {
	return props.EnableZooming && props.ChartID != ""
}

// renderArcs draws the sunburst arcs (and optional labels) via charts/arcs.
func renderArcs(props SunburstProps, result SunburstResult, theme *theming.Theme) string {
	arcGen := arcs.CreateArcGenerator(props.CornerRadius, 0)
	zoom := zoomEnabled(props)
	items := make([]arcs.ArcLayerItem, 0, len(result.Arcs))
	for i, a := range result.Arcs {
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
		if zoom {
			sp.HxGet = "/charts/" + props.ChartID + "/zoom?node=" + templ.EscapeString(a.ID)
			sp.HxTarget = "#chart-" + props.ChartID
			sp.HxSwap = "innerHTML"
		}
		if props.Animate {
			sp.Animate = true
			sp.AnimateBegin = core.StaggerBegin(i, props.MotionStagger)
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
	if zoom && len(result.Breadcrumb) > 0 {
		fill, fontSize, fontFamily := labelsTextStyle(theme)
		out += renderBreadcrumb(props.ChartID, result.Breadcrumb, fill, fontSize, fontFamily)
	}
	return out
}

// renderBreadcrumb draws the root→focus ancestor path as clickable text
// segments (top-left), each zooming to that ancestor. Only emitted when
// focused.
func renderBreadcrumb(chartID string, crumbs []Crumb, fill string, fontSize float64, fontFamily string) string {
	if fontSize <= 0 {
		fontSize = 11
	}
	var b strings.Builder
	x := 4.0
	y := fontSize + 2
	for i, c := range crumbs {
		if i > 0 {
			b.WriteString(`<text x="`)
			b.WriteString(fmtF(x))
			b.WriteString(`" y="`)
			b.WriteString(fmtF(y))
			b.WriteString(`" dominant-baseline="central" style="pointer-events:none;fill:`)
			b.WriteString(fill)
			b.WriteString(`;font-size:`)
			b.WriteString(fmtF(fontSize))
			b.WriteString(`px"> / </text>`)
			x += float64(len(" / ")) * fontSize * 0.5
		}
		b.WriteString(`<text x="`)
		b.WriteString(fmtF(x))
		b.WriteString(`" y="`)
		b.WriteString(fmtF(y))
		b.WriteString(`" dominant-baseline="central" hx-get="/charts/`)
		b.WriteString(chartID)
		b.WriteString(`/zoom?node=`)
		b.WriteString(templ.EscapeString(c.ID))
		b.WriteString(`" hx-target="#chart-`)
		b.WriteString(chartID)
		b.WriteString(`" hx-swap="innerHTML" style="cursor:pointer;pointer-events:auto;fill:`)
		b.WriteString(fill)
		b.WriteString(`;font-size:`)
		b.WriteString(fmtF(fontSize))
		if fontFamily != "" {
			b.WriteString(`;font-family:`)
			b.WriteString(fontFamily)
		}
		b.WriteString(`;text-decoration:underline">`)
		b.WriteString(templ.EscapeString(c.Label))
		b.WriteString(`</text>`)
		x += float64(len(c.Label)) * fontSize * 0.55
	}
	return b.String()
}

func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
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
