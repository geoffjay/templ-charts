// Package calendar mirrors @nivo/calendar: a day-grid heatmap over a date
// range, with weeks as columns (horizontal) or rows (vertical), per-day color
// from a quantized value scale, and month/year legends. Reuses charts/core,
// charts/theming, and the ported date math (no d3-time dependency).
//
// v2 scope: SVG only, static render. Day cells, day borders, and month/year
// text legends are supported; the month outline-path border is deferred.
package calendar

import (
	"time"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// CalendarDatum is one day's value. Day is "YYYY-MM-DD". Mirrors @nivo/calendar
// datum.
type CalendarDatum struct {
	Day   string
	Value float64
}

// ComputedDay is a positioned, colored day cell.
type ComputedDay struct {
	Date    time.Time
	Day     string // YYYY-MM-DD
	X       float64
	Y       float64
	Size    float64
	Color   string
	Value   *float64 // nil when no data for the day
	HasData bool
}

// MonthLegend is a month label positioned over its columns/rows.
type MonthLegend struct {
	Year  int
	Month int // 1..12
	X     float64
	Y     float64
	Label string
}

// YearLegend is a year label positioned beside its block.
type YearLegend struct {
	Year int
	X    float64
	Y    float64
}

// CalendarDirection is "horizontal" (weeks → columns) or "vertical".
type CalendarDirection string

const (
	DirectionHorizontal CalendarDirection = "horizontal"
	DirectionVertical   CalendarDirection = "vertical"
)

// CalendarProps is the input to the Calendar component. Mirrors @nivo/calendar
// commonDefaultProps (supported subset).
type CalendarProps struct {
	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool
	Data       []CalendarDatum

	// From/To bound the rendered range. When zero, derived from the data days.
	From time.Time
	To   time.Time

	Direction CalendarDirection

	// Colors is the quantize palette (value → bucket → color).
	Colors     []string
	EmptyColor string

	// MinValue/MaxValue bound the color domain; nil Min → 0, nil Max → data max.
	MinValue *float64
	MaxValue *float64

	YearSpacing  float64
	MonthSpacing float64
	DaySpacing   float64

	MonthBorderWidth float64
	MonthBorderColor string
	DayBorderWidth   float64
	DayBorderColor   string

	// EnableMonthLegends / EnableYearLegends gate the text legends. nil → true
	// (nivo default); pass a pointer to false to disable. Pointer-typed so the
	// default-on behavior survives Go's bool zero value.
	EnableMonthLegends *bool
	EnableYearLegends  *bool

	Theme       *theming.Theme
	Role        string
	IsFocusable bool

	core.MotionProps
}

// MonthLegendsEnabled resolves EnableMonthLegends: unset (nil) means true.
func (p CalendarProps) MonthLegendsEnabled() bool {
	return p.EnableMonthLegends == nil || *p.EnableMonthLegends
}

// YearLegendsEnabled resolves EnableYearLegends: unset (nil) means true.
func (p CalendarProps) YearLegendsEnabled() bool {
	return p.EnableYearLegends == nil || *p.EnableYearLegends
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// CalendarResult is the computed model produced by UseCalendar.
type CalendarResult struct {
	Days         []ComputedDay
	MonthLegends []MonthLegend
	YearLegends  []YearLegend
	CellSize     float64
	MinValue     float64
	MaxValue     float64
}
