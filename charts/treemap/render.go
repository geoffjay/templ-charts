package treemap

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
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

// renderNodes draws the node rects (parents before children so children sit on
// top), then leaf/parent labels.
func renderNodes(props TreemapProps, result TreemapResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)

	var rects, labels strings.Builder
	for _, n := range result.Nodes {
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
		rects.WriteString(`></rect>`)

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
	return rects.String() + labels.String()
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
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
