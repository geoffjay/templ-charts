package calendar

// Defaults mirrors @nivo/calendar commonDefaultProps. Fields left zero in a
// CalendarProps fall back to these via applyDefaults.
var Defaults = CalendarProps{
	Colors:             []string{"#61cdbb", "#97e3d5", "#e8c1a0", "#f47560"},
	Direction:          DirectionHorizontal,
	EmptyColor:         "#ffffff",
	YearSpacing:        30,
	MonthSpacing:       0,
	DaySpacing:         0,
	DayBorderWidth:     1,
	DayBorderColor:     "#000000",
	EnableMonthLegends: BoolPtr(true),
	EnableYearLegends:  BoolPtr(true),
	Role:               "img",
}
