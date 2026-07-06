package entries

import (
	"math"
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// calendarDemoData builds a deterministic year of daily "activity" values:
// a weekly rhythm (quiet weekends), slow seasonal swells, and occasional
// gaps, so the sequential color scale's full range is exercised.
func calendarDemoData(from time.Time, days int) []calendar.CalendarDatum {
	data := make([]calendar.CalendarDatum, 0, days)
	for i := 0; i < days; i++ {
		d := from.AddDate(0, 0, i)
		if (i*13)%17 == 0 {
			continue // leave some days empty, like real activity data
		}
		v := 38 + 26*math.Sin(float64(i)/28.0) + 14*math.Sin(float64(i)/5.2) + float64((i*7)%11)
		switch d.Weekday() {
		case time.Saturday, time.Sunday:
			v *= 0.4
		}
		if v < 0 {
			v = 0
		}
		data = append(data, calendar.CalendarDatum{Day: d.Format("2006-01-02"), Value: math.Round(v)})
	}
	return data
}

func init() {
	register(ChartEntry{
		Slug:        "calendar",
		Title:       "Calendar",
		Description: "Day-grid heatmap over a date range with quantized colors.",
		Snippet: `import (
    "time"

    "github.com/geoffjay/templ-charts/charts/calendar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(calendar.Calendar(calendar.CalendarProps{
    Width: 720, Height: 440,
    From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    To:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
    Data: []calendar.CalendarDatum{
        {Day: "2024-01-05", Value: 12},
        {Day: "2024-03-15", Value: 48},
        {Day: "2024-07-20", Value: 73},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := calendar.CalendarProps{
				Width: 720, Height: 440, Responsive: true,
				From:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
				Data:  calendarDemoData(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 366),
				Theme: theme,
			}
			// Colors is an explicit quantize list: sample gradient palettes into
			// 6 buckets; categorical picks keep the default ramp (arbitrary hues
			// make a meaningless value scale).
			if pal, ok := colors.LookupPalette(palette); ok && pal.Kind != colors.KindCategorical {
				p.Colors = pal.Swatch(6)
			}
			p.Animate = animate
			return render.String(calendar.Calendar(p))
		},
	})
}
