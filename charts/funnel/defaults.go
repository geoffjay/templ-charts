package funnel

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/funnel svgDefaultProps. Fields left zero in a
// FunnelProps fall back to these via applyDefaults.
var Defaults = FunnelProps{
	Layers:                 DefaultLayers,
	Direction:              FunnelDirectionVertical,
	Interpolation:          FunnelInterpolationSmooth,
	Spacing:                0,
	ShapeBlending:          0.66,
	Colors:                 colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	FillOpacity:            1,
	BorderWidth:            6,
	BorderColor:            colors.NewFromContextColor("color", nil),
	BorderOpacity:          0.66,
	EnableLabel:            core.BoolPtr(true),
	LabelColor:             colors.NewThemeColor("background"),
	EnableBeforeSeparators: core.BoolPtr(true),
	EnableAfterSeparators:  core.BoolPtr(true),
	Role:                   "img",
}
