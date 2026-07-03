package bullet

// Defaults mirrors @nivo/bullet defaultProps. Fields left zero in a BulletProps
// fall back to these via applyDefaults.
var Defaults = BulletProps{
	Layout:             BulletLayoutHorizontal,
	Reverse:            false,
	Spacing:            30,
	AxisPosition:       "after",
	TitlePosition:      "before",
	TitleAlign:         "middle",
	RangeColors:        "seq:cool",
	MeasureColors:      "seq:red_purple",
	MarkerColors:       "seq:red_purple",
	RangeBorderWidth:   0,
	MeasureSize:        0.4,
	MeasureBorderWidth: 0,
	MarkerSize:         0.6,
	Role:               "img",
}
