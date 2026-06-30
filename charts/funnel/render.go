package funnel

import (
	"fmt"
	"strings"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// applyDefaults fills zero-valued FunnelProps fields from Defaults.
func applyDefaults(p FunnelProps) FunnelProps {
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Direction == "" {
		p.Direction = Defaults.Direction
	}
	if p.Interpolation == "" {
		p.Interpolation = Defaults.Interpolation
	}
	if p.ShapeBlending == 0 {
		p.ShapeBlending = Defaults.ShapeBlending
	}
	if isZeroOrdinal(p.Colors) {
		p.Colors = Defaults.Colors
	}
	if p.FillOpacity == 0 {
		p.FillOpacity = Defaults.FillOpacity
	}
	if p.BorderWidth == 0 {
		p.BorderWidth = Defaults.BorderWidth
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if p.BorderOpacity == 0 {
		p.BorderOpacity = Defaults.BorderOpacity
	}
	if p.EnableLabel == nil {
		p.EnableLabel = Defaults.EnableLabel
	}
	if isZeroColor(p.LabelColor) {
		p.LabelColor = Defaults.LabelColor
	}
	if p.EnableBeforeSeparators == nil {
		p.EnableBeforeSeparators = Defaults.EnableBeforeSeparators
	}
	if p.EnableAfterSeparators == nil {
		p.EnableAfterSeparators = Defaults.EnableAfterSeparators
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

// renderLayers renders the enabled layers as an inner SVG string. Annotations
// are deferred in the v2 static pipeline.
func renderLayers(props FunnelProps, result FunnelResult, theme *theming.Theme) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case FunnelLayerSeparators:
			b.WriteString(renderSeparatorsLayer(props, result, theme))
		case FunnelLayerParts:
			b.WriteString(renderPartsLayer(props, result))
		case FunnelLayerLabels:
			if props.LabelEnabled() {
				b.WriteString(renderLabelsLayer(props, result, theme))
			}
		case FunnelLayerAnnotations:
			// Deferred in v2.
		}
	}
	return b.String()
}

func renderPartsLayer(props FunnelProps, result FunnelResult) string {
	var b strings.Builder
	for _, part := range result.Parts {
		// Filled trapezoid.
		fmt.Fprintf(&b, `<path d="%s" fill="%s" fill-opacity="%s"`, part.AreaPath, part.Color, fmtN(part.FillOpacity))
		if props.Interactive {
			fmt.Fprintf(&b, ` style="cursor:pointer" data-tc-tooltip="%s"`, interact.EscapeAttr(interact.TooltipHTML(part.Color, part.Label, part.FormattedValue)))
		}
		b.WriteString("></path>")
		// Side borders (no top/bottom join).
		if part.BorderWidth > 0 && part.BorderColor != "" {
			for _, d := range []string{part.BorderPathA, part.BorderPathB} {
				if d == "" {
					continue
				}
				fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="%s" stroke-width="%s" stroke-opacity="%s"></path>`,
					d, part.BorderColor, fmtN(part.BorderWidth), fmtN(part.BorderOpacity))
			}
		}
	}
	return b.String()
}

func renderLabelsLayer(props FunnelProps, result FunnelResult, theme *theming.Theme) string {
	_, fontSize, fontFamily := labelsTextStyle(theme)
	var b strings.Builder
	for _, part := range result.Parts {
		fmt.Fprintf(&b, `<text x="%s" y="%s" text-anchor="middle" dominant-baseline="central"`, fmtN(part.X), fmtN(part.Y))
		if part.LabelColor != "" {
			fmt.Fprintf(&b, ` fill="%s"`, part.LabelColor)
		}
		if fontSize != "" {
			fmt.Fprintf(&b, ` font-size="%s"`, fontSize)
		}
		if fontFamily != "" {
			fmt.Fprintf(&b, ` font-family="%s"`, fontFamily)
		}
		b.WriteString(` style="pointer-events:none">`)
		b.WriteString(escapeText(part.Label))
		b.WriteString("</text>")
	}
	return b.String()
}

func renderSeparatorsLayer(props FunnelProps, result FunnelResult, theme *theming.Theme) string {
	stroke, strokeWidth := gridLineStroke(theme)
	if stroke == "" {
		stroke = "#999999"
	}
	if strokeWidth == 0 {
		strokeWidth = 1
	}
	var b strings.Builder
	emit := func(seps []Separator) {
		for _, s := range seps {
			fmt.Fprintf(&b, `<line x1="%s" y1="%s" x2="%s" y2="%s" stroke="%s" stroke-width="%s"></line>`,
				fmtN(s.X0), fmtN(s.Y0), fmtN(s.X1), fmtN(s.Y1), stroke, fmtN(strokeWidth))
		}
	}
	emit(result.BeforeSeparators)
	emit(result.AfterSeparators)
	return b.String()
}

func labelsTextStyle(theme *theming.Theme) (fill, fontSize, fontFamily string) {
	if theme == nil {
		return "", "", ""
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
	if fs != nil {
		fontSize = fmt.Sprintf("%v", fs)
	}
	fontFamily = t.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	return fill, fontSize, fontFamily
}

func gridLineStroke(theme *theming.Theme) (string, float64) {
	if theme == nil {
		return "", 0
	}
	m := theme.Grid.Line.Extra
	var stroke string
	var sw float64
	if m != nil {
		if v, ok := m["stroke"].(string); ok {
			stroke = v
		}
		switch x := m["strokeWidth"].(type) {
		case float64:
			sw = x
		case int:
			sw = float64(x)
		}
	}
	return stroke, sw
}

func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func fmtN(v float64) string {
	if v != v || v > 1e308 || v < -1e308 {
		return "0"
	}
	s := fmt.Sprintf("%.3f", v)
	for len(s) > 1 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 1 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
