package calendar

import "github.com/geoffjay/templ-charts/charts/core"

// Defaults mirrors @nivo/calendar commonDefaultProps. Fields left zero in a
// CalendarProps fall back to these via applyDefaults.
var Defaults = CalendarProps{
	Colors:             []string{"#61cdbb", "#97e3d5", "#e8c1a0", "#f47560"},
	Direction:          DirectionHorizontal,
	EmptyColor:         "#ffffff",
	YearSpacing:        30,
	MonthSpacing:       0,
	DaySpacing:         0,
	MonthBorderWidth:   2,
	MonthBorderColor:   "#000000",
	DayBorderWidth:     1,
	DayBorderColor:     "#000000",
	EnableMonthLegends: BoolPtr(true),
	EnableYearLegends:  BoolPtr(true),
	Role:               "img",
	MotionProps:        core.MotionProps{Animate: true, MotionConfig: core.DefaultMotionConfig},
}
