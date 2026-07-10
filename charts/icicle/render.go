package icicle

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

func applyDefaults(p IcicleProps) IcicleProps {
	if p.Orientation == "" {
		p.Orientation = Defaults.Orientation
	}
	if p.GapX == 0 {
		p.GapX = Defaults.GapX
	}
	if p.GapY == 0 {
		p.GapY = Defaults.GapY
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.EnableLabels == nil {
		p.EnableLabels = Defaults.EnableLabels
	}
	if p.Label == "" {
		p.Label = Defaults.Label
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
// EnableZooming plus a ChartID (htmx mode). Static/standalone renders stay
// byte-identical.
func zoomEnabled(props IcicleProps) bool {
	return props.ZoomingEnabled() && props.ChartID != ""
}

// renderRects draws the icicle rects and optional labels.
func renderRects(props IcicleProps, result IcicleResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	zoom := zoomEnabled(props)

	var rects, labels strings.Builder
	for i, r := range result.Rects {
		if r.Width <= 0 || r.Height <= 0 {
			continue
		}
		rects.WriteString(`<rect x="`)
		rects.WriteString(fmtF(r.X))
		rects.WriteString(`" y="`)
		rects.WriteString(fmtF(r.Y))
		rects.WriteString(`" width="`)
		rects.WriteString(fmtF(r.Width))
		rects.WriteString(`" height="`)
		rects.WriteString(fmtF(r.Height))
		rects.WriteString(`" fill="`)
		rects.WriteString(r.Color)
		rects.WriteString(`"`)
		if props.BorderRadius > 0 {
			rects.WriteString(` rx="`)
			rects.WriteString(fmtF(props.BorderRadius))
			rects.WriteString(`"`)
		}
		if props.BorderWidth > 0 {
			stroke := getBorderColor(map[string]any{"color": r.Color})
			if stroke == "" {
				stroke = r.Color
			}
			rects.WriteString(` stroke="`)
			rects.WriteString(stroke)
			rects.WriteString(`" stroke-width="`)
			rects.WriteString(fmtF(props.BorderWidth))
			rects.WriteString(`"`)
		}
		if props.Interactive {
			rects.WriteString(` `)
			rects.WriteString(interact.TooltipAttrName)
			rects.WriteString(`="`)
			rects.WriteString(templ.EscapeString(interact.TooltipHTML(r.Color, r.ID, r.FormattedValue)))
			rects.WriteString(`" style="pointer-events:auto"`)
		}
		if zoom {
			writeZoomAttrs(&rects, props.ChartID, r.ID)
		}
		if props.Animate {
			rects.WriteString(`>`)
			rects.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
			rects.WriteString(`</rect>`)
		} else {
			rects.WriteString(`></rect>`)
		}

		if props.LabelsEnabled() {
			labels.WriteString(text(r.X+r.Width/2, r.Y+r.Height/2, labelText(props.Label, r), fill, fontSize, fontFamily))
		}
	}
	out := rects.String() + labels.String()
	if zoom && len(result.Breadcrumb) > 0 {
		out += renderBreadcrumb(props.ChartID, result.Breadcrumb, fill, fontSize, fontFamily)
	}
	return out
}

// writeZoomAttrs emits the htmx click-to-zoom wiring for a node: a GET to the
// chart's zoom verb that re-renders the full SVG into the chart container.
// Mirrors bar's activation/toggle target/swap idiom (#chart-<id> swap root).
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
// segments, each zooming to that ancestor. Rendered top-left, themed like the
// label layer. Only emitted when focused.
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
			b.WriteString(`px">`)
			b.WriteString(` / `)
			b.WriteString(`</text>`)
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

func labelText(path string, r ComputedRect) string {
	switch path {
	case "value":
		return strconv.FormatFloat(r.Value, 'g', -1, 64)
	case "formattedValue":
		return r.FormattedValue
	default:
		return r.ID
	}
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
