package calendar

import (
	"context"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
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
	if p.DayBorderColor == "" {
		p.DayBorderColor = Defaults.DayBorderColor
	}
	if p.DayBorderWidth == 0 {
		p.DayBorderWidth = Defaults.DayBorderWidth
	}
	if p.MonthBorderColor == "" {
		p.MonthBorderColor = Defaults.MonthBorderColor
	}
	if p.MonthBorderWidth == 0 {
		p.MonthBorderWidth = Defaults.MonthBorderWidth
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
	// Month outline sits on top of the day cells.
	b.WriteString(renderMonthBordersLayer(props, result))
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
	for i, day := range result.Days {
		cp := CalendarDayProps{
			Day:         day,
			BorderWidth: props.DayBorderWidth,
			BorderColor: props.DayBorderColor,
			Animate:     props.Animate,
		}
		if props.Animate {
			cp.AnimateBegin = core.StaggerBegin(i, props.MotionStagger)
		}
		if props.Interactive && day.HasData && day.Value != nil {
			cp.Tooltip = interact.TooltipHTML(day.Color, day.Day, strconv.FormatFloat(*day.Value, 'g', -1, 64))
		}
		s.WriteString(renderComponent(CalendarDay(cp)))
	}
	return s.String()
}

// renderMonthBordersLayer draws a per-month outline-path border around each
// calendar month's day cells. It runs only when props.MonthBorderWidth > 0.
// Days are grouped by calendar month (year+month), each group's contiguous
// blocks are merged into outline polygons via grid.GetCellsPolygons, and each
// month becomes one <path> whose d concatenates every polygon as a closed
// subpath ("M … L … Z"). Mirrors @nivo/calendar's month outline path.
func renderMonthBordersLayer(props CalendarProps, result CalendarResult) string {
	if props.MonthBorderWidth <= 0 {
		return ""
	}

	// Group days by calendar month, preserving first-seen order (days are
	// already chronological).
	type monthKey struct {
		year  int
		month int
	}
	var order []monthKey
	groups := make(map[monthKey][]ComputedDay)
	for _, d := range result.Days {
		k := monthKey{year: d.Date.Year(), month: int(d.Date.Month())}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], d)
	}

	var s strings.Builder
	for _, k := range order {
		polygons := grid.GetCellsPolygons(groups[k])
		var d strings.Builder
		for _, poly := range polygons {
			if len(poly) == 0 {
				continue
			}
			for i, v := range poly {
				if i == 0 {
					d.WriteString("M")
				} else {
					d.WriteString(" L")
				}
				d.WriteString(fmtC(v[0]))
				d.WriteString(",")
				d.WriteString(fmtC(v[1]))
			}
			d.WriteString(" Z")
		}
		if d.Len() == 0 {
			continue
		}
		s.WriteString(`<path fill="none" stroke="`)
		s.WriteString(props.MonthBorderColor)
		s.WriteString(`" stroke-width="`)
		s.WriteString(fmtC(props.MonthBorderWidth))
		s.WriteString(`" d="`)
		s.WriteString(d.String())
		s.WriteString(`"></path>`)
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
