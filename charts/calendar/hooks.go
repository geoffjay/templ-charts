package calendar

import (
	"math"
	"time"

	"github.com/geoffjay/templ-charts/charts/theming"
)

// UseCalendar mirrors @nivo/calendar's compute: it derives the date range, the
// per-year cell size, and the positioned/colored day cells plus month and year
// legends. props.Width/Height are the inner dimensions.
func UseCalendar(props CalendarProps) CalendarResult {
	from, to := resolveRange(props)
	if from.IsZero() || to.IsZero() || to.Before(from) {
		return CalendarResult{}
	}

	// Index data by day for O(1) lookup; track the value domain.
	byDay := make(map[string]float64, len(props.Data))
	minV := math.Inf(1)
	maxV := math.Inf(-1)
	for _, d := range props.Data {
		byDay[d.Day] = d.Value
		if d.Value < minV {
			minV = d.Value
		}
		if d.Value > maxV {
			maxV = d.Value
		}
	}
	if len(byDay) == 0 {
		minV, maxV = 0, 0
	}
	domainMin := 0.0
	if props.MinValue != nil {
		domainMin = *props.MinValue
	}
	domainMax := maxV
	if props.MaxValue != nil {
		domainMax = *props.MaxValue
	}

	yearRange := to.Year() - from.Year() + 1
	maxWeeks := 0
	for y := from.Year(); y <= to.Year(); y++ {
		if w := weeksInYear(y); w > maxWeeks {
			maxWeeks = w
		}
	}

	cellSize := computeCellSize(props.Width, props.Height, props.Direction, yearRange, maxWeeks,
		props.YearSpacing, props.MonthSpacing, props.DaySpacing)

	colorFor := quantizeColor(props.Colors, domainMin, domainMax)

	days := make([]ComputedDay, 0, 366)
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		yearIndex := d.Year() - from.Year()
		x, y := cellPosition(props.Direction, d, yearIndex, cellSize, props.YearSpacing, props.MonthSpacing, props.DaySpacing)
		cell := ComputedDay{Date: d, Day: key, X: x, Y: y, Size: cellSize, Color: props.EmptyColor}
		if v, ok := byDay[key]; ok {
			cell.Value = &v
			cell.HasData = true
			cell.Color = colorFor(v)
		}
		days = append(days, cell)
	}

	var monthLegends []MonthLegend
	if props.MonthLegendsEnabled() {
		monthLegends = computeMonthLegends(props, from, to, cellSize)
	}
	var yearLegends []YearLegend
	if props.YearLegendsEnabled() {
		yearLegends = computeYearLegends(props, from, to, cellSize)
	}

	return CalendarResult{
		Days:         days,
		MonthLegends: monthLegends,
		YearLegends:  yearLegends,
		CellSize:     cellSize,
		MinValue:     domainMin,
		MaxValue:     domainMax,
	}
}

// resolveRange returns the [from, to] range: explicit props, else derived from
// the data days (min/max).
func resolveRange(props CalendarProps) (time.Time, time.Time) {
	from, to := props.From, props.To
	if !from.IsZero() && !to.IsZero() {
		return dayFloor(from), dayFloor(to)
	}
	var minD, maxD time.Time
	for _, d := range props.Data {
		t, err := time.Parse("2006-01-02", d.Day)
		if err != nil {
			continue
		}
		if minD.IsZero() || t.Before(minD) {
			minD = t
		}
		if maxD.IsZero() || t.After(maxD) {
			maxD = t
		}
	}
	if from.IsZero() {
		from = minD
	}
	if to.IsZero() {
		to = maxD
	}
	return dayFloor(from), dayFloor(to)
}

