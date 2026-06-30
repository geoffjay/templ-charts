package bullet

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
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

// applyDefaults fills zero-valued BulletProps fields from Defaults.
func applyDefaults(p BulletProps) BulletProps {
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.Spacing == 0 {
		p.Spacing = Defaults.Spacing
	}
	if p.AxisPosition == "" {
		p.AxisPosition = Defaults.AxisPosition
	}
	if p.RangeColors == "" {
		p.RangeColors = Defaults.RangeColors
	}
	if p.MeasureColors == "" {
		p.MeasureColors = Defaults.MeasureColors
	}
	if p.MarkerColors == "" {
		p.MarkerColors = Defaults.MarkerColors
	}
	if p.MeasureSize == 0 {
		p.MeasureSize = Defaults.MeasureSize
	}
	if p.MarkerSize == 0 {
		p.MarkerSize = Defaults.MarkerSize
	}
	if p.TitlePosition == "" {
		p.TitlePosition = Defaults.TitlePosition
	}
	if p.TitleAlign == "" {
		p.TitleAlign = Defaults.TitleAlign
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderItems renders every bullet row as an inner SVG string.
func renderItems(props BulletProps, result BulletResult, theme *theming.Theme) string {
	horizontal := props.Layout != BulletLayoutVertical
	fill, fontSize, fontFamily := labelsTextStyle(theme)
	var b strings.Builder
	for _, item := range result.Items {
		fmt.Fprintf(&b, `<g transform="translate(%s,%s)">`, fmtN(item.OffsetX), fmtN(item.OffsetY))
		// Ranges (full band).
		for _, r := range item.Ranges {
			writeRect(&b, r)
		}
		// Measures (centered thin band).
		for _, m := range item.Measures {
			writeRect(&b, m)
		}
		// Axis.
		axisName := "x"
		if !horizontal {
			axisName = "y"
		}
		b.WriteString(renderComponent(axes.Axis(axes.AxisProps{
			Axis:          axisName,
			Scale:         item.Scale,
			X:             item.AxisX,
			Y:             item.AxisY,
			Length:        item.Width,
			TicksPosition: props.AxisPosition,
			TickSize:      5,
			TickPadding:   5,
			TextAlign:     "center",
			TextBaseline:  "middle",
		}, theme)))
		// Markers (lines).
		for _, mk := range item.Markers {
			fmt.Fprintf(&b, `<line x1="%s" y1="%s" x2="%s" y2="%s" stroke="%s" stroke-width="2"></line>`,
				fmtN(mk.X1), fmtN(mk.Y1), fmtN(mk.X2), fmtN(mk.Y2), mk.Color)
		}
		// Title.
		if item.Title != "" {
			fmt.Fprintf(&b, `<g transform="translate(%s,%s)`, fmtN(item.TitleX), fmtN(item.TitleY))
			if props.TitleRotation != 0 {
				fmt.Fprintf(&b, ` rotate(%s)`, fmtN(props.TitleRotation))
			}
			b.WriteString(`">`)
			fmt.Fprintf(&b, `<text dominant-baseline="central" text-anchor="%s"`, textAnchor(props.TitleAlign))
			if fill != "" {
				fmt.Fprintf(&b, ` fill="%s"`, fill)
			}
			if fontSize != "" {
				fmt.Fprintf(&b, ` font-size="%s"`, fontSize)
			}
			if fontFamily != "" {
				fmt.Fprintf(&b, ` font-family="%s"`, fontFamily)
			}
			fmt.Fprintf(&b, ">%s</text></g>", escapeText(item.Title))
		}
		b.WriteString("</g>")
	}
	return b.String()
}

func writeRect(b *strings.Builder, r ComputedRect) {
	fmt.Fprintf(b, `<rect x="%s" y="%s" width="%s" height="%s" fill="%s"></rect>`,
		fmtN(r.X), fmtN(r.Y), fmtN(maxF(r.Width, 0)), fmtN(maxF(r.Height, 0)), r.Color)
}

func textAnchor(align string) string {
	switch align {
	case "start":
		return "start"
	case "end":
		return "end"
	default:
		return "middle"
	}
}

// labelsTextStyle resolves the title text style from the theme labels block,
// falling back to the base text style.
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

func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
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
