package bump

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/bump svgDefaultProps. Fields left zero in a BumpProps
// fall back to these via applyDefaults.
var Defaults = BumpProps{
	Interpolation: InterpolationSmooth,
	XPadding:      0.6,
	XOuterPadding: 0.5,
	YOuterPadding: 0.5,

	LineWidth:         2,
	ActiveLineWidth:   4,
	InactiveLineWidth: 1,
	Opacity:           1,
	ActiveOpacity:     1,
	InactiveOpacity:   0.3,

	StartLabel: BoolPtr(false),
	EndLabel:   BoolPtr(true),

	PointSize:         6,
	ActivePointSize:   8,
	InactivePointSize: 4,

	Colors: colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},

	EnableGridX: true,
	EnableGridY: true,

	Layers:      DefaultLayers,
	Role:        "img",
	MotionProps: core.MotionProps{Animate: true, MotionConfig: "default"},
}
