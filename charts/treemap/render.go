package treemap

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

// applyDefaults fills zero-valued TreemapProps fields from Defaults.
func applyDefaults(p TreemapProps) TreemapProps {
	if p.Tile == "" {
		p.Tile = Defaults.Tile
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.NodeOpacity == 0 {
		p.NodeOpacity = Defaults.NodeOpacity
	}
	if p.EnableLabel == nil {
		p.EnableLabel = Defaults.EnableLabel
	}
	if p.EnableParentLabel == nil {
		p.EnableParentLabel = Defaults.EnableParentLabel
	}
	if p.ParentLabelSize == 0 {
		p.ParentLabelSize = Defaults.ParentLabelSize
	}
	if p.ParentLabelPadding == 0 {
		p.ParentLabelPadding = Defaults.ParentLabelPadding
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
func zoomEnabled(props TreemapProps) bool {
	return props.ZoomingEnabled() && props.ChartID != ""
}

// renderNodes draws the node rects (parents before children so children sit on
// top), then leaf/parent labels.
func renderNodes(props TreemapProps, result TreemapResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	zoom := zoomEnabled(props)

	var rects, labels strings.Builder
	for i, n := range result.Nodes {
		if n.Width <= 0 || n.Height <= 0 {
			continue
		}
		rects.WriteString(`<rect x="`)
		rects.WriteString(fmtF(n.X))
		rects.WriteString(`" y="`)
		rects.WriteString(fmtF(n.Y))
		rects.WriteString(`" width="`)
		rects.WriteString(fmtF(n.Width))
		rects.WriteString(`" height="`)
		rects.WriteString(fmtF(n.Height))
		rects.WriteString(`" fill="`)
		rects.WriteString(n.Color)
		rects.WriteString(`" fill-opacity="`)
		rects.WriteString(fmtF(props.NodeOpacity))
		rects.WriteString(`"`)
		if props.BorderWidth > 0 {
			rects.WriteString(` stroke="`)
			stroke := getBorderColor(map[string]any{"color": n.Color})
			if stroke == "" {
				stroke = n.Color
			}
			rects.WriteString(stroke)
			rects.WriteString(`" stroke-width="`)
			rects.WriteString(fmtF(props.BorderWidth))
			rects.WriteString(`"`)
		}
		if props.Interactive {
			rects.WriteString(` `)
			rects.WriteString(interact.TooltipAttrName)
			rects.WriteString(`="`)
			rects.WriteString(templ.EscapeString(interact.TooltipHTML(n.Color, n.ID, n.FormattedValue)))
			rects.WriteString(`" style="pointer-events:auto"`)
		}
		if zoom {
			writeZoomAttrs(&rects, props.ChartID, n.ID)
		}
		if props.Animate {
			rects.WriteString(`>`)
			rects.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
			rects.WriteString(`</rect>`)
		} else {
			rects.WriteString(`></rect>`)
		}

		// Labels.
		if n.IsParent {
			if props.ParentLabelEnabled() && n.Height >= props.ParentLabelSize {
				labels.WriteString(text(n.X+props.ParentLabelPadding, n.Y+props.ParentLabelSize/2, "start",
					n.ID, fill, fontSize, fontFamily))
			}
		} else if props.LabelEnabled() && math.Min(n.Width, n.Height) >= props.LabelSkipSize {
			labels.WriteString(text(n.X+n.Width/2, n.Y+n.Height/2, "middle",
				n.FormattedValue, fill, fontSize, fontFamily))
		}
	}
	out := rects.String() + labels.String()
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

func text(x, y float64, anchor, s, fill string, fontSize float64, fontFamily string) string {
	var b strings.Builder
	b.WriteString(`<text x="`)
	b.WriteString(fmtF(x))
	b.WriteString(`" y="`)
	b.WriteString(fmtF(y))
	b.WriteString(`" text-anchor="`)
	b.WriteString(anchor)
	b.WriteString(`" dominant-baseline="central" style="pointer-events:none;fill:`)
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
