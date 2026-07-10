package calendar

import "github.com/geoffjay/templ-charts/charts/core"

// Defaults mirrors @nivo/calendar commonDefaultProps. Fields left zero in a
// CalendarProps fall back to these via applyDefaults.
var Defaults = CalendarProps{
	Colors:         []string{"#61cdbb", "#97e3d5", "#e8c1a0", "#f47560"},
	Direction:      DirectionHorizontal,
	EmptyColor:     "#ffffff",
	YearSpacing:    30,
	MonthSpacing:   0,
	DaySpacing:     0,
	DayBorderWidth: 1,
	DayBorderColor: "#000000",
	// nivo's default monthBorderWidth is 2, but we default to 0 (outline off)
	// so the default render stays byte-identical to prior goldens. Set
	// MonthBorderWidth > 0 explicitly to enable the per-month outline.
	MonthBorderColor:   "#000000",
	MonthBorderWidth:   0,
	EnableMonthLegends: core.BoolPtr(true),
	EnableYearLegends:  core.BoolPtr(true),
	Role:               "img",
}
