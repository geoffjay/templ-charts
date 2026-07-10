package radar

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/radar svgDefaultProps. Fields left zero in a
// RadarProps fall back to these via applyDefaults.
var Defaults = RadarProps{
	Layers:          DefaultLayers,
	Rotation:        0,
	Curve:           core.CurveLinearClosed,
	BorderWidth:     2,
	BorderColor:     colors.NewFromContextColor("color", nil),
	GridLevels:      5,
	GridShape:       GridShapeCircular,
	GridLabelOffset: 16,
	EnableDots:      core.BoolPtr(true),
	DotSize:         6,
	DotColor:        colors.NewFromContextColor("color", nil),
	DotBorderWidth:  0,
	DotBorderColor:  colors.NewFromContextColor("color", nil),
	EnableDotLabel:  core.BoolPtr(false),
	DotLabelYOffset: -12,
	Colors:          colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	FillOpacity:     0.25,
	Legends:         nil,
	Role:            "img",
}
