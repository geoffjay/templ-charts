package entries

import (
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

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
				From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
				Data: []calendar.CalendarDatum{
					{Day: "2024-01-05", Value: 12},
					{Day: "2024-03-15", Value: 48},
					{Day: "2024-07-20", Value: 73},
				},
				Theme: theme,
			}
			p.Animate = animate
			return render.String(calendar.Calendar(p))
		},
	})
}
