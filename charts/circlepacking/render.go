package circlepacking

import (
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func applyDefaults(p CirclePackingProps) CirclePackingProps {
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.EnableLabels == nil {
		p.EnableLabels = Defaults.EnableLabels
	}
	if p.LabelsSkipRadius == 0 {
		p.LabelsSkipRadius = Defaults.LabelsSkipRadius
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
func zoomEnabled(props CirclePackingProps) bool {
	return props.EnableZooming && props.ChartID != ""
}

// renderCircles draws the packed circles (parents behind children) and labels.
func renderCircles(props CirclePackingProps, result CirclePackingResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	zoom := zoomEnabled(props)

	var circles, labels strings.Builder
	for i, c := range result.Circles {
		if c.R <= 0 {
			continue
		}
		circles.WriteString(`<circle cx="`)
		circles.WriteString(fmtF(c.X))
		circles.WriteString(`" cy="`)
		circles.WriteString(fmtF(c.Y))
		circles.WriteString(`" r="`)
		circles.WriteString(fmtF(c.R))
		circles.WriteString(`" fill="`)
		circles.WriteString(c.Color)
		circles.WriteString(`"`)
		if props.BorderWidth > 0 {
			stroke := getBorderColor(map[string]any{"color": c.Color})
			if stroke == "" {
				stroke = c.Color
			}
			circles.WriteString(` stroke="`)
			circles.WriteString(stroke)
			circles.WriteString(`" stroke-width="`)
			circles.WriteString(fmtF(props.BorderWidth))
			circles.WriteString(`"`)
		}
		if props.Interactive {
			circles.WriteString(` `)
			circles.WriteString(interact.TooltipAttrName)
			circles.WriteString(`="`)
			circles.WriteString(templ.EscapeString(interact.TooltipHTML(c.Color, c.ID, c.FormattedValue)))
			circles.WriteString(`" style="pointer-events:auto"`)
		}
		if zoom {
			writeZoomAttrs(&circles, props.ChartID, c.ID)
		}
		circles.WriteString(`>`)
		if props.Animate {
			circles.WriteString(core.SMILAnimate("r", "0", fmtF(c.R), core.StaggerBegin(i, props.MotionStagger)))
		}
		circles.WriteString(`</circle>`)

		if props.LabelsEnabled() && c.IsLeaf && c.R >= props.LabelsSkipRadius {
			labels.WriteString(text(c.X, c.Y, c.ID, fill, fontSize, fontFamily))
		}
	}
	out := circles.String() + labels.String()
	if zoom && len(result.Breadcrumb) > 0 {
		out += renderBreadcrumb(props.ChartID, result.Breadcrumb, fill, fontSize, fontFamily)
	}
	return out
}

// writeZoomAttrs emits the htmx click-to-zoom wiring for a node: a GET to the
// chart's zoom verb that re-renders the full SVG into the chart container
// (#chart-<id> swap root), mirroring bar's activation/toggle idiom.
func writeZoomAttrs(b *strings.Builder, chartID, nodeID string) {
	b.WriteString(` hx-get="/charts/`)
	b.WriteString(chartID)
	b.WriteString(`/zoom?node=`)
	b.WriteString(templ.EscapeString(nodeID))
	b.WriteString(`" hx-target="#chart-`)
	b.WriteString(chartID)
	b.WriteString(`" hx-swap="innerHTML" style="cursor:pointer;pointer-events:auto"`)
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
		b.WriteString(`" dominant-baseline="central"`)
		writeZoomAttrs(&b, chartID, c.ID)
		b.WriteString(`;fill:`)
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

func text(x, y float64, s, fill string, fontSize float64, fontFamily string) string {
	var b strings.Builder
	b.WriteString(`<text x="`)
	b.WriteString(fmtF(x))
	b.WriteString(`" y="`)
	b.WriteString(fmtF(y))
	b.WriteString(`" text-anchor="middle" dominant-baseline="central" style="pointer-events:none;fill:`)
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
	b.WriteString(templ.EscapeString(s))
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

func fmtF(v float64) string {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		r = 0 // normalize -0 to +0 for cross-platform-stable output
	}
	return strconv.FormatFloat(r, 'g', -1, 64)
}
