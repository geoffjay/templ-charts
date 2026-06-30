package calendar

import (
	"context"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/interact"
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

// applyDefaults fills zero-valued CalendarProps fields from Defaults.
func applyDefaults(p CalendarProps) CalendarProps {
	if len(p.Colors) == 0 {
		p.Colors = Defaults.Colors
	}
	if p.Direction == "" {
		p.Direction = Defaults.Direction
	}
	if p.EmptyColor == "" {
		p.EmptyColor = Defaults.EmptyColor
	}
	if p.YearSpacing == 0 {
		p.YearSpacing = Defaults.YearSpacing
	}
	if p.MonthBorderColor == "" {
		p.MonthBorderColor = Defaults.MonthBorderColor
	}
	if p.DayBorderColor == "" {
		p.DayBorderColor = Defaults.DayBorderColor
	}
	if p.DayBorderWidth == 0 {
		p.DayBorderWidth = Defaults.DayBorderWidth
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderLayers renders the day cells + month/year legends as an inner SVG string.
func renderLayers(props CalendarProps, result CalendarResult, theme *theming.Theme) string {
	var b strings.Builder
	b.WriteString(renderDaysLayer(props, result))
	if props.MonthLegendsEnabled() {
		b.WriteString(renderLegendText(result.MonthLegends, monthLegendItems, theme))
	}
	if props.YearLegendsEnabled() {
		b.WriteString(renderYearLegends(result.YearLegends, theme))
	}
	return b.String()
}

func renderDaysLayer(props CalendarProps, result CalendarResult) string {
	var s strings.Builder
	for _, day := range result.Days {
		cp := CalendarDayProps{
			Day:         day,
			BorderWidth: props.DayBorderWidth,
			BorderColor: props.DayBorderColor,
		}
		if props.Interactive && day.HasData && day.Value != nil {
			cp.Tooltip = interact.TooltipHTML(day.Color, day.Day, strconv.FormatFloat(*day.Value, 'g', -1, 64))
		}
		s.WriteString(renderComponent(CalendarDay(cp)))
	}
	return s.String()
}

// monthLegendItems / yearLegend rendering share a tiny text helper.
type legendText struct {
	X, Y  float64
	Label string
}

func monthLegendItems(ms []MonthLegend) []legendText {
	out := make([]legendText, len(ms))
	for i, m := range ms {
		out[i] = legendText{X: m.X, Y: m.Y, Label: m.Label}
	}
	return out
}

func renderLegendText(ms []MonthLegend, conv func([]MonthLegend) []legendText, theme *theming.Theme) string {
	var s strings.Builder
	fill := theme.Labels.Text.Fill
	for _, lt := range conv(ms) {
		s.WriteString(renderComponent(CalendarLegendLabel(CalendarLegendLabelProps{
			X: lt.X, Y: lt.Y, Label: lt.Label, Fill: fill,
		})))
	}
	return s.String()
}

func renderYearLegends(ys []YearLegend, theme *theming.Theme) string {
	var s strings.Builder
	fill := theme.Labels.Text.Fill
	for _, y := range ys {
		s.WriteString(renderComponent(CalendarLegendLabel(CalendarLegendLabelProps{
			X: y.X, Y: y.Y, Label: yearLabel(y.Year), Fill: fill, Bold: true,
		})))
	}
	return s.String()
}

func yearLabel(y int) string {
	// avoid strconv import churn; small int → string
	return itoa(y)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
