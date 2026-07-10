// Package calendar mirrors @nivo/calendar: a day-grid heatmap over a date
// range, with weeks as columns (horizontal) or rows (vertical), per-day color
// from a quantized value scale, and month/year legends. Reuses charts/core,
// charts/theming, and the ported date math (no d3-time dependency).
//
// SVG only, static render. Day cells, day borders, month/year
// text legends, and the per-month outline-path border are supported.
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

// GetX/GetY/GetWidth/GetHeight let ComputedDay satisfy grid.GridCellLike, so a
// month's day cells can be merged into an outline polygon via
// grid.GetCellsPolygons.
func (d ComputedDay) GetX() float64      { return d.X }
func (d ComputedDay) GetY() float64      { return d.Y }
func (d ComputedDay) GetWidth() float64  { return d.Size }
func (d ComputedDay) GetHeight() float64 { return d.Size }

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
	// Width and Height are the total SVG dimensions in pixels; the grid area is
	// these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the grid area (for legends).
	Margin core.Margin
	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool
	// Data holds one entry per day with a value; days are keyed "YYYY-MM-DD".
	Data []CalendarDatum

	// From/To bound the rendered range. When zero, derived from the data days.
	From time.Time
	To   time.Time

	// Direction lays out weeks as columns ("horizontal", default) or rows
	// ("vertical").
	Direction CalendarDirection

	// Interactive enables the client-side hover layer (charts/interact): each
	// data day emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Colors is the quantize palette (value → bucket → color).
	Colors []string
	// EmptyColor is the fill for days with no data; default "#ffffff".
	EmptyColor string

	// MinValue/MaxValue bound the color domain; nil Min → 0, nil Max → data max.
	MinValue *float64
	MaxValue *float64

	// YearSpacing is the gap in pixels between year blocks; default 30.
	// MonthSpacing and DaySpacing are the gaps in pixels between months and
	// between day cells; both default 0.
	YearSpacing  float64
	MonthSpacing float64
	DaySpacing   float64

	// DayBorderWidth is the day cell border width in pixels; default 1.
	// DayBorderColor is that border's color; default "#000000".
	DayBorderWidth float64
	DayBorderColor string

	// MonthBorderColor / MonthBorderWidth draw a per-month outline-path border
	// around each calendar month's day cells (see renderMonthBordersLayer). The
	// outline is emitted only when MonthBorderWidth > 0; it defaults to 0 (off)
	// so the default render carries no month outline. Mirrors @nivo/calendar's
	// monthBorderColor / monthBorderWidth.
	MonthBorderColor string
	MonthBorderWidth float64

	// EnableMonthLegends / EnableYearLegends gate the text legends. nil → true
	// (nivo default); pass a pointer to false to disable. Pointer-typed so the
	// default-on behavior survives Go's bool zero value.
	EnableMonthLegends *bool
	EnableYearLegends  *bool

	// Theme overrides chart styling; nil → theming.DefaultTheme.
	Theme *theming.Theme
	// Role is the root SVG ARIA role; default "img".
	Role string
	// AriaLabel, AriaLabelledBy, and AriaDescribedBy set the corresponding SVG
	// accessibility attributes.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title>/<desc> elements.
	Title string
	Desc  string
	// IsFocusable makes the SVG keyboard-focusable.
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// day cell (600ms). MotionStagger delays successive cells by that many
	// seconds (0 ⇒ all cells enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/calendar's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// MonthLegendsEnabled resolves EnableMonthLegends: unset (nil) means true.
func (p CalendarProps) MonthLegendsEnabled() bool {
	return p.EnableMonthLegends == nil || *p.EnableMonthLegends
}

// YearLegendsEnabled resolves EnableYearLegends: unset (nil) means true.
func (p CalendarProps) YearLegendsEnabled() bool {
	return p.EnableYearLegends == nil || *p.EnableYearLegends
}

// CalendarResult is the computed model produced by UseCalendar.
type CalendarResult struct {
	Days         []ComputedDay
	MonthLegends []MonthLegend
	YearLegends  []YearLegend
	CellSize     float64
	MinValue     float64
	MaxValue     float64
}