func dayFloor(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// weekOfYear mirrors d3's timeWeek.count(timeYear(d), d): the Sunday-based week
// index of d within its year (week 0 is the partial week before the first
// Sunday).
func weekOfYear(d time.Time) int {
	jan1 := time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
	jan1Weekday := int(jan1.Weekday()) // Sun=0
	dayOfYear := d.YearDay()           // 1-based
	return (dayOfYear - 1 + jan1Weekday) / 7
}

// weeksInYear returns the number of Sunday-aligned week rows a year spans.
func weeksInYear(year int) int {
	dec31 := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	return weekOfYear(dec31) + 1
}

// computeCellSize mirrors @nivo/calendar computeCellSize.
func computeCellSize(width, height float64, direction CalendarDirection, yearRange, maxWeeks int, yearSpacing, monthSpacing, daySpacing float64) float64 {
	var hCell, vCell float64
	yr := float64(yearRange)
	mw := float64(maxWeeks)
	if direction == DirectionVertical {
		hCell = (width - (yr-1)*yearSpacing - yr*(8*daySpacing)) / (yr * 7)
		vCell = (height - monthSpacing*12 - daySpacing*mw) / mw
	} else {
		hCell = (width - monthSpacing*12 - daySpacing*mw) / mw
		vCell = (height - (yr-1)*yearSpacing - yr*(8*daySpacing)) / (yr * 7)
	}
	return math.Max(0, math.Min(hCell, vCell))
}

// cellPosition mirrors @nivo/calendar cellPositionHorizontal/Vertical.
func cellPosition(direction CalendarDirection, d time.Time, yearIndex int, cellSize, yearSpacing, monthSpacing, daySpacing float64) (float64, float64) {
	week := float64(weekOfYear(d))
	weekday := float64(int(d.Weekday()))
	month := float64(int(d.Month()) - 1)
	step := cellSize + daySpacing
	yearOffset := float64(yearIndex) * (yearSpacing + 7*step)
	if direction == DirectionVertical {
		x := weekday*step + daySpacing/2 + yearOffset
		y := week*step + daySpacing/2 + month*monthSpacing
		return x, y
	}
	x := week*step + daySpacing/2 + month*monthSpacing
	y := weekday*step + daySpacing/2 + yearOffset
	return x, y
}

// quantizeColor returns a value→color function that buckets [min,max] across
// the colors slice (d3 scaleQuantize semantics).
func quantizeColor(palette []string, min, max float64) func(float64) string {
	if len(palette) == 0 {
		return func(float64) string { return "#000000" }
	}
	n := len(palette)
	span := max - min
	return func(v float64) string {
		if span <= 0 {
			return palette[0]
		}
		idx := int(math.Floor((v - min) / span * float64(n)))
		if idx < 0 {
			idx = 0
		}
		if idx >= n {
			idx = n - 1
		}
		return palette[idx]
	}
}

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// computeMonthLegends places one label per month at the position of its 1st day.
func computeMonthLegends(props CalendarProps, from, to time.Time, cellSize float64) []MonthLegend {
	var out []MonthLegend
	d := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	for !d.After(to) {
		if !d.Before(from) {
			yearIndex := d.Year() - from.Year()
			x, y := cellPosition(props.Direction, d, yearIndex, cellSize, props.YearSpacing, props.MonthSpacing, props.DaySpacing)
			label := monthNames[int(d.Month())-1]
			if props.Direction == DirectionVertical {
				out = append(out, MonthLegend{Year: d.Year(), Month: int(d.Month()), X: x - 8, Y: y + cellSize/2, Label: label})
			} else {
				out = append(out, MonthLegend{Year: d.Year(), Month: int(d.Month()), X: x + cellSize/2, Y: y - 6, Label: label})
			}
		}
		d = d.AddDate(0, 1, 0)
	}
	return out
}

// computeYearLegends places one label per year beside its block.
func computeYearLegends(props CalendarProps, from, to time.Time, cellSize float64) []YearLegend {
	var out []YearLegend
	step := cellSize + props.DaySpacing
	for y := from.Year(); y <= to.Year(); y++ {
		yearIndex := y - from.Year()
		yearOffset := float64(yearIndex) * (props.YearSpacing + 7*step)
		if props.Direction == DirectionVertical {
			out = append(out, YearLegend{Year: y, X: yearOffset + 3.5*step, Y: -18})
		} else {
			out = append(out, YearLegend{Year: y, X: -18, Y: yearOffset + 3.5*step})
		}
	}
	return out
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
