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

// renderCircles draws the packed circles (parents behind children) and labels.
func renderCircles(props CirclePackingProps, result CirclePackingResult, theme *theming.Theme) string {
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	fill, fontSize, fontFamily := labelsTextStyle(theme)

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
		circles.WriteString(`>`)
		if props.Animate {
			circles.WriteString(core.SMILAnimate("r", "0", fmtF(c.R), core.StaggerBegin(i, props.MotionStagger)))
		}
		circles.WriteString(`</circle>`)

		if props.LabelsEnabled() && c.IsLeaf && c.R >= props.LabelsSkipRadius {
			labels.WriteString(text(c.X, c.Y, c.ID, fill, fontSize, fontFamily))
		}
	}
	return circles.String() + labels.String()
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
