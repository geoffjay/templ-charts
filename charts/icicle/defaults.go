package icicle

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/icicle svgDefaultProps.
var Defaults = IcicleProps{
	Orientation:   OrientationBottom,
	GapX:          1,
	GapY:          1,
	Colors:        colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderWidth:   0,
	BorderRadius:  0,
	EnableLabels:  core.BoolPtr(false),
	EnableZooming: core.BoolPtr(true),
	Label:         "id",
	Role:          "img",
}
