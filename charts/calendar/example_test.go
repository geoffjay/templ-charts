package calendar_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleCalendar renders a small day-grid calendar heatmap to an SVG string
// with the charts/render helper.
func ExampleCalendar() {
	svg, err := render.String(calendar.Calendar(calendar.CalendarProps{
		Width:  600,
		Height: 200,
		From:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		Data: []calendar.CalendarDatum{
			{Day: "2024-01-05", Value: 12},
			{Day: "2024-03-15", Value: 48},
			{Day: "2024-07-20", Value: 73},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
