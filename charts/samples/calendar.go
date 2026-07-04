package samples

import (
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
)

// Calendar returns a full-year (2024) sample for the calendar chart: the from
// and to bounds plus a deterministic value on ~40% of the days in range. The
// three returned values map directly onto CalendarProps.From/To/Data.
func Calendar() (from, to time.Time, data []calendar.CalendarDatum) {
	from = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	i := 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		i++
		if i%5 == 0 || i%7 == 0 {
			data = append(data, calendar.CalendarDatum{
				Day:   d.Format("2006-01-02"),
				Value: float64((i*13)%100 + 1),
			})
		}
	}
	return from, to, data
}
