package demos

import (
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/core"
)

// CalendarDemo is one calendar tile on the /calendar page (static render).
type CalendarDemo struct {
	ID          string
	Title       string
	Description string
	Props       calendar.CalendarProps
}

// CalendarDemos returns the calendar demos for the /calendar page.
func CalendarDemos() []CalendarDemo {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	data := calendarData(from, to)
	return []CalendarDemo{
		{
			ID:          "calendar-year",
			Title:       "Year (horizontal)",
			Description: "A full 2024 with a 4-bucket quantized color scale, month + year legends.",
			Props: calendar.CalendarProps{
				Width: commonChartWidth, Height: 220,
				Margin: core.Margin{Top: 40, Right: 20, Bottom: 10, Left: 40},
				From:   from, To: to, Data: data,
				Interactive: true,
			},
		},
		{
			ID:          "calendar-vertical",
			Title:       "Vertical + custom colors",
			Description: "direction=vertical with a green sequential-style 5-color palette.",
			Props: calendar.CalendarProps{
				Width: commonChartWidth, Height: 520,
				Margin: core.Margin{Top: 30, Right: 20, Bottom: 10, Left: 40},
				From:   from, To: to, Data: data,
				Direction: calendar.DirectionVertical,
				Colors:    []string{"#d3f2a3", "#97e196", "#6cc08b", "#4c9b82", "#217a79"},
			},
		},
		{
			ID:          "calendar-month-border",
			Title:       "Month outlines",
			Description: "monthBorderWidth=2 draws a per-month outline-path border around each month's day cells.",
			Props: calendar.CalendarProps{
				Width: commonChartWidth, Height: 220,
				Margin: core.Margin{Top: 40, Right: 20, Bottom: 10, Left: 40},
				From:   from, To: to, Data: data,
				MonthBorderWidth: 2,
				MonthBorderColor: "#000000",
			},
		},
	}
}

// calendarData generates a deterministic value for ~40% of the days in range.
func calendarData(from, to time.Time) []calendar.CalendarDatum {
	var out []calendar.CalendarDatum
	i := 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		i++
		if i%5 == 0 || i%7 == 0 {
			out = append(out, calendar.CalendarDatum{
				Day:   d.Format("2006-01-02"),
				Value: float64((i*13)%100 + 1),
			})
		}
	}
	return out
}
