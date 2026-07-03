package icicle

import (
	"github.com/geoffjay/templ-charts/charts/colors"
)

// Defaults mirrors @nivo/icicle svgDefaultProps.
var Defaults = IcicleProps{
	Orientation:  OrientationBottom,
	GapX:         1,
	GapY:         1,
	Colors:       colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderWidth:  0,
	BorderRadius: 0,
	EnableLabels: BoolPtr(false),
	Label:        "id",
	Role:         "img",
}
