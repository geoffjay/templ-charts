package icicle

import (
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
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

// renderRects draws the icicle rects and optional labels.
func renderRects(props IcicleProps, result IcicleResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)

	var rects, labels strings.Builder
	for _, r := range result.Rects {
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
		rects.WriteString(`></rect>`)

		if props.LabelsEnabled() {
			labels.WriteString(text(r.X+r.Width/2, r.Y+r.Height/2, labelText(props.Label, r), fill, fontSize, fontFamily))
		}
	}
	return rects.String() + labels.String()
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
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
